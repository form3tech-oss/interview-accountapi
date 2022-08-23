package main

type AccountFetcher interface {
	FetchAccount(params FetchAccountParams) FetchAccountResult
}

type AccountCreator interface {
	CreateAccount()
}

type AccountRemover interface {
	DeleteAccount()
}

type Client struct {
	base_url string
}

type NewClientParams struct {
	BaseUrl string
}

func NewClient(params *NewClientParams) *Client {
	client := Client{
		base_url: params.BaseUrl,
	}

	return &client
}
