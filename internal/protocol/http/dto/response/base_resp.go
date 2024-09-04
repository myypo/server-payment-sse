package response

import httpErr "payment-sse/internal/protocol/http/error"

type BaseResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func NewSuccessResponse(data any) any {
	return data
}

func NewErrorResponse(err *httpErr.HttpError) BaseResponse {
	return BaseResponse{
		Error: err.Error(),
	}
}
