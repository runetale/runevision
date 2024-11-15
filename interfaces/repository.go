package interfaces

import "github.com/runetale/thor/domain/entity"

type DashboardRepository interface {
	Create(SQLExecuter, *entity.DashboardHistory) error
	GetHistories(SQLExecuter) ([]entity.DashboardHistory, error)
}
