package openai

import (
	"context"
	"net/http"
	"os"
)

type ImagesReuqest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"` // url or b64_json
	User           string `json:"user,omitempty"`
}

type ImagesResponse struct {
	Created int64                `json:"created,omitempty"`
	Data    []ImagesResponseData `json:"data,omitempty"`
}

type ImagesResponseData struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
}

func (c *Client) CreateImage(ctx context.Context, request *ImagesReuqest) (response ImagesResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL("images/generations"), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}

type ImagesEditReuqest struct {
	Image          *os.File `json:"image"`
	Mask           *os.File `json:"mask,omitempty"`
	Prompt         string   `json:"prompt"`
	N              int      `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` // url or b64_json
	User           string   `json:"user,omitempty"`
}

func (c *Client) CreateImageEdits(ctx context.Context, request *ImagesEditReuqest) (response ImagesResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL("images/edits"), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}

type ImagesVariationReuqest struct {
	Image          *os.File `json:"image"`
	N              int      `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` // url or b64_json
	User           string   `json:"user,omitempty"`
}

func (c *Client) CreateImageVariation(ctx context.Context, request *ImagesVariationReuqest) (response ImagesResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL("images/variations"), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
