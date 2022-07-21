package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//region COMMON MODELS

type Attributes struct {
	AlternativeNames []string `json:"alternative_names,omitempty"`
	BankID           string   `json:"bank_id,omitempty"`
	BankIDCode       string   `json:"bank_id_code,omitempty"`
	BaseCurrency     string   `json:"base_currency,omitempty"`
	Bic              string   `json:"bic,omitempty"`
	Country          string   `json:"country,omitempty"`
	Name             []string `json:"name,omitempty"`
}

type Data struct {
	Attributes     Attributes `json:"attributes,omitempty"`
	CreatedOn      time.Time  `json:"created_on,omitempty"`
	ID             string     `json:"id,omitempty"`
	ModifiedOn     time.Time  `json:"modified_on,omitempty"`
	OrganisationID string     `json:"organisation_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Version        float64    `json:"version,omitempty"`
}

type Links struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Self  string `json:"self,omitempty"`
}

//endregion

//region FETCH MODELS

type GetAccountByIdBackendResult struct {
	Data  Data `json:"data"`
	Links `json:"links"`
}

type GetAccountByIdResult struct {
	CreatedOn  time.Time  `json:"created_on"`
	Attributes Attributes `json:"attributes"`
}

//endregion

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		Fetch(w, r)
		return
	case http.MethodPut:
		Create(w, r)
		return
	case http.MethodDelete:
		Delete(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func Fetch(w http.ResponseWriter, r *http.Request) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	param := r.URL.Query().Get("account_id")

	url := "http://localhost:8080/v1/organisation/accounts/" + param
	resp, errGet := c.Get(url)

	if errGet != nil {
		fmt.Printf("Error %s", errGet.Error())
		return
	}

	defer resp.Body.Close()
	body, errReadAll := ioutil.ReadAll(resp.Body)

	if errReadAll != nil {
		fmt.Printf("Error %s", errReadAll.Error())
		return
	}

	if resp.StatusCode != 200 {
		var out bytes.Buffer

		json.Indent(&out, body, "", "  ")

		w.WriteHeader(resp.StatusCode)
		w.Write(out.Bytes())

		return
	}

	var backendResult GetAccountByIdBackendResult
	errJsonUnmarshal := json.Unmarshal(body, &backendResult)

	//map
	var result GetAccountByIdResult
	result.Attributes = backendResult.Data.Attributes
	result.CreatedOn = backendResult.Data.CreatedOn

	if errJsonUnmarshal != nil {
		fmt.Println("Failed on retrieve data", errJsonUnmarshal.Error())
		return
	}

	jsonBytes, errJsonMarshal := json.Marshal(result)
	if errJsonMarshal != nil {
		http.Error(w, errJsonMarshal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func Create(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Create")
	return
}

func Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Delete")
	return
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/accounts", ServeHTTP)
	http.ListenAndServe("localhost:8081", mux)
}
