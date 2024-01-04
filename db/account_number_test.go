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
		Balance:   utils.RandomInt(1000),
		AccountId: newAccount.Id,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id, "Id is nil")
	require.Equal(t, newAccount.Id, newAccountNumber.AccountId, "AccountId is not match")
}

func TestGetAccountNumber(t *testing.T) {
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
}

func TestUpdateBalanceAccountNumber(t *testing.T) {
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
		Balance:   utils.RandomInt(1000),
		AccountId: newAccount.Id,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id, "Id is nil")
	require.Equal(t, newAccount.Id, newAccountNumber.AccountId, "AccountId is not match")

	err = testQueries.UpdateBalanceAccountNumber(context.Background(), newAccountNumber.Id, 10)
	require.Nilf(t, err, "An error occur: %s\n", err)

	gotAccountNumber, err := testQueries.GetAccountNumber(context.Background(), newAccountNumber.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, gotAccountNumber.Id, "Id is nil")
	require.Equal(t, gotAccountNumber.Balance, newAccountNumber.Balance+10, "Balance is not match")
}
