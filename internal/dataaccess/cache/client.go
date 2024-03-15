package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tranHieuDev23/GoLoad/internal/configs"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataInSet(ctx context.Context, key string, data any) (bool, error)
}

type client struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewClient(
	cacheConfig configs.Cache,
	logger *zap.Logger,
) Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cacheConfig.Address,
		Username: cacheConfig.Username,
		Password: cacheConfig.Password,
	})

	return &client{
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r client) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data)).
		With(zap.Duration("ttl", ttl))

	if err := r.redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		logger.With(zap.Error(err)).Error("failed to set data into cache")
		return status.Errorf(codes.Internal, "failed to set data into cache: %+v", err)
	}

	return nil
}

func (r client) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key))

	data, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		logger.With(zap.Error(err)).Error("failed to get data from cache")
		return nil, status.Errorf(codes.Internal, "failed to get data from cache: %+v", err)
	}

	return data, nil
}

func (r client) AddToSet(ctx context.Context, key string, data ...any) error {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data))

	if err := r.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		logger.With(zap.Error(err)).Error("failed to set data into set inside cache")
		return status.Errorf(codes.Internal, "failed to set data into set inside cache: %+v", err)
	}

	return nil
}

func (r client) IsDataInSet(ctx context.Context, key string, data any) (bool, error) {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data))

	result, err := r.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check if data is member of set inside cache")
		return false, status.Errorf(codes.Internal, "failed to check if data is member of set inside cache: %+v", err)
	}

	return result, nil
}
