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

//The following code declares a function CreateEmbeddings for the client, which accepts a context, an EmbeddingsRequest struct and returns an EmbeddingsResponse and error.
//The function sends a POST request to the embeddings API of the server using Build method from c.requestFactory. It builds a request using the provided arguments;  a context, http.MethodPost and full URL for the embeddings API.
//It also checks for any errors that may occur in the process of building the request.
//If there is no error in creating the request, then this function sends the created request using sendRequest method of client c with both the request and response structs as its parameters.
//If there is no error encountered during the execution of the above statements, the function returns the result in the form of EmbdedingsResponse object.
func (c *Client) CreateEmbeddings(ctx context.Context, request *EmbeddingsRequest) (response EmbeddingsResponse, err error) {

	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(embeddings), request)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)

	return
}

