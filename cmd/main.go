package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/openai"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	weChatInfo *openai.WeChatInfo
	log        *zap.SugaredLogger
	client     *openai.Client
	cache      *redis.Client
)

func main() {
	viper.SetConfigFile("./conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("load config error: %s", err.Error()))
	}

	// init log
	log = defaultLogger().Sugar()

	// init redis
	if err = initRedis(); err != nil {
		log.Panicf("init redis error: %s", err.Error())
	}
	cache = getRedisClient()

	// init weChat info
	initWeChatInfo()

	// init openAI handler
	initOpenAIClient()

	engine := gin.Default()

	engine.GET("ai/chat", chat)
	engine.GET("ai/wx", weChatVerify)
	engine.POST("ai/wx", weChat)

	sv := &http.Server{
		Addr:    ":" + viper.GetString("port"),
		Handler: engine,
	}

	go func() {
		// 服务连接
		if err := sv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")
}

func initOpenAIClient() {
	apikey := viper.GetString("openai.apikey")
	org := viper.GetString("openai.org")
	proxyURL := viper.GetString("openai.proxy")

	client = openai.NewClient(apikey, org, proxyURL)
}

func chat(ctx *gin.Context) {

	request := &openai.ChatRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Errorf("Binding Lifecycle struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &openai.ChatCompletionRequset{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.Message{
			{Role: openai.ChatMessageRoleUser, Content: request.Content},
		},
		Temperature: 1,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}
