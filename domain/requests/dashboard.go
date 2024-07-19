package requests

type AddDashboardRequest struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Matches uint   `json:"matches"`
}
