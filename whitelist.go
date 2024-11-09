package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type SidraRequest struct {
	Headers map[string]string `json:"Headers"`
	Url     string            `json:"Url"`
	Method  string            `json:"Method"`
	Body    string            `json:"Body"`
}

type SidraResponse struct {
	StatusCode int    `json:"StatusCode"`
	Body       string `json:"Body"`
}

// Daftar IP yang diizinkan
var allowedIPs = map[string]bool{
	"192.168.1.1": true,
	"192.168.1.2": true,
}

// Fungsi untuk menangani permintaan dengan memeriksa whitelist IP
func whitelistHandler(req SidraRequest) SidraResponse {
	clientIP := req.Headers["X-Real-Ip"]
	if clientIP == "" {
		log.Println("Missing X-Real-IP header")
		return SidraResponse{
			StatusCode: 400,
			Body:       "Missing X-Real-IP header",
		}
	}

	if allowed, exists := allowedIPs[clientIP]; exists && allowed {
		log.Printf("IP allowed: %s", clientIP)
		return SidraResponse{
			StatusCode: 200,
			Body:       "IP allowed",
		}
	}

	log.Printf("IP not allowed: %s", clientIP)
	return SidraResponse{
		StatusCode: 403,
		Body:       "IP not allowed",
	}
}

func main() {
	// Set up the Unix domain socket listener
	socketPath := "/tmp/whitelist.sock"

	// Remove any existing socket file
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to remove existing socket file: %v", err)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Error setting up Unix domain socket: %v", err)
	}
	defer listener.Close()
	log.Printf("Whitelist plugin listening on %s", socketPath)

	// Accept and handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			decoder := json.NewDecoder(conn)
			var req SidraRequest
			if err := decoder.Decode(&req); err != nil {
				log.Printf("Error decoding request: %v", err)
				return
			}
			log.Printf("Received request: %+v\n", req)

			resp := whitelistHandler(req)

			encoder := json.NewEncoder(conn)
			if err := encoder.Encode(resp); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}(conn)
	}
}
