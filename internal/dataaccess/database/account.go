package database

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
)

type Account struct {
	AccountID   uint64 `sql:"account_id"`
	AccountName string `sql:"accountname"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByAccountName(ctx context.Context, accountName string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountDataAccessor struct {
	database Database
}

func NewAccountDataAccessor(
	database *goqu.Database,
) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}

// CreateAccount implements AccountDataAccessor.
func (a accountDataAccessor) CreateAccount(ctx context.Context, account Account) (uint64, error) {
	result, err := a.database.
		Insert("accounts").
		Rows(goqu.Record{
			"accountname": account.AccountName,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		log.Printf("failed to create account, err=%+v\n", err)
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed to get last inserted id, err=%+v\n", err)
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

// GetAccountByID implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	panic("unimplemented")
}

// GetAccountByAccountName implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (Account, error) {
	panic("unimplemented")
}

// WithDatabase implements AccountDataAccessor.
func (a *accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}
