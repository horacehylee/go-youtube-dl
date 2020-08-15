package client

import "net/http"

// Client for Youtube APIs
type Client struct {
	client *http.Client
}

// New returns HTTP client for Youtube APIs
func New() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}
