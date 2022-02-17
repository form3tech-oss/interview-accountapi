package accountapiclient

import (
	"accountapiclient"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var CurrnetVersion string
var apiAddress string

type Client interface {
	CreateAccount(account accountapiclient.AccountData) (accountapiclient.AccountData, []string)
	FetchAccount(accountId string) (accountapiclient.AccountData, []string)
	//DeleteAccount(accountId string) []string
}

type ApiClientV1 struct{}

func CreateClient(baseUrl string) Client {
	CurrnetVersion = "v1"
	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl = baseUrl + "/"
	}

	apiAddress = baseUrl + CurrnetVersion
	return ApiClientV1{}
}

func sendRequest(method string, endpoint string, requestBody string) ([]byte, int) {
	request, e := http.NewRequest(method, apiAddress+endpoint, bytes.NewBuffer([]byte(requestBody)))
	if e != nil {
		log.Print(e.Error())
		return nil, 400
	}

	httpclient := &http.Client{}
	response, e := httpclient.Do(request)
	if e != nil {
		log.Printf("error calling remote address: '%s'", e.Error())
		return nil, 503
	}

	defer response.Body.Close()
	responsecontent, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Printf("error reading http response: '%s'", e.Error())
	}
	return responsecontent, response.StatusCode
}

func (ApiClientV1) CreateAccount(account accountapiclient.AccountData) (accountapiclient.AccountData, []string) {
	if (account == accountapiclient.AccountData{}) {
		return accountapiclient.AccountData{}, []string{"invalid input"}
	}

	var errors []string
	requestbody, e := json.Marshal(accountapiclient.Account{
		AccountData: account,
	})
	if e != nil {
		log.Print(e.Error())
		return accountapiclient.AccountData{}, []string{"invalid input structure"}
	}

	var newaccount accountapiclient.AccountData
	response, statuccode := sendRequest("POST", "/organisation/accounts", string(requestbody))
	if statuccode >= 200 && statuccode <= 299 {
		json.NewDecoder(bytes.NewBuffer(response)).Decode(&newaccount)
	} else if statuccode >= 400 && statuccode <= 499 {
		var apierrors accountapiclient.ApiErrors
		json.NewDecoder(bytes.NewBuffer(response)).Decode(&apierrors)
		errors = append(errors, strings.Split(apierrors.ErrorMessage, "\n")...)
	} else {
		errors = append(errors, "somethign went wrong, try again")
	}

	if len(errors) > 0 {
		return accountapiclient.AccountData{}, errors
	}
	return account, nil
}

func (ApiClientV1) FetchAccount(accountId string) (accountapiclient.AccountData, []string) {
	var account accountapiclient.Account
	id, e := uuid.Parse(accountId)
	if e != nil || id == uuid.Nil {
		return accountapiclient.AccountData{}, []string{"invalid account id"}
	}

	responsebody, statuscode := sendRequest("GET", "/organisation/accounts/"+accountId, "")
	if statuscode >= 200 && statuscode <= 299 {
		json.NewDecoder(bytes.NewBuffer(responsebody)).Decode(&account)
	} else {
		var apiError accountapiclient.ApiErrors
		json.NewDecoder(bytes.NewBuffer(responsebody)).Decode(&apiError)
		return accountapiclient.AccountData{}, strings.Split(apiError.ErrorMessage, "\n")
	}

	return account.AccountData, nil
}
