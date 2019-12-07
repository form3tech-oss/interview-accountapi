package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lioda/interview-accountapi/model"
)

type responseDto struct {
	Data string `json:"data"`
}
type arrayDto struct {
	Data []accountDto `json:"data"`
}
type scalarDto struct {
	Data accountDto `json:"data"`
}

type accountDto struct {
	Type           string        `json:"type"`
	ID             string        `json:"id"`
	OrganizationID string        `json:"organisation_id"` //TODO
	Version        string        `json:"version"`
	Attributes     attributesDto `json:"attributes"`
}
type attributesDto struct {
	Country                     string    `json:"country"`
	BaseCurrency                string    `json:"base_currency"`
	BankID                      string    `json:"bank_id"`
	BankIDCode                  string    `json:"bank_id_code"`
	AccountNumber               string    `json:"account_number"`
	Bic                         string    `json:"bic"`
	Iban                        string    `json:"iban"`
	CustomerID                  string    `json:"customer_id"`
	Title                       string    `json:"title"`
	FirstName                   string    `json:"first_name"`
	BankAccountName             string    `json:"bank_account_name"`
	AlternativeBankAccountNames [3]string `json:"alternative_bank_account_names"`
	AccountClassification       string    `json:"account_classification"`
	JointAccount                bool      `json:"joint_account"`
	AccountMatchingOptOut       bool      `json:"account_matching_opt_out"`
	SecondaryIdentification     string    `json:"secondary_identification"`
}

// HTTPRemoteAPI implements RemoteApi using http / https protocol
type HTTPRemoteAPI struct {
	baseURL string
}

// NewHTTPRemote creates a new HTTPRemoteAPI for a base URL
func NewHTTPRemote(baseURL string) HTTPRemoteAPI {
	return HTTPRemoteAPI{baseURL: baseURL}
}

// Get run a GET HTTP request on a path and with queryParams and returns JSON
func (h HTTPRemoteAPI) Get(path string, queryParams string) string {
	// resp, _ := http.Get(h.baseURL + path + "?" + queryParams)
	// dto := responseDto{}
	// decoder := json.NewDecoder(resp.Body)
	// decoder.Decode(&dto)
	// return dto.Data
	panic("not implemented")
}

// GetArray returns an array of account
func (h HTTPRemoteAPI) GetArray(path string, queryParams string) ([]model.Account, error) { //TODO pagination
	resp, err := http.Get(h.baseURL + path + "?" + queryParams)
	fmt.Printf("%v", err)
	if resp.StatusCode != 200 {
		reason, _ := ioutil.ReadAll(resp.Body)
		return []model.Account{}, errors.New(string(reason))
	}

	dtos := arrayDto{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&dtos)

	result := make([]model.Account, 0)
	for _, dto := range dtos.Data {
		account := model.Account{
			ID:                          uuid.MustParse(dto.ID),
			AccountNumber:               dto.Attributes.AccountNumber,
			AccountClassification:       dto.Attributes.AccountClassification,
			AccountMatchingOptOut:       dto.Attributes.AccountMatchingOptOut,
			BankID:                      dto.Attributes.BankID,
			BankIDCode:                  dto.Attributes.BankIDCode,
			BaseCurrency:                dto.Attributes.BaseCurrency,
			BankAccountName:             dto.Attributes.BankAccountName,
			Bic:                         dto.Attributes.Bic,
			Country:                     dto.Attributes.Country,
			CustomerID:                  dto.Attributes.CustomerID,
			FirstName:                   dto.Attributes.FirstName,
			Iban:                        dto.Attributes.Iban,
			JointAccount:                dto.Attributes.JointAccount,
			SecondaryIdentification:     dto.Attributes.SecondaryIdentification,
			Title:                       dto.Attributes.Title,
			AlternativeBankAccountNames: dto.Attributes.AlternativeBankAccountNames,
		}
		result = append(result, account)
	}
	return result, nil
}

// Post run a POST HTTP request on a path with an account
func (h HTTPRemoteAPI) Post(path string, organizationID uuid.UUID, account model.Account) (uuid.UUID, error) {
	toJSON := scalarDto{
		Data: accountDto{
			ID:             account.ID.String(),
			Type:           "accounts",
			Version:        "0",
			OrganizationID: organizationID.String(),
			Attributes: attributesDto{
				AccountClassification:       account.AccountClassification,
				AccountMatchingOptOut:       account.AccountMatchingOptOut,
				AlternativeBankAccountNames: account.AlternativeBankAccountNames,
				AccountNumber:               account.AccountNumber,
				BankAccountName:             account.BankAccountName,
				BankID:                      account.BankID,
				BankIDCode:                  account.BankIDCode,
				BaseCurrency:                account.BaseCurrency,
				Bic:                         account.Bic,
				Country:                     account.Country,
				CustomerID:                  account.CustomerID,
				FirstName:                   account.FirstName,
				Iban:                        account.Iban,
				JointAccount:                account.JointAccount,
				SecondaryIdentification:     account.SecondaryIdentification,
				Title:                       account.Title,
			},
		},
	}
	bytes, _ := json.Marshal(toJSON)
	resp, _ := http.Post(h.baseURL+path, "application/json", strings.NewReader(string(bytes)))
	if resp.StatusCode != 201 {
		reason, _ := ioutil.ReadAll(resp.Body)
		return uuid.Nil, errors.New(string(reason))
	}
	return account.ID, nil
}
