package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/advena/interview-accountapi/cmd/app/account"
	"github.com/advena/interview-accountapi/cmd/app/account/validation"
	"github.com/advena/interview-accountapi/cmd/app/downstream"
)

type form3AccountsHandler struct {
	repository account.AccountsRepository
	validator  validation.Validator
}

func (a *form3AccountsHandler) Create(newAccount account.Account) (account.Account, error) {
	validationResult := a.validator.Validate(newAccount)
	if validationResult.IsValid() {
		return a.repository.Create(newAccount)
	} else {
		return account.Account{}, errors.New(validationResult.Message())
	}
}

func (a *form3AccountsHandler) Delete(accountID string) (bool, error) {
	return a.repository.Delete(accountID)
}

func (a *form3AccountsHandler) Fetch(accountID string) (account.Account, error) {
	return a.repository.Fetch(accountID)
}

func NewForm3AccountHandler(url string) account.AccountsRepository {
	client := http.Client{Timeout: 10 * time.Second}
	validators := []validation.Validator{validation.BankIDValidator(), validation.BICValidator()}
	accountRepository := downstream.Form3RestRespository(client, url)
	accountValidator := validation.AccountValidator{Validators: validators}

	return &form3AccountsHandler{accountRepository, accountValidator}
}
