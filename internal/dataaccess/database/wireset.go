package database

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	InitializeAndMigrateUpDB,
	InitializeGoquDB,
	NewMigrator,
	NewAccountDataAccessor,
	NewAccountPasswordDataAccessor,
	NewDownloadTaskDataAccessor,
	NewTokenPublicKeyDataAccessor,
)
