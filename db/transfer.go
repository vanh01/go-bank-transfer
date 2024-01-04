package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	Id        uuid.UUID `json:"id"`
	From      uuid.UUID `json:"from"`
	To        uuid.UUID `json:"to"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `json:"-"`
}

const createTransfer = `
INSERT INTO bank.transfers ("from", "to", "amount")
VALUES($1, $2, $3)
RETURNING id, "from", "to", "amount", created_at, updated_at, is_deleted
`

func (q *Queries) CreateTransfer(ctx context.Context, t Transfer) (Transfer, error) {
	row := q.DB.QueryRowContext(ctx, createTransfer, t.From, t.To, t.Amount)
	if row.Err() != nil {
		return Transfer{}, row.Err()
	}
	var transfer Transfer
	err := row.Scan(
		&transfer.Id,
		&transfer.From,
		&transfer.To,
		&transfer.Amount,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.IsDeleted,
	)
	return transfer, err
}

const getTransfer = `
SELECT id, from, to, amount, created_at, updated_at, is_deleted
FROM bank.transfers
WHERE id = $1
`

func (q *Queries) GetTransfer(ctx context.Context, id uuid.UUID) (Transfer, error) {
	row := q.DB.QueryRowContext(ctx, getTransfer, id)
	if row.Err() != nil {
		return Transfer{}, row.Err()
	}
	var transfer Transfer
	err := row.Scan(
		&transfer.Id,
		&transfer.From,
		&transfer.To,
		&transfer.Amount,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.IsDeleted,
	)
	return transfer, err
}
