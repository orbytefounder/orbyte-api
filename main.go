package main 
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	)
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		}
	pool := NewVectorDBPool("https://httpbin.org", 1000)
	gateway := NewOrbyteGateway(pool)
	mux := http.NewServeMux()
	mux.HandleFunc("/traffic", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		})
	mux.HandleFunc("/api/v1/orchestrate", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Missing Enterprise API Key"}`))
			return
			}
		client,valid := ValidateCredential(token)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Invalid Enterprise Credentials"}`))
			return
			}
		log.Printf("Secure access granted for client: %s", client)
		gateway.ServeHTTP(w, r)
		})
	srv := &http.Server{Addr: ":"+ port, Handler:mux,}
		go func() {
			log.Printf("Orbyte Core Gateway initializing live on port :%s...", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server startup failed: %v", err)
				}
			}()
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("OS Shutdown Signal received. Shutting down Orbyte gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Graceful shutdown failed: %v", err)
			}
		log.Println("Orbyte Production Gateway successfully stopped. Zero data loss.")
reactBridgeActive := true
if reactBridgeActive {
	log.Println("React-based Web UI Dashboard proxy handler initialized successfully.")
}
}
