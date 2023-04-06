package models

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"sort"
	"strings"

	"github.com/mengbin92/openai/log"
)

var (
	logger = log.DefaultLogger().Sugar()
)

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Requset struct {
	Model            string                 `json:"model,omitempty"`
	Messages         []Message              `json:"messages,omitempty"`
	Temperature      float64                `json:"temperature,omitempty"`
	TopP             float64                `json:"top_p,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	Stop             string                 `json:"stop,omitempty"`
	MaxTokens        int                    `json:"max_tokens,omitempty"`
	PresencePenalty  float64                `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64                `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
}

type Choice struct {
	Index        int `json:"index,omitempty"`
	Message      `json:"message,omitempty"`
	FinishReason string `json:"finish_reason,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Response struct {
	Id      string   `json:"id,omitempty"`
	Object  string   `json:"object,omitempty"`
	Created int64    `json:"created,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
	Usage   `json:"usage,omitempty"`
}

type ChatRequest struct {
	Content string `json:"content" form:"content"`
	Tokens  int    `json:"tokens" form:"tokens"`
}

type WeChatVerify struct {
	Signature string `json:"signature" form:"signature"`
	Timestamp string `json:"timestamp" form:"timestamp"`
	Nonce     string `json:"nonce" form:"nonce"`
	Echostr   string `json:"echostr" form:"echostr"`
}

type WeChatMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
}

func (p *WeChatVerify) Verify(token string) bool {
	s := []string{token, p.Timestamp, p.Nonce}
	sort.Strings(s)
	str := strings.Join(s, "")
	hashs := sha1.New()
	hashs.Write([]byte(str))

	signature := hex.EncodeToString(hashs.Sum(nil))
	logger.Infof("calc signature on local: %s", signature)
	if signature == p.Signature {
		return true
	} else {
		return false
	}
}
