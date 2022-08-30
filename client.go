package main

import (
	"net/http"
	"time"
)

const BASE_URL_V1 = "http://localhost:8080/v1/organisation"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: BASE_URL_V1,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
