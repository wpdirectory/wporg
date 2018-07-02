package wporg

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

// Client contains data required for making requests to the API
type Client struct {
	userAgent  string
	httpClient *http.Client
}

// NewClient returns a new client for accessing the WordPress.org APIs
func NewClient(options ...func(c *Client)) *Client {
	c := &Client{}
	for _, option := range options {
		option(c)
	}

	// Set default user-agent if not set
	if c.userAgent == "" {
		c.userAgent = "wporg/1.0"
	}

	// Set default client if not set
	if c.httpClient == nil {
		c.httpClient = getDefaultClient()
	}

	return c
}

func getDefaultClient() *http.Client {
	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

	httpClient := &http.Client{
		Timeout:   time.Second * time.Duration(30),
		Transport: netTransport,
	}

	return httpClient
}
