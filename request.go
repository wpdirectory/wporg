package wporg

import (
	"fmt"
	"net/http"
)

// getRequest performs the HTTP GET request, using the provided URL
func (c *Client) getRequest(URL string) (*http.Response, error) {

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	// Set User Agent
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// Check status code is 2XX
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("Invalid HTTP status code: %d", resp.StatusCode)
	}

	return resp, nil

}

// postRequest performs the HTTP POST request, using the provided URL and body
func (c *Client) postRequest(URL string, body []byte) (*http.Response, error) {

	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return nil, err
	}

	// Set User Agent
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// Check status code is 2XX
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("Invalid HTTP status code: %d", resp.StatusCode)
	}

	return resp, nil

}
