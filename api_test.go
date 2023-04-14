package openai

import (
	"context"
	"os"
	"testing"
)

func TestAPI(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	org := os.Getenv("ORG")
	proxyUrl := os.Getenv("PROXYURL")
	if apiKey == "" {
		t.Fatal("no openai token got from environment. Try later after set API_KEY environment variable")
	}

	client := NewClient(apiKey, org, proxyUrl)

	var err error

	// test model list
	_, err = client.ModelList(context.TODO())
	if err != nil {
		t.Fatalf("test model list error: %s", err.Error())
	}

	// test model
	_, err = client.ModelInfo(context.TODO(), GPT3Dot5Turbo)
	if err != nil {
		t.Fatalf("test model info error: %s", err.Error())
	}

	// test chat
	request := &ChatCompletionRequset{
		Model: GPT3Dot5Turbo,
		Messages: []Message{
			{Role: ChatMessageRoleUser, Content: "hello"},
		},
		Temperature: 0.2,
	}
	_, err = client.CreateChatCompletion(context.TODO(), request)
	if err != nil {
		t.Fatalf("test CreateChatCompletion error: %s", err.Error())
	}
}
