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

// The CreateEdits method creates a new edits request with the specified parameters and sends it to the server.
// It takes in parameters:
// - ctx, of type context.Context which is the execution context of the function.
// - request, of type *EditsRequest which contains the data to be sent in the edits request body.
// It returns:
// - response, an object of type EditsResponse, which contains the response from the server.
// - err, an error object which will contain any errors that occurred during the creation or sending of the request.
func (c *Client) CreateEdits(ctx context.Context, request *EditsRequest) (response EditsResponse, err error) {
	//build the edits request using the requestFactory instance.
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(edits), request)
	if err != nil {
		return
	}
	//send the edits request to the server using sendRequest method
	err = c.sendRequest(req,  &response)
	return
}
