package server

import (
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/mengbin92/openai/client"
	"github.com/mengbin92/openai/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger
	weChatInfo *models.WeChatInfo
	handler    *client.Client
)

func goChat(msg string, tokens int) (*models.Response, error) {
	req := &models.Requset{
		Model: "gpt-3.5-turbo",
		Messages: []models.Message{
			{Role: "user", Content: msg},
		},
		Temperature: 0.2,
	}
	reqByte, err := sonic.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal chat request error")
	}
	return handler.Do(http.MethodPost, strings.NewReader(string(reqByte)))
}
