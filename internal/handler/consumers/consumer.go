package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/mq/consumer"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/mq/producer"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

type Root interface {
	Start(ctx context.Context) error
}

type root struct {
	downloadTaskCreatedHandler DownloadTaskCreated
	mqConsumer                 consumer.Consumer
	logger                     *zap.Logger
}

func NewRoot(
	downloadTaskCreatedHandler DownloadTaskCreated,
	mqConsumer consumer.Consumer,
	logger *zap.Logger,
) Root {
	return &root{
		downloadTaskCreatedHandler: downloadTaskCreatedHandler,
		mqConsumer:                 mqConsumer,
		logger:                     logger,
	}
}

func (r root) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, r.logger)

	if err := r.mqConsumer.RegisterHandler(
		producer.MessageQueueDownloadTaskCreated,
		func(ctx context.Context, queueName string, payload []byte) error {
			var event producer.DownloadTaskCreated
			if err := json.Unmarshal(payload, &event); err != nil {
				return err
			}

			return r.downloadTaskCreatedHandler.Handle(ctx, event)
		}); err != nil {
		logger.With(zap.Error(err)).Error("failed to register download task created handler")
		return fmt.Errorf("failed to register download task created handler: %w", err)
	}

	return r.mqConsumer.Start(ctx)
}
