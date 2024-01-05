package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"main/db"
	"main/utils"
)

func TestCreateHistory(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	account := createAccountRandom(t)

	accountNumber := createAccountNumberRandom(t, account.Id)
	amount := utils.RandomInt(accountNumber.Balance)
	balance := accountNumber.Balance
	h := db.History{
		Balance:         accountNumber.Balance,
		Amount:          amount,
		AccountNumberId: accountNumber.Id,
	}
	newHistory, err := testQueries.CreateHistory(context.Background(), h)

	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newHistory.Id, "Id is nil")
	require.Equal(t, accountNumber.Id, newHistory.AccountNumberId)
	require.Equal(t, amount, newHistory.Amount)
	require.Equal(t, balance, newHistory.Balance)
}
