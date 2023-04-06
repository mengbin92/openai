package server

import (
	"context"
	"os"

	"github.com/mengbin92/openai/log"
	openai "github.com/sashabaranov/go-openai"
)

var (
	logger = log.DefaultLogger().Sugar()
	token  = os.Getenv("TOKEN")
	client *openai.Client
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
