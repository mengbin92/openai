package openai

import (
	"context"
	"net/http"
)

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

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

func (c *Client) CreateChatCompletion(ctx context.Context, request *ChatCompletionRequset) (response ChatCompletionResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(chatCompletion), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
