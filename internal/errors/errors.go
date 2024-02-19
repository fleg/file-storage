package errors

import (
	"fmt"
	"net/http"
)

type BaseError struct {
	Title    string
	Detail   string
	Status   int
	Internal error
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Internal)
}

func (e *BaseError) Unwrap() error {
	return e.Internal
}

func NewNotFoundError(detail string, err error) error {
	return &BaseError{
		Title:    "NotFoundError",
		Detail:   detail,
		Status:   http.StatusNotFound,
		Internal: err,
	}
}

func NewBadRequestError(detail string, err error) error {
	return &BaseError{
		Title:    "BadRequestError",
		Detail:   detail,
		Status:   http.StatusBadRequest,
		Internal: err,
	}
}

func NewUnsupportedMediaTypeError(detail string, err error) error {
	return &BaseError{
		Title:    "UnsupportedMediaTypeError",
		Detail:   detail,
		Status:   http.StatusUnsupportedMediaType,
		Internal: err,
	}
}
