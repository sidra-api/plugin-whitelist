package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Daftar IP yang diizinkan
var allowedIPs = map[string]bool{
	"192.168.1.1": true,
	"192.168.1.2": true,
}

// Struktur request yang diterima
type SidraRequest struct {
	IP string `json:"ip"`
}

// Struktur response yang akan dikirim
type SidraResponse struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

// Fungsi utama untuk menjalankan plugin whitelist
func main() {
	socketPath := "/tmp/whitelist.sock"

	// Hapus file socket lama jika ada
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	// Membuat listener di socket UNIX
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Whitelist plugin listening on", socketPath)

	// Loop untuk menerima koneksi dan memproses request
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleRequest(conn)
	}
}

// Fungsi untuk memproses setiap request yang masuk
func handleRequest(conn net.Conn) {
	defer conn.Close()

	var request SidraRequest
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println("Failed to decode request:", err)
		response := SidraResponse{
			StatusCode: 400,
			Body:       "Failed to decode request.",
		}
		json.NewEncoder(conn).Encode(response)
		return
	}

	// Proses whitelist menggunakan IP yang diterima
	response := whitelistHandler(request.IP)

	// Kirim response kembali ke klien
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(response)
	if err != nil {
		fmt.Println("Failed to encode response:", err)
	}
}

// Fungsi untuk mengecek apakah IP diizinkan atau tidak
func whitelistHandler(ip string) SidraResponse {
	if ip == "" {
		return SidraResponse{
			StatusCode: 400,
			Body:       "IP is missing.",
		}
	}

	if allowed, exists := allowedIPs[ip]; exists && allowed {
		return SidraResponse{
			StatusCode: 200,
			Body:       "Allowed",
		}
	}

	return SidraResponse{
		StatusCode: 403,
		Body:       "IP not allowed",
	}
}
