package grpc

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"google.golang.org/grpc"

	"github.com/tranHieuDev23/GoLoad/internal/generated/grpc/go_load"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler go_load.GoLoadServiceServer
}

func NewServer(
	handler go_load.GoLoadServiceServer,
) Server {
	return &server{
		handler: handler,
	}
}

func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}

	defer listener.Close()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	)
	go_load.RegisterGoLoadServiceServer(server, s.handler)
	return server.Serve(listener)
}
