package httpErr

import (
	"net/http"
	dom "payment-sse/internal/domain"
)

type HttpError struct {
	userErr error
	Code    int
}

func NewBadRequest(err error) *HttpError {
	return &HttpError{
		userErr: err,
		Code:    http.StatusBadRequest,
	}
}

func NewTimeout(err error) *HttpError {
	return &HttpError{
		userErr: err,
		Code:    http.StatusRequestTimeout,
	}
}

func HttpErrorFromDom(de dom.DomError) *HttpError {
	return &HttpError{
		userErr: de.User(),
		Code: (func() int {
			switch de.Type() {
			case dom.NotFound:
				return http.StatusNotFound
			case dom.Conflict:
				return http.StatusConflict
			case dom.Gone:
				return http.StatusGone
			case dom.BadRequest:
				return http.StatusBadRequest
			case dom.Forbidden:
				return http.StatusForbidden
			case dom.Internal:
				return http.StatusInternalServerError
			}

			return http.StatusInternalServerError
		})(),
	}
}

func (e *HttpError) Error() string {
	return e.userErr.Error()
}
