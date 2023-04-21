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
	engine.POST("ai/audio", audio)
	engine.GET("ai/edits", edits)
	engine.GET("ai/images", imageGen)
	engine.GET("ai/embedding", embedding)
	engine.POST("ai/file", upload)
	engine.DELETE("ai/file/:id", deleteFile)
	engine.GET("ai/file/:id", retrieveFile)
	engine.GET("ai/file/content/:id", retrieveFileContent)

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
