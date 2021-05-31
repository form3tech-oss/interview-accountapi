package app

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/advena/interview-accountapi/cmd/app/account"
	"github.com/advena/interview-accountapi/cmd/app/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateFetchDeleteAccount(t *testing.T) {
	url := "http://accountapi:8080/v1/organisation/accounts/"
	client := http.Client{Timeout: 10 * time.Second}

	//create handler

	accountHandler := handler.Handler(client, url)

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
	createdAccount, err := accountHandler.Create(newAccount)
	validate(err)

	//get account
	fetchedAccount, err := accountHandler.Fetch(newAccount.ID)
	validate(err)

	assert.True(t, createdAccount.ID == fetchedAccount.ID)

	//delete account
	wasDeleted, err := accountHandler.Delete(newAccount.ID)
	validate(err)

	//validate created account is deleted
	exists, err := accountHandler.Fetch(newAccount.ID)
	assert.NotNil(t, err)
	assert.Equal(t, exists.ID, "")

	assert.True(t, wasDeleted)

}

func validate(error error) {
	if error != nil {
		fmt.Printf("Error: %s \n", error.Error())
	}
}
