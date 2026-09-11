package main 
import (
	"fmt"
	"net/http"
)
func trafficHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
func main() {
	fmt.Println("--- Orbyte: Orchestration Engine Ready ---")
	http.HandleFunc("/traffic", trafficHandler)
	http.ListenAndServe(":8080", nil)
}