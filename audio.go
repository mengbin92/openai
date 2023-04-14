package openai

import (
	"context"
	"net/http"
	"os"
)

type TranscriptionsRequest struct {
	File           *os.File `json:"file"`
	Model          string   `json:"model"` // only whisper-1
	Prompt         string   `json:"prompt,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` // json, text, srt, verbose_json, or vtt
	Temperature    float64  `json:"temperature,omitempty"`
	Language       string   `json:"language,omitempty"`
}

type TranscriptionsResponse struct {
	Text string `json:"text"`
}

func (c *Client) CreateTranscriptions(ctx context.Context, request *TranscriptionsRequest) (response TranscriptionsResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(audioTranscriptions), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}

func (c *Client) CreateTranslations(ctx context.Context, request *TranscriptionsRequest) (response TranscriptionsResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(audioTranslations), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
