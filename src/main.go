package main

type AccountFetcher interface {
	FetchAccount(params FetchAccountParams) FetchAccountResult
}

