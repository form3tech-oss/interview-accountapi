package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"models"
	"net/http"
)

type ResponseData struct {
	Data  *models.AccountData `json:"data,omitempty"`
	Links struct {
		Self string `json:"self"`
	} `json:"links,omitempty"`
}

var URL = "https://api.form3.tech/"
var ACCOUNTS_URL_V1 = URL + "/v1/organisation/accounts/"

// Check if the API is available (server down, authentication required, etc.)
func Check() (err error) {
	response, err := http.Get(URL + "/v1/health")
	if err == nil {
		if response.StatusCode == http.StatusOK {
			responseData, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			var result map[string]interface{}
			json.Unmarshal([]byte(responseData), &result)
			if result["status"] == "up" {
				err = nil
				return
			}
		}
	}
	err = errors.New("API Unavailable")
	return
}

// Create an organizing account provided in input
// Will fail if the account already exists
func Create(account *models.AccountData) (string, error) {
	err := Check()
	if err == nil {
		type PayloadData struct {
			Data *models.AccountData `json:"data,omitempty"`
		}
		payload := PayloadData{
			Data: account,
		}
		payloadBytes, _ := json.Marshal(payload)
		body := bytes.NewReader(payloadBytes)
		req, err := http.NewRequest(http.MethodPost, ACCOUNTS_URL_V1, body)
		req.Header.Set("Content-Type", "application/vnd.api+json")
		if err == nil {
			response, _ := http.DefaultClient.Do(req)
			if response.StatusCode == http.StatusCreated {
				defer response.Body.Close()
				responseData, _ := ioutil.ReadAll(response.Body)
				var result ResponseData
				json.Unmarshal([]byte(responseData), &result)
				return result.Data.ID, nil
			}
			return "", errors.New(fmt.Sprintf("Record \"%s\" not created", account.ID))
		}
	}
	return "", err
}

// Fetch an existing account providing its id
func Fetch(id string) (models.AccountData, error) {
	var invalid models.AccountData
	err := Check()
	if err == nil {
		response, _ := http.Get(ACCOUNTS_URL_V1 + id)
		if response.StatusCode == http.StatusOK {
			defer response.Body.Close()
			responseData, _ := ioutil.ReadAll(response.Body)
			var result ResponseData
			json.Unmarshal([]byte(responseData), &result)
			return *result.Data, nil
		}
		err = errors.New(fmt.Sprintf("Record \"%s\" not found", id))
	}
	return invalid, err
}

// Delete an existing account providing its id and version
func Delete(id string, version int64) error {
	err := Check()
	if err == nil {
		var url = fmt.Sprintf("%s%s?version=%v", ACCOUNTS_URL_V1, id, version)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err == nil {
			response, _ := http.DefaultClient.Do(request)
			if response.StatusCode == http.StatusNoContent {
				return nil
			} else if response.StatusCode == http.StatusConflict {
				account, _ := Fetch(id)
				suggested_version := *account.Version
				return errors.New(fmt.Sprintf("Invalid version \"%v\" for record \"%s\". Hint: Try with version \"%v\"", version, id, suggested_version))
			}
		}
		return errors.New(fmt.Sprintf("Record \"%s\" not found", id))
	}
	return err
}
