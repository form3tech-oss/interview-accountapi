package handler

import (
	"net/http"
	"testing"

	"github.com/advena/interview-accountapi/cmd/app/account"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const url = "http://accounts/api/"
const validAccountID = "account-id"

func Account() account.Account {
	return account.Account{
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
}

func TestFetchAccount(t *testing.T) {
	//given
	client := http.Client{}

	httpmock.ActivateNonDefault(&client)
	httpmock.RegisterResponder("GET", url+validAccountID, httpmock.NewStringResponder(200, validRespone()))

	handler := Handler(client, url)

	//when
	account, err := handler.Fetch(validAccountID)

	//then
	assert.NoError(t, err)
	assert.Equal(t, account.ID, validAccountID)
}

func TestCreateAccount(t *testing.T) {
	//given
	client := http.Client{}

	httpmock.ActivateNonDefault(&client)
	httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(200, validRespone()))

	handler := Handler(client, url)

	//when
	createdAccount, err := handler.Create(Account())

	//then
	assert.NoError(t, err)
	assert.IsType(t, createdAccount, account.Account{})
}

func TestDeleteAccount(t *testing.T) {
	//given
	client := http.Client{}

	httpmock.ActivateNonDefault(&client)

	handler := Handler(client, url)

	httpmock.RegisterResponder("DELETE", url+validAccountID+"?version=0", httpmock.NewStringResponder(204, ""))

	//when
	account, err := handler.Delete(validAccountID)

	//then
	assert.NoError(t, err)
	assert.True(t, account)
}

func validRespone() string {
	return `{
		"data": {
			"type": "accounts",
			"id": "account-id",
			"version": 0,
			"organisation_id": "organisation-id",
			"attributes": {
				"country": "GB",
				"base_currency": "GBP",
				"account_number": "41426819",
				"bank_id": "400300",
				"bank_id_code": "GBDSC",
				"bic": "NWBKGB22"
			}
		}
	}`
}
