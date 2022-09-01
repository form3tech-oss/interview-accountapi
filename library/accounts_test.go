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

// after one insertion
func TestGetAccounts(t *testing.T) {
	bodyResponse := models.AccountBodyResponse{}
	err := library.GetAccounts(&bodyResponse)
	if err != nil {
		t.Error(err.Error())
	} else {
		fmt.Printf("%+v\n", bodyResponse)
		if len(bodyResponse.Data) != 1 {
			t.Errorf("expected %d, but got: %d", 1, len(bodyResponse.Data))
		}
		expectedType := "accounts"
		first := bodyResponse.Data[0]
		if first.Type != expectedType {
			t.Errorf("expected %s, but got: %s", expectedType, first.Type)
		}

		expectedCurrency := "GBP"
		if first.Attributes.BaseCurrency != expectedCurrency {
			t.Errorf("expected %s, but got: %s", expectedType, first.Type)
		}
	}
}

// with no insertion
func TestCreateAccountSuccess(t *testing.T) {

	// get account request
	requestFromFile, err := os.ReadFile("samples/account_request.json")
	if err != nil {
		t.Error(err.Error())
	}

	var accountRequest models.AccountBodyRequest
	marshallError := json.Unmarshal(requestFromFile, &accountRequest)
	if marshallError != nil {
		t.Error(marshallError.Error())
	}
	//fmt.Printf("%+v %+v \n", accountRequest, accountRequest.Data.Attributes)

	// make actual request
	createResponse, createError := library.CreateAccount(&accountRequest)
	if createError != nil {
		t.Error(errors.New(createError.Message))
	}
	fmt.Printf("%+v\n", createResponse)

	// test
	accountsResponse := models.AccountBodyResponse{}
	err = library.GetAccounts(&accountsResponse)
	if err != nil {
		t.Error(err.Error())
	}
	if len(accountsResponse.Data) < 1 {
		t.Errorf("expected greater or equal to 1 received: %d", len(accountsResponse.Data))
	}

	// filter
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
