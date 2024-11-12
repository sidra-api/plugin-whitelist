package main

import (
	"log"
	"net"
	"strings"

	"github.com/sidra-gateway/go-pdk/server"
)

// Daftar IP yang diizinkan
var allowedIPs = map[string]bool{
	"192.168.1.1": true,
	"192.168.1.2": true,
}

// Fungsi untuk menangani permintaan dengan memeriksa whitelist IP
func whitelistHandler(request server.Request) server.Response {
	clientIP := request.Headers["X-Real-Ip"]
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(request.Headers["X-Forwarded-For"])
	}
	if clientIP == "" {
		clientIP = request.Headers["Remote-Addr"]
	}

	ipAllowed := false
	for ip := range allowedIPs {
		if strings.TrimSpace(ip) == clientIP {
			ipAllowed = true
			break
		}
	}

	var response server.Response
	if ipAllowed {
		log.Println("IP allowed:", clientIP)
		response.StatusCode = 200
		response.Body = "Allowed"
	} else {
		log.Println("IP not allowed:", clientIP)
		response.StatusCode = 403
		response.Body = "IP not allowed"
	}
	return response
}

func main() {
	// Menggunakan server dari go-pdk
	srv := server.NewServer("whitelist", whitelistHandler)
	log.Println("Whitelist plugin using go-pdk server.")
	srv.Start()
}
