package entity

import "time"

// tempolary (snt)
type ScanResponse struct {
	ID           uint      `db:"id" json:"-"`
	Name         string    `db:"name" json:"name"`
	SequentialID string    `db:"sequential_id" json:"sequential_id"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func NewScanResponse(
	name, sid, status string,
) *ScanResponse {
	return &ScanResponse{
		Name:         name,
		SequentialID: sid,
		Status:       status,
		CreatedAt:    time.Now(),
	}
}
