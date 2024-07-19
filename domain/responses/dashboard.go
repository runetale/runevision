package responses

import "github.com/runetale/runevision/domain/entity"

type GetDashboardResponse struct {
	ActiveVolnerablities uint                    `json:"active_vulnerablities"`
	RecentScans          GetDashboardRecentScans `json:"recent_scans"`
}

type GetDashboardRecentScans struct {
	Count     uint                      `json:"count"`
	Histories []entity.DashboardHistory `json:"histories"`
}

func NewGetDashboardResponse(
	activeVulnerablities, count uint, histories []entity.DashboardHistory,
) GetDashboardResponse {
	return GetDashboardResponse{
		ActiveVolnerablities: activeVulnerablities,
		RecentScans: GetDashboardRecentScans{
			Count:     count,
			Histories: histories,
		},
	}
}
