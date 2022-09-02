package library_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsebasct/account-api-lib/library"
	"github.com/jsebasct/account-api-lib/models"
	"os"
	"testing"
)

func getAccountFromFile(filePathName string) (*models.AccountBodyRequest, error) {
	requestFromFile, err := os.ReadFile(filePathName)
	if err != nil {
		return nil, err
	}

	var accountRequest models.AccountBodyRequest
	err = json.Unmarshal(requestFromFile, &accountRequest)
	if err != nil {
		return nil, err
	}
	return &accountRequest, nil
}

func TestCreateAccountSuccess(t *testing.T) {
	accountRequest, err := getAccountFromFile("samples/account_request.json")
	if err != nil {
		t.Error(err)
	}

	createResponse, createError := library.CreateAccount(accountRequest)
	if createError != nil {
		t.Error(errors.New(createError.Message))
	}
	fmt.Printf("%+v\n", createResponse)

	// test amount
	accountsResponse := models.AccountListResponse{}
	err = library.ListAccounts(&accountsResponse)
	if err != nil {
		t.Error(err.Error())
	}
	if len(accountsResponse.Data) < 1 {
		t.Errorf("expected greater or equal to 1 received: %d", len(accountsResponse.Data))
	}

	// test id
	sameIds := 0
	for _, acc := range accountsResponse.Data {
		if acc.Id == accountRequest.Data.Id {
			sameIds += 1
		}
	}
	if sameIds != 1 {
		t.Errorf("expected equal to 1, but received: %d", sameIds)
	}
}

func TestFetchAccount(t *testing.T) {
	accountRequest, err := getAccountFromFile("samples/account_request.json")
	if err != nil {
		t.Error(err)
	}

	account, errFetch := library.FetchAccount(accountRequest.Data.Id)
	if errFetch != nil {
		t.Error(errors.New(errFetch.Message))
	}

	// check id
	if account.Data.Id != accountRequest.Data.Id {
		t.Error(errors.New("account does not have same id"))
	}
	fmt.Println(account)
}

func TestDeleteAccount(t *testing.T) {
	accountFromFile, err := getAccountFromFile("samples/account_request.json")
	if err != nil {
		t.Error(err)
	}

	deleteError := library.DeleteAccount(accountFromFile.Data.Id, accountFromFile.Data.Version)
	if deleteError != nil {
		t.Errorf("Account: %s with version: %d, can't be deleted", accountFromFile.Data.Id, accountFromFile.Data.Version)
	}
}
