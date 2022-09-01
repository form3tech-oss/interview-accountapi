package utils

import (
	"bytes"
	"encoding/json"
	"github.com/jsebasct/account-api-lib/models"
	"io"
	"net/http"
	"time"
)

const HTTP_STATUS_CODE_CREATED = 201

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetDecodedRequest(url string, target interface{}) error {
	response, getError := myClient.Get(url)
	if getError != nil {
		ShowError("GetDecodedRequest", getError)
		return getError
	}
	defer response.Body.Close()

	decodeError := json.NewDecoder(response.Body).Decode(&target)
	return decodeError
}

func GetUnmarshalledJson(url string, target interface{}) error {

	response, getError := myClient.Get(url)
	if getError != nil {
		ShowError("GetUnmarshalledJson", getError)
		return getError
	}
	defer response.Body.Close()

	bodyByte, readError := io.ReadAll(response.Body)
	if readError != nil {
		ShowError("GetUnmarshalledJson while reading", readError)
		return readError
	}

	res := json.Unmarshal(bodyByte, target)
	return res
}

func PostAccountRequest(url string, bodyRequest *models.AccountBodyRequest) (resp *http.Response, err error) {
	jsonData, marshallErr := json.Marshal(bodyRequest)

	if marshallErr != nil {
		ShowError("PostAccountRequest while reading", marshallErr)
		return nil, marshallErr
	}

	responsePost, postError := myClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return responsePost, postError
}

func EvaluatePostAccountResponse(responsePost *http.Response, postError error) (*models.AccountData, *models.ErrorResponse) {

	if postError != nil {
		ShowError("PostAccountRequest", postError)
		return nil, &models.ErrorResponse{Code: responsePost.StatusCode, Message: postError.Error()}
	}
	defer responsePost.Body.Close()

	if responsePost.StatusCode != HTTP_STATUS_CODE_CREATED {
		bodyErrorByte, readError := io.ReadAll(responsePost.Body)
		if readError != nil {
			ShowError("PostAccountRequest while reading ERROR", readError)
			return nil, &models.ErrorResponse{Code: responsePost.StatusCode, Message: readError.Error()}
		}

		errorMessage := models.ErrorResponse{}
		err := json.Unmarshal(bodyErrorByte, &errorMessage)
		if err != nil {
			ShowError("PostAccountRequest while unmarshalling", err)
			return nil, &models.ErrorResponse{Code: responsePost.StatusCode, Message: err.Error()}
		} else {
			return nil, &errorMessage
		}
	}

	var accountResponse models.AccountData
	decodeError := json.NewDecoder(responsePost.Body).Decode(&accountResponse)
	return &accountResponse, &models.ErrorResponse{Code: responsePost.StatusCode, Message: decodeError.Error()}
}
