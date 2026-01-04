package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/idlistic/go-backend-api-sample/internal/model"
)

type TimetableRepository struct {
	db *sql.DB
}

func NewTimetableRepository(db *sql.DB) *TimetableRepository {
	return &TimetableRepository{db: db}
}

func (r *TimetableRepository) GetOrdersTimetable(
	ctx context.Context,
	branchID int64,
	date string, // YYYY-MM-DD
) ([]model.TimetableItem, error) {

	// LEFT JOIN เพื่อให้ timeslot ที่ไม่มี order ก็ยังออกมา (orders = [])
	const q = `
SELECT
  t.id,
  t.start_time,
  t.end_time,
  t.capacity,
  t.reserved,
  t.is_active,

  o.id AS order_id,
  o.customer_name,
  o.status,
  o.created_at
FROM timeslots t
LEFT JOIN orders o
  ON o.timeslot_id = t.id
 AND o.branch_id = t.branch_id
 AND o.status = 'created'
WHERE t.branch_id = $1
  AND t.service_date = $2::date
ORDER BY t.start_time ASC, o.created_at ASC;
`

	rows, err := r.db.QueryContext(ctx, q, branchID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// group by timeslot_id
	items := make([]model.TimetableItem, 0, 16)
	indexByTimeslot := make(map[int64]int, 16)

	for rows.Next() {
		var (
			tsID      int64
			startTime string
			endTime   string
			capacity  int
			reserved  int
			isActive  bool

			orderID       sql.NullInt64
			customerName  sql.NullString
			status        sql.NullString
			createdAtTime sql.NullTime
		)

		if err := rows.Scan(
			&tsID,
			&startTime,
			&endTime,
			&capacity,
			&reserved,
			&isActive,
			&orderID,
			&customerName,
			&status,
			&createdAtTime,
		); err != nil {
			return nil, err
		}

		pos, ok := indexByTimeslot[tsID]
		if !ok {
			items = append(items, model.TimetableItem{
				Timeslot: model.TimetableTimeslot{
					ID:        tsID,
					StartTime: startTime,
					EndTime:   endTime,
					Capacity:  capacity,
					Reserved:  reserved,
					IsActive:  isActive,
				},
				Orders: make([]model.TimetableOrder, 0, 4),
			})
			pos = len(items) - 1
			indexByTimeslot[tsID] = pos
		}

		// order อาจเป็น NULL (เพราะ LEFT JOIN)
		if orderID.Valid {
			createdAt := ""
			if createdAtTime.Valid {
				createdAt = createdAtTime.Time.UTC().Format(time.RFC3339)
			}

			items[pos].Orders = append(items[pos].Orders, model.TimetableOrder{
				ID:           orderID.Int64,
				CustomerName: customerName.String,
				Status:       status.String,
				CreatedAt:    createdAt,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
