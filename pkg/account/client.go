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
	Host    string
	Port    int
	Version string
}

// AccountClient is a client for the account service.
type AccountClient struct {
	Host       string
	Port       int
	Version    string
	HttpClient *http.Client
}

// NewAccountClient creates an AccountClient using a Config struct and returning a pointer to it.
func NewAccountClient(c *Config) *AccountClient {
	return &AccountClient{
		Host:       c.Host,
		Port:       c.Port,
		Version:    c.Version,
		HttpClient: &http.Client{},
	}
}

func (a *AccountClient) Greet() string {
	return "Hello! the URL is: " + a.getUrl()
}

func (a *AccountClient) getUrl() string {
	return fmt.Sprintf("http://%s:%v/%s/organisation/accounts", a.Host, a.Port, a.Version)
}

func (a *AccountClient) CreateAccount(ctx context.Context, accountData *AccountData) (*AccountData, error) {
	body, err := json.Marshal(Request{Data: accountData})
	if err != nil {
		return nil, err
	}
	data := &AccountData{}
	return data, a.ExecuteRequest(ctx, http.MethodPost, a.getUrl(), body, data)
}

func (a *AccountClient) DeleteAccount(ctx context.Context, accountId string, version int64) error {
	return a.ExecuteRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s?version=%d", a.getUrl(), accountId, version), nil, nil)
}

func (a *AccountClient) FetchAccount(ctx context.Context, accountId string) (*AccountData, error) {
	data := &AccountData{}
	return data, a.ExecuteRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", a.getUrl(), accountId), nil, data)
}

func (a *AccountClient) ExecuteRequest(ctx context.Context, method, url string, body []byte, i interface{}) error {

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err
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

	res, err := a.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	result := &Response{
		Data: i,
	}

	if len(b) > 0 {
		err = json.Unmarshal(b, result)
		if err != nil {
			return err
		}
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return newErrorResponse(res.StatusCode, result.ErrorMessage)
	}
	return nil
}

//TODO: Timeouts, Rate Limiting and Retry Strategy
//Should a request to the Form3 API respond with a status code indicating a temporary error (429, 500, 503 or 504, see above) or no response is received at all, wait and retry the request using an exponential back-off algorithm. See the code panel on the right for a simple example implementation in pseudo code.

type ErrorResponse struct {
	Code    int
	Message string
}

func newErrorResponse(code int, message *string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: *message,
	}
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("error: (%d) message: %s", er.Code, er.Message)
}
