package main

import (
	"log"

	"github.com/nebulabox/nebulabox/internal/api"
	"github.com/nebulabox/nebulabox/internal/database"
)

func main() {
	// Initialize databases
	if err := database.Init(); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		log.Println("Continuing with in-memory storage...")
	}
	defer database.Close()

	// Create API server
	server, err := api.NewServer()
	if err != nil {
		log.Fatalf("Failed to create API server: %v", err)
	}
	defer server.Close()

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
