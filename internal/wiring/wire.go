//go:build wireinject
// +build wireinject

//
//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/GoLoad/internal/app"
	"github.com/tranHieuDev23/GoLoad/internal/configs"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess"
	"github.com/tranHieuDev23/GoLoad/internal/handler"
	"github.com/tranHieuDev23/GoLoad/internal/logic"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	utils.WireSet,
	dataaccess.WireSet,
	logic.WireSet,
	handler.WireSet,
	app.WireSet,
)

func InitializeStandaloneServer(configFilePath configs.ConfigFilePath) (*app.StandaloneServer, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}
