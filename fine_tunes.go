package openai

import (
	"context"
	"fmt"
	"net/http"
)

type FineTuneRequest struct {
	TrainingFile                 string    `json:"training_file"`
	ValidationFile               string    `json:"validation_file"`
	Model                        string    `json:"model"`
	NEpochs                      int       `json:"n_epochs"`
	BatchSize                    int       `json:"batch_size"`
	LearningRateMultiplier       float32   `json:"learning_rate_multiplier"`
	PromptLossWeight             float32   `json:"prompt_loss_weight"`
	ComputeClassificationMetrics bool      `json:"compute_classification_metrics"`
	ClassificationNClasses       int       `json:"classification_n_classes"`
	ClassificationPositiveClass  string    `json:"classification_positive_class"`
	ClassificationBetas          []float32 `json:"classification_betas"`
	Suffix                       string    `json:"suffix"`
}

type FineTune struct {
	ID              string              `json:"id"`
	Object          string              `json:"object"`
	Model           string              `json:"model"`
	CreatedAt       int64               `json:"created_at"`
	Events          []FineTuneEvent     `json:"events"`
	FineTunedModel  string              `json:"fine_tuned_model"`
	Hyperparams     FineTuneHyperparams `json:"hyperparams"`
	OrganizationID  string              `json:"organization_id"`
	ResultFiles     []File              `json:"result_files"`
	Status          string              `json:"status"`
	ValidationFiles []File              `json:"validation_files"`
	TrainingFiles   []File              `json:"training_files"`
	UpdatedAt       int64               `json:"updated_at"`
}

type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type FineTuneHyperparams struct {
	BatchSize              int64   `json:"batch_size"`
	LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
	NEpochs                int64   `json:"n_epochs"`
	PromptLossWeight       float64 `json:"prompt_loss_weight"`
}

type FineTuneList struct {
	Object string     `json:"object"`
	Data   []FineTune `json:"data"`
}

type FineTuneEventList struct {
	Object string          `json:"object"`
	Data   []FineTuneEvent `json:"data"`
}

type FineTuneDeleteResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func (c *Client) CreateFineTune(ctx context.Context, request FineTuneRequest) (response FineTune, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL(fineTunes), request)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) ListFineTune(ctx context.Context) (response FineTuneList, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fullURL(fineTunes), nil)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) RetrieveFineTune(ctx context.Context, id string) (response FineTune, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fmt.Sprintf("%s/%s", fullURL(fineTunes), id), nil)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) CancelFineTune(ctx context.Context, id string) (response FineTune, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fmt.Sprintf("%s/%s/cancel", fullURL(fineTunes), id), nil)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) ListFineTuneEvents(ctx context.Context, id string) (response FineTuneEventList, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fmt.Sprintf("%s/%s/events", fullURL(fineTunes), id), nil)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) DeleteFineTuneModel(ctx context.Context, model string) (response FineTuneDeleteResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", fullURL(models), model), nil)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
