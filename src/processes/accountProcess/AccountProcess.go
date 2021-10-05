package accountProcess

import (
	"fmt"

	"../../../src/endpoint"
	"../../../src/validators/accountValidator"
	"../../messages"
)

func GetAccounts() {

	endpoint.GetAccounts()
}

func PostAccount(accountData map[string]interface{}) {
	validationResult := accountValidator.Validate(accountData)

	if !validationResult {
		fmt.Printf(messages.ACCOUNT_WITH_ERRORS + " Aborting POST request.\n")
	} else {
		endpoint.PostAccount(accountData)
	}
}

func DeleteAccount(id string) {
	endpoint.DeleteAccount(id)
}
