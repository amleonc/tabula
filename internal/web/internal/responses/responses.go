package responses

type SuccessResponse struct {
	Code int `json:"code,omitempty"`
	Body any `json:"body,omitempty"`
}

type ErrorResponse struct {
	Code  int `json:"code,omitempty"`
	Error any `json:"error,omitempty"`
}

func NewSuccessResponse(code int, body any) SuccessResponse {
	return SuccessResponse{code, body}
}

func NewErrorResponse(code int, err error) ErrorResponse {
	return ErrorResponse{code, err.Error()}
}
