package interfaces

import (
	"github.com/runetale/runevision/domain/entity"
)

type DashboardInteractor interface {
	Get() ([]entity.DashboardHistory, error)
	Add(*entity.DashboardHistory) error
}
