package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// configuration to create an AccountClient
type Config struct {
	BaseUrl    string
	Version    string
	MaxRetries int
	Wait       int
}

// AccountClient is a client for the organization/account API.
// BaseUrl the base URL for the API calls.
// Version the version of the API to use.
// HttpClient the http client to use for the API calls.
// LimitRateAndRetry the implementation for the retry / exponential backoff strategy.
type AccountClient struct {
	BaseUrl           string
	Version           string
	HttpClient        *http.Client
	LimitRateAndRetry *LimitRateAndRetry
}

// NewAccountClient Constructor. Creates an AccountClient with the Config struct and returning a pointer to it.
// baseURL defaults to http://api.form3.tech
// version defaults to v1
// httpClient defaults to http.Client{}
// LimitRateAndRetry defaults to LimitRateAndRetry{MaxRetries: 3, Wait: 500}
func NewAccountClient(c *Config) *AccountClient {
	if c == nil {
		return &AccountClient{HttpClient: &http.Client{}, LimitRateAndRetry: &LimitRateAndRetry{}}
	}

	httpClient := &http.Client{}
	return &AccountClient{
		BaseUrl:    c.BaseUrl,
		Version:    c.Version,
		HttpClient: httpClient,
		LimitRateAndRetry: &LimitRateAndRetry{
			MaxRetries: &c.MaxRetries,
			Wait:       &c.Wait,
		},
	}
}

// Returns the base URL for the API calls.
// Defaults to http://api.form3.tech/v1/organisation/accounts
func (a *AccountClient) GetUrl() string {
	url := "http://api.form3.tech"
	version := "v1"

	if len(a.BaseUrl) > 0 {
		url = a.BaseUrl
	}

	if len(a.Version) > 0 {
		version = a.Version
	}

	return fmt.Sprintf("%s/%s/organisation/accounts", url, version)
}

// Creates a new account using the AccountData struct and returns a pointer to the created account.
// This same context will be used for any retry, allowing cancel, timeout, and error handling.
// Return error if the account already exists or if any validation fails for the input data.
func (a *AccountClient) CreateAccount(ctx context.Context, accountData *AccountData) (*AccountData, error) {
	body, err := json.Marshal(Request{Data: accountData})
	if err != nil {
		return nil, err
	}
	data := &AccountData{}
	if err := a.ExecuteRequest(ctx, http.MethodPost, a.GetUrl(), body, data); err != nil {
		return nil, err
	}
	return data, nil
}

// Delete an existing account with the account ID and version.
// This same context will be used for any retry, allowing cancel, timeout, and error handling.
// Returns error if the account is not found.
func (a *AccountClient) DeleteAccount(ctx context.Context, accountId string, version int64) error {
	return a.ExecuteRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s?version=%d", a.GetUrl(), accountId, version), nil, nil)
}

// Fetch an existing account with the account ID.
// This same context will be used for any retry, allowing cancel, timeout, and error handling.
// Returns error if the account is not found.
func (a *AccountClient) FetchAccount(ctx context.Context, accountId string) (*AccountData, error) {
	data := &AccountData{}
	if err := a.ExecuteRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", a.GetUrl(), accountId), nil, data); err != nil {
		return nil, err
	}
	return data, nil
}

// Function to execute the request and handle the response.
// Returns error if the request fails or if the response status code is not 2xx / 3xx.
// The response body is decoded into the given interface.
func (a *AccountClient) ExecuteRequest(ctx context.Context, method, url string, body []byte, i interface{}) error {

	// build request with the given context
	req, err := buildRequest(ctx, method, url, body)
	if err != nil {
		return err
	}

	// the request execution is delegated to the LimitRateAndRetry implementation.
	res, err := a.LimitRateAndRetry.ExponentialBackOff(a.HttpClient, req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	result := &Response{
		Data: i,
	}

	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil && err != io.EOF {
		return err
	}

	// read the body and return an error if the status code is not 2xx /3xx
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return newErrorResponse(res.StatusCode, result.ErrorMessage)
	}
	return nil
}

// Function to build the request with the given context and the necessary headers.
func buildRequest(ctx context.Context, method, url string, body []byte) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Host = "api.form3.tech"
	req.Header.Set("Host", "api.form3.tech")
	req.Header.Set("Date", time.Now().Format(time.RFC3339))
	req.Header.Set("Accept", "vnd.api+json")
	req.Header.Set("Accept-Encoding", "gzip")

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/vnd.api+json")
		req.Header.Set("Content-Length", fmt.Sprint(len(body)))
	}
	return req, nil
}

// Struct to hold error responses from the API.
type ErrorResponse struct {
	Code    int
	Message string
}

// ErrorResponse constructor.
func newErrorResponse(code int, message *string) *ErrorResponse {
	var value string
	if message != nil {
		value = *message
	}
	return &ErrorResponse{
		Code:    code,
		Message: value,
	}
}

// Error function to implement the error interface.
func (er ErrorResponse) Error() string {
	return fmt.Sprintf("error: (%d) message: %s", er.Code, er.Message)
}
