package library_test

import (
	"fmt"
	"github.com/jsebasct/account-api-lib/library"
	"github.com/jsebasct/account-api-lib/models"
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
