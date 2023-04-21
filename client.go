// This is an openAI library implemented in Go.
package openai

import (
	"io"
	"net/http"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

// Clinet is the openAI API client
type Client struct {
	apiKey   string
	org      string
	proxyUrl string

	HttpClient     *http.Client
	requestFactory RequestFactory
	formFactory    func(body io.Writer) FormFactory
}

// NewClient creates and returns a new instance of the Client struct.
// It takes in an API key, organization name and proxy URL as strings
func NewClient(apikey, org, proxyUrl string) *Client {
	// Create a new instance of the Client struct with the given parameters
	client := &Client{
		apiKey:         apikey,
		org:            org,
		proxyUrl:       proxyUrl,
		requestFactory: newDefaultRequestFcatory(),
		formFactory: func(body io.Writer) FormFactory {
			return newDefaultForm(body)
		},
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

// sendRequest sends an HTTP request to a given URL and handles the response
// It takes 2 parameters: req *http.Request - a pointer to an http.Request object and v interface{} - an Response object such as ChatCompletionResponse CompletionResponse ...
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	// Set the headers for the request being sent
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	if c.org != "" {
		req.Header.Add("OpenAI-Organization", c.org)
	}

	// Send the HTTP request and handle errors
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send ChatCompletion request error")
	}
	defer res.Body.Close()

	// If HTTP status code is not in success range, parse the error response and return the error message
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

	// Handle the response from openAI
	if v != nil {
		if err := sonic.ConfigDefault.NewDecoder(res.Body).Decode(v); err != nil {
			return errors.Wrap(err, "unmarshal ChatCompletionResponse error")
		}
	}
	return nil
}