package account

import "errors"

type AccountsRepository interface {
	Create(account Account) (Account, error)
	Delete(accountID string) (bool, error)
	Fetch(accountID string) (Account, error)
}

type AccountsHandler struct {
	Repository AccountsRepository
	Validator  Validator
}

func (a *AccountsHandler) Create(account Account) (Account, error) {
	validationResult := a.Validator.Validate(account)
	if validationResult.IsValid() {
		return a.Repository.Create(account)
	} else {
		return Account{}, errors.New(validationResult.message())
	}
}

func (a *AccountsHandler) Delete(accountID string) (bool, error) {
	return a.Repository.Delete(accountID)
}

func (a *AccountsHandler) Fetch(accountID string) (Account, error) {
	return a.Repository.Fetch(accountID)
}

type Account struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
}

type Attributes struct {
	AlternativeNames        []string              `json:"alternative_names"`
	BankID                  string                `json:"bank_id"`
	BankIDCode              string                `json:"bank_id_code"`
	BaseCurrency            string                `json:"base_currency"`
	BIC                     string                `json:"bic"`
	Country                 string                `json:"country"`
	Name                    []string              `json:"name"`
	JointAccount            bool                  `json:"joint_account"`
	AccountMatchingOptOut   bool                  `json:"account_matching_opt_out"`
	SecondaryIdentification string                `json:"secondary_identification"`
	PrivateIdentification   PrivateIdentification `json:"private_identification"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date"`
	BirthCountry   string   `json:"birth_country"`
	Identification string   `json:"identificaiton"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}
