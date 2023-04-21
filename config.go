package openai

import (
	"io"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	APIKey string
	Org    string
	Proxy  string

	HttpRequestFactory  RequestFactory
	HttpFielFormFactory func(body io.Writer) FormFactory
}

func NewClientWithConfig(config *ClientConfig) *Client {
	// Create a new instance of the Client struct with the given parameters
	client := &Client{
		apiKey:         config.APIKey,
		org:            config.Org,
		proxyUrl:       config.Proxy,
		requestFactory: config.HttpRequestFactory,
		formFactory:    config.HttpFielFormFactory,
	}

	// If a proxy URL was provided, set up an HTTP client with the proxy details
	if len(client.proxyUrl) != 0 {
		// Parse the proxy URL into a url.URL struct
		proxyUrl, err := url.Parse(client.proxyUrl)
		if err != nil {
			// If there was an error parsing the URL, panic
			panic(err)
		}
		// Set up the HTTP client with the proxy details
		client.HttpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		// Otherwise, set up a regular HTTP client
		client.HttpClient = &http.Client{}
	}
	// Return the newly created client instance
	return client
}
