package main

type CreateParams struct{}

type CreateResult struct{}

func (client Client) Create(params *CreateParams) *CreateResult {
	return &CreateResult{}
}
