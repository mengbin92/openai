package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/openai"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("load config error: %s", err.Error()))
	}

	engine := gin.Default()

	engine.GET("ai/chat", chat)

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
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func chat(ctx *gin.Context) {
	apikey := viper.GetString("openai.apikey")
	org := viper.GetString("openai.org")
	proxyURL := viper.GetString("openai.proxy")

	client := openai.NewClient(apikey, org, proxyURL)

	request := &openai.ChatRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Printf("Binding Lifecycle struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &openai.ChatCompletionRequset{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.Message{
			{Role: "user", Content: request.Content},
		},
		Temperature: 0.2,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}
