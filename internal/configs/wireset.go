package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "GRPC"),
	wire.FieldsOf(new(Config), "HTTP"),
	wire.FieldsOf(new(Config), "Log"),
	wire.FieldsOf(new(Config), "Auth"),
	wire.FieldsOf(new(Config), "Database"),
	wire.FieldsOf(new(Config), "Cache"),
	wire.FieldsOf(new(Config), "MQ"),
	wire.FieldsOf(new(Config), "Cron"),
	wire.FieldsOf(new(Config), "Download"),
)
