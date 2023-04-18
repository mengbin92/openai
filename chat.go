package openai

import (
	"context"
	"net/http"
)

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ChatCompletionRequest struct defines the openAI ChatCompletionRequest type.
// The purpose of this struct is to define the parameters required to make a request to an OpenAI API endpoint to perform chat completion.
type ChatCompletionRequset struct {
	Model            string                 `json:"model,omitempty"`
	Messages         []Message              `json:"messages,omitempty"`
	Temperature      float64                `json:"temperature,omitempty"`
	TopP             float64                `json:"top_p,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	Stop             string                 `json:"stop,omitempty"`
	MaxTokens        int                    `json:"max_tokens,omitempty"`
	PresencePenalty  float64                `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64                `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
}

type Choice struct {
	Index        int `json:"index,omitempty"`
	Message      `json:"message,omitempty"`
	FinishReason string `json:"finish_reason,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// The function CreateChatCompletion creates a new completed chat.
// It takes in parameters: 
// - ctx, of type context.Context which is the execution context of the function.
// - request, of type *ChatCompletionRequset which contains the necessary information to create a new chat completion.
// It returns:
// - response, of type ChatCompletionResponse which contains the response from the server.
// - err, of type error which contains any errors that occurred. If there were no errors, it will be nil.
func (c *Client) CreateChatCompletion(ctx context.Context, request *ChatCompletionRequset) (response ChatCompletionResponse, err error) {
	// Build the request using the request factory.
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(chatCompletion), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
