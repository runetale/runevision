package localclient

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewLocalClient,
)
