package form3client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
	CreatedOn      string             `json:"created_on,omitempty"`
	ModifiedOn     string             `json:"modified_on,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

var ErrNoAccountData error = errors.New("no data provided to create an account")
var ErrParsingResponse error = errors.New("unable to parse JSON response")

func ClientError(err error) error {
	return fmt.Errorf("client error: %v", err)
}

// Create creates a new account
func CreateAccount(jsonAccountData string) (string, error) {

	// If no accountData was given, return an error with a message.
	if jsonAccountData == "" {
		return "", ClientError(ErrNoAccountData)
	}

	// we parse the req to send it
	jsonData := []byte(jsonAccountData)

	// We post the request
	resp, err := http.Post("http://localhost:8080/v1/organisation/accounts", "application/vnd.api+json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", ClientError(err)
	}

	// if we return => we close the response
	defer resp.Body.Close()

	log.Println(resp)
	// we read the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ClientError(err)
	}

	// we return the response
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("API error: %v; %w", resp.StatusCode, errors.New(string(bodyBytes)))
	} else {
		return string(bodyBytes), nil
	}

}

// Delete an existing account
func DeleteAccount(accountID, version string) (string, error) {
	return "", nil
}
