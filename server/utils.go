package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/mengbin92/openai/client"
	"github.com/mengbin92/openai/common/cache"
	"github.com/mengbin92/openai/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger
	weChatInfo *models.WeChatInfo
	handler    *client.Client
)

func goChat(msg string) (*models.Response, error) {
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

func goChatWithChan(reqCache *models.WeChatCache, respChan chan string, errChan chan error) {
	req := &models.Requset{
		Model: "gpt-3.5-turbo",
		Messages: []models.Message{
			{Role: "user", Content: reqCache.Content},
		},
		Temperature: 0.2,
	}
	reqByte, err := sonic.Marshal(req)
	if err != nil {
		errChan <- errors.Wrap(err, "marshal chat request error")
		return
	}

	resp, err := handler.Do(http.MethodPost, strings.NewReader(string(reqByte)))
	if err != nil {
		errChan <- errors.Wrap(err, "get chat response from openai error")
		return
	}
	go func() {
	LOOP:
		err = cache.Get().Set(context.Background(), reqCache.Key(), []byte(resp.Choices[0].Message.Content), viper.GetDuration("redis.expire")*time.Second).Err()
		if err != nil {
			logger.Debugf("Set data error: %s", err.Error())
			goto LOOP
		}
	}()

	respChan <- resp.Choices[0].Message.Content
}
