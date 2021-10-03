package endpoint

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"../../src/messages"
)

func GetHealth() {

	response, err := http.Get(BASE_URL + HEALTH)

	if err != nil {
		fmt.Printf(messages.REQ_NOT_FOUND+" with error %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
