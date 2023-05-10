package openai

import (
	"bytes"
	"context"
	"net/http"

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

// ListFiles takes a context, send a request to OpenAI  and return a list of files as response
func (c *Client) ListFiles(ctx context.Context) (response Files, err error) {
	// Create a request using the provided context and URL
	req, err := c.requestFactory.Build(ctx, http.MethodGet, fullURL(files), nil)
	if err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "create list files request error")
		return
	}
	// Send the request and receive the response
	if err = c.sendRequest(req, &response); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "get list files from openai error")
		return
	}
	// Return the successful response and no errors
	return
}

// CreateFile takes a context and a file request, creates a new file, and returns the created file as response
func (c *Client) CreateFile(ctx context.Context, request *FileRequest) (response File, err error) {
	// Create an empty buffer to write form data to
	buf := &bytes.Buffer{}
	// Create a new form using the empty buffer
	factory := c.formFactory(buf)
	
	// Add the file data to the form data with the key "file"
	if err = factory.CreateFormFile("file", request.File); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "load file error")
		return
	}
	// Add the purpose value to the form data with the key "purpose"
	if err = factory.WriteField("purpose", request.Purpose); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "set purpose error")
		return
	}

	// close file multipart
	if err = factory.Close(); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "close file multipart error")
		return
	}

	// Create a new POST request using the provided context, URL and buffer of form data
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL(files), buf)
	if err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "create upload file request error")
		return
	}
	// Set the Content-Type header of the request to match the type of form data being sent
	req.Header.Add("Content-Type", factory.FormDataContentType())
	// Send the request and receive the response
	if err = c.sendRequest(req, &response); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "upload file to openai error")
		return
	}
	// Return the successful response and no errors
	return
}

// DeleteFile takes a context and a file, creates a new DELETE request for deleting the file with provided ID,
// and returns the response after sending the request
func (c *Client) DeleteFile(ctx context.Context, request *File) (response File, err error) {
	// Create a new DELETE request using the provided context and URL of the file to be deleted
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL(files+"/"+request.Id), nil)
	if err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "create delete file request error")
		return
	}
	// Send the DELETE request and receive the response
	if err = c.sendRequest(req, &response); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrapf(err, "delete file with id: %s error", request.Id)
		return
	}
	// Return the successful response and no errors
	return
}

// RetrieveFile takes a context and a file, creates a new GET request for retrieving the file with provided ID,
// and returns the response after sending the request
func (c *Client) RetrieveFile(ctx context.Context, request *File) (response File, err error) {
	// Create a new GET request using the provided context and URL of the file to be retrieved
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL(files+"/"+request.Id), nil)
	if err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "create retrieve file request error")
		return
	}
	// Send the GET request and receive the response
	if err = c.sendRequest(req, &response); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrapf(err, "retrieve file with id: %s error", request.Id)
		return
	}
	// Return the successful response and no errors
	return
}

// RetrieveFileContent takes a context and a file, creates a new GET request for retrieving the content of the file with provided ID,
// and returns the response after sending the request
func (c *Client) RetrieveFileContent(ctx context.Context, request *File) (response File, err error) {
	// Create a new GET request using the provided context and URL of the file content to be retrieved
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL(files+"/"+request.Id+"/content"), nil)
	if err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrap(err, "create retrieve file request error")
		return
	}
	// Send the GET request and receive the response
	if err = c.sendRequest(req, &response); err != nil {
		// Wrap the error with additional context for easier debugging and return it
		errors.Wrapf(err, "retrieve file with id: %s error", request.Id)
		return
	}
	// Return the successful response and no errors
	return
}
