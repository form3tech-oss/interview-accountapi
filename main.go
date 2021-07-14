package main

import (
	"api"
	"fmt"
	"models"
)

func main() {
	data := models.CreateAccountFromJSONFile("account_gb.json")
	fmt.Println(data.ID)
	id, err := api.Create(&data)
	fmt.Println(id)
	_, err = api.Create(&data)
	fmt.Println(err)
	account, err := api.Fetch(data.ID)
	if err == nil {
		fmt.Println(account.ID, *account.Version)
		err = api.Delete(id, 1)
		fmt.Println(err)
		err = api.Delete("id", *account.Version)
		fmt.Println(err)
		api.Delete(id, *account.Version)
		_, err = api.Fetch(id)
		fmt.Println(err)
	}
}
