package openai

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"encoding/xml"
	"sort"
	"strings"
)

type ChatRequest struct {
	Content string `json:"content" form:"content"`
	Tokens  int    `json:"tokens" form:"tokens"`
}

type WeChatInfo struct {
	Token     string
	AppID     string
	Appsecret string
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
	if signature == p.Signature {
		return true
	} else {
		return false
	}
}

type WeChatCache struct {
	OpenID  string `json:"open_id"`
	Content string `json:"content"`
}

func (cache *WeChatCache) Key() string {
	hash := sha512.New384()
	hash.Write([]byte(cache.OpenID + "-" + cache.Content))
	return hex.EncodeToString(hash.Sum(nil))
}
