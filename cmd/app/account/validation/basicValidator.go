package validation

import "github.com/advena/interview-accountapi/cmd/app/account"

type bankIDValidator struct {
}

func (c bankIDValidator) Validate(account account.Account) (result ValidationResult) {
	validationResult := ValidationResult{}
	country := account.Attributes.Country

	validateEmptyBankID(country, account.Attributes.BankID, &validationResult)
	validateNonEmptyBankID(country, account.Attributes.BankID, &validationResult)
	validateBankIDLength(country, account.Attributes.BankID, &validationResult)

	return validationResult

}

func validateEmptyBankID(country string, bankID string, result *ValidationResult) {
	if country == "NL" && bankID != "" {
		result.fail("BankId for NL must be empty")
	}
}

func validateNonEmptyBankID(country string, bankID string, result *ValidationResult) {
	switch country {
	case
		"GB",
		"BE":
		if bankID == "" {
			result.fail("BankID for " + country + " must not be empty")
		}
	}
}

func validateBankIDLength(country string, bankID string, result *ValidationResult) {
	switch country {
	case "GB":
		if len(bankID) != 6 {
			result.fail("BankID for GB must be 6 characters length")
		}
	case "AU":
		if bankID != "" && len(bankID) != 6 {
			result.fail("BankID for AU must be 6 characters length")
		}
	case "BE":
		if len(bankID) != 3 {
			result.fail("BankID for AU must be 3 characters length")
		}
	}
}

func BankIDValidator() (validator Validator) {
	return bankIDValidator{}
}

type bicValidator struct {
}

func (b bicValidator) Validate(account account.Account) ValidationResult {
	validationResult := ValidationResult{}
	country := account.Attributes.Country

	validateBIC(country, account.Attributes.BIC, &validationResult)

	return validationResult
}

func validateBIC(country string, BIC string, result *ValidationResult) {
	switch country {
	case
		"GB",
		"AU",
		"BE":
		if BIC == "" {
			result.fail("BIC for " + country + " cannot be empty")
		}
	}
}

func BICValidator() (validator Validator) {
	return bicValidator{}
}

type AccountValidator struct {
	Validators []Validator
}

func (a AccountValidator) Validate(account account.Account) ValidationResult {
	result := ValidationResult{}
	for _, validator := range a.Validators {
		result.concat(validator.Validate(account))
	}
	return result
}
