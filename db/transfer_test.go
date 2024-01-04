package db_test

import (
	"context"
	"fmt"
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

func createAccountRandom(t *testing.T) db.Account {
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
	return newAccount
}

func createAccountNumberRandom(t *testing.T, accountId uuid.UUID) db.AccountNumber {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	accountNumberB := db.AccountNumber{
		Number:    utils.RandomString(20),
		AccountId: accountId,
		Balance:   1000,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumberB)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id, "Id is nil")
	return newAccountNumber
}

func TestTransferMoney(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	accountA := createAccountRandom(t)
	accountB := createAccountRandom(t)

	accountNumberA := createAccountNumberRandom(t, accountA.Id)
	accountNumberB := createAccountNumberRandom(t, accountB.Id)

	amount := int64(10)
	numTx := 90
	errchan := make(chan error, numTx)
	reschan := make(chan db.TransferMoneyResponse, numTx)

	for i := 0; i < numTx; i++ {
		go func() {
			response, err := testQueries.TransferMoney(context.Background(), accountNumberA.Id, accountNumberB.Id, amount)
			errchan <- err
			reschan <- *response
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < numTx; i++ {
		err := <-errchan
		require.NoError(t, err)

		res := <-reschan
		require.NotEmpty(t, res)
		require.NotEmpty(t, res.Transfer)
		require.NotEmpty(t, res.From)
		require.NotEmpty(t, res.To)

		fmt.Println(">>tx: ", res.From.Balance, res.To.Balance)

		require.Equal(t, amount, res.Transfer.Amount)
		require.Equal(t, accountNumberA.Id, res.Transfer.From)

		sub := (i + 1) * int(amount) * 2
		require.Equal(t, int64(sub), res.To.Balance-res.From.Balance)

		dif1 := accountNumberA.Balance - res.From.Balance
		dif2 := res.To.Balance - accountNumberB.Balance
		require.Equal(t, dif1, dif2)
		require.True(t, dif1 > 0)
		require.True(t, dif1%amount == 0)

		k := int(dif1 / amount)
		require.True(t, k >= 1 && k <= numTx)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
}
