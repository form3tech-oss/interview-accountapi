package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/advena/interview-accountapi/cmd/app/account"
)

type AccountsHandler interface {
	Create(account account.Account) (account.Account, error)
	Delete(accountID string) (bool, error)
	Fetch(accountID string) (account.Account, error)
}

type Data struct {
	Body account.Account `json:"data"`
}

type Form3AccountsHandler struct {
	url    string
	client http.Client
}

func (handler Form3AccountsHandler) Create(a account.Account) (account.Account, error) {
	jsonData, jsonErr := json.Marshal(&Data{a})
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

	createdAccount := account.Account{}
	json.NewDecoder(createAccountResponse.Body).Decode(&createdAccount)

	return createdAccount, nil
}

func (handler Form3AccountsHandler) Delete(accountID string) (bool, error) {
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

func (handler Form3AccountsHandler) Fetch(accountID string) (account.Account, error) {

	fetchedAccountResponse, fetchErr := handler.client.Get(handler.url + accountID)
	if fetchErr != nil {
		log.Fatal(fetchErr)
		return account.Account{}, fetchErr
	}

	defer fetchedAccountResponse.Body.Close()

	fetchedAccount := new(Data)

	json.NewDecoder(fetchedAccountResponse.Body).Decode(&fetchedAccount)

	if fetchedAccountResponse.StatusCode != 200 {
		return account.Account{}, errors.New("not existing account for " + accountID)
	}

	return fetchedAccount.Body, nil
}

func Handler(url string) AccountsHandler {
	return Form3AccountsHandler{url, http.Client{Timeout: 10 * time.Second}}
}
