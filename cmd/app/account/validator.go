package account

type AccountValidator interface {
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
