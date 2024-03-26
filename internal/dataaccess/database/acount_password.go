package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

var (
	TabNameAccountPasswords = goqu.T("account_passwords")
)

const (
	ColNameAccountPasswordsOfAccountID = "of_account_id"
	ColNameAccountPasswordsHash        = "hash"
)

type AccountPassword struct {
	OfAccountID uint64 `db:"of_account_id" goqu:"skipupdate"`
	Hash        string `db:"hash"`
}

type AccountPasswordDataAccessor interface {
	CreateAccountPassword(ctx context.Context, accountPassword AccountPassword) error
	GetAccountPassword(ctx context.Context, ofAccountID uint64) (AccountPassword, error)
	UpdateAccountPassword(ctx context.Context, accountPassword AccountPassword) error
	WithDatabase(database Database) AccountPasswordDataAccessor
}

type accountPasswordDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewAccountPasswordDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a accountPasswordDataAccessor) CreateAccountPassword(ctx context.Context, accountPassword AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	_, err := a.database.
		Insert(TabNameAccountPasswords).
		Rows(goqu.Record{
			ColNameAccountPasswordsOfAccountID: accountPassword.OfAccountID,
			ColNameAccountPasswordsHash:        accountPassword.Hash,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create account password")
		return status.Error(codes.Internal, "failed to create account password")
	}

	return nil
}

func (a accountPasswordDataAccessor) GetAccountPassword(
	ctx context.Context,
	ofAccountID uint64,
) (AccountPassword, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("of_account_id", ofAccountID))
	accountPassword := AccountPassword{}
	found, err := a.database.
		From(TabNameAccountPasswords).
		Where(goqu.Ex{ColNameAccountPasswordsOfAccountID: ofAccountID}).
		ScanStructContext(ctx, &accountPassword)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account password by id")
		return AccountPassword{}, status.Error(codes.Internal, "failed to get account password by id")
	}

	if !found {
		logger.Warn("cannot find account by id")
		return AccountPassword{}, sql.ErrNoRows
	}

	return accountPassword, nil
}

func (a accountPasswordDataAccessor) UpdateAccountPassword(ctx context.Context, accountPassword AccountPassword) error {
	logger := utils.LoggerWithContext(ctx, a.logger)
	_, err := a.database.
		Update(TabNameAccountPasswords).
		Set(goqu.Record{ColNameAccountPasswordsHash: accountPassword.Hash}).
		Where(goqu.Ex{ColNameAccountPasswordsOfAccountID: accountPassword.OfAccountID}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to update account password")
		return status.Error(codes.Internal, "failed to update account password")
	}

	return nil
}

func (a accountPasswordDataAccessor) WithDatabase(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
		logger:   a.logger,
	}
}
