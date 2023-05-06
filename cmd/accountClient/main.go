package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/asiless/interview-accountapi/pkg/account"
)

func main() {

	cfg := account.Config{
		BaseUrl: "localhost:8080",
		Version: "v1",
	}

	client := account.NewAccountClient(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.FetchAccount(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	if err == nil {
		printAccount(res)
		err = client.DeleteAccount(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", *res.Version)
		if err == nil {
			fmt.Printf("Deleted account %s %d\n", res.ID, res.Version)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	res, err = client.CreateAccount(ctx, createRequest())
	if err == nil {
		printAccount(res)
	} else {
		fmt.Println(err)
	}

	var errorResponse *account.ErrorResponse

	res, err = client.CreateAccount(ctx, createRequest())
	if err == nil {
		printAccount(res)
	} else if errors.As(err, &errorResponse) {
		fmt.Printf("ESTA OK code => %d\n", errorResponse.Code)
		fmt.Printf("ESTA OK message => %s\n", errorResponse.Message)
		fmt.Printf("ESTA OK err => %v", err)
		fmt.Printf("ESTA OK errorResponse => %v", errorResponse)

	} else {
		fmt.Println("ESTA MAL: %w", err)
	}
	res, err = client.FetchAccount(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	if err != nil {
		printAccount(res)
	} else {
		fmt.Println(err)
	}
}

func createRequest() *account.AccountData {
	version := int64(0)
	country := "GB"
	accountClassification := "Personal"
	jointAccount := false
	accountMatchingOptOut := false
	switched := false
	status := "confirmed"

	return &account.AccountData{
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
		Version:        &version,
		CreatedOn:      "2021-01-01T00:00:00Z",
		Attributes: &account.AccountAttributes{
			Country:                 &country,
			BaseCurrency:            "GBP",
			BankID:                  "123456",
			BankIDCode:              "GBDSC",
			Bic:                     "EXMPLGB2XXX",
			AccountNumber:           "12345678",
			Name:                    []string{"BETO", "SILESS"},
			AlternativeNames:        []string{"BETO", "SILESS"},
			AccountClassification:   &accountClassification,
			JointAccount:            &jointAccount,
			AccountMatchingOptOut:   &accountMatchingOptOut,
			SecondaryIdentification: "A1B2C3D4",
			Switched:                &switched,
			Status:                  &status,
		},
	}
}

func printAccount(res *account.AccountData) {
	if res == nil {
		fmt.Println("<<nil>>")
		return
	}

	j, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("error: %w", err)
		return
	}
	fmt.Printf("%s\n", string(j))
}

// {
// 	"data": {
// 	  "type": "accounts",
// 	  "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
// 	  "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
// 	  "attributes": {
// 		"name": ["BETO", "SILESS"],
// 		"country": "GB",
// 		"base_currency": "GBP",
// 		"bank_id": "123456",
// 		"bank_id_code": "GBDSC",
// 		"bic": "EXMPLGB2XXX",
// 		"user_defined_data": [
// 		  {
// 			"key": "account_related_key",
// 			"value": "account_related_value"
// 		  }
// 		],
// 		"validation_type": "card",
// 		"reference_mask": "############",
// 		"acceptance_qualifier": "same_day",
// 		"switched_account_details": {
// 		  "switched_effective_date": "2022-07-23",
// 		  "account_number": "12345678",
// 		  "account_with": {
// 			"bank_id": "123456",
// 			"bank_id_code": "GBDSC"
// 		  },
// 		  "account_number_code": "BBAN",
// 		  "account_type": 0
// 		}
// 	  }
// 	}
//   }
