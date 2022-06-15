package grpc

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
}

type server struct {
	UnimplementedProductServer
}

func NewServer() Server {
	s := grpc.NewServer()
	RegisterProductServer(s, &server{})
	return Server{
		grpcServer: s,
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
