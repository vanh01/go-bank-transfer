package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"main/db"
	"main/utils"
)

func TestCreateAccount(t *testing.T) {
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
}

func TestGetAccount(t *testing.T) {
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

	account, err := testQueries.GetAccount(context.Background(), newAccount.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.Equal(t, newAccount.Username, account.Username)
	require.Equal(t, newAccount.Password, account.Password)
}

func TestGetFirstAccount(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}

	account, err := testQueries.GetFirstAccount(context.Background())
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, account.Id)
	require.NotEmpty(t, account.Username)
	require.NotEmpty(t, account.Password)
}

func TestGetSecondAccount(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}

	account, err := testQueries.GetSecondAccount(context.Background())
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, account.Id)
	require.NotEmpty(t, account.Username)
	require.NotEmpty(t, account.Password)
}
