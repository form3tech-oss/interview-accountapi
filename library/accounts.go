package library

import (
	"errors"
	"fmt"
	"github.com/jsebasct/account-api-lib/models"
	"github.com/jsebasct/account-api-lib/utils"
)

const SERVER = "http://localhost:8080"

func GetAccounts(bodyResponse *models.AccountBodyResponse) error {
	const getPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", SERVER, getPath)
	//decodeError := utils.GetDecodedRequest(requestURL, &bodyResponse)
	decodeError := utils.GetUnmarshalledJson(requestURL, &bodyResponse)

	if decodeError != nil {
		utils.ShowError("GetAccounts", decodeError)
		return errors.New(fmt.Sprintf("Can't retrieve %s resource", requestURL))
	}

	return nil
}

func CreateAccount(bodyRequest *models.AccountBodyRequest) (*models.AccountData, *models.ErrorResponse) {
	const postPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", SERVER, postPath)
	response, errorResponse := utils.EvaluatePostAccountResponse(utils.PostAccountRequest(requestURL, bodyRequest))
	return response, errorResponse
}
