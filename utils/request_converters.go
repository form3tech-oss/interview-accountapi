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

func GetRequest(url string) (resp *http.Response, err error) {
	response, getError := myClient.Get(url)
	if getError != nil {
		ShowError("GetAccountRequest", getError)
		return nil, getError
	}

	return response, nil
}

func EvaluateGetAccountResponse(getResponse *http.Response, err error) (*models.AccountBodyResponse, *models.ErrorResponse) {
	defer getResponse.Body.Close()

	// evaluate status code
	if getResponse.StatusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		err = UnmarshalTo(getResponse.Body, &errorResponse)
		if err != nil {
			ShowError("EvaluateGetAccountResponse", err)
		}
		errorResponse.Code = getResponse.StatusCode

		return nil, &errorResponse
	}

	var accountBodyRespone models.AccountBodyResponse
	err = UnmarshalTo(getResponse.Body, &accountBodyRespone)
	if err != nil {
		ShowError("EvaluateGetAccountResponse while unmarshalling", err)
		return nil, &models.ErrorResponse{Message: err.Error()}
	}

	return &accountBodyRespone, nil
}

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

	res := UnmarshalTo(response.Body, target)
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
		var errorMessage models.ErrorResponse
		err := UnmarshalTo(responsePost.Body, &errorMessage)
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
	decodeError := UnmarshalTo(responsePost.Body, &accountResponse)
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

	if responseDelete.StatusCode == http.StatusNotFound {
		return &models.ErrorResponse{Code: http.StatusNotFound, Message: "Not Found"}
	}

	if responseDelete.StatusCode != http.StatusNoContent {
		var errorResponse models.ErrorResponse
		unmarshallErr := UnmarshalTo(responseDelete.Body, &errorResponse)
		if unmarshallErr != nil {
			ShowError("EvaluateDeleteAccountResponse", unmarshallErr)
			return &models.ErrorResponse{Message: unmarshallErr.Error()}
		}
		errorResponse.Code = responseDelete.StatusCode
		return &errorResponse
	}

	return nil
}

func UnmarshalTo(body io.ReadCloser, target interface{}) error {
	bodyByte, readError := io.ReadAll(body)
	if readError != nil {
		ShowError("UnmarshalTo while reading", readError)
		return readError
	}

	res := json.Unmarshal(bodyByte, target)
	return res
}
