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
INSERT INTO public.accounts (username, "password")
VALUES($1, $2)
RETURNING id, username, password, created_at, updated_at, is_deleted
`

func (q *Queries) CreateAccount(ctx context.Context, a Account) (Account, error) {
	row := q.DB.QueryRowContext(ctx, createAccount, a.Username, a.Password)
	var newAccount Account
	err := row.Scan(&newAccount.Id, &newAccount.Username, &newAccount.Password, &newAccount.CreatedAt, &newAccount.UpdatedAt, &newAccount.IsDeleted)
	return newAccount, err
}

const getAccount = `
SELECT *
FROM accounts
WHERE id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id uuid.UUID) (Account, error) {
	row := q.DB.QueryRowContext(ctx, getAccount, id)
	var a Account
	err := row.Scan(&a.Id, &a.Username, &a.Password, &a.CreatedAt, &a.UpdatedAt, &a.IsDeleted)
	return a, err
}
