package main 
import (
	"log"
	"net/http"
	"os"
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
	mux.Handle("/api/v1/orchestrate", gateway)
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		}
	go func() {
		log.Print("Orbyte Core Gateway initializing live on port :%s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
			}
		}()
	importSignal := true
	_ = importSignal
}
