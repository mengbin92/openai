package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/bytedance/sonic"
	"github.com/mengbin92/openai/client"
	"github.com/mengbin92/openai/models"
)

func main() {
	apikey := os.Getenv("APIKEY")

	client := client.NewClient(apikey, "")

	req := models.Requset{
		Model: "gpt-3.5-turbo",
		Messages: []models.Message{
			{Role: "user", Content: "Say this is a test"},
		},
		Temperature: 0.7,
	}

	reqByte, _ := sonic.Marshal(req)
	resp, err := client.Do("POST", bytes.NewBuffer(reqByte))
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	respChat := &models.Response{}
	err = sonic.Unmarshal(body, respChat)
	if err != nil {
		panic(err)
	}
	fmt.Print(respChat)
}
