package errors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Status() int {
	return e.Code
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func MapGORMError(err error) *AppError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewAppError(http.StatusNotFound, "Record not found")
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return NewAppError(http.StatusConflict, "Duplicate key error")
	}
	// Add more GORM error mappings as needed
	return NewAppError(http.StatusInternalServerError, "Database error: "+err.Error())
}

func NewValidationError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, message)
}
