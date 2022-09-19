package tests

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/giannimassi/interview-accountapi/accountapi"
	"github.com/giannimassi/interview-accountapi/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_integration(t *testing.T) {

	// The following code is an integration test for the accountapi package.
	// It is not a unit test, but it is a test that can be run against the actual api backend to verify that it works as expected.
	// The test assumes to be run as part of the docker compose provided in the root of the project (see README.md for more info).
	var testRequestTimeout = time.Second * 10

	host := os.Getenv("ACCOUNTAPI_HOST")
	require.NotEmpty(t, host, "ACCOUNTAPI_HOST is not set")

	t.Run("When creating a valid account, the account is created and the account can be fetched and deleted by ID", func(t *testing.T) {
		// Setup library
		c := accountapi.NewAccountAPIClient(accountapi.WithHTTPClient(http.DefaultClient), accountapi.WithHost(host))
		ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
		defer cancel()

		id := uuid.New().String()
		a := accountData(id, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")

		// Create account
		created, err := c.Create(ctx, a)
		require.NoError(t, err)

		// Fetch account
		fetched, err := c.Fetch(ctx, *created)
		require.NoError(t, err)
		require.Equal(t, created, fetched, "Accounts created and fetched do not match: %v != %v", created, fetched)

		// Delete account
		err = c.Delete(ctx, *fetched)
		require.NoError(t, err)

		// Fetch account again
		fetchedAgain, err := c.Fetch(ctx, *created)
		require.Equal(t, errors.New("unexpected status code: 404"), err)
		require.Nil(t, fetchedAgain, "Account should have been deleted, but it was not")
	})

	t.Run("When creating account with existing id, an error is returned", func(t *testing.T) {
		// Setup library
		c := accountapi.NewAccountAPIClient(accountapi.WithHTTPClient(http.DefaultClient), accountapi.WithHost(host))
		ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
		defer cancel()

		id := uuid.New().String()
		a := accountData(id, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")

		// Create account
		_, err := c.Create(ctx, a)
		require.NoError(t, err)

		b := accountData(id, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")

		// Fail creating account with same id
		_, err = c.Create(ctx, b)
		require.EqualError(t, err, "unexpected status code: 409")
	})

	t.Run("When fetching account by unexisting ID, an account is not returned", func(t *testing.T) {
		c := accountapi.NewAccountAPIClient(accountapi.WithHTTPClient(http.DefaultClient), accountapi.WithHost(host))
		ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
		defer cancel()

		id := uuid.New().String()
		a := accountData(id, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")

		// Fetch account by unexistent ID
		_, err := c.Fetch(ctx, a)
		require.EqualError(t, err, "unexpected status code: 404")
	})

	t.Run("When deleting account by unexisting ID, an account is not deleted", func(t *testing.T) {
		c := accountapi.NewAccountAPIClient(accountapi.WithHTTPClient(http.DefaultClient), accountapi.WithHost(host))
		ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
		defer cancel()

		id := uuid.New().String()
		a := accountData(id, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")
		a.Version = new(int64)

		// Delete account by unexistent ID
		err := c.Delete(ctx, a)
		require.EqualError(t, err, "specified resource does not exist")
	})

	t.Run("When deleting account by incorrect version, an account is not deleted", func(t *testing.T) {
		// Setup library
		c := accountapi.NewAccountAPIClient(accountapi.WithHTTPClient(http.DefaultClient), accountapi.WithHost(host))
		ctx, cancel := context.WithTimeout(context.Background(), testRequestTimeout)
		defer cancel()

		id := uuid.New().String()
		a := accountData(id, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")

		// Create account
		created, err := c.Create(ctx, a)
		require.NoError(t, err)

		// Overwrite version with wrong value
		*created.Version++

		// Fetch account
		err = c.Delete(ctx, *created)
		require.EqualError(t, err, "specified version incorrect")
	})
}

func accountData(id, orgID string) model.AccountData {
	return model.AccountData{
		ID:             id,
		OrganisationID: orgID,
		Type:           "accounts",
		Attributes: &model.AccountAttributes{
			Country:               stringP("GB"),
			BaseCurrency:          "GBP",
			BankID:                "400302",
			BankIDCode:            "GBDSC",
			AccountNumber:         "10000004",
			AccountClassification: stringP("Personal"),
			Name:                  []string{"Elizabeth", "I", "Tudor"},
		},
	}
}

func stringP(s string) *string {
	return &s
}
