package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Initialise router & starts server
func StartServer() {
	router := Routes()
	if err := Server(router); err != nil {
		log.Fatalf("Server encountered an error: %v", err)
	}
}

// Starts the HTTP server
func Server(handler http.Handler) error {
	var port string
	// Import specified port
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return err
	}
	defer listener.Close()

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)
	return server.Serve(listener)
}
