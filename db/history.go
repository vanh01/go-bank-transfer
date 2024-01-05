package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type History struct {
	Id              uuid.UUID `json:"id"`
	Balance         int64     `json:"balance"`
	Amount          int64     `json:"amount"`
	AccountNumberId uuid.UUID `json:"account_number_id"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	IsDeleted       bool      `json:"-"`
}

const createHistory = `
INSERT INTO bank.histories (balance, amount, account_number_id)
VALUES($2, $3, $1)
RETURNING id, balance, amount, account_number_id, created_at, updated_at, is_deleted
`

type CreateHistoryParams struct {
	Balance         int64     `json:"balance"`
	Amount          int64     `json:"amount"`
	AccountNumberId uuid.UUID `json:"account_number_id"`
}

func (q *Queries) CreateHistory(ctx context.Context, param CreateHistoryParams) (History, error) {
	row := q.DB.QueryRowContext(ctx, createHistory, param.AccountNumberId, param.Balance, param.Amount)
	if row.Err() != nil {
		return History{}, row.Err()
	}
	var newHistory History

	err := row.Scan(
		&newHistory.Id,
		&newHistory.Balance,
		&newHistory.Amount,
		&newHistory.AccountNumberId,
		&newHistory.CreatedAt,
		&newHistory.UpdatedAt,
		&newHistory.IsDeleted,
	)

	return newHistory, err
}
