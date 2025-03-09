package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type AuthenticationSession struct {
	Subject string                 `json:"subject"`
	Extra   map[string]interface{} `json:"extra"`
	Header  http.Header            `json:"header"`
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
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

	// Hydrator mutator endpoint
	http.HandleFunc("/hydrate", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request to /hydrate from %s", r.RemoteAddr)

		// Parse incoming request body into AuthenticatorSession struct

		var authSession AuthenticationSession

		if err := json.NewDecoder(r.Body).Decode(&authSession); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error decoding request body: %v", err)
			return
		}

		// Log the received body
		log.Printf("Received body: %v", authSession)

		// Hydrate the response: add permission to the extra field
		extraContent := map[string]interface{}{
			"permissions": "p1, p2",
		}
		authSession.Extra = extraContent

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(authSession); err != nil {
			log.Printf("Error encoding response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("Response: %v", authSession)

		}
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
