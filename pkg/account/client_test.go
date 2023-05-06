package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var account = createAccount()
var message = "some error!"

func TestExecuteRequest(t *testing.T) {

	type testCase struct {
		name         string
		context      func() context.Context
		headers      map[string]string
		code         int
		requestBody  interface{}
		responseBody interface{}
		data         interface{}
		err          error
	}

	testCases := []testCase{
		{
			name:         "request ok",
			context:      func() context.Context { return context.Background() },
			code:         200,
			responseBody: nil,
			data:         nil,
			err:          nil,
		},
		{
			name:         "request ok nil context",
			context:      func() context.Context { return nil },
			code:         200,
			responseBody: nil,
			data:         nil,
			err:          errors.New("net/http: nil Context"),
		},
		{
			name:    "request ok with content",
			context: func() context.Context { return context.Background() },
			code:    200,
			requestBody: Request{
				Data: account,
			},
			responseBody: Response{
				Data: account,
			},
			data: account,
			err:  nil,
		},
		{
			name: "request timeout",
			context: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				return ctx
			},
			headers:      nil,
			code:         200,
			requestBody:  nil,
			responseBody: nil,
			data:         nil,
			err:          fmt.Errorf("context canceled"),
		},
		{
			name:    "request failed",
			context: func() context.Context { return context.Background() },
			code:    500,
			responseBody: Response{
				ErrorMessage: &message,
			},
			data: nil,
			err: &ErrorResponse{
				Code:    500,
				Message: message,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			requestBytes, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			responseBytes, err := json.Marshal(tc.responseBody)
			if err != nil {
				t.Fatal(err)
			}

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

				assert.Equal(t, "vnd.api+json", req.Header.Get("Accept"))
				assert.Equal(t, "gzip", req.Header.Get("Accept-Encoding"))
				assert.NotEmpty(t, req.Header.Get("Date"))
				assert.Equal(t, "api.form3.tech", req.Host)

				if tc.requestBody != nil {
					assert.Equal(t, "application/vnd.api+json", req.Header.Get("Content-Type"))
					assert.Equal(t, fmt.Sprint(len(requestBytes)), req.Header.Get("Content-Length"))

					body, err := io.ReadAll(req.Body)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, requestBytes, body)
				}
				rw.WriteHeader(tc.code)
				rw.Write(responseBytes)
			}))
			defer server.Close()

			ac := &AccountClient{
				BaseUrl:    server.URL,
				Version:    "v1",
				HttpClient: server.Client(),
			}

			var data *AccountData
			if tc.data != nil {
				data = &AccountData{}
			}

			err = ac.ExecuteRequest(tc.context(), http.MethodGet, server.URL, requestBytes, data)

			if tc.err != nil {
				fmt.Println(err)
				if uerr := errors.Unwrap(err); uerr != nil {
					assert.Equal(t, tc.err, uerr)
				} else {
					assert.Equal(t, tc.err.Error(), err.Error())
				}
			} else {
				assert.Nil(t, err)
			}

			if tc.data != nil {
				assert.Equal(t, tc.data, data)
			} else {
				assert.Nil(t, data)
			}
		})
	}
}

func TestExecuteRequestFailToMarshal(t *testing.T) {

	type testCase struct {
		name          string
		code          int
		responseBytes []byte
		err           error
	}

	testCases := []testCase{
		{
			name:          "request ok",
			code:          200,
			responseBytes: []byte(`{"key":value}`),
			err:           fmt.Errorf("invalid character 'v' looking for beginning of value"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

				rw.WriteHeader(tc.code)
				rw.Write(tc.responseBytes)
			}))
			defer server.Close()

			ac := &AccountClient{
				BaseUrl:    server.URL,
				Version:    "v1",
				HttpClient: server.Client(),
			}

			data := &AccountData{}

			err := ac.ExecuteRequest(context.Background(), http.MethodGet, server.URL, nil, data)

			if tc.err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFetchAccount(t *testing.T) {

	type testCase struct {
		name         string
		code         int
		responseBody interface{}
		expectedUrl  string
		data         interface{}
		err          error
	}

	testCases := []testCase{
		{
			name: "request ok",
			code: 200,
			responseBody: Response{
				Data: account,
			},
			expectedUrl: "/v1/organisation/accounts/abc123",
			data:        account,
			err:         nil,
		},
		{
			name: "account not found",
			code: 404,
			responseBody: Response{
				ErrorMessage: &message,
			},
			expectedUrl: "/v1/organisation/accounts/abc123",
			data:        nil,
			err:         &ErrorResponse{Code: 404, Message: "some error!"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			responseBytes, err := json.Marshal(tc.responseBody)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tc.expectedUrl, req.RequestURI)
				rw.WriteHeader(tc.code)
				rw.Write(responseBytes)
			}))
			defer server.Close()

			ac := &AccountClient{
				BaseUrl:    server.URL,
				Version:    "v1",
				HttpClient: server.Client(),
			}

			data, err := ac.FetchAccount(context.Background(), "abc123")

			if tc.err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			if tc.data != nil {
				assert.Equal(t, tc.data, data)
			} else {
				assert.Nil(t, data)
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {

	type testCase struct {
		name         string
		code         int
		expectedUrl  string
		responseBody interface{}
		err          error
	}

	testCases := []testCase{
		{
			name:         "request ok",
			code:         200,
			expectedUrl:  "/v1/organisation/accounts/abc123?version=0",
			responseBody: nil,
			err:          nil,
		},
		{
			name:        "account not found",
			code:        404,
			expectedUrl: "/v1/organisation/accounts/abc123?version=0",
			responseBody: Response{
				ErrorMessage: &message,
			},
			err: &ErrorResponse{Code: 404, Message: "some error!"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			responseBytes, err := json.Marshal(tc.responseBody)
			if err != nil {
				t.Fatal(err)
			}

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tc.expectedUrl, req.RequestURI)
				rw.WriteHeader(tc.code)
				rw.Write(responseBytes)
			}))
			defer server.Close()

			ac := &AccountClient{
				BaseUrl:    server.URL,
				Version:    "v1",
				HttpClient: server.Client(),
			}

			err = ac.DeleteAccount(context.Background(), "abc123", 0)

			if tc.err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {

	type testCase struct {
		name         string
		code         int
		requestBody  interface{}
		responseBody interface{}
		expectedUrl  string
		data         interface{}
		err          error
	}

	testCases := []testCase{
		{
			name: "request ok",
			code: 200,
			requestBody: Request{
				Data: account,
			},
			responseBody: Response{
				Data: account,
			},
			expectedUrl: "/v1/organisation/accounts",
			data:        account,
			err:         nil,
		},
		{
			name: "account not found",
			code: 404,
			requestBody: Request{
				Data: account,
			},
			responseBody: Response{
				ErrorMessage: &message,
			},
			expectedUrl: "/v1/organisation/accounts",
			data:        nil,
			err:         &ErrorResponse{Code: 404, Message: "some error!"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			responseBytes, err := json.Marshal(tc.responseBody)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tc.expectedUrl, req.RequestURI)
				rw.WriteHeader(tc.code)
				rw.Write(responseBytes)
			}))
			defer server.Close()

			ac := &AccountClient{
				BaseUrl:    server.URL,
				Version:    "v1",
				HttpClient: server.Client(),
			}

			data, err := ac.CreateAccount(context.Background(), account)

			if tc.err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			if tc.data != nil {
				assert.Equal(t, tc.data, data)
			} else {
				assert.Nil(t, data)
			}
		})
	}
}

func TestNewAccountClient(t *testing.T) {

	type testCase struct {
		name   string
		cfg    *Config
		client *AccountClient
		url    string
	}

	maxRetries := 2
	zeroRetries := 0
	testCases := []testCase{
		{
			name: "create with config",
			cfg: &Config{
				BaseUrl:    "http://localhost:8080",
				Version:    "v1",
				MaxRetries: maxRetries,
			},
			client: &AccountClient{
				BaseUrl:           "http://localhost:8080",
				Version:           "v1",
				HttpClient:        &http.Client{},
				LimitRateAndRetry: &LimitRateAndRetry{MaxRetries: &maxRetries},
			},
			url: "http://localhost:8080/v1/organisation/accounts",
		},
		{
			name:   "create with nil config",
			cfg:    nil,
			client: &AccountClient{HttpClient: &http.Client{}, LimitRateAndRetry: &LimitRateAndRetry{}},
			url:    "http://api.form3.tech/v1/organisation/accounts",
		},
		{
			name:   "create with empty config",
			cfg:    &Config{},
			client: &AccountClient{HttpClient: &http.Client{}, LimitRateAndRetry: &LimitRateAndRetry{MaxRetries: &zeroRetries}},
			url:    "http://api.form3.tech/v1/organisation/accounts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			client := NewAccountClient(tc.cfg)
			assert.Equal(t, tc.client, client)
			assert.Equal(t, tc.url, client.GetUrl())
		})
	}
}

func createAccount() *AccountData {
	version := int64(0)
	country := "GB"
	accountClassification := "Personal"
	jointAccount := false
	accountMatchingOptOut := false
	switched := false
	status := "confirmed"

	return &AccountData{
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
		Version:        &version,
		CreatedOn:      "2023-05-06T17:25:05.311Z",
		Attributes: &AccountAttributes{
			Country:                 &country,
			BaseCurrency:            "GBP",
			BankID:                  "123456",
			BankIDCode:              "GBDSC",
			Bic:                     "EXMPLGB2XXX",
			AccountNumber:           "12345678",
			Name:                    []string{"FIRST", "LAST"},
			AlternativeNames:        nil,
			AccountClassification:   &accountClassification,
			JointAccount:            &jointAccount,
			AccountMatchingOptOut:   &accountMatchingOptOut,
			SecondaryIdentification: "A1B2C3D4",
			Switched:                &switched,
			Status:                  &status,
		},
	}
}
