package app

import (
	"context"
	"syscall"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/database"
	"github.com/tranHieuDev23/GoLoad/internal/handler/grpc"
	"github.com/tranHieuDev23/GoLoad/internal/handler/http"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

type Server struct {
	databaseMigrator database.Migrator
	grpcServer       grpc.Server
	httpServer       http.Server
	logger           *zap.Logger
}

func NewServer(
	databaseMigrator database.Migrator,
	grpcServer grpc.Server,
	httpServer http.Server,
	logger *zap.Logger,
) *Server {
	return &Server{
		databaseMigrator: databaseMigrator,
		grpcServer:       grpcServer,
		httpServer:       httpServer,
		logger:           logger,
	}
}

func (s Server) Start() error {
	if err := s.databaseMigrator.Up(context.Background()); err != nil {
		s.logger.With(zap.Error(err)).Error("failed to execute database up migration")
		return err
	}

	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("grpc server stopped")
	}()

	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("http server stopped")
	}()

	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
	return nil
}
