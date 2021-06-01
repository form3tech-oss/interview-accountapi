package account

import "strings"

type Validator interface {
	Validate(account Account) (result ValidationResult)
}

type ValidationResult struct {
	errorMessages []string
}

func (r *ValidationResult) fail(error string) {
	r.errorMessages = append(r.errorMessages, error)
}

func (r ValidationResult) IsValid() (valid bool) {
	return len(r.errorMessages) == 0
}

func (r *ValidationResult) concat(other ValidationResult) (output ValidationResult) {
	r.errorMessages = append(r.errorMessages, other.errorMessages...)
	return *r
}

func (r ValidationResult) message() (message string) {
	return strings.Join(r.errorMessages[:], ", ")
}
