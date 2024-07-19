package interfaces

import "github.com/runetale/runevision/domain/entity"

type DashboardRepository interface {
	Create(SQLExecuter, *entity.DashboardHistory) error
	GetHistories(SQLExecuter) ([]entity.DashboardHistory, error)
}
