package accountapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/giannimassi/interview-accountapi/model"
	"github.com/stretchr/testify/require"
)

func TestNewAccountAPIClient(t *testing.T) {
	c := NewAccountAPIClient(WithHTTPClient(http.DefaultClient))

	require.Equal(t, http.DefaultClient, c.client, "client do not match: expected %s, got %s", http.DefaultClient, c.client)
}

func TestAccountAPIClient_Create(t *testing.T) {
	var testAccountData1 = model.AccountData{ID: "fakeid"}

	tests := []struct {
		name string

		// Input
		account model.AccountData

		// Mocking
		responseBody   []byte
		responseStatus int
		ctx            context.Context

		// Output
		wantReturn *model.AccountData
		wantErr    error
	}{
		{
			name:    "ok",
			account: testAccountData1,

			responseBody:   []byte(`{"data":{"ID": "fakeid"}}`),
			responseStatus: http.StatusCreated,
			ctx:            context.Background(),

			wantReturn: &testAccountData1,
			wantErr:    nil,
		},
		{
			name:    "bad request",
			account: model.AccountData{},

			responseBody:   []byte(``),
			responseStatus: http.StatusBadRequest,
			ctx:            context.Background(),

			wantReturn: nil,
			wantErr:    fmt.Errorf("unexpected status code: %d", http.StatusBadRequest),
		},
		{
			name:    "error decoding response",
			account: model.AccountData{},

			responseBody:   []byte(``),
			responseStatus: http.StatusCreated,
			ctx:            context.Background(),

			wantReturn: nil,
			wantErr:    fmt.Errorf("while decoding response: EOF"),
		},
		{
			name:    "nil context",
			account: model.AccountData{},

			responseBody:   []byte(``),
			responseStatus: http.StatusCreated,
			ctx:            nil,

			wantErr: fmt.Errorf("while formatting request: net/http: nil Context"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup mocked client and server
			mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Write(tt.responseBody)
			}))
			defer mockedServer.Close()

			a := NewAccountAPIClient(WithHTTPClient(mockedServer.Client()), WithHost(mockedServer.URL))
			got, err := a.Create(tt.ctx, tt.account)
			if tt.wantErr == nil {
				require.NoError(t, err, "unexpected error: %s", err)
			} else {
				require.EqualError(t, err, tt.wantErr.Error(), "error do not match: expected %s, got %s", tt.wantErr, err)
			}
			require.Equal(t, tt.wantReturn, got, "response do not match: expected %s, got %s", tt.wantReturn, got)
		})
	}
}

func TestAccountAPIClient_Fetch(t *testing.T) {
	var testAccountData1 = model.AccountData{ID: "fakeid"}

	tests := []struct {
		name string

		// Input
		account model.AccountData

		// Mocking
		responseBody   []byte
		responseStatus int
		ctx            context.Context

		// Output
		wantReturn *model.AccountData
		wantErr    error
	}{
		{
			name:    "ok",
			account: testAccountData1,

			responseBody:   []byte(`{"data":{"ID": "fakeid"}}`),
			responseStatus: http.StatusOK,
			ctx:            context.Background(),

			wantReturn: &testAccountData1,
			wantErr:    nil,
		},
		{
			name:    "bad request",
			account: model.AccountData{},

			responseBody:   []byte(``),
			responseStatus: http.StatusBadRequest,
			ctx:            context.Background(),

			wantReturn: nil,
			wantErr:    fmt.Errorf("unexpected status code: %d", http.StatusBadRequest),
		},
		{
			name:    "error decoding response",
			account: model.AccountData{},

			responseBody:   []byte(``),
			responseStatus: http.StatusOK,
			ctx:            context.Background(),

			wantReturn: nil,
			wantErr:    fmt.Errorf("while decoding response: EOF"),
		},
		{
			name: "nil context",

			account: model.AccountData{},

			responseBody:   []byte(``),
			responseStatus: http.StatusOK,
			ctx:            nil,

			wantErr: fmt.Errorf("while formatting request: net/http: nil Context"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup mocked client and server
			mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Write(tt.responseBody)
			}))
			defer mockedServer.Close()

			a := NewAccountAPIClient(WithHTTPClient(mockedServer.Client()), WithHost(mockedServer.URL))
			got, err := a.Fetch(tt.ctx, tt.account)
			if tt.wantErr == nil {
				require.NoError(t, err, "unexpected error: %s", err)
			} else {
				require.EqualError(t, err, tt.wantErr.Error(), "error do not match: expected %s, got %s", tt.wantErr, err)
			}
			require.Equal(t, tt.wantReturn, got, "response do not match: expected %s, got %s", tt.wantReturn, got)
		})
	}
}

func TestAccountAPIClient_Delete(t *testing.T) {
	var testAccountData1 = model.AccountData{ID: "fakeid", Version: new(int64)}

	tests := []struct {
		name string

		// Input
		account model.AccountData

		// Mocking
		responseStatus int
		ctx            context.Context

		// Output
		wantErr error
	}{
		{
			name:    "ok",
			account: testAccountData1,

			responseStatus: http.StatusNoContent,
			ctx:            context.Background(),

			wantErr: nil,
		},
		{
			name:    "resource not found",
			account: testAccountData1,

			responseStatus: http.StatusNotFound,
			ctx:            context.Background(),

			wantErr: fmt.Errorf("specified resource does not exist"),
		},
		{
			name:    "resources incorrect",
			account: testAccountData1,

			responseStatus: http.StatusConflict,
			ctx:            context.Background(),

			wantErr: fmt.Errorf("specified version incorrect"),
		},

		{
			name:    "unexpected error",
			account: testAccountData1,

			responseStatus: http.StatusBadRequest,
			ctx:            context.Background(),

			wantErr: fmt.Errorf("unexpected status code: %d", http.StatusBadRequest),
		},
		{
			name:    "nil context",
			account: testAccountData1,

			responseStatus: http.StatusNoContent,
			ctx:            nil,

			wantErr: fmt.Errorf("while formatting request: net/http: nil Context"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup mocked client and server
			mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseStatus)
			}))
			defer mockedServer.Close()

			a := NewAccountAPIClient(WithHTTPClient(mockedServer.Client()), WithHost(mockedServer.URL))
			err := a.Delete(tt.ctx, tt.account)
			if tt.wantErr == nil {
				require.NoError(t, err, "unexpected error: %s", err)
			} else {
				require.EqualError(t, err, tt.wantErr.Error(), "error do not match: expected %s, got %s", tt.wantErr, err)
			}
		})
	}
}
