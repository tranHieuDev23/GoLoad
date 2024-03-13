package cache

import (
	"context"

	"go.uber.org/zap"
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
	if err := c.client.AddToSet(ctx, setKeyNameTakenAccountName, accountName); err != nil {
		return err
	}

	return nil
}

func (c takenAccountName) Has(ctx context.Context, accountName string) (bool, error) {
	return c.client.IsDataInSet(ctx, setKeyNameTakenAccountName, accountName)
}
