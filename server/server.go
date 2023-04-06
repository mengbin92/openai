package server

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type Server struct {
	client *openai.Client
	sv     *http.Server
}

func NewServer() *Server {
	token := os.Getenv("APIKEY")
	org := os.Getenv("ORG")
	proxyURL := os.Getenv("PROXY")

	defaultConfig := openai.DefaultConfig(token)
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

	return &Server{
		client: openai.NewClientWithConfig(defaultConfig),
	}
}

func (s *Server) Run(port string) error {
	engine := gin.Default()

	engine.GET("/chat", s.chat)
	engine.GET("/wxChat", s.wxChat)

	s.sv = &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	return s.sv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.sv.Shutdown(ctx)
}
