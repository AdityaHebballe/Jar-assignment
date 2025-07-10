package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			apiKey := os.Getenv("API_KEY")
			if r.Header.Get("X-API-KEY") != strings.TrimSpace(apiKey) {
				http.Error(w, "Wrong password", http.StatusForbidden)
				return
			}
			fmt.Fprintf(w, "This is a secret message from the admin service!")
			return
		}
		http.ServeFile(w, r, "admin.html")
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong from admin-service")
	})

	log.Println("Admin service starting on port 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
