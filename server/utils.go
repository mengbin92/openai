package server

import (
	"context"

	"github.com/mengbin92/openai/models"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger
	weChatInfo *models.WeChatInfo
	client     *openai.Client
)

func goChat(msg string, tokens int) (openai.ChatCompletionResponse, error) {
	return client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
}
