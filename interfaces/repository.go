package interfaces

import "github.com/runetale/runevision/domain/entity"

type DashboardRepository interface {
	Create(*entity.DashboardHistory) error
	GetHistories() ([]entity.DashboardHistory, error)
}
