//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/runetale/runevision/handler"
	"github.com/runetale/runevision/interactor"
	"github.com/runetale/runevision/interfaces"
	"github.com/runetale/runevision/repository"
)

var wireSet = wire.NewSet(
	handler.WireSet,
	interactor.WireSet,
	repository.WireSet,
)

func InitializeDashboardHandler(db interfaces.SQLExecuter) (h interfaces.DashboardHandler) {
	wire.Build(wireSet)
	return
}
