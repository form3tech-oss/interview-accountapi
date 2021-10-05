package endpoint

import "strconv"

var PORT int = 8080

var BASE_URL = "http://localhost:" + strconv.Itoa(PORT)

var HEALTH_URL string = BASE_URL + "/v1/health"
var ACCOUNT_URL = BASE_URL + "/v1/organisation/accounts"
var APP_JSON = BASE_URL + "application/json"
var DELETE_URL = ACCOUNT_URL + "/"
