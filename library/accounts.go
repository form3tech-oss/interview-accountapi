package library

import (
	"errors"
	"fmt"
	"github.com/jsebasct/account-api-lib/models"
	"github.com/jsebasct/account-api-lib/utils"
	"os"
)

func GetServerURL() string {
	var server = os.Getenv("SERVER_URL")
	if server == "" {
		server = "http://localhost:8080"
	}
	return server
}

func ListAccounts(bodyResponse *models.AccountListResponse) error {
	server := GetServerURL()
	fmt.Println("SERVER", server)
	const listPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", server, listPath)

	decodeError := utils.GetUnmarshalledJson(requestURL, &bodyResponse)
	//account, err := utils.EvaluateGetAccountResponse(utils.GetRequest(requestURL))

	if decodeError != nil {
		utils.ShowError("ListAccounts", decodeError)
		return errors.New(fmt.Sprintf("Can't retrieve %s resource", requestURL))
	}

	return nil
}

func FetchAccount(accountId string) (*models.AccountBodyResponse, *models.ErrorResponse) {
	server := GetServerURL()
	fetchPath := fmt.Sprintf("/v1/organisation/accounts/%s", accountId)
	requestURL := fmt.Sprintf("%s%s", server, fetchPath)

	account, err := utils.EvaluateGetAccountResponse(utils.GetRequest(requestURL))
	return account, err
}

func CreateAccount(bodyRequest *models.AccountBodyRequest) (*models.AccountBodyResponse, *models.ErrorResponse) {
	server := GetServerURL()
	const postPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", server, postPath)
	fmt.Println("create account to URL: ", requestURL)

	response, errorResponse := utils.EvaluatePostAccountResponse(utils.PostAccountRequest(requestURL, bodyRequest))
	return response, errorResponse
}

func DeleteAccount(accountId string, accountVersion int) *models.ErrorResponse {
	server := GetServerURL()
	var deletePath = fmt.Sprintf("/v1/organisation/accounts/%s?version=%d", accountId, accountVersion)
	requestURL := fmt.Sprintf("%s%s", server, deletePath)

	response := utils.EvaluateDeleteAccountResponse(utils.DeleteAccountRequest(requestURL))
	return response
}
