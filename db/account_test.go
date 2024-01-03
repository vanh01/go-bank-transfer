package db_test

import (
	"context"
	"main/db"
	"main/utils"
	"testing"

	"github.com/google/uuid"
)

var testQueries *db.Queries

func TestCreateAccount(t *testing.T) {
	testQueries = &db.Queries{
		DB: db.NewDB(),
	}
	a := db.Account{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if newAccount.Id == uuid.Nil {
		t.Fatal("ID is nil")
	}
	if newAccount.Username != a.Username || newAccount.Password != a.Password {
		t.Fatalf("Old username:%s\nNew username%s\nOld password:%s\nNew password%s\n", a.Username, newAccount.Username, a.Password, newAccount.Password)
	}
}

func TestGetAccount(t *testing.T) {
	testQueries = &db.Queries{
		DB: db.NewDB(),
	}
	a := db.Account{
		Username: utils.RandomString(7),
		Password: utils.RandomString(7),
	}
	newAccount, err := testQueries.CreateAccount(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if newAccount.Id == uuid.Nil {
		t.Fatal("ID is nil")
	}
	if newAccount.Username != a.Username || newAccount.Password != a.Password {
		t.Fatalf("Old username:%s\nNew username%s\nOld password:%s\nNew password%s\n", a.Username, newAccount.Username, a.Password, newAccount.Password)
	}
	account, err := testQueries.GetAccount(context.Background(), newAccount.Id)
	if err != nil {
		t.Fatal(err)
	}
	if account.Username != newAccount.Username {
		t.Fatalf("Username is incorrect!\nCreated username:%s\nGot username%s\n", newAccount.Username, account.Username)
	}
	if account.Password != newAccount.Password {
		t.Fatalf("Password is incorrect!\nCreated password:%s\nGot password%s\n", newAccount.Password, account.Password)
	}
}
