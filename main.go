package main

import (
	"context"
	"github.com/law-a-1/product-service/ent"
	"github.com/law-a-1/product-service/grpc"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger, err := NewLogger()
	if err != nil {
		logger.Fatalf("failed to create logger: %v", err)
	}
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			logger.Fatalf("failed to sync logger: %v", err)
		}
	}(logger)
	logger.Info("logger created")

	persistent, err := NewPersistent()
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	defer func(persistent *ent.Client) {
		err := persistent.Close()
		if err != nil {
			logger.Fatalf("failed to close database: %v", err)
		}
	}(persistent)
	logger.Info("database connected")

	// Migrate database
	if err := persistent.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	logger.Info("database migrated")

	cache := NewCache()

	server := NewServer(logger, persistent, cache)
	server.SetupMiddlewares()
	server.SetupRoutes()

	grpcServer := grpc.NewServer()

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	if err := grpcServer.Start(); err != nil {
		logger.Fatalf("failed to start grpc server: %v", err)
	}
}
