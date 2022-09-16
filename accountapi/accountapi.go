package accountapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/giannimassi/interview-accountapi/model"
)

type AccountAPIClientOptions func(*AccountAPIClient) *AccountAPIClient

// WithHTTPClient returns an AccountAPIClientOptions setting up the API client with the provided http client.
func WithHTTPClient(client *http.Client) AccountAPIClientOptions {
	return func(c *AccountAPIClient) *AccountAPIClient {
		c.client = client
		return c
	}
}

// WithHost returns an AccountAPIClientOptions setting up the API client with the provided host.
func WithHost(host string) AccountAPIClientOptions {
	return func(c *AccountAPIClient) *AccountAPIClient {
		c.host = host
		return c
	}
}

// AccountAPIClient allows to perform requests to the account api client
type AccountAPIClient struct {
	client *http.Client // Unexported in order to avoid direct access to the client
	host   string
}

// NewAccountAPIClient returns a new AccountAPIClient with all provided options applied and host set
func NewAccountAPIClient(options ...AccountAPIClientOptions) *AccountAPIClient {
	a := new(AccountAPIClient)
	for _, o := range options {
		o(a)
	}
	return a
}

const (
	createAccountURLV1 = "/v1/organisation/accounts"

	contentTypeJSON = "application/vnd.api+json"
)

type accountDataBody struct {
	Data model.AccountData `json:"data"`
}

// Create performs the CREATE operation on the AccountData entity.
func (a *AccountAPIClient) Create(ctx context.Context, account model.AccountData) (*model.AccountData, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(accountDataBody{Data: account}); err != nil {
		return nil, fmt.Errorf("while encoding account data: %w", err)
	}

	// TODO: use url.URL to build the URL
	req, err := http.NewRequestWithContext(ctx, "POST", a.host+createAccountURLV1, body)
	if err != nil {
		return nil, fmt.Errorf("while formatting request: %w", err)
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	response, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making request: %w", err)
	}

	if response.StatusCode != http.StatusCreated {
		// Handle unexpected errors
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Handle Successful creation
	var responseAccount accountDataBody
	if err := json.NewDecoder(response.Body).Decode(&responseAccount); err != nil {
		return nil, fmt.Errorf("while decoding response: %w", err)
	}

	return &responseAccount.Data, nil
}

// Fetch performs the FETCH operation on the AccountData entity.
func (a *AccountAPIClient) Fetch(ctx context.Context, account model.AccountData) (*model.AccountData, error) {
	// TODO: use url.URL to build the URL
	req, err := http.NewRequestWithContext(ctx, "GET", a.host+createAccountURLV1+"/"+account.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("while formatting request: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making request: %w", err)
	}

	// TODO: use sentinel errors for other meaningful cases (e.g. 404)
	if response.StatusCode != http.StatusOK {
		// Handle unexpected errors
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Handle account found
	var responseAccount accountDataBody
	if err := json.NewDecoder(response.Body).Decode(&responseAccount); err != nil {
		return nil, fmt.Errorf("while decoding response: %w", err)
	}

	return &responseAccount.Data, nil
}

// Delete performs the DELETE operation on the AccountData entity.
func (a *AccountAPIClient) Delete(ctx context.Context, account model.AccountData) error {
	// TODO: use url.URL to build the URL
	req, err := http.NewRequestWithContext(ctx, "DELETE", a.host+createAccountURLV1+"/"+account.ID+"?version="+strconv.FormatInt(*account.Version, 10), nil)
	if err != nil {
		return fmt.Errorf("while formatting request: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("while making request: %w", err)
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		// Handle successful deletion
		return nil

	// TODO: use sentinel errors for other meaningful cases (e.g. 404)
	case http.StatusNotFound:
		// Handle resource not found
		return errors.New("specified resource does not exist")
	case http.StatusConflict:
		// Handle incorrect version
		return errors.New("specified version incorrect")
	default:
		// Handle unexpected errors
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}
