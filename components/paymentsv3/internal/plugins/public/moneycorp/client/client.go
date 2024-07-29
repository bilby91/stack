package client

import (
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
	pageSize   int
}

func newHTTPClient(clientID, apiKey, endpoint string) *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &apiTransport{
			clientID:   clientID,
			apiKey:     apiKey,
			endpoint:   endpoint,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func NewClient(clientID, apiKey, endpoint string, pageSize int) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	c := &Client{
		httpClient: newHTTPClient(clientID, apiKey, endpoint),
		endpoint:   endpoint,
		pageSize:   pageSize,
	}

	return c, nil
}

func (c Client) PageSize() int {
	return c.pageSize
}
