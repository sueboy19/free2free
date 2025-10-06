package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"free2free/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// formatValidationErrors formats validator errors into a readable string
func formatValidationErrors(errors validator.ValidationErrors) string {
	var msgs []string
	for _, err := range errors {
		field := strings.ToLower(err.Field())
		tag := err.Tag()
		param := err.Param()

		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", field, param)
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", field, param)
		case "oneof":
			msg = fmt.Sprintf("%s must be one of %s", field, param)
		case "url":
			msg = fmt.Sprintf("%s must be a valid URL", field)
		default:
			msg = fmt.Sprintf("%s validation failed on %s", tag, field)
		}
		msgs = append(msgs, msg)
	}
	return strings.Join(msgs, "; ")
}

// CustomRecovery 捕捉 panic 並轉換為 error
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Error(errors.NewAppError(http.StatusInternalServerError, "Internal server error"))
			}
		}()
		c.Next()
	}
}

// ErrorHandler 統一處理錯誤
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			var status int
			var message string

			if appErr, ok := lastErr.Err.(*errors.AppError); ok {
				status = appErr.Status()
				message = appErr.Message
			} else if lastErr.Type == gin.ErrorTypeBind {
				// Handle binding errors, including validator errors
				if ve, ok := lastErr.Err.(validator.ValidationErrors); ok {
					message = formatValidationErrors(ve)
				} else {
					// Other binding errors like JSON parse
					message = lastErr.Err.Error()
				}
				status = http.StatusBadRequest
			} else if httpErr, ok := lastErr.Err.(interface{ Status() int }); ok {
				status = httpErr.Status()
				message = lastErr.Error()
			} else {
				status = http.StatusInternalServerError
				message = "Internal server error"
			}

			resp := ErrorResponse{
				Error: message,
				Code:  status,
			}

			// 確保響應是 JSON
			c.Header("Content-Type", "application/json")
			c.JSON(status, resp)
			c.Abort()
		}
	}
}
