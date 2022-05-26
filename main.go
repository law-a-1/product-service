package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	pg, err := NewPostgres()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pg.Close()

	rd := NewRedis()

	// Migrate database
	if err := pg.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	server := NewServer(pg, rd)
	server.SetupMiddlewares()
	server.SetupRoutes()

	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
