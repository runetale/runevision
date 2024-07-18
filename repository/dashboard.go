package repository

import (
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/interfaces"
)

type DashboardRepository struct {
	db interfaces.SQLExecuter
}

func NewDashboardRepository(
	db interfaces.SQLExecuter,
) interfaces.DashboardRepository {
	return &DashboardRepository{
		db: db,
	}
}

func (r *DashboardRepository) GetHistories() ([]entity.DashboardHistory, error) {
	var histories []entity.DashboardHistory
	err := r.db.Select(&histories, `SELECT * FROM dashboard_histories`)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *DashboardRepository) Create(obj *entity.DashboardHistory) error {
	err := r.db.NameExec(`INSERT INTO dashboard_histories (name, status, matches, created_at) VALUES (:name, :status, :matches, :created_at)`, obj)
	if err != nil {
		return err
	}
	return nil
}
