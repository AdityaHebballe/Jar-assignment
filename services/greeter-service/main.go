package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong from greeter-service")
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Service discovery via Kubernetes DNS
		echoSvcHost := os.Getenv("ECHO_SERVICE_HOST")
		if echoSvcHost == "" {
			echoSvcHost = "echo-service"
		}
		echoSvcPort := os.Getenv("ECHO_SERVICE_PORT")
		if echoSvcPort == "" {
			echoSvcPort = "8081"
		}

		payload := map[string]string{"message": "Hello from greeter-service"}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "Failed to create payload for echo-service", http.StatusInternalServerError)
			log.Printf("Error creating payload: %v", err)
			return
		}

		url := fmt.Sprintf("http://%s:%s/echo", echoSvcHost, echoSvcPort)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			http.Error(w, "Failed to call echo-service", http.StatusInternalServerError)
			log.Printf("Error calling echo-service: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response from echo-service", http.StatusInternalServerError)
			log.Printf("Error reading response from echo-service: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	log.Println("Greeter service starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
