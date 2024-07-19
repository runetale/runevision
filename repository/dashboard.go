package repository

import (
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/interfaces"
)

type DashboardRepository struct {
}

func NewDashboardRepository() interfaces.DashboardRepository {
	return &DashboardRepository{}
}

func (r *DashboardRepository) GetHistories(db interfaces.SQLExecuter) ([]entity.DashboardHistory, error) {
	var histories []entity.DashboardHistory
	err := db.Select(&histories, `SELECT * FROM dashboard_histories`)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *DashboardRepository) Create(db interfaces.SQLExecuter, obj *entity.DashboardHistory) error {
	err := db.NameExec(`INSERT INTO dashboard_histories (name, status, matches, created_at) VALUES (:name, :status, :matches, :created_at)`, obj)
	if err != nil {
		return err
	}
	return nil
}
