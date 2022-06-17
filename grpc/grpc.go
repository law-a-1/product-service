package grpc

import (
	"context"
	"github.com/law-a-1/product-service/ent/product"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	grpcServer := grpc.NewServer()
	s := Server{
		grpcServer: grpcServer,
		db:         db,
		logger:     logger,
	}
	RegisterProductServer(grpcServer, s)
	return s
}

func (s Server) Start() error {
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		return err
	}
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

func (s *Server) DecreaseStock(ctx context.Context, in *DecreaseStockRequest) (*DecreaseStockResponse, error) {
	p, err := s.db.Product.Query().Where(product.ID(int(in.ID))).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			s.logger.Fatalf("Product not found")
			st := status.New(codes.NotFound, "Product with given ID not found.")
			return &DecreaseStockResponse{}, st.Err()
		}
		s.logger.Fatalf("")
		st := status.New(codes.Internal, "Failed to find product with given ID.")
		return &DecreaseStockResponse{}, st.Err()
	}
	if p.Stock < int(in.Amount) {
		s.logger.Fatalf("")
		st := status.New(codes.InvalidArgument, "Stock is not enough.")
		return &DecreaseStockResponse{}, st.Err()
	}
	_, err = s.db.Product.Update().SetStock(int(in.Amount)).Save(ctx)
	if err != nil {
		s.logger.Fatalf("")
		st := status.New(codes.Internal, "Failed to update stock data.")
		return &DecreaseStockResponse{}, st.Err()
	}

	return &DecreaseStockResponse{}, nil
}
