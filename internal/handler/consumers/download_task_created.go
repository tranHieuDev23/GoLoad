package consumers

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/mq/producer"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

type DownloadTaskCreated interface {
	Handle(ctx context.Context, event producer.DownloadTaskCreated) error
}

type downloadTaskCreated struct {
	logger *zap.Logger
}

func NewDownloadTaskCreated(
	logger *zap.Logger,
) DownloadTaskCreated {
	return &downloadTaskCreated{
		logger: logger,
	}
}

func (d downloadTaskCreated) Handle(ctx context.Context, event producer.DownloadTaskCreated) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("event", event))
	logger.Info("download task created event received")

	return nil
}
