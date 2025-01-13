package helper

import (
	"graded-challenge-2-client/dto"
	"net/http"
)

type ErrorResponse dto.WebErrorResponse

var (
	ErrBadRequest = ErrorResponse{
		Status: http.StatusBadRequest,
		Type:   "Bad Request",
	}

	ErrInternalServer = ErrorResponse{
		Status: http.StatusInternalServerError,
		Type:   "Internal Server",
	}

	ErrNotFound = ErrorResponse{
		Status: http.StatusNotFound,
		Type:   "Not Found",
	}

	ErrUnauthorized = ErrorResponse{
		Status: http.StatusUnauthorized,
		Type:   "Unauthorized",
	}
)

func (er *ErrorResponse) ErrorFormat(detail interface{}) (int, *ErrorResponse) {
	er.Detail = detail

	return er.Status, er
}
