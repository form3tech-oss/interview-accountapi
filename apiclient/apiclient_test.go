package accountapiclient

import (
	"accountapiclient"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAccount_ReturnsOneError_WhenInputIsEmpty(t *testing.T) {
	client := CreateClient("http://localhost:8080")
	_, errors := client.CreateAccount(accountapiclient.AccountData{})
	if len(errors) != 1 {
		t.Fatalf("expected one error, received '%s'", fmt.Sprint(len(errors)))
	}
}

func TestCreateAccount_ReturnsNoErrors_WhenCrationIsSuccessful(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{
			"data": {
				"attributes": {
					"alternative_names": null,
					"country": "CA",
					"name": [
						"M"
					]
				},
				"created_on": "2022-02-14T14:11:46.906Z",
				"id": "49dac5ee-6ffb-4bb3-a24d-9c36d4f4ca36",
				"modified_on": "2022-02-14T14:11:46.906Z",
				"organisation_id": "78398917-e6bd-4671-bc99-666c5015af99",
				"type": "accounts",
				"version": 0
			},
			"links": {
				"self": "/v1/organisation/accounts/49dac5ee-6ffb-4bb3-a24d-9c36d4f4ca36"
			}
		}`))
	}))
	defer server.Close()

	client := CreateClient(server.URL)
	country := "CA"
	_, errors := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if len(errors) > 0 {
		t.Fatalf("expected no errors, received '%s'", fmt.Sprint(len(errors)))
	}
}

func TestCreateAccount_ReturnsCreatedAccount_WhenCreationIsSuccessful(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{
			"data": {
				"attributes": {
					"alternative_names": null,
					"country": "CA",
					"name": [
						"M"
					]
				},
				"created_on": "2022-02-14T14:11:46.906Z",
				"id": "49dac5ee-6ffb-4bb3-a24d-9c36d4f4ca36",
				"modified_on": "2022-02-14T14:11:46.906Z",
				"organisation_id": "78398917-e6bd-4671-bc99-666c5015af99",
				"type": "accounts",
				"version": 0
			},
			"links": {
				"self": "/v1/organisation/accounts/49dac5ee-6ffb-4bb3-a24d-9c36d4f4ca36"
			}
		}`))
	}))
	defer server.Close()

	client := CreateClient(server.URL)
	country := "CA"
	account, _ := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if (account == accountapiclient.AccountData{}) {
		t.Fatalf("expected account with data, received '%s'", fmt.Sprint(account))
	}
}

func TestCreateAccount_ReturnsValidationErrors_WhenCreationFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`{
			"error_message": "validation failure list:\nvalidation failure list:\nvalidation failure list:\ncountry in body should match '^[A-Z]{2}$'"
		}`))
	}))
	defer server.Close()

	client := CreateClient(server.URL)
	country := "CA"
	_, errors := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if len(errors) < 2 {
		t.Fatalf("expected two errors, received '%s'", fmt.Sprint(len(errors)))
	}
}

func TestCreateAccount_ReturnsOneError_WhenCreationCrashes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(500)
		rw.Write([]byte(""))
	}))
	defer server.Close()

	client := CreateClient(server.URL)
	country := "CA"
	_, errors := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if len(errors) != 1 {
		t.Fatalf("expected one error, received '%s'", fmt.Sprint(len(errors)))
	}
}
