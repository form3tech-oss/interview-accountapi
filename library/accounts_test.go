package library_test

import (
	"fmt"
	"github.com/jsebasct/account-api-lib/library"
	"github.com/jsebasct/account-api-lib/models"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	bodyResponse := models.AccountBodyResponse{}
	err := library.GetAccounts(&bodyResponse)
	if err != nil {
		t.Error(err.Error())
	} else {
		fmt.Printf("%+v\n", bodyResponse)
	}
}
