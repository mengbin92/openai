package openai

import (
	"context"
	"net/http"
)

type EmbeddingsRequest struct {
	Model string      `json:"model"`
	Input interface{} `json:"input"`
	User  string      `json:"user,omitempty"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingsResponse struct {
	Model  string      `json:"model"`
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Usage  Usage       `json:"usage"`
}

func (c *Client) CreateEmbeddings(ctx context.Context, request *EmbeddingsRequest) (response EmbeddingsResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(embeddings), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
