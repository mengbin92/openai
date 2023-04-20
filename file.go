package openai

import (
	"bytes"
	"context"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type FileRequest struct {
	File    string `json:"file" form:"file"`
	Purpose string `json:"purpose" form:"purpose"`
}

type File struct {
	Id        string `json:"id" form:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
	Deleted   bool   `json:"deleted"`
}

type Files struct {
	Data []File `json:"data"`
}

func (c *Client) ListFiles(ctx context.Context) (response Files, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fullURL(files), nil)
	if err != nil {
		errors.Wrap(err, "create list files request error")
		return
	}
	if err = c.sendRequest(req, &response); err != nil {
		errors.Wrap(err, "get list files from openai error")
		return
	}
	return
}

func (c *Client) CreateFile(ctx context.Context, request *FileRequest) (response File, err error) {
	buf := &bytes.Buffer{}
	factory := c.formFactory(buf)

	file, err := os.Open(request.File)
	if err != nil {
		errors.Wrapf(err, "load file: %s error", request.File)
	}
	defer file.Close()
	if err = factory.CreateFormFile("file", file); err != nil {
		errors.Wrap(err, "load file error")
		return
	}

	if err = factory.WriteField("purpose", request.Purpose); err != nil {
		errors.Wrap(err, "set purpose error")
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL(files), buf)
	if err != nil {
		errors.Wrap(err, "create upload file request error")
		return
	}
	req.Header.Add("Content-Type", factory.FormDataContentType())
	if err = c.sendRequest(req, &response); err != nil {
		errors.Wrap(err, "upload file to openai error")
		return
	}
	return
}

func (c *Client) DeleteFile(ctx context.Context, request *File) (response File, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL(files+"/"+request.Id), nil)
	if err != nil {
		errors.Wrap(err, "create delete file request error")
		return
	}
	if err = c.sendRequest(req, &response); err != nil {
		errors.Wrapf(err, "delete file with id: %s error", request.Id)
		return
	}
	return
}

func (c *Client) RetrieveFile(ctx context.Context, request *File) (response File, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL(files+"/"+request.Id), nil)
	if err != nil {
		errors.Wrap(err, "create retrieve file request error")
		return
	}
	if err = c.sendRequest(req, &response); err != nil {
		errors.Wrapf(err, "retrieve file with id: %s error", request.Id)
		return
	}
	return
}

func (c *Client) RetrieveFileContent(ctx context.Context, request *File) (response File, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL(files+"/"+request.Id+"/content"), nil)
	if err != nil {
		errors.Wrap(err, "create retrieve file request error")
		return
	}
	if err = c.sendRequest(req, &response); err != nil {
		errors.Wrapf(err, "retrieve file with id: %s error", request.Id)
		return
	}
	return
}
