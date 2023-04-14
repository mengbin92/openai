package openai

import (
	"context"
	"fmt"
	"net/http"
)

type Model struct {
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	Created    int64        `json:"created"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
	Root       string       `json:"root"`
	Parent     interface{}  `json:"parent"`
}

type Permission struct {
	ID                 string      `json:"id"`
	Object             string      `json:"object"`
	Created            int64       `json:"created"`
	AllowCreateEngine  bool        `json:"allow_create_engine"`
	AllowSampling      bool        `json:"allow_sampling"`
	AllowLogprobs      bool        `json:"allow_logprobs"`
	AllowSearchIndices bool        `json:"allow_search_indices"`
	AllowView          bool        `json:"allow_view"`
	AllowFineTuning    bool        `json:"allow_fine_tuning"`
	Organization       string      `json:"organization"`
	Group              interface{} `json:"group"`
	IsBlocking         bool        `json:"is_blocking"`
}

type ModelList struct {
	Models []Model `json:"data"`
}

func (c *Client) ModelList(ctx context.Context) (list ModelList, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fullURL(models), nil)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &list)
	return
}

func (c *Client) ModelInfo(ctx context.Context, name string) (model Model, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fmt.Sprintf("%s/%s", fullURL(models), name), nil)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &model)
	return
}
