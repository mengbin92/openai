package openai

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type ImagesReuqest struct {
	Prompt         string `json:"prompt" form:"prompt"`
	N              int    `json:"n,omitempty" form:"n,omitempty"`
	Size           string `json:"size,omitempty" form:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty" form:"response_format,omitempty"` // url or b64_json
	User           string `json:"user,omitempty" form:"user,omitempty"`
}

type ImagesResponse struct {
	Created int64                `json:"created,omitempty"`
	Data    []ImagesResponseData `json:"data,omitempty"`
}

type ImagesResponseData struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
}

func (c *Client) CreateImage(ctx context.Context, request *ImagesReuqest) (response ImagesResponse, err error) {
	req, err := c.requestFactory.Build(ctx, http.MethodPost, fullURL("images/generations"), request)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}

type ImagesEditReuqest struct {
	Image          *os.File `json:"image"`
	Mask           *os.File `json:"mask,omitempty"`
	Prompt         string   `json:"prompt"`
	N              int      `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` // url or b64_json
	User           string   `json:"user,omitempty"`
}

func (c *Client) CreateImageEdits(ctx context.Context, request *ImagesEditReuqest) (response ImagesResponse, err error) {
	buf := &bytes.Buffer{}
	factory := c.formFactory(buf)

	// load image
	if err = factory.CreateFormFile("image", request.Image); err != nil {
		errors.Wrap(err, "load origin image error")
		return
	}

	// load mask
	if request.Mask != nil {
		if err = factory.CreateFormFile("mask", request.Mask); err != nil {
			errors.Wrap(err, "load mask image error")
			return
		}
	}

	// write param
	if err = factory.WriteField("prompt", request.Prompt); err != nil {
		errors.Wrap(err, "write prompt error")
		return
	}

	if request.N != 0 {
		if err = factory.WriteField("n", fmt.Sprintf("%d", request.N)); err != nil {
			errors.Wrap(err, "write n error")
			return
		}
	}

	if request.Size != "" {
		if err = factory.WriteField("size", request.Size); err != nil {
			errors.Wrap(err, "write size error")
			return
		}
	}

	if request.ResponseFormat != "" {
		if err = factory.WriteField("response_format", request.ResponseFormat); err != nil {
			errors.Wrap(err, "write response_format error")
			return
		}
	}

	if request.User != "" {
		if err = factory.WriteField("user", request.User); err != nil {
			errors.Wrap(err, "write user error")
			return
		}
	}

	if err = factory.Close(); err != nil {
		errors.Wrap(err, "write close error")
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL("images/edits"), buf)
	if err != nil {
		errors.Wrap(err, "NewRequestWithContext error")
		return
	}
	err = c.sendRequest(req, &response)
	return
}

type ImagesVariationReuqest struct {
	Image          *os.File `json:"image"`
	N              int      `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` // url or b64_json
	User           string   `json:"user,omitempty"`
}

func (c *Client) CreateImageVariation(ctx context.Context, request *ImagesVariationReuqest) (response ImagesResponse, err error) {
	buf := &bytes.Buffer{}
	factory := c.formFactory(buf)

	// load image
	if err = factory.CreateFormFile("image", request.Image); err != nil {
		errors.Wrap(err, "load origin image error")
		return
	}

	if request.N != 0 {
		if err = factory.WriteField("n", fmt.Sprintf("%d", request.N)); err != nil {
			errors.Wrap(err, "write n error")
			return
		}
	}

	if request.Size != "" {
		if err = factory.WriteField("size", request.Size); err != nil {
			errors.Wrap(err, "write size error")
			return
		}
	}

	if request.ResponseFormat != "" {
		if err = factory.WriteField("response_format", request.ResponseFormat); err != nil {
			errors.Wrap(err, "write response_format error")
			return
		}
	}

	if request.User != "" {
		if err = factory.WriteField("user", request.User); err != nil {
			errors.Wrap(err, "write user error")
			return
		}
	}

	if err = factory.Close(); err != nil {
		errors.Wrap(err, "write close error")
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL("images/edits"), buf)
	if err != nil {
		errors.Wrap(err, "NewRequestWithContext error")
		return
	}
	err = c.sendRequest(req, &response)
	return
}
