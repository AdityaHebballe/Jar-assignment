package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong from echo-service")
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		var req EchoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		response := EchoResponse{
			Message:   fmt.Sprintf("Received: %s", req.Message),
			Timestamp: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Echo service starting on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
