package hltv

import (
	"net/http"
)

const Version = "0.0.1"

type Client struct {
	http    *http.Client
	baseURL string
}

func New(httpClient *http.Client) *Client {
	c := &Client{
		http:    httpClient,
		baseURL: "https://www.hltv.org/",
	}

	return c
}
