package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type SidraRequest struct {
	IP string `json:"ip"`
}

type SidraResponse struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request in proxy...")

		// Parsing IP dari request Postman
		var request SidraRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, `{"status_code": 400, "body": "Failed to decode request."}`, http.StatusBadRequest)
			log.Println("Failed to decode request:", err)
			return
		}
		log.Println("Decoded request:", request)

		// Arahkan ke socket plugin whitelist
		socketPath := "/tmp/whitelist.sock"
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			http.Error(w, `{"status_code": 500, "body": "Failed to connect to socket."}`, http.StatusInternalServerError)
			log.Println("Failed to connect to socket:", err)
			return
		}
		defer conn.Close()
		log.Println("Connected to whitelist socket at", socketPath)

		// Kirim request JSON ke socket whitelist
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(request)
		if err != nil {
			http.Error(w, `{"status_code": 500, "body": "Failed to encode request."}`, http.StatusInternalServerError)
			log.Println("Failed to encode request:", err)
			return
		}
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			http.Error(w, `{"status_code": 500, "body": "Failed to write request to socket."}`, http.StatusInternalServerError)
			log.Println("Failed to write request to socket:", err)
			return
		}
		log.Println("Request sent to whitelist plugin:", request)

		// Baca respons dari socket whitelist
		var response SidraResponse
		respBytes, err := io.ReadAll(conn)
		if err != nil {
			http.Error(w, `{"status_code": 500, "body": "Failed to read response from socket."}`, http.StatusInternalServerError)
			log.Println("Failed to read response from socket:", err)
			return
		}
		err = json.Unmarshal(respBytes, &response)
		if err != nil {
			http.Error(w, `{"status_code": 500, "body": "Failed to parse response from socket."}`, http.StatusInternalServerError)
			log.Println("Failed to parse response from socket:", err)
			return
		}
		log.Println("Response received from whitelist plugin:", response)

		// Kirim respons plugin langsung ke Postman
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Proxy HTTP server running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
