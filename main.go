package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Echo handler - echoes back whatever is POSTed to it
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request to /echo from %s", r.RemoteAddr)

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Only POST method is allowed")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error reading request body: %v", err)
			return
		}
		defer r.Body.Close()

		// Log the received body
		log.Printf("Received body: %s", string(body))

		// Echo back the body
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	// OAuth introspection endpoint
	http.HandleFunc("/introspect", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request to /introspect from %s", r.RemoteAddr)

		// Simple response with active=true
		response := map[string]interface{}{
			"active": true,
			"sub":    "test-user",
			"scope":  "read write",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Root endpoint for basic info
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Echo Server is running. Use /echo for POST requests and /introspect for OAuth introspection.")
	})

	// Start the server
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
