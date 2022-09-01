package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jsebasct/account-api-lib/models"
	"io"
	"net/http"
	"time"
)

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

func PostAccountRequest(url string, bodyRequest *models.AccountBodyRequest) (*models.AccountData, error) {
	//requestBodyBytes := new(bytes.Buffer)
	//json.NewEncoder(requestBodyBytes).Encode(&bodyRequest)
	jsonData, marshallErr := json.Marshal(bodyRequest)

	if marshallErr != nil {
		ShowError("PostAccountRequest while reading", marshallErr)
		return nil, marshallErr
	}

	responsePost, postError := myClient.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if postError != nil {
		ShowError("PostAccountRequest", postError)
		return nil, postError
	}
	defer responsePost.Body.Close()

	if responsePost.StatusCode != 201 {

		//extracted
		bodyErrorByte, readError := io.ReadAll(responsePost.Body)
		if readError != nil {
			ShowError("PostAccountRequest while reading ERROR", readError)
			return nil, readError
		}

		errorMessage := models.ErrorResponse{}
		err := json.Unmarshal(bodyErrorByte, &errorMessage)
		if err != nil {
			ShowError("PostAccountRequest while marshalling", err)
			return nil, err
		} else {
			//fmt.Printf("%+v\n", errorMessage)
			return nil, errors.New(errorMessage.Message)
		}
		//
	}

	var accountResponse models.AccountData
	decodeError := json.NewDecoder(responsePost.Body).Decode(&accountResponse)
	return &accountResponse, decodeError
}
