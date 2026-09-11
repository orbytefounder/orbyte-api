package traffic

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Config controls concurrency limits and basic rate shaping.
type Config struct {
	MaxConcurrent int
	RequestsPerSec float64
	Burst          int
}

// DefaultConfig returns production-friendly traffic limits.
func DefaultConfig() Config {
	return Config{
		MaxConcurrent:  1024,
		RequestsPerSec: 500,
		Burst:          100,
	}
}

// Engine limits in-flight work and smooths request bursts.
type Engine struct {
	semaphore chan struct{}

	mu          sync.Mutex
	tokens      float64
	maxTokens   float64
	refillRate  float64
	lastRefill  time.Time

	activeRequests int64
	totalAccepted  int64
	totalRejected  int64
}

// NewEngine creates a traffic controller from cfg.
func NewEngine(cfg Config) *Engine {
	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = 1024
	}
	if cfg.RequestsPerSec <= 0 {
		cfg.RequestsPerSec = 500
	}
	if cfg.Burst <= 0 {
		cfg.Burst = int(cfg.RequestsPerSec)
	}

	return &Engine{
		semaphore:  make(chan struct{}, cfg.MaxConcurrent),
		tokens:     float64(cfg.Burst),
		maxTokens:  float64(cfg.Burst),
		refillRate: cfg.RequestsPerSec,
		lastRefill: time.Now(),
	}
}

// Middleware wraps handlers with concurrency + rate-limit gates.
func (e *Engine) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !e.allow() {
			atomic.AddInt64(&e.totalRejected, 1)
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		select {
		case e.semaphore <- struct{}{}:
		default:
			atomic.AddInt64(&e.totalRejected, 1)
			http.Error(w, "server busy", http.StatusServiceUnavailable)
			return
		}

		atomic.AddInt64(&e.activeRequests, 1)
		atomic.AddInt64(&e.totalAccepted, 1)
		defer func() {
			<-e.semaphore
			atomic.AddInt64(&e.activeRequests, -1)
		}()

		next.ServeHTTP(w, r)
	})
}

func (e *Engine) allow() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(e.lastRefill).Seconds()
	e.lastRefill = now

	e.tokens += elapsed * e.refillRate
	if e.tokens > e.maxTokens {
		e.tokens = e.maxTokens
	}

	if e.tokens < 1 {
		return false
	}

	e.tokens--
	return true
}

// Stats exposes live counters for health/metrics endpoints.
func (e *Engine) Stats() (active, accepted, rejected int64) {
	return atomic.LoadInt64(&e.activeRequests),
		atomic.LoadInt64(&e.totalAccepted),
		atomic.LoadInt64(&e.totalRejected)
}
