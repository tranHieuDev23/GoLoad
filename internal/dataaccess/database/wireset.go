package database

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	InitializeDB,
	InitializeGoquDB,
	NewAccountDataAccessor,
	NewAccountPasswordDataAccessor,
	NewDownloadTaskDataAccessor,
	NewTokenPublicKeyDataAccessor,
)
