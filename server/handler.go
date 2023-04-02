package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/openai/models"
	openai "github.com/sashabaranov/go-openai"
)

func (s *Server) chat(ctx *gin.Context) {
	chat := &models.ChatRequest{}
	if err := ctx.Bind(chat); err != nil {
		logger.Errorf("Binding Lifecycle struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens := 5
	if chat.Tokens != 0 {
		tokens = chat.Tokens
	}
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: chat.Content,
				},
			},
			MaxTokens: tokens,
		},
	)
	if err != nil {
		logger.Errorf("get chat response from openai error: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}
