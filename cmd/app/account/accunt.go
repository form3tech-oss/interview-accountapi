package account

type Account struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
}

type Attributes struct {
	AlternativeNames []string `json:"alternative_names"`
	BankID           string   `json:"bank_id"`
	BankIDCode       string   `json:"bank_id_code"`
	BaseCurrency     string   `json:"base_currency"`
	BIC              string   `json:"bic"`
	Country          string   `json:"country"`
	Name             []string `json:"name"`
	// AccountClassification   string                `json:"account_classification"`
	JointAccount            bool   `json:"joint_account"`
	AccountMatchingOptOut   bool   `json:"account_matching_opt_out"`
	SecondaryIdentification string `json:"secondary_identification"`
	// Status                  string                `json:"status"`
	PrivateIdentification PrivateIdentification `json:"private_identification"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date"`
	BirthCountry   string   `json:"birth_country"`
	Identification string   `json:"identificaiton"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}
