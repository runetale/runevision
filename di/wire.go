//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/runetale/thor/database"
	"github.com/runetale/thor/domain/config"
	"github.com/runetale/thor/handler"
	"github.com/runetale/thor/interactor"
	"github.com/runetale/thor/interfaces"
	"github.com/runetale/thor/localclient"
	"github.com/runetale/thor/repository"
	"github.com/runetale/thor/utility"
)

var wireSet = wire.NewSet(
	utility.WireSet,
	database.WireSet,
	handler.WireSet,
	interactor.WireSet,
	repository.WireSet,
	localclient.WireSet,
)

func InitializeLogger(logConfig config.Log) (l *utility.Logger) {
	wire.Build(wireSet)
	return
}

func InitializePostgres(dbConfig config.Postgres, logConfig config.Log) (db *database.Postgres) {
	wire.Build(wireSet)
	return
}

func InitializeDashboardRepository() (repo interfaces.DashboardRepository) {
	wire.Build(wireSet)
	return
}

func InitializeDashboardHandler(dbConfig config.Postgres, logConfig config.Log) (h interfaces.DashboardHandler) {
	wire.Build(wireSet)
	return
}

func InitializeHackHandler(dbConfig config.Postgres, logConfig config.Log) (h interfaces.HackHandler) {
	wire.Build(wireSet)
	return
}
