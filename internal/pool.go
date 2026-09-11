package internal
import (
	"errors"
	"sync"
	"time"
)
type ConnectionPool struct {
	mu sync.Mutex
	connections chan string
	maxSize int
}
func NewConnectionPool(size int) *ConnectionPool {
	pool := &ConnectionPool{
		connections: make(chan string, size),
		maxSize: size,
	}
	for i := 0; i< size; i++ {
		pool.connections <- "db-connection-active"
	}
	return pool
}
func (p *ConnectionPool) Acquire(timeout time.Duration) (string, error) {
	select {
	case conn := <-p.connections:
		return conn, nil
	case <-time.After(timeout):
		return "", errors.New("database connection timeout: pool is exhausted")
	}
}
func (p *ConnectionPool) Release(conn string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connections <- conn
}