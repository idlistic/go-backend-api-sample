package model

import "time"

type Order struct {
	ID           int64     `json:"id"`
	BranchID     int64     `json:"branch_id"`
	TimeslotID   int64     `json:"timeslot_id"`
	CustomerName string    `json:"customer_name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
