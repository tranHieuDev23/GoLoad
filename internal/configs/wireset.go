package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "Account"),
	wire.FieldsOf(new(Config), "Database"),
)
