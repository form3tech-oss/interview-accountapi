package form3client

// type form3Response struct {
// 	Data          *AccountData `json:"data,omitempty"`
// 	Relationships interface{}  `json:"relationships,omitempty"`
// }

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
	CreatedOn      string             `json:"created_on,omitempty"`
	ModifiedOn     string             `json:"modified_on,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

const AccountsEndpoint string = "/v1/organisation/accounts"

// Get FullAccountUrl
func getFullAccountUrl() string {
	return GetBaseUrl() + AccountsEndpoint + "/"
}

// Create creates a new account
func CreateAccount(jsonAccountData string) (ClientResponse, error) {
	return PostRequest(getFullAccountUrl(), jsonAccountData)
}

// Delete an existing account
func DeleteAccount(accountID, version string) (ClientResponse, error) {
	return DeletionRequest(getFullAccountUrl() + accountID + "?version=" + version)
}

// Fetch an existing account
func FetchAccount(accountID string) (ClientResponse, error) {
	return FetchRequest(getFullAccountUrl() + accountID)
}
