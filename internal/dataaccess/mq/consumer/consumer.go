package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/configs"
)

type HandlerFunc func(ctx context.Context, queueName string, payload []byte) error

type Consumer interface {
	RegisterHandler(queueName string, handlerFunc HandlerFunc) error
	Start(ctx context.Context) error
}

type partitionConsumerAndHandlerFunc struct {
	queueName         string
	partitionConsumer sarama.PartitionConsumer
	handlerFunc       HandlerFunc
}

type consumer struct {
	saramaConsumer                      sarama.Consumer
	partitionConsumerAndHandlerFuncList []partitionConsumerAndHandlerFunc
	logger                              *zap.Logger
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = mqConfig.ClientID
	saramaConfig.Metadata.Full = true
	return saramaConfig
}

func NewConsumer(
	mqConfig configs.MQ,
	logger *zap.Logger,
) (Consumer, error) {
	saramaConsumer, err := sarama.NewConsumer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama consumer: %w", err)
	}

	return &consumer{
		saramaConsumer: saramaConsumer,
		logger:         logger,
	}, nil
}

func (c *consumer) RegisterHandler(queueName string, handlerFunc HandlerFunc) error {
	partitionConsumer, err := c.saramaConsumer.ConsumePartition(queueName, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("failed to create sarama partition consumer: %w", err)
	}

	c.partitionConsumerAndHandlerFuncList = append(
		c.partitionConsumerAndHandlerFuncList,
		partitionConsumerAndHandlerFunc{
			queueName:         queueName,
			partitionConsumer: partitionConsumer,
			handlerFunc:       handlerFunc,
		})

	return nil
}

func (c consumer) Start(_ context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for i := range c.partitionConsumerAndHandlerFuncList {
		go func(i int) {
			queueName := c.partitionConsumerAndHandlerFuncList[i].queueName
			partitionConsumer := c.partitionConsumerAndHandlerFuncList[i].partitionConsumer
			handlerFunc := c.partitionConsumerAndHandlerFuncList[i].handlerFunc
			logger := c.logger.With(zap.String("queue_name", queueName))

			for {
				select {
				case message := <-partitionConsumer.Messages():
					if err := handlerFunc(context.Background(), queueName, message.Value); err != nil {
						logger.With(zap.Error(err)).Error("failed to handle message")
					}

				case <-signals:
					break
				}
			}
		}(i)
	}

	<-signals
	return nil
}
