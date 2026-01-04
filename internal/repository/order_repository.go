package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/idlistic/go-backend-api-sample/internal/model"
)

var (
	ErrTimeslotNotFound    = errors.New("timeslot not found")
	ErrTimeslotInactive    = errors.New("timeslot is inactive")
	ErrTimeslotFullyBooked = errors.New("timeslot is fully booked")
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderNotCancellable = errors.New("order is not cancellable")
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateWithTimeslotReservation(
	ctx context.Context,
	branchID int64,
	timeslotID int64,
	customerName string,
) (model.Order, error) {

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		// PostgreSQL default is Read Committed; fine for this flow when we lock the row
	})
	if err != nil {
		return model.Order{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// 1) Lock the timeslot row
	var capacity, reserved int
	var isActive bool

	const lockQ = `
SELECT capacity, reserved, is_active
FROM timeslots
WHERE id = $1 AND branch_id = $2
FOR UPDATE;
`
	err = tx.QueryRowContext(ctx, lockQ, timeslotID, branchID).Scan(&capacity, &reserved, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrTimeslotNotFound
		}
		return model.Order{}, err
	}

	if !isActive {
		return model.Order{}, ErrTimeslotInactive
	}
	if reserved >= capacity {
		return model.Order{}, ErrTimeslotFullyBooked
	}

	// 2) Reserve: reserved + 1
	const reserveQ = `
UPDATE timeslots
SET reserved = reserved + 1,
    updated_at = now()
WHERE id = $1 AND branch_id = $2;
`
	if _, err := tx.ExecContext(ctx, reserveQ, timeslotID, branchID); err != nil {
		return model.Order{}, err
	}

	// 3) Create order
	const insertQ = `
INSERT INTO orders (branch_id, timeslot_id, customer_name, status)
VALUES ($1, $2, $3, 'created')
RETURNING id, branch_id, timeslot_id, customer_name, status, created_at, updated_at;
`
	var out model.Order
	if err := tx.QueryRowContext(ctx, insertQ, branchID, timeslotID, customerName).Scan(
		&out.ID,
		&out.BranchID,
		&out.TimeslotID,
		&out.CustomerName,
		&out.Status,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		return model.Order{}, err
	}

	// 4) Commit
	if err := tx.Commit(); err != nil {
		return model.Order{}, err
	}

	return out, nil
}

func (r *OrderRepository) CancelAndReleaseTimeslot(
	ctx context.Context,
	orderID int64,
) (model.Order, error) {

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return model.Order{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// 1) Lock order row
	var out model.Order
	const lockOrderQ = `
SELECT id, branch_id, timeslot_id, customer_name, status, created_at, updated_at
FROM orders
WHERE id = $1
FOR UPDATE;
`
	if err := tx.QueryRowContext(ctx, lockOrderQ, orderID).Scan(
		&out.ID,
		&out.BranchID,
		&out.TimeslotID,
		&out.CustomerName,
		&out.Status,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrOrderNotFound
		}
		return model.Order{}, err
	}

	// only "created" can be cancelled
	if out.Status != "created" {
		return model.Order{}, ErrOrderNotCancellable
	}

	// 2) Lock timeslot row and ensure reserved > 0
	var reserved int
	const lockTimeslotQ = `
SELECT reserved
FROM timeslots
WHERE id = $1 AND branch_id = $2
FOR UPDATE;
`
	if err := tx.QueryRowContext(ctx, lockTimeslotQ, out.TimeslotID, out.BranchID).Scan(&reserved); err != nil {
		// timeslot missing shouldn't happen in demo, but treat as not found timeslot
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrTimeslotNotFound
		}
		return model.Order{}, err
	}

	// 3) Update order -> cancelled
	const cancelOrderQ = `
UPDATE orders
SET status = 'cancelled',
    updated_at = now()
WHERE id = $1;
`
	if _, err := tx.ExecContext(ctx, cancelOrderQ, orderID); err != nil {
		return model.Order{}, err
	}

	// 4) Release reserved (guard: never below 0)
	const releaseQ = `
UPDATE timeslots
SET reserved = CASE WHEN reserved > 0 THEN reserved - 1 ELSE 0 END,
    updated_at = now()
WHERE id = $1 AND branch_id = $2;
`
	if _, err := tx.ExecContext(ctx, releaseQ, out.TimeslotID, out.BranchID); err != nil {
		return model.Order{}, err
	}

	// refresh order status in response
	out.Status = "cancelled"

	if err := tx.Commit(); err != nil {
		return model.Order{}, err
	}

	return out, nil
}

func (r *OrderRepository) ListByBranchAndDate(
	ctx context.Context,
	branchID int64,
	date string, // YYYY-MM-DD
) ([]model.Order, error) {

	// orders ไม่มี service_date -> join timeslots เพื่อ filter ตามวันที่
	const q = `
SELECT
  o.id, o.branch_id, o.timeslot_id, o.customer_name, o.status, o.created_at, o.updated_at
FROM orders o
JOIN timeslots t
  ON t.id = o.timeslot_id
 AND t.branch_id = o.branch_id
WHERE o.branch_id = $1
  AND t.service_date = $2::date
ORDER BY t.start_time ASC, o.created_at ASC;
`

	rows, err := r.db.QueryContext(ctx, q, branchID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Order, 0, 32)
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.ID,
			&o.BranchID,
			&o.TimeslotID,
			&o.CustomerName,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
