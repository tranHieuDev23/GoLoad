package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

var (
	TabNameAccounts = goqu.T("accounts")

	ErrAccountNotFound = status.Error(codes.NotFound, "account not found")
)

const (
	ColNameAccountsID          = "id"
	ColNameAccountsAccountName = "account_name"
)

type Account struct {
	ID          uint64 `db:"id" goqu:"skipinsert,skipupdate"`
	AccountName string `db:"account_name"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByAccountName(ctx context.Context, accountName string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewAccountDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, account Account) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Any("account", account))

	result, err := a.database.
		Insert(TabNameAccounts).
		Rows(goqu.Record{
			ColNameAccountsAccountName: account.AccountName,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create account")
		return 0, status.Error(codes.Internal, "failed to create account")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, status.Error(codes.Internal, "failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)

	account := Account{}
	found, err := a.database.
		From(TabNameAccounts).
		Where(goqu.C(ColNameAccountsID).Eq(id)).
		ScanStructContext(ctx, &account)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account by id")
		return Account{}, status.Error(codes.Internal, "failed to get account by id")
	}

	if !found {
		logger.Warn("cannot find account by id")
		return Account{}, ErrAccountNotFound
	}

	return account, nil
}

func (a accountDataAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (Account, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", accountName))

	account := Account{}
	found, err := a.database.
		From(TabNameAccounts).
		Where(goqu.C(ColNameAccountsAccountName).Eq(accountName)).
		ScanStructContext(ctx, &account)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get account by name")
		return Account{}, status.Error(codes.Internal, "failed to get account by name")
	}

	if !found {
		logger.Warn("cannot find account by name")
		return Account{}, ErrAccountNotFound
	}

	return account, nil
}

func (a accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   a.logger,
	}
}
