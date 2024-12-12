package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/sidra-gateway/go-pdk/server"
)

var allowedIPs map[string]bool

// Fungsi untuk inisialisasi daftar IP yang diizinkan dari environment variable
func initAllowedIPs() {
	allowedIPs = make(map[string]bool) 
	allowedIPsEnv := os.Getenv("ALLOWED_IPS") 
	if allowedIPsEnv == "" {
		// Jika environment variable tidak ada, gunakan daftar IP default
		allowedIPsEnv = "192.168.1.1,192.168.1.2"
	}
	for _, ip := range strings.Split(allowedIPsEnv, ",") {
		allowedIPs[strings.TrimSpace(ip)] = true // Menambahkan setiap IP ke map
	}
}

// Fungsi handler untuk memproses permintaan dan memeriksa apakah IP klien diizinkan
func whitelistHandler(request server.Request) server.Response {
	clientIP := request.Headers["X-Real-Ip"]
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(request.Headers["X-Forwarded-For"])
	}
	if clientIP == "" {
		clientIP = request.Headers["Remote-Addr"]
	}

	// Membuat respons berdasarkan apakah IP diizinkan atau tidak
	var response server.Response
	if allowedIPs[clientIP] {
		// Jika IP ada di daftar allowedIPs
		log.Println("IP allowed:", clientIP) 
		response.StatusCode = 200            
		response.Body = "Allowed"            
	} else {
		// Jika IP tidak ada di daftar allowedIPs
		log.Println("IP not allowed:", clientIP) 
		response.StatusCode = 403                
		response.Body = "IP not allowed"         
	}
	return response 
}

func main() {
	initAllowedIPs() // Memanggil fungsi untuk inisialisasi daftar IP yang diizinkan
	srv := server.NewServer("whitelist", whitelistHandler)
	log.Println("Whitelist plugin using go-pdk server.")
	srv.Start()
}
