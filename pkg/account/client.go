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

func (a *AccountClient) CreateAccount(ctx context.Context, account AccountAttributes) (*Response, error) {
	body, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}
	return a.ExecuteRequest(ctx, http.MethodPost, a.getUrl(), body)
}

func (a *AccountClient) DeleteAccount(ctx context.Context, accountId, version string) (*Response, error) {
	return a.ExecuteRequest(ctx, http.MethodDelete, a.getUrl(), nil)
}

func (a *AccountClient) FetchAccount(ctx context.Context, accountId string) (*Response, error) {
	//add parameter to url to fetch account
	return a.ExecuteRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", a.getUrl(), accountId), nil)
}

func (a *AccountClient) ExecuteRequest(ctx context.Context, method, url string, body []byte) (*Response, error) {

	reader, length := getBodyReaderAndLength(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, reader)

	req.Header.Set("Host", "api.form3.tech")
	req.Header.Set("Date", time.Now().Format(time.RFC3339))
	req.Header.Set("Accept", "vnd.api+json")
	req.Header.Set("Accept-Encoding", "gzip")

	if length > 0 {
		req.Header.Set("Content-Type", "application/vnd.api+json")
		req.Header.Set("Content-Length", fmt.Sprint(length))
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
	return &result, nil
}

func getBodyReaderAndLength(body []byte) (io.Reader, int) {
	if len(body) > 0 {
		return bytes.NewReader(body), len(body)
	}
	return nil, 0
}

//TODO: Timeouts, Rate Limiting and Retry Strategy
//Should a request to the Form3 API respond with a status code indicating a temporary error (429, 500, 503 or 504, see above) or no response is received at all, wait and retry the request using an exponential back-off algorithm. See the code panel on the right for a simple example implementation in pseudo code.
