# OpenAI

## Overview

This is a library implemented in Go for [OpenAI API](https://platform.openai.com/docs/api-reference). It supports:  

* [Models](https://platform.openai.com/docs/api-reference/models)
* [Completions](https://platform.openai.com/docs/api-reference/completions)
* [Chat](https://platform.openai.com/docs/api-reference/chat)
* [Edits](https://platform.openai.com/docs/api-reference/edits)
* [Images](https://platform.openai.com/docs/api-reference/images)
* [Embeddings](https://platform.openai.com/docs/api-reference/embeddings)
* [Audio](https://platform.openai.com/docs/api-reference/audio)
* [Files](https://platform.openai.com/docs/api-reference/files)
* [Fine-tunes](https://platform.openai.com/docs/api-reference/fine-tunes)
* [Moderations](https://platform.openai.com/docs/api-reference/moderations) 

## Install

```bash
go get github.com/mengbin92/openai
```

## ChatGPT example code  

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mengbin92/openai"
)

func main() {
	client := openai.NewClient("your token", "your org", "proxy")

	resp, err := client.CreateChatCompletion(
		context.Background(),
		&openai.ChatCompletionRequset{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.Message{
				{Role: openai.ChatMessageRoleUser, Content: "hi!"},
			},
		},
	)
	if err != nil {
		fmt.Printf("CreateChatCompletion error: %s\n", err.Error())
		os.Exit(-1)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
```