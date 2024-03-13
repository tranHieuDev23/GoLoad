package dataaccess

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/cache"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/database"
)

var WireSet = wire.NewSet(
	cache.WireSet,
	database.WireSet,
)
