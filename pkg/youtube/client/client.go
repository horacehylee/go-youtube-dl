package client

import "net/http"

// Client for Youtube APIs
type Client struct {
	client *http.Client
}

// NewClient for Youtube APIs
func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}
