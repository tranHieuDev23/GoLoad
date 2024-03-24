package http

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tranHieuDev23/GoLoad/internal/configs"
	"github.com/tranHieuDev23/GoLoad/internal/generated/grpc/go_load"
	handlerGRPC "github.com/tranHieuDev23/GoLoad/internal/handler/grpc"
	"github.com/tranHieuDev23/GoLoad/internal/handler/http/servemuxoptions"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

const (
	//nolint:gosec // This is just to specify the cookie name
	AuthTokenCookieName = "GOLOAD_AUTH"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	grpcConfig configs.GRPC
	httpConfig configs.HTTP
	authConfig configs.Auth
	logger     *zap.Logger
}

func NewServer(
	grpcConfig configs.GRPC,
	httpConfig configs.HTTP,
	authConfig configs.Auth,
	logger *zap.Logger,
) Server {
	return &server{
		grpcConfig: grpcConfig,
		httpConfig: httpConfig,
		authConfig: authConfig,
		logger:     logger,
	}
}

func (s server) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, s.logger)

	tokenExpiresInDuration, err := s.authConfig.Token.GetExpiresInDuration()
	if err != nil {
		return err
	}

	grpcMux := runtime.NewServeMux(
		servemuxoptions.WithAuthCookieToAuthMetadata(AuthTokenCookieName, handlerGRPC.AuthTokenMetadataName),
		servemuxoptions.WithAuthMetadataToAuthCookie(
			handlerGRPC.AuthTokenMetadataName, AuthTokenCookieName, tokenExpiresInDuration),
	)
	err = go_load.RegisterGoLoadServiceHandlerFromEndpoint(
		ctx,
		grpcMux,
		s.grpcConfig.Address,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		})
	if err != nil {
		return err
	}

	httpServer := http.Server{
		Addr:              s.httpConfig.Address,
		ReadHeaderTimeout: time.Minute,
		Handler:           grpcMux,
	}

	logger.With(zap.String("address", s.httpConfig.Address)).Info("starting http server")
	return httpServer.ListenAndServe()
}
