package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `json:"-"`
}

const createAccount = `
INSERT INTO bank.accounts (username, "password")
VALUES($1, $2)
RETURNING id, username, password, created_at, updated_at, is_deleted
`

func (q *Queries) CreateAccount(ctx context.Context, a Account) (Account, error) {
	row := q.DB.QueryRowContext(ctx, createAccount, a.Username, a.Password)
	if row.Err() != nil {
		return Account{}, row.Err()
	}
	var newAccount Account
	err := row.Scan(&newAccount.Id, &newAccount.Username, &newAccount.Password, &newAccount.CreatedAt, &newAccount.UpdatedAt, &newAccount.IsDeleted)
	return newAccount, err
}

const getAccount = `
SELECT id, username, password, created_at, updated_at, is_deleted
FROM bank.accounts
WHERE id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id uuid.UUID) (Account, error) {
	row := q.DB.QueryRowContext(ctx, getAccount, id)
	if row.Err() != nil {
		return Account{}, row.Err()
	}
	var a Account
	err := row.Scan(&a.Id, &a.Username, &a.Password, &a.CreatedAt, &a.UpdatedAt, &a.IsDeleted)
	return a, err
}

const getFirstAccount = `
SELECT id, username, password, created_at, updated_at, is_deleted
FROM bank.accounts
ORDER BY created_at
LIMIT 1
`

func (q *Queries) GetFirstAccount(ctx context.Context) (Account, error) {
	row := q.DB.QueryRowContext(ctx, getFirstAccount)
	if row.Err() != nil {
		return Account{}, row.Err()
	}
	var a Account
	err := row.Scan(&a.Id, &a.Username, &a.Password, &a.CreatedAt, &a.UpdatedAt, &a.IsDeleted)
	return a, err
}

const getSecondAccount = `
SELECT id, username, password, created_at, updated_at, is_deleted
FROM bank.accounts
ORDER BY created_at
OFFSET 1
LIMIT 1
`

func (q *Queries) GetSecondAccount(ctx context.Context) (Account, error) {
	row := q.DB.QueryRowContext(ctx, getSecondAccount)
	if row.Err() != nil {
		return Account{}, row.Err()
	}
	var a Account
	err := row.Scan(&a.Id, &a.Username, &a.Password, &a.CreatedAt, &a.UpdatedAt, &a.IsDeleted)
	return a, err
}
