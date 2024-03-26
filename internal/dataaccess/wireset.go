package dataaccess

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/cache"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/database"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/file"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/mq"
)

var WireSet = wire.NewSet(
	cache.WireSet,
	database.WireSet,
	mq.WireSet,
	file.WireSet,
)
