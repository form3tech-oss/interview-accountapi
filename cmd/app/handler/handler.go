package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/advena/interview-accountapi/cmd/app/account"
)

type AccountsHandler interface {
	Create(account account.Account) (account.Account, error)
	Delete(accountID string) (bool, error)
	Fetch(accountID string) (account.Account, error)
}

type data struct {
	Body account.Account `json:"data"`
}

type form3AccountsHandler struct {
	url    string
	client http.Client
}

func (handler form3AccountsHandler) Create(a account.Account) (account.Account, error) {
	jsonData, jsonErr := json.Marshal(&data{a})
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return account.Account{}, jsonErr
	}

	createAccountResponse, createErr := handler.client.Post(handler.url, "application/json", bytes.NewBuffer(jsonData))
	if createErr != nil {
		log.Fatal(createErr)
		return account.Account{}, createErr
	}

	defer createAccountResponse.Body.Close()

	createdAccount := data{}
	json.NewDecoder(createAccountResponse.Body).Decode(&createdAccount)

	return createdAccount.Body, nil
}

func (handler form3AccountsHandler) Delete(accountID string) (bool, error) {
	deleteAccountRequest, reqErr := http.NewRequest("DELETE", handler.url+accountID+"?version=0", nil)
	if reqErr != nil {
		log.Fatal(reqErr)
		return false, reqErr
	}

	deleteAccountResponse, delErr := handler.client.Do(deleteAccountRequest)

	if delErr != nil {
		log.Fatal(delErr)
		return false, delErr
	}

	defer deleteAccountResponse.Body.Close()

	return deleteAccountResponse.StatusCode == 204, nil
}

func (handler form3AccountsHandler) Fetch(accountID string) (account.Account, error) {

	fetchedAccountResponse, fetchErr := handler.client.Get(handler.url + accountID)
	if fetchErr != nil {
		log.Fatal(fetchErr)
		return account.Account{}, fetchErr
	}

	defer fetchedAccountResponse.Body.Close()

	fetchedAccount := data{}

	if fetchedAccountResponse.StatusCode != 200 {
		return account.Account{}, errors.New("not existing account for " + accountID)
	}

	decodeErr := json.NewDecoder(fetchedAccountResponse.Body).Decode(&fetchedAccount)
	if decodeErr != nil {
		return account.Account{}, decodeErr
	}

	return fetchedAccount.Body, nil
}

func Handler(client http.Client, url string) AccountsHandler {
	return form3AccountsHandler{url, client}
}
