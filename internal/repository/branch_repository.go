package repository

import (
	"context"
	"database/sql"

	"github.com/idlistic/go-backend-api-sample/internal/model"
)

type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) List(ctx context.Context) ([]model.Branch, error) {
	const q = `
SELECT id, name, created_at
FROM branches
ORDER BY id ASC;
`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Branch, 0, 16)
	for rows.Next() {
		var b model.Branch
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
