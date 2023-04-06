package server

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
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

func (s *Server) wxChat(ctx *gin.Context) {
	logger.Info("Get Msg from wechat")
	verify := &models.WeChatVerify{
		Signature: ctx.Query("signature"),
		Timestamp: ctx.Query("timestamp"),
		Nonce:     ctx.Query("nonce"),
		Echostr:   ctx.Query("echostr"),
	}
	if !verify.Verify(token) {
		logger.Error("WeChat Verify failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "WeChat Verify failed"})
		return
	}
	logger.Info("verify pass")

	reqBody := &models.WeChatMsg{}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Errorf("read request body error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "read request body error"})
		return
	}
	xml.Unmarshal(body, reqBody)
	reqBytes, _ := sonic.Marshal(reqBody)
	logger.Infof("Get requset from wechat: %s", string(reqBytes))

	switch reqBody.MsgType {
	case "text":
		resp := &models.WeChatMsg{}
		resp.FromUserName = reqBody.ToUserName
		resp.ToUserName = reqBody.FromUserName
		resp.CreateTime = time.Now().Unix()
		resp.MsgType = "text"
		resp.Content = reqBody.Content
		respBytes, _ := xml.Marshal(resp)
		logger.Infof("return msg to wechat: %s", string(respBytes))
		ctx.Writer.Header().Set("Content-Type", "text/xml")
		ctx.Writer.WriteString(string(respBytes))
	default:
		logger.Errorf("unknow MsgType: %s", reqBody.MsgType)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unknow MsgType: %s", reqBody.MsgType)})
		return
	}
}
