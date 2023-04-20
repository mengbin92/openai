package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/mengbin92/openai"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func initWeChatInfo() {
	weChatInfo = &openai.WeChatInfo{
		AppID:     viper.GetString("wechat.appID"),
		Appsecret: viper.GetString("wechat.appsecret"),
		Token:     viper.GetString("wechat.token"),
	}
}

func weChatVerify(ctx *gin.Context) {
	log.Info("Get Msg from wechat")
	verify := &openai.WeChatVerify{
		Signature: ctx.Query("signature"),
		Timestamp: ctx.Query("timestamp"),
		Nonce:     ctx.Query("nonce"),
		Echostr:   ctx.Query("echostr"),
	}
	if !verify.Verify(weChatInfo.Token) {
		log.Error("WeChat Verify failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "WeChat Verify failed"})
		return
	}
	ctx.Writer.WriteString(verify.Echostr)
}

func weChat(ctx *gin.Context) {
	log.Info("Get Msg from wechat")
	verify := &openai.WeChatVerify{
		Signature: ctx.Query("signature"),
		Timestamp: ctx.Query("timestamp"),
		Nonce:     ctx.Query("nonce"),
		Echostr:   ctx.Query("echostr"),
	}
	if !verify.Verify(weChatInfo.Token) {
		log.Error("WeChat Verify failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "WeChat Verify failed"})
		return
	}
	log.Info("verify pass")

	reqBody := &openai.WeChatMsg{}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Errorf("read request body error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "read request body error"})
		return
	}
	xml.Unmarshal(body, reqBody)
	reqBytes, _ := sonic.Marshal(reqBody)
	log.Infof("Get requset from wechat: %s", string(reqBytes))

	reqCache := &openai.WeChatCache{
		OpenID:  reqBody.FromUserName,
		Content: reqBody.Content,
	}

	resp := &openai.WeChatMsg{}
	resp.FromUserName = reqBody.ToUserName
	resp.ToUserName = reqBody.FromUserName
	resp.CreateTime = time.Now().Unix()
	resp.MsgType = reqBody.MsgType

	respChan := make(chan string)
	errChan := make(chan error)

	switch reqBody.MsgType {
	case "text":
		reply, err := cache.Get(context.Background(), reqCache.Key()).Bytes()
		if err != nil && len(reply) == 0 {
			log.Info("get nothing from local cache,now get data from openai")

			go goChatWithChan(reqCache, respChan, errChan)

			select {
			case resp.Content = <-respChan:
			case err := <-errChan:
				resp.Content = err.Error()
			case <-time.After(4900 * time.Millisecond):
				resp.Content = "前方网络拥堵....\n等待是为了更好的相遇，稍后请重新发送上面的问题来获取答案，感谢理解"
				// default:
				// 	resp.Content = "答案整理中，请30s稍后重试"
			}
		} else {
			resp.Content = string(reply)
		}
		respBytes, _ := xml.Marshal(resp)
		log.Infof("return msg to wechat: %s", string(respBytes))
		ctx.Writer.Header().Set("Content-Type", "text/xml")
		ctx.Writer.WriteString(string(respBytes))
	default:
		log.Errorf("unknow MsgType: %s", reqBody.MsgType)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unknow MsgType: %s", reqBody.MsgType)})
		return
	}
}

func goChatWithChan(reqCache *openai.WeChatCache, respChan chan string, errChan chan error) {
	req := &openai.ChatCompletionRequset{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.Message{
			{Role: openai.ChatMessageRoleUser, Content: reqCache.Content},
		},
		Temperature: 0.9,
	}
	log.Info(req)

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		errChan <- errors.Wrap(err, "get chat response from openai error")
		return
	}
	go func() {
	LOOP:
		err = cache.Set(context.Background(), reqCache.Key(), []byte(resp.Choices[0].Message.Content), viper.GetDuration("redis.expire")*time.Second).Err()
		if err != nil {
			log.Debugf("Set data error: %s", err.Error())
			goto LOOP
		}
	}()

	respChan <- resp.Choices[0].Message.Content
}
