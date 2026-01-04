package model

import "time"

type Timeslot struct {
	ID          int64     `json:"id"`
	BranchID    int64     `json:"branch_id"`
	ServiceDate string    `json:"service_date"` // YYYY-MM-DD
	StartTime   string    `json:"start_time"`   // HH:MM:SS
	EndTime     string    `json:"end_time"`     // HH:MM:SS
	Capacity    int       `json:"capacity"`
	Reserved    int       `json:"reserved"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
