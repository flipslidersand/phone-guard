package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting Phone Guard API on port %s...\n", port)

	// TODO: Initialize HTTP server
	// TODO: Setup database connection (PostgreSQL)
	// TODO: Setup LINE Notify client
	// TODO: Register routes
	//   - GET /api/numbers/:number (lookup)
	//   - POST /api/calls (log incoming call)
	//   - GET /api/whitelist/:device_id
	//   - POST /api/whitelist (add entry)
	//   - POST /api/notify (send LINE message)

	log.Println("Server started")
}
