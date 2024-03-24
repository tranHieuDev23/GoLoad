package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tranHieuDev23/GoLoad/internal/configs"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

type Client interface {
	Produce(ctx context.Context, queueName string, payload []byte) error
}

type client struct {
	saramaSyncProducer sarama.SyncProducer
	logger             *zap.Logger
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 1
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = mqConfig.ClientID
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func NewClient(
	mqConfig configs.MQ,
	logger *zap.Logger,
) (Client, error) {
	saramaSyncProducer, err := sarama.NewSyncProducer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return &client{
		saramaSyncProducer: saramaSyncProducer,
		logger:             logger,
	}, nil
}

func (c client) Produce(ctx context.Context, queueName string, payload []byte) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("queue_name", queueName)).
		With(zap.ByteString("payload", payload))

	if _, _, err := c.saramaSyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: queueName,
		Value: sarama.ByteEncoder(payload),
	}); err != nil {
		logger.With(zap.Error(err)).Error("failed to produce message")
		return status.Error(codes.Internal, "failed to produce message")
	}

	return nil
}
