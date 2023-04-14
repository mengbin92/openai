package openai

import (
	"bytes"
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type RequestFactory interface {
	Build(ctx context.Context, method, url string, request any) (*http.Request, error)
}

type httpRequestFactory struct {
}

func (f *httpRequestFactory) Build(ctx context.Context, method, url string, request any) (*http.Request, error) {
	if request == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}
	requestBytes, err := sonic.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshal request error")
	}
	return http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBytes))
}
