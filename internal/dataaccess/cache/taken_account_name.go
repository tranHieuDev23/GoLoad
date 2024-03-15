package cache

import (
	"context"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

const (
	setKeyNameTakenAccountName = "taken_account_name_set"
)

type TakenAccountName interface {
	Add(ctx context.Context, accountName string) error
	Has(ctx context.Context, accountName string) (bool, error)
}

type takenAccountName struct {
	client Client
	logger *zap.Logger
}

func NewTakenAccountName(
	client Client,
	logger *zap.Logger,
) TakenAccountName {
	return &takenAccountName{
		client: client,
		logger: logger,
	}
}

func (c takenAccountName) Add(ctx context.Context, accountName string) error {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("account_name", accountName))

	if err := c.client.AddToSet(ctx, setKeyNameTakenAccountName, accountName); err != nil {
		logger.With(zap.Error(err)).Error("failed to add account name to set in cache")
		return err
	}

	return nil
}

func (c takenAccountName) Has(ctx context.Context, accountName string) (bool, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).With(zap.String("account_name", accountName))
	result, err := c.client.IsDataInSet(ctx, setKeyNameTakenAccountName, accountName)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check if account name is in set in cache")
		return false, err
	}

	return result, nil
}
