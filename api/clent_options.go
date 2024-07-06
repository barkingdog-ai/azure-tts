package api

import (
	"net/http"
	"time"
)

type ClientOption func(*AzureTTSClient) error

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *AzureTTSClient) error {
		c.HTTPClient = httpClient
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client. The default is 30 seconds. If you are overriding the http client as well, just include
// the timeout there.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *AzureTTSClient) error {
		c.HTTPClient.Timeout = timeout
		return nil
	}
}
