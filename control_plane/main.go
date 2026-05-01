package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleWorkloadSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Workload submitted successfully\n")
}

func main() {
	http.HandleFunc("/api/v1/workloads", handleWorkloadSubmission)
	
	port := 8080
	log.Printf("Control Plane API starting on port %d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
