package entity

import (
	"time"

	"github.com/runetale/runevision/domain/requests"
)

type DashboardHistory struct {
	ID        uint      `db:"id" json:"-"`
	Name      string    `db:"name" json:"name"`
	Status    string    `db:"status" json:"status"`
	Matches   uint      `db:"matches" json:"matches"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewDashboardHistory(
	name, status string, matches uint,
) *DashboardHistory {
	return &DashboardHistory{
		Name:      name,
		Status:    status,
		Matches:   matches,
		CreatedAt: time.Now(),
	}
}

func NewDashboardHistoryFromRequest(req requests.AddDashboardRequest) *DashboardHistory {
	return NewDashboardHistory(req.Name, req.Status, req.Matches)
}
