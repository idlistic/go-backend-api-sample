package model

type TimetableOrder struct {
	ID           int64  `json:"id"`
	CustomerName string `json:"customer_name"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

type TimetableTimeslot struct {
	ID        int64  `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Capacity  int    `json:"capacity"`
	Reserved  int    `json:"reserved"`
	IsActive  bool   `json:"is_active"`
}

type TimetableItem struct {
	Timeslot TimetableTimeslot `json:"timeslot"`
	Orders   []TimetableOrder  `json:"orders"`
}
