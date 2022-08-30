package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	//r, err := http.Get(url)
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println("HUbo un error al hacer el GET")
		return err
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("HUbo un error al Lllevarlo a cadena")
		log.Fatalln(err)
	}

	//Convert the body to type string
	bodyAsString := string(body)
	fmt.Println("bodyAsString", bodyAsString)

	//var bodyAsStruct AccountBodyResponse
	res := json.Unmarshal([]byte(bodyAsString), target)

	return res
}
