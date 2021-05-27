package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func main() {
	//url
	url := "http://localhost:8080/v1/organisation/accounts/"

	//create client
	client := http.Client{Timeout: 10 * time.Second}

	//create Account data
	newAccount := Account{
		ID:             uuid.New().String(),
		OrganisationID: uuid.New().String(),
		Type:           "accounts",
		Version:        0,
		Attributes: Attributes{
			BankID:       "400400",
			BankIDCode:   "GBDSC",
			BaseCurrency: "GBP",
			BIC:          "NWBKGB22",
			Country:      "GB",
			Name:         []string{"some name"},
		},
	}

	//save account
	jsonData, jsonErr := json.Marshal(&Data{newAccount})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	createAccountResponse, createErr := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if createErr != nil {
		log.Fatal(createErr)
	}

	defer createAccountResponse.Body.Close()

	// createdAccount := Account

	//get account

	fetchedAccountResponse, fetchErr := client.Get(url + newAccount.ID)
	if fetchErr != nil {
		log.Fatal(fetchErr)
	}

	defer fetchedAccountResponse.Body.Close()

	fetchedAccount := new(Data)

	json.NewDecoder(createAccountResponse.Body).Decode(fetchedAccount)

	fmt.Println(fetchedAccount)

	//delete account
	deleteAccountRequest, err := http.NewRequest("DELETE", url+newAccount.ID+"?version="+strconv.Itoa(newAccount.Version), nil)
	if err != nil {
		log.Fatal(err)
	}

	deleteAccountResponse, delErr := client.Do(deleteAccountRequest)

	if delErr != nil {
		log.Fatal(delErr)
	}

	defer deleteAccountResponse.Body.Close()

	//validate created account is deleted
	fetchDeletedAccountResponse, err := client.Get(url + newAccount.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fetchDeletedAccountResponse.Status)

}

type Data struct {
	Body Account `json:"data"`
}

type Account struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
}

type Attributes struct {
	AlternativeNames []string `json:"alternative_names"`
	BankID           string   `json:"bank_id"`
	BankIDCode       string   `json:"bank_id_code"`
	BaseCurrency     string   `json:"base_currency"`
	BIC              string   `json:"bic"`
	Country          string   `json:"country"`
	Name             []string `json:"name"`
	// AccountClassification   string                `json:"account_classification"`
	JointAccount            bool   `json:"joint_account"`
	AccountMatchingOptOut   bool   `json:"account_matching_opt_out"`
	SecondaryIdentification string `json:"secondary_identification"`
	// Status                  string                `json:"status"`
	PrivateIdentification PrivateIdentification `json:"private_identification"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date"`
	BirthCountry   string   `json:"birth_country"`
	Identification string   `json:"identificaiton"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}
