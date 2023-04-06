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
	"github.com/sashabaranov/go-openai"
)

func chat(ctx *gin.Context) {
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
	resp, err := goChat(chat.Content, tokens)
	if err != nil {
		logger.Errorf("get chat response from openai error: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func wxChat(ctx *gin.Context) {
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
		chatResp, err := goChat(reqBody.Content, 0)
		if err != nil {
			logger.Errorf("call chatGPT got error: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("call chatGPT got error: %s", err.Error())})
			return
		}
		resp.Content = chatResp.Choices[0].Message.Content
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

func listModels(ctx *gin.Context) {
	logger.Info("List and describe the various models available in the API")
	resp, err := client.ListModels(ctx)
	if err != nil {
		logger.Errorf("call chatGPT got error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("call ListModels got error: %s", err.Error())})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func completion(ctx *gin.Context) {
	chat := &models.ChatRequest{}
	if err := ctx.Bind(chat); err != nil {
		logger.Errorf("Binding Lifecycle struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// tokens := 5
	// if chat.Tokens != 0 {
	// 	tokens = chat.Tokens
	// }
	resp, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:  openai.GPT3TextDavinci003,
			Prompt: chat.Content,
		},
	)
	if err != nil {
		logger.Errorf("get completion response from openai error: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}
