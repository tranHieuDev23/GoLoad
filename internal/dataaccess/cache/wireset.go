package cache

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewRedisClient,
	NewTokenPublicKey,
	NewTakenAccountName,
)
