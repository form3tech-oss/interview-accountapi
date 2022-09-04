package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsebasct/account-api-lib/models"
	"io"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetUnmarshalledJson(url string, target interface{}) error {

	response, getError := myClient.Get(url)
	if getError != nil {
		ShowError("GetUnmarshalledJson", getError)
		return getError
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("Expected 200 but got %d", response.StatusCode))
		ShowError("GetUnmarshalledJson", err)
		return err
	}

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

	responsePost, postError := myClient.Post(url, "application/vnd.api+json", bytes.NewBuffer(jsonData))
	return responsePost, postError
}

func EvaluatePostAccountResponse(responsePost *http.Response, postError error) (*models.AccountBodyResponse, *models.ErrorResponse) {

	if postError != nil {
		ShowError("EvaluatePostAccountResponse", postError)
		return nil, &models.ErrorResponse{Code: responsePost.StatusCode, Message: postError.Error()}
	}

	defer responsePost.Body.Close()

	if responsePost.StatusCode != http.StatusCreated {
		bodyErrorByte, readError := io.ReadAll(responsePost.Body)
		if readError != nil {
			ShowError("EvaluatePostAccountResponse while reading ERROR", readError)
			return nil, &models.ErrorResponse{Code: responsePost.StatusCode, Message: readError.Error()}
		}

		fmt.Println("bad body as string", string(bodyErrorByte))

		var errorMessage models.ErrorResponse
		err := json.Unmarshal(bodyErrorByte, &errorMessage)
		if err != nil {
			ShowError("EvaluatePostAccountResponse while unmarshalling", err)
			return nil, &models.ErrorResponse{Code: responsePost.StatusCode, Message: err.Error()}
		} else {
			internalErrorMsg := fmt.Sprintf("EvaluatePostAccountResponse with status code: %d", responsePost.StatusCode)
			ShowError(internalErrorMsg, errors.New(errorMessage.Message))
			return nil, &errorMessage
		}
	}

	var accountResponse models.AccountBodyResponse
	decodeError := json.NewDecoder(responsePost.Body).Decode(&accountResponse)
	if decodeError != nil {
		return nil, &models.ErrorResponse{Message: decodeError.Error()}
	}
	return &accountResponse, nil
}

func DeleteAccountRequest(url string) (responseDelete *http.Response, err error) {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseDelete, err = myClient.Do(req)
	return responseDelete, err
}

func EvaluateDeleteAccountResponse(responseDelete *http.Response, deleteError error) *models.ErrorResponse {
	if deleteError != nil {
		ShowError("EvaluateDeleteAccountResponse", deleteError)
		return &models.ErrorResponse{Message: deleteError.Error()}
	}

	defer responseDelete.Body.Close()

	if responseDelete.StatusCode != http.StatusNoContent {
		statusCodeError := errors.New(fmt.Sprintf("Status Code %d different than %d", responseDelete.StatusCode, http.StatusNoContent))
		ShowError("EvaluateDeleteAccountResponse", statusCodeError)
		return &models.ErrorResponse{Message: statusCodeError.Error()}
	}

	return nil
}
