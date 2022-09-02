package main

import (
	"fmt"
	"github.com/jsebasct/account-api-lib/library"
	"github.com/jsebasct/account-api-lib/models"
	"github.com/jsebasct/account-api-lib/utils"
)

func main() {
	fmt.Println("---- Start ----")

	bodyResponse := models.AccountListResponse{}
	err := library.ListAccounts(&bodyResponse)

	if err != nil {
		utils.ShowError("Error trying to get account", err)
		return
	}

	fmt.Printf("%+v\n", bodyResponse)
	// fmt.Printf("%+v\n", *bodyResponse.Data[0].Attributes.Country)
	fmt.Println("==== End ====")
}
