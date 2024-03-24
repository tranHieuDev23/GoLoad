package database

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/doug-martin/goqu/v9"

	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

var (
	TabNameTokenPublicKeys = goqu.T("token_public_keys")
)

const (
	ColNameTokenPublicKeysID        = "id"
	ColNameTokenPublicKeysPublicKey = "public_key"
)

type TokenPublicKey struct {
	ID        uint64 `db:"id" goqu:"skipinsert,skipupdate"`
	PublicKey string `db:"public_key"`
}

type TokenPublicKeyDataAccessor interface {
	CreatePublicKey(ctx context.Context, tokenPublicKey TokenPublicKey) (uint64, error)
	GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error)
	WithDatabase(database Database) TokenPublicKeyDataAccessor
}

type tokenPublicKeyDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewTokenPublicKeyDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a tokenPublicKeyDataAccessor) CreatePublicKey(
	ctx context.Context,
	tokenPublicKey TokenPublicKey,
) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, a.logger)
	result, err := a.database.
		Insert(TabNameTokenPublicKeys).
		Rows(goqu.Record{
			ColNameTokenPublicKeysPublicKey: tokenPublicKey.PublicKey,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create token public key")
		return 0, status.Error(codes.Internal, "failed to create token public key")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, status.Error(codes.Internal, "failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (a tokenPublicKeyDataAccessor) GetPublicKey(ctx context.Context, id uint64) (TokenPublicKey, error) {
	logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("id", id))

	tokenPublicKey := TokenPublicKey{}
	found, err := a.database.
		Select().
		From(TabNameTokenPublicKeys).
		Where(goqu.Ex{
			ColNameTokenPublicKeysID: id,
		}).
		Executor().
		ScanStructContext(ctx, &tokenPublicKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get public key")
		return TokenPublicKey{}, status.Error(codes.Internal, "failed to get public key")
	}

	if !found {
		logger.Warn("public key not found")
		return TokenPublicKey{}, sql.ErrNoRows
	}

	return tokenPublicKey, nil
}

func (a tokenPublicKeyDataAccessor) WithDatabase(database Database) TokenPublicKeyDataAccessor {
	a.database = database
	return a
}
