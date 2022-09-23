package main

import (
	"net/http"

	"github.com/pkg/errors"
)

// Error types for API Gateway Response
var (
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized request")
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
)

// errorCodeMap is to map error with status code
var errorCodeMap = map[error]int{
	ErrBadRequest:          http.StatusBadRequest,
	ErrUnauthorized:        http.StatusUnauthorized,
	ErrNotFound:            http.StatusNotFound,
	ErrInternalServerError: http.StatusInternalServerError,
}

// GetCodeFromError gets the code of the error in the dict.
// Returns 500 if the error is not present
func GetCodeFromError(err error) int {
	if val, ok := errorCodeMap[err]; ok {
		return val
	}
	return 500
}
