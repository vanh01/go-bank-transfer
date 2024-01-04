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
	a := db.Account{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccount.Id, "Id is nil")
	require.Equal(t, newAccount.Username, a.Username, "Username is not match")
	require.Equal(t, newAccount.Password, a.Password, "Password is not match")
}

func TestGetAccount(t *testing.T) {
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

	account, err := testQueries.GetAccount(context.Background(), newAccount.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.Equal(t, newAccount.Username, account.Username, "Username is not match")
	require.Equal(t, newAccount.Password, account.Password, "Password is not match")
}

func TestGetFirstAccount(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}

	account, err := testQueries.GetFirstAccount(context.Background())
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, account.Id, "Id is nil")
	require.NotEmpty(t, account.Username, "", "Username is empty")
	require.NotEmpty(t, account.Password, "", "Password is empty")
}

func TestGetSecondAccount(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}

	account, err := testQueries.GetSecondAccount(context.Background())
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, account.Id, "Id is nil")
	require.NotEmpty(t, account.Username, "", "Username is empty")
	require.NotEmpty(t, account.Password, "", "Password is empty")
}
