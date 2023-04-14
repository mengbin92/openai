package openai

import (
	"context"
	"net/http"
)

type EditsRequest struct {
	Model       string  `json:"model,omitempty"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction,omitempty"`
	N           int     `json:"n,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
}

type EditsChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

type EditsResponse struct {
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Usage   Usage         `json:"usage"`
	Choices []EditsChoice `json:"choices"`
}

func (c *Client) CreateEdits(ctx context.Context, request *EditsRequest) (response EditsResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(edits), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
