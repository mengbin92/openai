package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/openai"
)

func chat(ctx *gin.Context) {

	request := &openai.ChatRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Errorf("Binding ChatRequest struct error: %s\n", err.Error())
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

func audio(ctx *gin.Context) {
	log.Info("call Transcriptions from audio to text")
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Errorf("获取file文件失败, %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	src, _ := os.Create(fileHeader.Filename)
	defer src.Close()

	defer os.Remove(fileHeader.Filename)

	io.Copy(src, file)

	req := &openai.TranscriptionsRequest{
		File:  src,
		Model: openai.GPT3Whisper1,
	}
	resp, err := client.CreateTranscriptions(ctx, req)
	if err != nil {
		log.Errorf("CreateTranscriptions from openai error: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": resp, "code": http.StatusOK})
}

type editsRequest struct {
	Input       string `json:"input" form:"input"`
	Instruction string `json:"instruction" form:"instruction"`
}

func edits(ctx *gin.Context) {
	log.Info("call edits from openai")
	request := &editsRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Errorf("Binding editsRequest struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &openai.EditsRequest{
		Model:       openai.GPT3TextDavincEdit001,
		Input:       request.Input,
		Instruction: request.Instruction,
	}
	resp, err := client.CreateEdits(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func imageGen(ctx *gin.Context) {
	log.Info("call images generations function")
	req := &openai.ImagesReuqest{}
	if err := ctx.Bind(req); err != nil {
		log.Errorf("Binding ImagesReuqest struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.CreateImage(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func embedding(ctx *gin.Context) {
	log.Info("call embedding  function")
	req := &openai.EmbeddingsRequest{}
	if err := ctx.Bind(req); err != nil {
		log.Errorf("Binding EmbeddingsRequest struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Model == "" {
		req.Model = openai.TextEmbeddingAda002
	}

	resp, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func upload(ctx *gin.Context) {
	log.Info("upload file to openai")
	purpose := ctx.PostForm("purpose")
	if purpose == "" {
		purpose = "fine-tune"
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Errorf("获取file文件失败, %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	src, _ := os.Create(fileHeader.Filename)
	defer src.Close()

	defer os.Remove(fileHeader.Filename)

	io.Copy(src, file)

	req := &openai.FileRequest{
		File:    fileHeader.Filename,
		Purpose: purpose,
	}

	resp, err := client.CreateFile(ctx, req)
	if err != nil {
		log.Errorf("get error when up load file: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}
