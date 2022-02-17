package accountapiclient

import (
	"accountapiclient"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var CurrnetVersion string
var apiAddress string

type Client interface {
	CreateAccount(account accountapiclient.AccountData) (accountapiclient.AccountData, []string)
	//FetchAccount(accountId string) (accountapiclient.AccountData, []string)
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
		log.Fatal(e.Error())
	}

	httpclient := &http.Client{}
	response, e := httpclient.Do(request)
	if e != nil {
		log.Fatal(e.Error())
	}

	defer response.Body.Close()
	responsecontent, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Fatal(e.Error())
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
		log.Fatal(e.Error())
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
