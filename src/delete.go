package main

type DeleteParams struct{}

type DeleteResult struct{}

func (client Client) Delete(CreateParams) *DeleteResult {
	return &DeleteResult{}
}
