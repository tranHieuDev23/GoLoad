package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "Log"),
	wire.FieldsOf(new(Config), "Auth"),
	wire.FieldsOf(new(Config), "Database"),
)
