package errorx

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	PageErrNotFound      = errors.New("not found")
	PageErrInternalError = errors.New("not found")
	PageErrFail          = errors.New("fail")
)

type PageError struct {
	err  error
	code int
}

func (e *PageError) Error() string {
	return e.err.Error()
}

func (e *PageError) Unwrap() error {
	return e.err
}

func (e *PageError) IsError() bool {
	if e != nil {
		return true
	}
	return false
}

func (e *PageError) GetStatusCode() int {
	if e == nil {
		return http.StatusOK
	}
	return e.code
}

func NewPageError(c int, f string, a ...any) *PageError {
	return &PageError{
		err:  fmt.Errorf(f, a...),
		code: c,
	}
}

func NewInternalErrPageError() *PageError {
	return &PageError{
		err:  PageErrInternalError,
		code: http.StatusInternalServerError,
	}
}

func NewNotFoundErrPageError() *PageError {
	return &PageError{
		err:  PageErrNotFound,
		code: http.StatusNotFound,
	}
}

func NewFailErrPageError() *PageError {
	return &PageError{
		err:  PageErrFail,
		code: http.StatusOK,
	}
}
