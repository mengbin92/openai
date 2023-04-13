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
	"github.com/mengbin92/openai/common/cache"
	"github.com/mengbin92/openai/models"
)

// func chat(ctx *gin.Context) {
// 	logger.Infof("Get msg from openai")
// 	chat := &models.ChatRequest{}
// 	if err := ctx.Bind(chat); err != nil {
// 		logger.Errorf("Binding Lifecycle struct error: %s\n", err.Error())
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	tokens := 5
// 	if chat.Tokens != 0 {
// 		tokens = chat.Tokens
// 	}
// 	resp, err := goChat(chat.Content, tokens)
// 	if err != nil {
// 		logger.Errorf("get chat response from openai error: %s", err.Error())
// 		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
// }

func weChatVerify(ctx *gin.Context) {
	logger.Info("Get Msg from wechat")
	verify := &models.WeChatVerify{
		Signature: ctx.Query("signature"),
		Timestamp: ctx.Query("timestamp"),
		Nonce:     ctx.Query("nonce"),
		Echostr:   ctx.Query("echostr"),
	}
	if !verify.Verify(weChatInfo.Token) {
		logger.Error("WeChat Verify failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "WeChat Verify failed"})
		return
	}
	ctx.Writer.WriteString(verify.Echostr)
}

func weChat(ctx *gin.Context) {
	logger.Info("Get Msg from wechat")
	verify := &models.WeChatVerify{
		Signature: ctx.Query("signature"),
		Timestamp: ctx.Query("timestamp"),
		Nonce:     ctx.Query("nonce"),
		Echostr:   ctx.Query("echostr"),
	}
	if !verify.Verify(weChatInfo.Token) {
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
		reqCache := &models.WeChatCache{
			OpenID:  reqBody.FromUserName,
			Content: reqBody.Content,
		}

		resp := &models.WeChatMsg{}
		resp.FromUserName = reqBody.ToUserName
		resp.ToUserName = reqBody.FromUserName
		resp.CreateTime = time.Now().Unix()
		resp.MsgType = "text"

		respChan := make(chan string)
		errChan := make(chan error)

		reply, err := cache.Get().Get(context.Background(), reqCache.Key()).Bytes()
		if err != nil && len(reply) == 0 {
			logger.Info("get nothing from local cache,now get data from openai")

			go goChatWithChan(reqCache, respChan, errChan)

			select {
			case resp.Content = <-respChan:
			case err := <-errChan:
				resp.Content = err.Error()
			case <-time.After(4 * time.Second):
				resp.Content = "前方网络拥堵....\n等待是为了更好的相遇，稍后请重新发送上面的问题来获取答案，感谢理解"
			default:
				resp.Content = "答案整理中，请30s稍后重试"
			}
		} else {
			resp.Content = string(reply)
		}
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

func chat(ctx *gin.Context) {
	logger.Infof("Get msg from openai")
	chat := &models.ChatRequest{}
	if err := ctx.Bind(chat); err != nil {
		logger.Errorf("Binding Lifecycle struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := goChat(chat.Content)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": response})
}
