package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

func initOpenAIClient() {
	apikey := viper.GetString("openai.apikey")
	org := viper.GetString("openai.org")
	proxyURL := viper.GetString("openai.proxy")

	defaultConfig := openai.DefaultConfig(apikey)
	defaultConfig.OrgID = org

	var httpClient *http.Client

	if len(proxyURL) != 0 {
		proxyUrl, err := url.Parse(proxyURL)
		if err != nil {
			panic(err)
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		httpClient = &http.Client{}
	}
	defaultConfig.HTTPClient = httpClient
	client = openai.NewClientWithConfig(defaultConfig)
}

type Server struct {
	sv *http.Server
}

func NewServer() *Server {
	initOpenAIClient()
	return &Server{}
}

func (s *Server) Run(port string) error {
	engine := gin.Default()

	engine.GET("/chat", chat)
	engine.GET("/wxChat", wxChat)
	engine.GET("/models", listModels)
	engine.GET("/completion", completion)

	s.sv = &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	return s.sv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.sv.Shutdown(ctx)
}
