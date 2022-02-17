package accountapiclient

import (
	"accountapiclient"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAccount_ReturnsCreatedAccount(t *testing.T) {
	client := CreateClient("http://localhost:8080")

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

	t.Run("TestFetchAccount", func(t *testing.T) {
		foundAccount, e := client.FetchAccount(account.ID)
		if (len(e) > 0 || foundAccount == accountapiclient.AccountData{}) {
			t.Fatalf("expected account with no errors, recieved '%s'", fmt.Sprint(len(e)))
		}
	})
}

func TestCreateAccount_ReturnsError_WhenApiIsDown(t *testing.T) {
	client := CreateClient("http://myhost:8080")
	country := "CA"
	version := int64(1)
	_, e := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Version:        &version,
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if len(e) <= 0 {
		t.Fatalf("expected remote host error, received '%s'", fmt.Sprint(e))
	}
}
