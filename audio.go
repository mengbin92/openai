package openai

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type TranscriptionsRequest struct {
	FilePath       string  `json:"file_path"`
	Model          string  `json:"model"` // only whisper-1
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"` // json, text, srt, verbose_json, or vtt
	Temperature    float64 `json:"temperature,omitempty"`
	Language       string  `json:"language,omitempty"`
}

type TranscriptionsResponse struct {
	Text string `json:"text"`
}

func (c *Client) CreateTranscriptions(ctx context.Context, request *TranscriptionsRequest) (response TranscriptionsResponse, err error) {
	response, err = c.callAudio(ctx, request, audioTranscriptions)
	return
}

func (c *Client) CreateTranslations(ctx context.Context, request *TranscriptionsRequest) (response TranscriptionsResponse, err error) {
	response, err = c.callAudio(ctx, request, audioTranslations)
	return
}

func (c *Client) callAudio(ctx context.Context, request *TranscriptionsRequest, url string) (response TranscriptionsResponse, err error) {
	buf := &bytes.Buffer{}
	factory := c.formFactory(buf)

	// read audio file
	err = factory.CreateFormFile("file", request.FilePath)
	if err != nil {
		errors.Wrap(err, "load audio file error")
		return
	}

	if err = factory.WriteField("model", request.Model); err != nil {
		errors.Wrap(err, "write model error")
		return
	}

	if request.Language != "" {
		if err = factory.WriteField("language", request.Language); err != nil {
			errors.Wrap(err, "write language error")
			return
		}
	}
	if request.Prompt != "" {
		if err = factory.WriteField("prompt", request.Prompt); err != nil {
			errors.Wrap(err, "write prompt error")
			return
		}
	}
	if request.ResponseFormat != "" {
		if err = factory.WriteField("response_format", request.ResponseFormat); err != nil {
			errors.Wrap(err, "write response_format error")
			return
		}
	}
	if request.Temperature != 0 {
		if err = factory.WriteField("temperature", fmt.Sprintf("%f", request.Temperature)); err != nil {
			errors.Wrap(err, "write temperature error")
			return
		}
	}

	if err = factory.Close(); err != nil {
		errors.Wrap(err, "close multipart error")
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL(url), buf)
	if err != nil {
		errors.Wrap(err, "NewRequestWithContext error")
		return
	}
	req.Header.Add("Content-Type", factory.FormDataContentType())
	err = c.sendRequest(req, &response)
	return
}
