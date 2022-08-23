package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type FetchAccountParams struct {
	ID             string `json:"id,omitempty"`
	OrganisationID string `json:"organizationId,omitempty"`
	Type           string `json:"type,omitempty"`
	Version        *int64 `json:"version,omitempty"`
}

type FetchAccountResult struct {
	// accountData AccountData
}

func handleOk(response *http.Response) (*AccountData, error) {
	accountData := new(AccountData)
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(accountData)
	return accountData, err
}

func handleBadRequest(response *http.Response) (*AccountData, error) {
	var error_messages map[string]string
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&error_messages)
	if err != nil {
		return nil, err
	}
	return nil, errors.New(error_messages["error_message"])
}

func defaultHandler() (*AccountData, error) {
	// TODO: Client side validation not requested for the submission
	panic("Not implemented!")
}

var handlers map[int]func(*http.Response) (*AccountData, error)

func init() {
	handlers = map[int]func(*http.Response) (*AccountData, error){
		200: handleOk,
		400: handleBadRequest,
	}
}

// Avoided the pointer as to make it clear that the private fields
// are meant to remain unchanged
func (client Client) FetchAccount(params FetchAccountParams) (*AccountData, error) {
	url := fmt.Sprintf("%s/v1/organisation/accounts/%s", client.base_url, params.ID)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseHandler, ok := handlers[resp.StatusCode]
	if !ok {
		return defaultHandler()
	} else {
		accountData, err := responseHandler(resp)
		return accountData, err
	}
}
