package openai

import (
	"context"
	"net/http"
)

type CompletionRequest struct {
	Model            string         `json:"model"`
	Prompt           interface{}    `json:"prompt,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type CompletionChoice struct {
	Text         string        `json:"text"`
	Index        int           `json:"index"`
	FinishReason string        `json:"finish_reason"`
	LogProbs     LogprobResult `json:"logprobs"`
}

type LogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// The CreateCompletion method creates a new completion request with the specified parameters and sends it to the server.
// It takes in parameters:
// - ctx, of type context.Context which is the execution context of the function.
// - request, of type *CompletionRequest which contains the data to be sent in the completion request body.
// It returns:
// - response, an object of type CompletionResponse, which contains the response from the server.
// - err, an error object which will contain any errors that occurred during the creation or sending of the request.
func (c *Client) CreateCompletion(ctx context.Context, request *CompletionRequest) (response CompletionResponse, err error) {
	//build the completion request using the requestFactory instance.
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(completions), request)
	if err != nil {
		return
	}
	//send the completion request to the server using sendRequest method 
	err = c.sendRequest(req,  &response)
	return
}

