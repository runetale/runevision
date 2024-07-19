package interactor

import (
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/interfaces"
)

type DashboardInteractor struct {
	dashboardRepository interfaces.DashboardRepository
}

func NewDashboardInteractor(
	dashboardRepository interfaces.DashboardRepository,
) interfaces.DashboardInteractor {
	return &DashboardInteractor{
		dashboardRepository: dashboardRepository,
	}
}

func (i *DashboardInteractor) Get() ([]entity.DashboardHistory, error) {
	histories, err := i.dashboardRepository.GetHistories()
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (i *DashboardInteractor) Add(history *entity.DashboardHistory) error {
	err := i.dashboardRepository.Create(history)
	if err != nil {
		return err
	}
	return nil
}
