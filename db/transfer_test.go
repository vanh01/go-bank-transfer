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

func createAccountRandom(t *testing.T) db.Account {
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
	return newAccount
}

func createAccountNumberRandom(t *testing.T, accountId uuid.UUID) db.AccountNumber {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	accountNumber := db.CreateAccountNumberParams{
		Number:    utils.RandomString(20),
		AccountId: accountId,
		Balance:   1000,
	}
	newAccountNumber, err := testQueries.CreateAccountNumber(context.Background(), accountNumber)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumber.Id)
	require.Equal(t, accountNumber.AccountId, newAccountNumber.AccountId)
	return newAccountNumber
}

func TestCreateTransfer(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	accountA := createAccountRandom(t)
	accountB := createAccountRandom(t)

	accountNumberA := createAccountNumberRandom(t, accountA.Id)
	accountNumberB := createAccountNumberRandom(t, accountB.Id)

	transfer := db.CreateTransferParams{
		From:   accountNumberA.Id,
		To:     accountNumberB.Id,
		Amount: utils.RandomInt(accountNumberA.Balance),
	}
	newTransfer, err := testQueries.CreateTransfer(context.Background(), transfer)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newTransfer.Id)
	require.Equal(t, accountNumberA.Id, transfer.From)
	require.Equal(t, accountNumberB.Id, transfer.To)
}

func TestGetTransfer(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	accountA := createAccountRandom(t)
	accountB := createAccountRandom(t)

	accountNumberA := createAccountNumberRandom(t, accountA.Id)
	accountNumberB := createAccountNumberRandom(t, accountB.Id)

	transfer := db.CreateTransferParams{
		From:   accountNumberA.Id,
		To:     accountNumberB.Id,
		Amount: utils.RandomInt(accountNumberA.Balance),
	}
	newTransfer, err := testQueries.CreateTransfer(context.Background(), transfer)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newTransfer.Id)
	require.Equal(t, accountNumberA.Id, transfer.From)
	require.Equal(t, accountNumberB.Id, transfer.To)

	gotTransfer, err := testQueries.GetTransfer(context.Background(), newTransfer.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, gotTransfer.Id)
	require.Equal(t, newTransfer.From, gotTransfer.From)
	require.Equal(t, newTransfer.To, gotTransfer.To)
	require.Equal(t, newTransfer.Amount, gotTransfer.Amount)
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
	numTx := 50
	errchan := make(chan error, numTx)
	reschan := make(chan db.TransferMoneyResponse, numTx)

	for i := 0; i < numTx; i++ {
		go func() {
			param := db.TransferMoneyParams{
				From:   accountNumberA.Id,
				To:     accountNumberB.Id,
				Amount: amount,
			}
			response, err := testQueries.TransferMoney(context.Background(), param)
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

	newAccountNumberA, err := testQueries.GetAccountNumber(context.Background(), accountNumberA.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumberA.Id)

	newAccountNumberB, err := testQueries.GetAccountNumber(context.Background(), accountNumberB.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumberB.Id)

	require.Equal(t, 2*amount*int64(numTx), newAccountNumberB.Balance-newAccountNumberA.Balance)
}

func TestTransferMoneyTwoWays(t *testing.T) {
	testQueries := &db.Queries{
		DB: db.NewDB(),
	}
	accountA := createAccountRandom(t)
	accountB := createAccountRandom(t)

	accountNumberA := createAccountNumberRandom(t, accountA.Id)
	accountNumberB := createAccountNumberRandom(t, accountB.Id)

	amount := int64(10)
	numTx := 10
	errchan := make(chan error, numTx)
	reschan := make(chan db.TransferMoneyResponse, 2*numTx)

	for i := 0; i < 2*numTx; i++ {
		from := accountNumberA.Id
		to := accountNumberB.Id
		if i%2 == 0 {
			from = accountNumberB.Id
			to = accountNumberA.Id
		}
		go func() {
			param := db.TransferMoneyParams{
				From:   from,
				To:     to,
				Amount: amount,
			}
			response, err := testQueries.TransferMoney(context.Background(), param)
			errchan <- err
			reschan <- *response
		}()
	}

	for i := 0; i < 2*numTx; i++ {
		err := <-errchan
		require.NoError(t, err)

		res := <-reschan
		require.NotEmpty(t, res)
		require.NotEmpty(t, res.Transfer)
		require.NotEmpty(t, res.From)
		require.NotEmpty(t, res.To)

		fmt.Println(">>tx: ", res.From.Balance, res.To.Balance)

		require.Equal(t, amount, res.Transfer.Amount)

		dif1 := accountNumberA.Balance - res.From.Balance
		dif2 := res.To.Balance - accountNumberB.Balance
		require.Equal(t, dif1, dif2)
		require.True(t, dif1%amount == 0)
	}

	newAccountNumberA, err := testQueries.GetAccountNumber(context.Background(), accountNumberA.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumberA.Id)

	newAccountNumberB, err := testQueries.GetAccountNumber(context.Background(), accountNumberB.Id)
	require.Nilf(t, err, "An error occur: %s\n", err)
	require.NotEqual(t, uuid.Nil, newAccountNumberB.Id)

	require.Equal(t, accountNumberA.Balance, newAccountNumberA.Balance)
	require.Equal(t, accountNumberB.Balance, newAccountNumberB.Balance)
}
