package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	// Specific API routes
	mux.HandleFunc("/admin/secret", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiKey := os.Getenv("API_KEY")
		if r.Header.Get("X-API-KEY") != strings.TrimSpace(apiKey) {
			http.Error(w, "Wrong password", http.StatusForbidden)
			return
		}
		fmt.Fprintf(w, "This is a secret message from the admin service!")
	})

	mux.HandleFunc("/admin/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong from admin-service")
	})

	mux.HandleFunc("/admin/debug", func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")
		fmt.Fprintf(w, "API_KEY: [%s]", apiKey)
	})

	// Root handler for the single-page application
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "admin.html")
	})

	log.Println("Admin service starting on port 8082...")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
