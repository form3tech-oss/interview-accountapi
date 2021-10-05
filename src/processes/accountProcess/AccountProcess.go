package accountProcess

import (
	"fmt"

	"../../../src/endpoint"
	"../../../src/validators/accountValidator"
	"../../messages"
)

func GetAccounts() {
	fmt.Println("\nGetting Accounts:")
	endpoint.GetAccounts()
}

func PostAccount(accountData map[string]interface{}) {
	validationResult := accountValidator.Validate(accountData)

	if !validationResult {
		fmt.Printf(messages.ACCOUNT_WITH_ERRORS + " Aborting POST request.\n")
	} else {
		fmt.Println("\nPosting an Account:")
		endpoint.PostAccount(accountData)
	}
}

func DeleteAccount(id string) {
	fmt.Println("\nRemoving account number: " + id)
	endpoint.DeleteAccount(id)
}
