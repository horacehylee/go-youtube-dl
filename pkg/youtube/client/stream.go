package client

import (
	"fmt"
	"io"
	"net/http"
)

// StreamLength of the url provided
func (c *Client) StreamLength(url string) (int64, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get content length, status code returned: %v", resp.StatusCode)
	}
	return resp.ContentLength, nil
}

// Stream used for getting chunk of stream data
func (c *Client) Stream(url string, from int64, to int64) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// TODO: can split length into multiple goroutine to speed up http calls, each with around 2MB
	req.Header.Set("range", fmt.Sprintf("bytes=%v-%v", from, to))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusPartialContent {
		return nil, fmt.Errorf("failed to get stream, status code returned: %v", resp.StatusCode)
	}
	return resp.Body, nil
}
