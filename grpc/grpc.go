package grpc

import (
	"log"
	"net"
	"os"

	"github.com/law-a-1/product-service/ent"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	db         *ent.Client
	logger     *zap.SugaredLogger
	UnimplementedProductServer
}

func NewServer(logger *zap.SugaredLogger, db *ent.Client) Server {
	s := grpc.NewServer()
	RegisterProductServer(s, &Server{})
	return Server{
		grpcServer: s,
		db:         db,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		return err
	}
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
