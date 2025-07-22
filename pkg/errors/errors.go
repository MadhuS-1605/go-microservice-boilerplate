package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
func ErrNotFound(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func ErrAlreadyExists(resource string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("%s already exists", resource),
	}
}

func ErrInvalidInput(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func ErrInternalServer(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Err:     err,
	}
}

func ErrUnauthorized(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func ErrForbidden(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}
