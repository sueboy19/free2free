package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"free2free/errors"
	"free2free/middleware"
)

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		code         int
		message      string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "成功發送錯誤響應",
			code:         http.StatusBadRequest,
			message:      "無效的請求",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"無效的請求","code":400,"code_error":"INVALID_INPUT"}`,
		},
		{
			name:         "內部伺服器錯誤",
			code:         http.StatusInternalServerError,
			message:      "伺服器錯誤",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"伺服器錯誤","code":500,"code_error":"INTERNAL_ERROR"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			// Add the error handling middleware
			router.Use(middleware.ErrorHandler())

			// Create a simple handler that adds an error
			router.GET("/test", func(c *gin.Context) {
				c.Error(errors.NewAppError(tt.code, tt.message))
				// Don't send any response here; let the error handler do it
			})

			// Make a request to trigger the error handling
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestErrorResponseStruct(t *testing.T) {
	errResp := middleware.ErrorResponse{
		Error:     "測試錯誤",
		Code:      400,
		ErrorCode: "TEST_ERROR",
	}

	data, err := json.Marshal(errResp)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"測試錯誤","code":400,"code_error":"TEST_ERROR"}`, string(data))
}
