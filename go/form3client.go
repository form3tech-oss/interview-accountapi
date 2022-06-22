package form3client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var ErrNoAccountData error = errors.New("expected JSON string not provided")
var ErrParsingResponse error = errors.New("unable to parse JSON response")

const DefaultBaseUrl string = "http://localhost:8080"

// Get baseUrl
func GetBaseUrl() string {
	url := os.Getenv("ACCOUNT_API_BASE_URL")
	if len(url) == 0 {
		return DefaultBaseUrl
	}
	return url
}

type ClientResponse struct {
	StatusCode int
	Body       string
}

func ClientError(err error) error {
	return fmt.Errorf("client error: %v", err)
}

// Create client
func Client() *http.Client {
	return &http.Client{}
}

func ParseResponse(resp *http.Response, okCode int) (ClientResponse, error) {
	// Read Response Body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}

	// We return the response
	bodyString := string(bodyBytes)
	if resp.StatusCode != okCode {
		apiError := errors.New(bodyString)
		errorString := fmt.Errorf("API error: %v; %w", resp.StatusCode, apiError)
		return ClientResponse{resp.StatusCode, bodyString}, errorString
	} else {
		return ClientResponse{resp.StatusCode, bodyString}, nil
	}
}

func DeletionRequest(url string) (ClientResponse, error) {
	client := Client()

	// Create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}
	defer resp.Body.Close()

	// Parse the response
	return ParseResponse(resp, http.StatusNoContent)
}

func FetchRequest(url string) (ClientResponse, error) {
	client := Client()

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}
	defer resp.Body.Close()

	// Parse the response
	return ParseResponse(resp, http.StatusOK)
}

func PostRequest(url string, jsonAccountData string) (ClientResponse, error) {
	client := Client()

	// If no accountData was given, return an error with a message.
	if jsonAccountData == "" {
		return ClientResponse{}, ClientError(ErrNoAccountData)
	}

	// we parse the req to send it
	jsonData := []byte(jsonAccountData)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return ClientResponse{}, ClientError(err)
	}
	defer resp.Body.Close()

	// Parse the response
	return ParseResponse(resp, http.StatusCreated)
}
