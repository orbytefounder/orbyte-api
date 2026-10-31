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
	log.Printf("Orbyte Core Gateway initializing live on port :%s...", port)
	if err:=http.ListenAndServe(":"=+port,mux); err != nil {
		log.Fatalf(Server startup failed: %v", err)
				   }
				   }
