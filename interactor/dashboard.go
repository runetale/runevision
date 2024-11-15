package interactor

import (
	"github.com/runetale/thor/database"
	"github.com/runetale/thor/domain/entity"
	"github.com/runetale/thor/interfaces"
)

type DashboardInteractor struct {
	db                  *database.Postgres
	dashboardRepository interfaces.DashboardRepository
}

func NewDashboardInteractor(
	db *database.Postgres,
	dashboardRepository interfaces.DashboardRepository,
) interfaces.DashboardInteractor {
	return &DashboardInteractor{
		db:                  db,
		dashboardRepository: dashboardRepository,
	}
}

func (i *DashboardInteractor) Get() ([]entity.DashboardHistory, error) {
	histories, err := i.dashboardRepository.GetHistories(i.db)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (i *DashboardInteractor) Add(history *entity.DashboardHistory) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	err = i.dashboardRepository.Create(i.db, history)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
