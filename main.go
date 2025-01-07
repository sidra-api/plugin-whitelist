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
	} else {
		log.Printf("INFO: Allowed IPs loaded from environment: %s\n", allowedIPsEnv)
	}
	for _, ip := range strings.Split(allowedIPsEnv, ",") {
		allowedIPs[strings.TrimSpace(ip)] = true // Menambahkan setiap IP ke map
	}
	log.Printf("INFO: Initialized allowed IPs: %v\n", allowedIPs)
}

// Fungsi handler untuk memproses permintaan dan memeriksa apakah IP klien diizinkan
func whitelistHandler(request server.Request) server.Response {
	clientIP := request.Headers["X-Real-IP"]
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(request.Headers["X-Forwarded-For"])
	}
	if clientIP == "" {
		clientIP = request.Headers["Remote-Addr"]
	}

	if clientIP == "" {
		log.Println("ERROR: Unable to determine client IP")
		return server.Response{
			StatusCode: 400,
			Body:       "Bad Request - Unable to determine IP",
		}
	}

	// Membuat respons berdasarkan apakah IP diizinkan atau tidak
	var response server.Response
	if allowedIPs[clientIP] {
		// Jika IP ada di daftar allowedIPs
		log.Printf("INFO: IP allowed: %s\n", clientIP) 
		response.StatusCode = 200            
		response.Body = "Allowed"            
	} else {
		// Jika IP tidak ada di daftar allowedIPs
		log.Printf("WARNING: IP not allowed: %s\n", clientIP) 
		response.StatusCode = 403                
		response.Body = "IP not allowed"         
	}
	return response 
}

func main() {
	initAllowedIPs() // Memanggil fungsi untuk inisialisasi daftar IP yang diizinkan

	// Mengambil nama plugin dari environment variable, atau menggunakan default
	pluginName := os.Getenv("PLUGIN_NAME")
	if pluginName == "" {
		pluginName = "whitelist"
	}

	srv := server.NewServer(pluginName, whitelistHandler)
	log.Printf("Plugin '%s' using go-pdk server.\n", pluginName)
	srv.Start()
}
