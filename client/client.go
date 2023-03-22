package client

import (
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	CHAT_URL = `https://api.openai.com/v1/chat/completions`
)

type Client struct {
	apiKey       string
	Organization string
}

func NewClient(apikey, org string) *Client {
	return &Client{
		apiKey:       apikey,
		Organization: org,
	}
}

func (c *Client) Do(method string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, CHAT_URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	if c.Organization != "" {
		req.Header.Add("OpenAI-Organization", c.Organization)
	}

	proxy := os.Getenv("http_proxy")

	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	return client.Do(req)
}
