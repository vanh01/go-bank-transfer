package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"main/db"
	"main/utils"
)

func TestCreateTransfer(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	a := db.Account{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccount.Id, "Id is nil")
	require.Equal(t, newAccount.Username, a.Username, "Username is not match")
	require.Equal(t, newAccount.Password, a.Password, "Password is not match")

	accountNumber := db.AccountNumber{
		Number:    utils.RandomString(20),
		AccountId: newAccount.Id,
		Balance:   utils.RandomInt(100),
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id, "Id is nil")
	require.Equal(t, newAccount.Id, newAccountNumber.AccountId, "AccountId is not match")

	gotAccountNumber, err := testQueries.GetAccountNumber(context.Background(), newAccountNumber.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, gotAccountNumber.Id, "Id is nil")
	require.Equal(t, newAccountNumber.Id, gotAccountNumber.Id, "Id is not match")
	require.Equal(t, newAccountNumber.Number, gotAccountNumber.Number, "Number is not match")
	require.Equal(t, newAccountNumber.Balance, gotAccountNumber.Balance, "Balance is not match")
	require.Equal(t, newAccountNumber.AccountId, gotAccountNumber.AccountId, "AccountId is not match")

	b := db.Account{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccountB, err := testQueries.CreateAccount(context.Background(), b)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountB.Id, "Id is nil")
	require.Equal(t, newAccountB.Username, b.Username, "Username is not match")
	require.Equal(t, newAccountB.Password, b.Password, "Password is not match")

	accountNumberB := db.AccountNumber{
		Number:    utils.RandomString(20),
		AccountId: newAccountB.Id,
		Balance:   utils.RandomInt(100),
	}
	newAccountNumberB, err := testQueries.CreateAccountNumber(context.Background(), accountNumberB)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumberB.Id, "Id is nil")
	require.Equal(t, newAccountB.Id, newAccountNumberB.AccountId, "AccountId is not match")

	gotAccountNumberB, err := testQueries.GetAccountNumber(context.Background(), newAccountNumberB.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, gotAccountNumberB.Id, "Id is nil")
	require.Equal(t, newAccountNumberB.Id, gotAccountNumberB.Id, "Id is not match")
	require.Equal(t, newAccountNumberB.Number, gotAccountNumberB.Number, "Number is not match")
	require.Equal(t, newAccountNumberB.Balance, gotAccountNumberB.Balance, "Balance is not match")
	require.Equal(t, newAccountNumberB.AccountId, gotAccountNumberB.AccountId, "AccountId is not match")

	transfer := db.Transfer{
		From:   gotAccountNumber.Id,
		To:     gotAccountNumberB.Id,
		Amount: utils.RandomInt(gotAccountNumber.Balance),
	}
	newTransfer, err := testQueries.CreateTransfer(context.Background(), transfer)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newTransfer.Id, "Id is nil")
	require.Equal(t, gotAccountNumber.Id, transfer.From, "From Id is not match")
	require.Equal(t, gotAccountNumberB.Id, transfer.To, "To Id is not match")
}
