package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountNumber struct {
	Id        uuid.UUID `json:"id"`
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	AccountId uuid.UUID `json:"account_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `json:"-"`
}

const createAccountNumber = `
INSERT INTO bank.account_numbers (number, account_id, balance)
VALUES($1, $2, $3)
RETURNING id, number, balance, account_id, created_at, updated_at, is_deleted
`

type CreateAccountNumberParams struct {
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	AccountId uuid.UUID `json:"account_id"`
}

func (q *Queries) CreateAccountNumber(ctx context.Context, param CreateAccountNumberParams) (AccountNumber, error) {
	row := q.DB.QueryRowContext(ctx, createAccountNumber, param.Number, param.AccountId, param.Balance)
	if row.Err() != nil {
		return AccountNumber{}, row.Err()
	}
	var newAccountNumber AccountNumber
	err := row.Scan(
		&newAccountNumber.Id,
		&newAccountNumber.Number,
		&newAccountNumber.Balance,
		&newAccountNumber.AccountId,
		&newAccountNumber.CreatedAt,
		&newAccountNumber.UpdatedAt,
		&newAccountNumber.IsDeleted)
	return newAccountNumber, err
}

const getAccountNumber = `
SELECT id, number, balance, account_id, created_at, updated_at, is_deleted
FROM bank.account_numbers
WHERE id = $1
`

func (q *Queries) GetAccountNumber(ctx context.Context, id uuid.UUID) (AccountNumber, error) {
	row := q.DB.QueryRowContext(ctx, getAccountNumber, id)
	if row.Err() != nil {
		return AccountNumber{}, row.Err()
	}
	var accountNumber AccountNumber
	err := row.Scan(
		&accountNumber.Id,
		&accountNumber.Number,
		&accountNumber.Balance,
		&accountNumber.AccountId,
		&accountNumber.CreatedAt,
		&accountNumber.UpdatedAt,
		&accountNumber.IsDeleted)
	return accountNumber, err
}

const updateBalanceAccountNumber = `
UPDATE bank.account_numbers
SET balance = balance + $2
WHERE id = $1
RETURNING id, number, balance, account_id, created_at, updated_at, is_deleted
`

type UpdateBalanceAccountNumberParams struct {
	Id     uuid.UUID `json:"id"`
	Amount int64     `json:"amount"`
}

func (q *Queries) UpdateBalanceAccountNumber(ctx context.Context, param UpdateBalanceAccountNumberParams) (AccountNumber, error) {
	row := q.DB.QueryRowContext(ctx, updateBalanceAccountNumber, param.Id, param.Amount)
	if row.Err() != nil {
		return AccountNumber{}, row.Err()
	}
	var accountNumber AccountNumber
	err := row.Scan(
		&accountNumber.Id,
		&accountNumber.Number,
		&accountNumber.Balance,
		&accountNumber.AccountId,
		&accountNumber.CreatedAt,
		&accountNumber.UpdatedAt,
		&accountNumber.IsDeleted)
	return accountNumber, err
}
