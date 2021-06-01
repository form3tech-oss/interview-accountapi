package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatesAccountWithGivenCountryAndBIC(t *testing.T) {
	params := []struct {
		country        string
		bic            string
		expectedResult bool
	}{
		{"GB", "NWBKGB22", true},
		{"AU", "NWBKGB22", true},
		{"BE", "NWBKGB22", true},
		{"NL", "NWBKGB22", true},
		{"GB", "", false},
		{"AU", "", false},
		{"BE", "", false},
		{"NL", "", true},
	}

	for _, tt := range params {
		t.Run(tt.country, func(t *testing.T) {
			account := Account{}
			account.Attributes.BIC = tt.bic
			account.Attributes.Country = tt.country
			result := BICValidator().Validate(account)
			assert.Equal(t, result.IsValid(), tt.expectedResult)
		})
	}

}

func TestValidatesAccountWithGivenCountryAndBankID(t *testing.T) {
	params := []struct {
		country        string
		bankID         string
		expectedResult bool
	}{
		{"GB", "222222", true},
		{"GB", "2222222", false},
		{"GB", "22222", false},
		{"GB", "", false},
		{"AU", "222222", true},
		{"AU", "22222", false},
		{"AU", "2222222", false},
		{"AU", "", true},
		{"BE", "333", true},
		{"BE", "33", false},
		{"BE", "3333", false},
		{"BE", "", false},
		{"NL", "non-empty", false},
		{"NL", "", true},
	}

	for _, tt := range params {
		t.Run(tt.country, func(t *testing.T) {
			account := Account{}
			account.Attributes.BankID = tt.bankID
			account.Attributes.Country = tt.country
			result := BankIDValidator().Validate(account)
			assert.Equal(t, result.IsValid(), tt.expectedResult)
		})
	}

}

func TestValidatesAccountWithAllRequiredFiledsValid(t *testing.T) {
	account := Account{}
	account.Attributes.BankID = "222222"
	account.Attributes.Country = "GB"
	account.Attributes.BIC = "NWBKGB22"

	result := AccountValidator{[]Validator{BankIDValidator(), BICValidator()}}.Validate(account)

	assert.True(t, result.IsValid())
}

func TestValidatesAccountWithAllRequiredFieldsInvalid(t *testing.T) {
	account := Account{}
	account.Attributes.BankID = "22222"
	account.Attributes.Country = "GB"
	account.Attributes.BIC = ""

	result := AccountValidator{[]Validator{BankIDValidator(), BICValidator()}}.Validate(account)

	assert.False(t, result.IsValid())
}
