package accountapiclient

import (
	"accountapiclient"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAccount_ReturnsCreatedAccount(t *testing.T) {
	client := CreateClient("http://accountapi:8080")

	country := "CA"
	version := int64(1)
	account, _ := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Version:        &version,
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if (account == accountapiclient.AccountData{}) {
		t.Fatalf("expected valid account type, received '%s'", fmt.Sprint(account))
	}
}
