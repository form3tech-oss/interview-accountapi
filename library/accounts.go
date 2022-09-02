package library

import (
	"errors"
	"fmt"
	"github.com/jsebasct/account-api-lib/models"
	"github.com/jsebasct/account-api-lib/utils"
)

const SERVER = "http://localhost:8080"

func ListAccounts(bodyResponse *models.AccountListResponse) error {
	const listPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", SERVER, listPath)

	//decodeError := utils.GetDecodedRequest(requestURL, &bodyResponse)
	decodeError := utils.GetUnmarshalledJson(requestURL, &bodyResponse)

	if decodeError != nil {
		utils.ShowError("ListAccounts", decodeError)
		return errors.New(fmt.Sprintf("Can't retrieve %s resource", requestURL))
	}

	return nil
}

func FetchAccount(accountId string) (*models.AccountBodyResponse, *models.ErrorResponse) {

	fetchPath := fmt.Sprintf("/v1/organisation/accounts/%s", accountId)
	requestURL := fmt.Sprintf("%s%s", SERVER, fetchPath)

	account := models.AccountBodyResponse{}
	decodeError := utils.GetUnmarshalledJson(requestURL, &account)

	if decodeError != nil {
		//utils.ShowError("FetchAccount", decodeError)
		return nil, &models.ErrorResponse{Message: decodeError.Error()}
	}

	return &account, nil
}

func CreateAccount(bodyRequest *models.AccountBodyRequest) (*models.AccountData, *models.ErrorResponse) {
	const postPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", SERVER, postPath)
	response, errorResponse := utils.EvaluatePostAccountResponse(utils.PostAccountRequest(requestURL, bodyRequest))
	return response, errorResponse
}

func DeleteAccount(accountId string, accountVersion int) *models.ErrorResponse {
	var deletePath = fmt.Sprintf("/v1/organisation/accounts/%s?version=%d", accountId, accountVersion)
	requestURL := fmt.Sprintf("%s%s", SERVER, deletePath)

	response := utils.EvaluateDeleteAccountResponse(utils.DeleteAccountRequest(requestURL))
	return response
}
