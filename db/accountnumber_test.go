package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"main/db"
	"main/utils"
)

func TestCreateAccountNumber(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	a := db.CreateAccountParams{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccount.Id)
	require.Equal(t, newAccount.Username, a.Username)
	require.Equal(t, newAccount.Password, a.Password)

	accountNumber := db.CreateAccountNumberParams{
		Number:    utils.RandomString(20),
		Balance:   utils.RandomInt(1000),
		AccountId: newAccount.Id,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id)
	require.Equal(t, newAccount.Id, newAccountNumber.AccountId)
}

func TestGetAccountNumber(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	a := db.CreateAccountParams{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccount.Id)
	require.Equal(t, newAccount.Username, a.Username)
	require.Equal(t, newAccount.Password, a.Password)

	accountNumber := db.CreateAccountNumberParams{
		Number:    utils.RandomString(20),
		AccountId: newAccount.Id,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id)
	require.Equal(t, newAccount.Id, newAccountNumber.AccountId)

	gotAccountNumber, err := testQueries.GetAccountNumber(context.Background(), newAccountNumber.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, gotAccountNumber.Id)
	require.Equal(t, newAccountNumber.Id, gotAccountNumber.Id)
	require.Equal(t, newAccountNumber.Number, gotAccountNumber.Number)
	require.Equal(t, newAccountNumber.Balance, gotAccountNumber.Balance)
	require.Equal(t, newAccountNumber.AccountId, gotAccountNumber.AccountId)
}

func TestUpdateBalanceAccountNumber(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	a := db.CreateAccountParams{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccount.Id)
	require.Equal(t, newAccount.Username, a.Username)
	require.Equal(t, newAccount.Password, a.Password)

	accountNumber := db.CreateAccountNumberParams{
		Number:    utils.RandomString(20),
		Balance:   utils.RandomInt(1000),
		AccountId: newAccount.Id,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id)
	require.Equal(t, newAccount.Id, newAccountNumber.AccountId)

	updateBalanceAccountNumberParams := db.UpdateBalanceAccountNumberParams{
		Id:     newAccountNumber.Id,
		Amount: 10,
	}
	gotAccountNumber, err := testQueries.UpdateBalanceAccountNumber(context.Background(), updateBalanceAccountNumberParams)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.Equal(t, gotAccountNumber.Balance, newAccountNumber.Balance+10)
}
