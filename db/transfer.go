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

type TransferMoneyResponse struct {
	From     AccountNumber
	To       AccountNumber
	Transfer Transfer
}

func (q *Queries) TransferMoney(ctx context.Context, from, to uuid.UUID, amount int64) (*TransferMoneyResponse, error) {
	response := &TransferMoneyResponse{}
	err := q.ExecuteTx(context.Background(), func(q *Queries) error {
		// get from account number
		fromAccountNumber, err := q.GetAccountNumber(context.Background(), from)
		if err != nil {
			return err
		}

		// get to account number
		toAccountNumber, err := q.GetAccountNumber(context.Background(), to)
		if err != nil {
			return err
		}

		// create a transfer
		transfer := Transfer{
			From:   fromAccountNumber.Id,
			To:     toAccountNumber.Id,
			Amount: amount,
		}
		newTransfer, err := q.CreateTransfer(context.Background(), transfer)
		if err != nil {
			return err
		}
		response.Transfer = newTransfer

		// update balance of first account number
		response.From, err = q.UpdateBalanceAccountNumber(context.Background(), fromAccountNumber.Id, -amount)
		if err != nil {
			return err
		}

		// update balance of second account number
		response.To, err = q.UpdateBalanceAccountNumber(context.Background(), toAccountNumber.Id, amount)
		if err != nil {
			return err
		}

		return nil
	})
	return response, err
}
