package openai

import (
	"context"
	"net/http"
)

type ModerationRequest struct {
	Input string `json:"input" form:"input"`
	Model string `json:"model,omitempty" form:"model,omitempty"`
}

type ModerationResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Results []Result `json:"results"`
}

type Result struct {
	Categories     Categories     `json:"categories"`
	CategoryScores CategoryScores `json:"category_scores"`
	Flagged        bool           `json:"flagged"`
}

type Categories struct {
	Hate            bool `json:"hate"`
	HateThreatening bool `json:"hate/threatening"`
	SelfHarm        bool `json:"self-harm"`
	Sexual          bool `json:"sexual"`
	SexualMinors    bool `json:"sexual/minors"`
	Violence        bool `json:"violence"`
	ViolenceGraphic bool `json:"violence/graphic"`
}

type CategoryScores struct {
	Hate            float64 `json:"hate"`
	HateThreatening float64 `json:"hate/threatening"`
	SelfHarm        float64 `json:"self-harm"`
	Sexual          float64 `json:"sexual"`
	SexualMinors    float64 `json:"sexual/minors"`
	Violence        float64 `json:"violence"`
	ViolenceGraphic float64 `json:"violence/graphic"`
}

func (c *Client) CreateModeration(ctx context.Context, request ModerationRequest) (response ModerationResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(moderations), request)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
