package app

import (
	"fmt"
	"testing"

	"github.com/advena/interview-accountapi/cmd/app/account"
	"github.com/advena/interview-accountapi/cmd/app/account/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateFetchDeleteAccount(t *testing.T) {
	url := "http://accountapi:8080/v1/organisation/accounts/"

	//create handler
	accountHandler := handler.NewForm3AccountHandler(url)

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

	//invalid account is not created
	invalidAccount := account.Account{}
	invalidAccount.Attributes.Country = "GB"

	_, err = accountHandler.Create(invalidAccount)

	assert.NotNil(t, err)

}

func validate(error error) {
	if error != nil {
		fmt.Printf("Error: %s \n", error.Error())
	}
}
