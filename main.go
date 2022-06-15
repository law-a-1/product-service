package main

import (
	"context"
	"os"

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

	persistent, err := NewPersistent(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
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
		logger.Fatalf("failed creating schema resources: %v", err)
	}
	logger.Info("database migrated")

	server := NewServer(logger, persistent)
	server.SetupMiddlewares()
	server.SetupRoutes()

	grpcServer := grpc.NewServer(logger, persistent)

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	if err := grpcServer.Start(); err != nil {
		logger.Fatalf("failed to start grpc server: %v", err)
	}
}
