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

// configurations for the account client
type Config struct {
	BaseUrl string
	Version string
}

// AccountClient is a client for the account service.
type AccountClient struct {
	BaseUrl    string
	Version    string
	HttpClient *http.Client
}

// NewAccountClient creates an AccountClient using a Config struct and returning a pointer to it.
func NewAccountClient(c *Config) *AccountClient {
	if c == nil {
		return &AccountClient{HttpClient: &http.Client{}}
	}

	return &AccountClient{
		BaseUrl:    c.BaseUrl,
		Version:    c.Version,
		HttpClient: &http.Client{},
	}
}

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

func (a *AccountClient) DeleteAccount(ctx context.Context, accountId string, version int64) error {
	return a.ExecuteRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s?version=%d", a.GetUrl(), accountId, version), nil, nil)
}

func (a *AccountClient) FetchAccount(ctx context.Context, accountId string) (*AccountData, error) {
	data := &AccountData{}
	if err := a.ExecuteRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", a.GetUrl(), accountId), nil, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (a *AccountClient) ExecuteRequest(ctx context.Context, method, url string, body []byte, i interface{}) error {

	req, err := buildRequest(ctx, method, url, body)
	if err != nil {
		return err
	}

	res, err := a.HttpClient.Do(req)
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

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return newErrorResponse(res.StatusCode, result.ErrorMessage)
	}
	return nil
}

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

//TODO: Timeouts, Rate Limiting and Retry Strategy
//Should a request to the Form3 API respond with a status code indicating a temporary error (429, 500, 503 or 504, see above) or no response is received at all, wait and retry the request using an exponential back-off algorithm. See the code panel on the right for a simple example implementation in pseudo code.

type ErrorResponse struct {
	Code    int
	Message string
}

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

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("error: (%d) message: %s", er.Code, er.Message)
}
