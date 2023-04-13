package client

import (
	"io"
	"net/http"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/mengbin92/openai/models"
	"github.com/pkg/errors"
)

var (
	CHAT_URL = `https://api.openai.com/v1/chat/completions`
)

type Client struct {
	apiKey   string
	org      string
	proxyUrl string
}

func NewClient(apikey, org, proxyUrl string) *Client {
	return &Client{
		apiKey:   apikey,
		org:      org,
		proxyUrl: proxyUrl,
	}
}

func (c *Client) Do(method string, body io.Reader) (*models.Response, error) {
	req, err := http.NewRequest(method, CHAT_URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	if c.org != "" {
		req.Header.Add("OpenAI-Organization", c.org)
	}

	var client *http.Client
	if len(c.proxyUrl) != 0 {
		proxyUrl, err := url.Parse(c.proxyUrl)
		if err != nil {
			panic(err)
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		client = &http.Client{}
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "get response from openai")
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body error")
	}
	resp := &models.Response{}
	err = sonic.Unmarshal(data, resp)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal openai response error")
	}
	return resp, nil

}
