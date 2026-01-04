package repository

import (
	"context"
	"database/sql"

	"github.com/idlistic/go-backend-api-sample/internal/model"
)

type TimeslotRepository struct {
	db *sql.DB
}

func NewTimeslotRepository(db *sql.DB) *TimeslotRepository {
	return &TimeslotRepository{db: db}
}

func (r *TimeslotRepository) ListByBranchAndDate(
	ctx context.Context,
	branchID int64,
	date string, // YYYY-MM-DD
) ([]model.Timeslot, error) {

	const q = `
SELECT
  id, branch_id, service_date, start_time, end_time,
  capacity, reserved, is_active, created_at, updated_at
FROM timeslots
WHERE branch_id = $1
  AND service_date = $2::date
ORDER BY start_time ASC;
`

	rows, err := r.db.QueryContext(ctx, q, branchID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Timeslot, 0, 16)

	for rows.Next() {
		var t model.Timeslot
		if err := rows.Scan(
			&t.ID,
			&t.BranchID,
			&t.ServiceDate,
			&t.StartTime,
			&t.EndTime,
			&t.Capacity,
			&t.Reserved,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
