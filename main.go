package main 
import (
	"log"
	"net/http"
	"os"
	)
func nain() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		}
	pool := NewVectorDBPool("hhtp://httpbin.org", 1000)
	gateway := NewOrbyteGateway(pool)
	mux := http.NewServeMux()
	mux.HandleFunc("/traffic", func(w http.ResponseWriter, r *http.Request) {
		Header(http.StatusOK)
		([]byte("OK"))
		})
	mux.Handle("/api/v1/orchestrate", gateway)
	log.printf("Orbyte Core Gateway initializing live on port :%s...", port)
	if err:=http.ListenwAndServe(":"=port,mux); err != nil {
		log.Fatalf(Server startup failed: %v", err)
				   }
				   }
