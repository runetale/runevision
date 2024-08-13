//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/runetale/runevision/database"
	"github.com/runetale/runevision/domain/config"
	"github.com/runetale/runevision/handler"
	"github.com/runetale/runevision/interactor"
	"github.com/runetale/runevision/interfaces"
	"github.com/runetale/runevision/repository"
	"github.com/runetale/runevision/utility"
)

var wireSet = wire.NewSet(
	utility.WireSet,
	database.WireSet,
	handler.WireSet,
	interactor.WireSet,
	repository.WireSet,
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
