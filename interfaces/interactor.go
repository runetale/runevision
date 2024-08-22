package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/domain/requests"
)

type DashboardInteractor interface {
	Get() ([]entity.DashboardHistory, error)
	Add(*entity.DashboardHistory) error
}

type HackInteractor interface {
	Scan(*requests.HackDoScanRequest, echo.Context) (*entity.HackHistory, error)
}
