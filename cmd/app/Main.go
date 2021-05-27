package main

import (
	"fmt"

	"github.com/advena/interview-accountapi/cmd/app/account"
	"github.com/advena/interview-accountapi/cmd/app/handler"
	"github.com/google/uuid"
)

func main() {
	//url
	url := "http://localhost:8080/v1/organisation/accounts/"

	//create handler
	accountHandler := handler.Handler(url)

	//create Account data
	newAccount := account.Account{
		ID:             uuid.New().String(),
		OrganisationID: uuid.New().String(),
		Type:           "accounts",
		Version:        0,
		Attributes: account.Attributes{
			BankID:       "400400",
			BankIDCode:   "GBDSC",
			BaseCurrency: "GBP",
			BIC:          "NWBKGB22",
			Country:      "GB",
			Name:         []string{"some name"},
		},
	}

	//save account
	accountHandler.Create(newAccount)

	//get account
	accountHandler.Fetch(newAccount.ID)

	//delete account
	accountHandler.Delete(newAccount.ID)

	//validate created account is deleted
	exists, err := accountHandler.Fetch(newAccount.ID)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf(exists.ID)

}
