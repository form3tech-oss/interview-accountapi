package library_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsebasct/account-api-lib/library"
	"github.com/jsebasct/account-api-lib/models"
	"net/http"
	"os"
	"strings"
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
		fmt.Printf("Error while creating: %+v\n", createError)
		t.Error(errors.New(createError.Message))
	}
	fmt.Printf("Created response: %+v\n", createResponse)

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

func TestFetchAccountSuccess(t *testing.T) {
	accountRequest, err := getAccountFromFile("samples/account_request.json")
	if err != nil {
		t.Error(err)
	}

	account, errFetch := library.FetchAccount(accountRequest.Data.Id)
	if errFetch != nil {
		t.Error(errors.New(errFetch.Message))
	}

	if account.Data.Id != accountRequest.Data.Id {
		t.Error(errors.New("account does not have same id"))
	}
}

func TestFetchUnexistingAccount(t *testing.T) {
	unexistingAccountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4cc"

	account, errFetch := library.FetchAccount(unexistingAccountId)
	if account != nil {
		fmt.Println("account", account)
		t.Error(errors.New("the account should be nil since ID doesn't exist"))
	}
	if errFetch != nil {
		if errFetch.Code != http.StatusNotFound {
			t.Error(errors.New("the response code should be 404"))
		}
		if !strings.ContainsAny(errFetch.Message, "does not exist") {
			t.Error(errors.New("wrong message for not found"))
		}
		fmt.Println("message", errFetch)
	}
}

func TestFetchInvalidId(t *testing.T) {
	unexistingAccountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4cc_11111"

	account, errFetch := library.FetchAccount(unexistingAccountId)
	if account != nil {
		fmt.Println("account", account)
		t.Error(errors.New("the account should be nil since ID doesn't exist"))
	}
	if errFetch != nil {
		if errFetch.Code != http.StatusBadRequest {
			t.Error(errors.New("the response code should be 400"))
		}
		if !strings.ContainsAny(errFetch.Message, "does not exist") {
			t.Error(errors.New("wrong message for not found"))
		}
		fmt.Println("message", errFetch)
	}
}

func TestDeleteAccountInvalidVersion(t *testing.T) {
	accountFromFile, err := getAccountFromFile("samples/account_request.json")
	if err != nil {
		t.Error(err)
	}

	deleteError := library.DeleteAccount(accountFromFile.Data.Id, accountFromFile.Data.Version+10)

	if deleteError == nil {
		t.Error("should have returned an error")
	}
	fmt.Printf("%+v", deleteError)

	if deleteError.Code != http.StatusConflict {
		t.Errorf("should have returned a status %d", http.StatusConflict)
	}

	expectedErrMsg := "invalid version"
	if expectedErrMsg != deleteError.Message {
		t.Errorf("message should be: %s", expectedErrMsg)
	}
}

func TestDeleteAccountSuccess(t *testing.T) {
	accountFromFile, err := getAccountFromFile("samples/account_request.json")
	if err != nil {
		t.Error(err)
	}

	deleteError := library.DeleteAccount(accountFromFile.Data.Id, accountFromFile.Data.Version)
	if deleteError != nil {
		t.Errorf("Account: %s with version: %d, can't be deleted", accountFromFile.Data.Id, accountFromFile.Data.Version)
	}
}

func TestDeleteAccountNotFound(t *testing.T) {
	accountId := "ad27e265-9605-4b4b-a0e5-3003ea9cc4ff"
	version := 0
	deleteError := library.DeleteAccount(accountId, version)

	if deleteError == nil {
		t.Error("should have returned an error")
	}
	fmt.Printf("%+v", deleteError)

	if deleteError.Code != http.StatusNotFound {
		t.Errorf("should have returned a status %d", http.StatusNotFound)
	}
}
