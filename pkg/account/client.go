package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrDuplicatedAccount = errors.New("duplicated account")

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
	return a.ExecuteRequest(ctx, http.MethodPost, a.getUrl(), body)
}

func (a *AccountClient) DeleteAccount(ctx context.Context, accountId string, version int64) (*AccountData, error) {
	return a.ExecuteRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s?version=%d", a.getUrl(), accountId, version), nil)
}

func (a *AccountClient) FetchAccount(ctx context.Context, accountId string) (*AccountData, error) {
	//add parameter to url to fetch account
	return a.ExecuteRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", a.getUrl(), accountId), nil)
}

func (a *AccountClient) ExecuteRequest(ctx context.Context, method, url string, body []byte) (*AccountData, error) {

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)

	req.Header.Set("Host", "api.form3.tech")
	req.Header.Set("Date", time.Now().Format(time.RFC3339))
	req.Header.Set("Accept", "vnd.api+json")
	req.Header.Set("Accept-Encoding", "gzip")

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/vnd.api+json")
		req.Header.Set("Content-Length", fmt.Sprint(len(body)))
	}

	if err != nil {
		return nil, err
	}

	res, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := Response{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrorMessage != nil {
		msg := fmt.Sprintf("account error: (%d) %s", res.StatusCode, *result.ErrorMessage)

		if res.StatusCode == 409 {
			return nil, fmt.Errorf("%s: %w", msg, ErrDuplicatedAccount)
		} else {
			return nil, fmt.Errorf(msg)
		}
	}
	return result.Data, nil
}

//TODO: Timeouts, Rate Limiting and Retry Strategy
//Should a request to the Form3 API respond with a status code indicating a temporary error (429, 500, 503 or 504, see above) or no response is received at all, wait and retry the request using an exponential back-off algorithm. See the code panel on the right for a simple example implementation in pseudo code.
