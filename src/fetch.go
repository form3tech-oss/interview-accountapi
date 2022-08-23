type FetchAccountParams struct {
	ID             string `json:"id,omitempty"`
	OrganisationID string `json:"organizationId,omitempty"`
	Type           string `json:"type,omitempty"`
	Version        *int64 `json:"version,omitempty"`
}

type FetchAccountResult struct {
	AccountData
}
// Avoided the pointer as to make it clear that the private fields
// are meant to remain unchanged
func (client Client) FetchAccount(params FetchAccountParams) (*FetchAccountResult, error) {
	return &FetchAccountResult{
		AccountData: AccountData{},
	}, nil
}
