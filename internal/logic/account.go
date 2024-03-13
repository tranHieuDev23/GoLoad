package logic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/cache"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/database"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

type CreateAccountParams struct {
	AccountName string
	Password    string
}

type CreateAccountOutput struct {
	ID          uint64
	AccountName string
}

type CreateSessionParams struct {
	AccountName string
	Password    string
}

type Account interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (token string, err error)
}

type account struct {
	goquDatabase                *goqu.Database
	takenAccountNameCache       cache.TakenAccountName
	accountDataAccessor         database.AccountDataAccessor
	accountPasswordDataAccessor database.AccountPasswordDataAccessor
	hashLogic                   Hash
	tokenLogic                  Token
	logger                      *zap.Logger
}

func NewAccount(
	goquDatabase *goqu.Database,
	takenAccountNameCache cache.TakenAccountName,
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordDataAccessor database.AccountPasswordDataAccessor,
	hashLogic Hash,
	tokenLogic Token,
	logger *zap.Logger,
) Account {
	return &account{
		goquDatabase:                goquDatabase,
		takenAccountNameCache:       takenAccountNameCache,
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		hashLogic:                   hashLogic,
		tokenLogic:                  tokenLogic,
		logger:                      logger,
	}
}

func (a account) isAccountAccountNameTaken(ctx context.Context, accountName string) (bool, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("account_name", accountName))

	accountNameTaken, err := a.takenAccountNameCache.Has(ctx, accountName)
	if err != nil {
		logger.With(zap.Error(err)).Warn("failed to get account name from taken set in cache, will fall back to database")
	} else {
		return accountNameTaken, nil
	}

	if _, err := a.accountDataAccessor.GetAccountByAccountName(ctx, accountName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	if err := a.takenAccountNameCache.Add(ctx, accountName); err != nil {
		logger.With(zap.Error(err)).Warn("failed to set account name into taken set in cache")
	}

	return true, nil
}

func (a account) CreateAccount(ctx context.Context, params CreateAccountParams) (CreateAccountOutput, error) {
	var accountID uint64

	txErr := a.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		accountNameTaken, err := a.isAccountAccountNameTaken(ctx, params.AccountName)
		if err != nil {
			return err
		}

		if accountNameTaken {
			return errors.New("account name is already taken")
		}

		accountID, err = a.accountDataAccessor.WithDatabase(td).CreateAccount(ctx, database.Account{
			AccountName: params.AccountName,
		})
		if err != nil {
			return err
		}

		hashedPassword, err := a.hashLogic.Hash(ctx, params.Password)
		if err != nil {
			return err
		}

		if err := a.accountPasswordDataAccessor.WithDatabase(td).CreateAccountPassword(ctx, database.AccountPassword{
			OfAccountID: accountID,
			Hash:        hashedPassword,
		}); err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return CreateAccountOutput{}, txErr
	}

	return CreateAccountOutput{
		ID:          accountID,
		AccountName: params.AccountName,
	}, nil
}

func (a account) CreateSession(ctx context.Context, params CreateSessionParams) (token string, err error) {
	existingAccount, err := a.accountDataAccessor.GetAccountByAccountName(ctx, params.AccountName)
	if err != nil {
		return "", err
	}

	existingAccountPassword, err := a.accountPasswordDataAccessor.GetAccountPassword(ctx, existingAccount.ID)
	if err != nil {
		return "", err
	}

	isHashEqual, err := a.hashLogic.IsHashEqual(ctx, params.Password, existingAccountPassword.Hash)
	if err != nil {
		return "", err
	}

	if !isHashEqual {
		return "", errors.New("incorrect password")
	}

	return "", nil
}
