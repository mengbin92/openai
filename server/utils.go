package server

import (
	"context"

	"github.com/mengbin92/openai/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

var (
	logger = log.DefaultLogger().Sugar()
	token  = viper.GetString("openai.token")
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
