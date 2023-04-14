package openai

import (
	"net/http"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type Client struct {
	apiKey   string
	org      string
	proxyUrl string

	HttpClient     *http.Client
	requestFactory RequestFactory
}

func NewClient(apikey, org, proxyUrl string) *Client {
	client := &Client{
		apiKey:         apikey,
		org:            org,
		proxyUrl:       proxyUrl,
		requestFactory: &httpRequestFactory{},
	}
	if len(client.proxyUrl) != 0 {
		proxyUrl, err := url.Parse(client.proxyUrl)
		if err != nil {
			panic(err)
		}
		client.HttpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		client.HttpClient = &http.Client{}
	}
	return client
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	if c.org != "" {
		req.Header.Add("OpenAI-Organization", c.org)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send ChatCompletion request error")
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrResponse
		err = sonic.ConfigDefault.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			reqErr := &RequestError{
				Code: res.StatusCode,
				Err:  err,
			}
			return errors.Wrap(reqErr, "request sent to openai error")
		}
		errRes.Err.StatusCode = res.StatusCode
		return errors.Wrap(&errRes, "error response from openai")
	}

	if v != nil {
		if err := sonic.ConfigDefault.NewDecoder(res.Body).Decode(v); err != nil {
			return errors.Wrap(err, "unmarshal ChatCompletionResponse error")
		}
	}
	return nil
}
