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

func CreateAccount(bodyRequest *models.AccountBodyRequest) (*models.AccountData, error) {
	const postPath = "/v1/organisation/accounts/"
	requestURL := fmt.Sprintf("%s%s", SERVER, postPath)
	response, decodeError := utils.PostAccountRequest(requestURL, bodyRequest)

	if decodeError != nil {
		utils.ShowError("CreateAccount", decodeError)
		return nil, errors.New(fmt.Sprintf("Can't do a Post to %s ", requestURL))
	}

	return response, nil
}
