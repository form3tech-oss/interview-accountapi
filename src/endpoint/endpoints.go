package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../src/messages"
)

func GetHealth() {

	response, err := http.Get(HEALTH_URL)

	if err != nil {
		fmt.Printf(messages.REQ_NOT_FOUND+" with error %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		defer response.Body.Close()
	}
}

func GetAccounts() {
	response, err := http.Get(ACCOUNT_URL)

	if err != nil {
		fmt.Printf(messages.REQ_NOT_FOUND+" with error %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func PostAccount(accountData map[string]interface{}) {

	fmt.Println(accountData)
	jsonValue, err := json.Marshal(accountData)

	response, err := http.Post(ACCOUNT_URL, APP_JSON, bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf(messages.REQ_NOT_FOUND+" with error %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func DeleteAccount(id string) {
	client := &http.Client{}

	request, err := http.NewRequest("DELETE", DELETE_URL+id, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()
}
