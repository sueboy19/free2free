package errors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AppError struct {
	Code      int    `json:"code"`
	Message   string `json:"error"`
	ErrorCode string `json:"code_error,omitempty"` // Standardized error code
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

func NewAppErrorWithCode(code int, message string, errorCode string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		ErrorCode: errorCode,
	}
}

func MapGORMError(err error) *AppError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewAppErrorWithCode(http.StatusNotFound, "resource not found", "NOT_FOUND")
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return NewAppErrorWithCode(http.StatusConflict, "duplicate resource", "DUPLICATE_RESOURCE")
	}
	// Add more GORM error mappings as needed
	return NewAppErrorWithCode(http.StatusInternalServerError, "database error", "DATABASE_ERROR")
}

func NewValidationError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusBadRequest, message, "VALIDATION_ERROR")
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusUnauthorized, message, "AUTH_REQUIRED")
}

func NewForbiddenError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusForbidden, message, "FORBIDDEN")
}

// NewAuthenticationError creates standardized authentication error
func NewAuthenticationError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusUnauthorized, "authentication failed", "AUTH_FAILED")
}

// NewTokenExpiredError creates standardized token expired error
func NewTokenExpiredError() *AppError {
	return NewAppErrorWithCode(http.StatusUnauthorized, "token expired", "TOKEN_EXPIRED")
}

// NewInvalidTokenError creates standardized invalid token error
func NewInvalidTokenError() *AppError {
	return NewAppErrorWithCode(http.StatusUnauthorized, "invalid token format", "INVALID_FORMAT")
}

// NewOAuthError creates standardized OAuth error
func NewOAuthError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusBadRequest, message, "OAUTH_FAILED")
}

// NewInternalError creates standardized internal server error
func NewInternalError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusInternalServerError, message, "INTERNAL_ERROR")
}

// NewBadRequestError creates standardized bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppErrorWithCode(http.StatusBadRequest, message, "INVALID_INPUT")
}
