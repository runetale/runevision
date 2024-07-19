package utility

import "github.com/google/wire"

var WireSet = wire.NewSet(
	MustNewLoggerFromConfig,
)
