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

type CreateTransferParams struct {
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Amount int64     `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, param CreateTransferParams) (Transfer, error) {
	row := q.DB.QueryRowContext(ctx, createTransfer, param.From, param.To, param.Amount)
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
SELECT id, "from", "to", "amount", created_at, updated_at, is_deleted
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

type TransferMoneyParams struct {
	From   uuid.UUID
	To     uuid.UUID
	Amount int64
}

type TransferMoneyResponse struct {
	From     AccountNumber
	To       AccountNumber
	Transfer Transfer
}

func (q *Queries) TransferMoney(ctx context.Context, param TransferMoneyParams) (*TransferMoneyResponse, error) {
	response := &TransferMoneyResponse{}
	err := q.ExecuteTx(context.Background(), func(q *Queries) error {
		// get from account number
		fromAccountNumber, err := q.GetAccountNumber(context.Background(), param.From)
		if err != nil {
			return err
		}

		// get to account number
		toAccountNumber, err := q.GetAccountNumber(context.Background(), param.To)
		if err != nil {
			return err
		}

		// create history
		fromHistory := CreateHistoryParams{
			Balance:         fromAccountNumber.Balance,
			Amount:          -param.Amount,
			AccountNumberId: fromAccountNumber.Id,
		}
		_, err = q.CreateHistory(context.Background(), fromHistory)
		if err != nil {
			return err
		}

		// create history
		toHistory := CreateHistoryParams{
			Balance:         toAccountNumber.Balance,
			Amount:          param.Amount,
			AccountNumberId: toAccountNumber.Id,
		}
		_, err = q.CreateHistory(context.Background(), toHistory)
		if err != nil {
			return err
		}

		// create a createTransfer
		createTransfer := CreateTransferParams{
			From:   fromAccountNumber.Id,
			To:     toAccountNumber.Id,
			Amount: param.Amount,
		}
		newTransfer, err := q.CreateTransfer(context.Background(), createTransfer)
		if err != nil {
			return err
		}
		response.Transfer = newTransfer

		updatebalance := func(ctx context.Context, from, to UpdateBalanceAccountNumberParams) (AccountNumber, AccountNumber, error) {
			var account1, account2 AccountNumber
			if from.Id.String() < to.Id.String() {
				account1, err = q.UpdateBalanceAccountNumber(context.Background(), from)
				if err != nil {
					return account1, account2, err
				}
				account2, err = q.UpdateBalanceAccountNumber(context.Background(), to)
				if err != nil {
					return account1, account2, err
				}
			} else {
				account2, err = q.UpdateBalanceAccountNumber(context.Background(), to)
				if err != nil {
					return account1, account2, err
				}
				account1, err = q.UpdateBalanceAccountNumber(context.Background(), from)
				if err != nil {
					return account1, account2, err
				}
			}
			return account1, account2, nil
		}
		updateBalanceAccountNumberParamsA := UpdateBalanceAccountNumberParams{
			Id:     fromAccountNumber.Id,
			Amount: -param.Amount,
		}
		updateBalanceAccountNumberParamsB := UpdateBalanceAccountNumberParams{
			Id:     toAccountNumber.Id,
			Amount: param.Amount,
		}
		response.From, response.To, err = updatebalance(context.Background(), updateBalanceAccountNumberParamsA, updateBalanceAccountNumberParamsB)

		return nil
	})
	return response, err
}
