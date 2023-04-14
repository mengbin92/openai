package openai

// APIError provides error information returned by the OpenAI API.
type APIError struct {
	Code       string `json:"code,omitempty"`
	Message    string `json:"message"`
	Param      string `json:"param,omitempty"`
	Type       string `json:"type"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

type ErrResponse struct {
	Err *APIError `json:"error,omitempty"`
}

func (e *ErrResponse) Error() string {
	return e.Err.Message
}

type RequestError struct {
	Code int
	Err  error
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}
