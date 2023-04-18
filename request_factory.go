package openai

import (
	"bytes"
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)
// The RequestFactory interface defines the behavior for creating an http.Request object given certain parameters.
type RequestFactory interface {
	Build(ctx context.Context, method, url string, request any) (*http.Request, error)
}

// The httpRequestFactory struct implements the RequestFactory interface.
type httpRequestFactory struct{}

// The Build method of the httpRequestFactory struct builds and returns an http.Request object given the specified parameters.
// It takes in parameters: 
// - ctx, of type context.Context which is the execution context of the function.
// - method, of type string which represents the HTTP method to be used for the request.
// - url, of type string which is the URL that the request should be sent to.
// - request, of type any, which specifies the data to be sent in the request body. It uses Sonic to marshal the request data into bytes.
// It returns:
// - A pointer to the created *http.Request object which may consist of the specified context, method, url and request data.
// - err, an error object which will contain any errors that occurred during assembly of http.Request object.
func (f *httpRequestFactory) Build(ctx context.Context, method, url string, request any) (*http.Request, error) {
	// Check if request data is nil.
	if request == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	//Marshal the request data using Sonic marshaler.
	requestBytes, err := sonic.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshal request error")
	}
	//Use the marshaled bytes to create a new request with the specified context, method and URL.
	return http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBytes))
}

