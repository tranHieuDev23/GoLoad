package file

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewClient,
)
