package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/openai"
)

// The ChatRequest struct defines the structure for requests to be sent to the chat service. It contains two fields:
// - Content: the content of the message
// - Tokens: the number of tokens to generate in the response.
type ChatRequest struct {
	Content string `json:"content" form:"content"`
	Tokens  int    `json:"tokens,omitempty" form:"tokens,omitempty"`
}

func chat(ctx *gin.Context) {
	request := &ChatRequest{}
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
		FilePath: fileHeader.Filename,
		Model:    openai.GPT3Whisper1,
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

func deleteFile(ctx *gin.Context) {
	log.Info("delete file from openai")

	id := ctx.Param("id")

	resp, err := client.DeleteFile(ctx, &openai.File{Id: id})
	if err != nil {
		log.Errorf("get error when delete file from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func retrieveFile(ctx *gin.Context) {
	log.Info("retrieve file from openai")

	id := ctx.Param("id")

	resp, err := client.RetrieveFile(ctx, &openai.File{Id: id})
	if err != nil {
		log.Errorf("get error when retrieve file from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func retrieveFileContent(ctx *gin.Context) {
	log.Info("retrieve file content from openai")

	id := ctx.Param("id")

	resp, err := client.RetrieveFileContent(ctx, &openai.File{Id: id})
	if err != nil {
		log.Errorf("get error when retrieve file content from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func retrieveFiles(ctx *gin.Context) {
	log.Info("retrieve files from openai")

	resp, err := client.ListFiles(ctx)
	if err != nil {
		log.Errorf("get error when retrieve file from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func createFineTune(ctx *gin.Context) {
	log.Info("create fine-tunes from openai")

	id := ctx.Param("id")

	req := openai.FineTuneRequest{
		TrainingFile: id,
	}

	resp, err := client.CreateFineTune(ctx, req)
	if err != nil {
		log.Errorf("get error when CreateFineTune from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func listFineTune(ctx *gin.Context) {
	log.Info("list fine-tunes from openai")

	resp, err := client.ListFineTunes(ctx)
	if err != nil {
		log.Errorf("get error when ListFineTunes from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func retrieveFineTune(ctx *gin.Context) {
	log.Info("retrieve fine-tunes from openai")

	id := ctx.Param("id")

	resp, err := client.RetrieveFineTune(ctx, id)
	if err != nil {
		log.Errorf("get error when RetrieveFineTune from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func cancleFineTune(ctx *gin.Context) {
	log.Info("cancle fine-tunes from openai")

	id := ctx.Param("id")

	resp, err := client.CancelFineTune(ctx, id)
	if err != nil {
		log.Errorf("get error when CancelFineTune from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func listFineTuneEvents(ctx *gin.Context) {
	log.Info("cancle fine-tunes from openai")

	id := ctx.Param("id")

	resp, err := client.ListFineTuneEvents(ctx, id)
	if err != nil {
		log.Errorf("get error when ListFineTuneEvents from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func deleteFineTuneModel(ctx *gin.Context) {
	log.Info("delete fine-tunes model from openai")

	models := ctx.Param("models")

	resp, err := client.DeleteFineTuneModel(ctx, models)
	if err != nil {
		log.Errorf("get error when DeleteFineTuneModel from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}

func createModerations(ctx *gin.Context) {
	log.Info("create moderations from openai")

	var req openai.ModerationRequest
	if err := ctx.Bind(&req); err != nil {
		log.Errorf("Binding ModerationRequest struct error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.CreateModeration(ctx, req)
	if err != nil {
		log.Errorf("get error when DeleteFineTuneModel from openai: %s", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": resp})
}
