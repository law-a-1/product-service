package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	client, err := NewPostgres()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()

	// Migrate database
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	server := NewServer(client)
	server.SetupMiddlewares()
	server.SetupRoutes()

	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
