package grpc

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHandler,
	NewServer,
)
