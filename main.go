package main

import (
	"encoding/json"

	"./src/processes/accountProcess"
	"./src/processes/healthProcess"
)

func main() {
	healthProcess.GetHealth()

	accountProcess.PostAccount(GetAccountToPost())

	accountProcess.GetAccounts()

	id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	accountProcess.DeleteAccount(id)
}

func GetAccountToPost() map[string]interface{} {
	str := `{"attributes":{"account_classification":"Personal","account_matching_opt_out":false,"account_number":"41426819","alternative_names":["Sam Holder"],"bank_id":"400300","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB22","country":"GB","iban":"GB11NWBK40030041426819","joint_account":false,"name":["Samantha Holder"],"secondary_identification":"A1B2C3D4","status":"confirmed","switched":false},"id":"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc","organisation_id":"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c","type":"accounts","version":0}`

	var data map[string]interface{}
	json.Unmarshal([]byte(str), &data)

	return data
}
