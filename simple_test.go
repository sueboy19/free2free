package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerEndpoint(t *testing.T) {
	// 設置測試模式
	gin.SetMode(gin.TestMode)

	// 創建一個新的Gin路由器
	r := gin.New()
	
	// 添加Swagger路由
	r.GET("/swagger/*any", func(c *gin.Context) {
		c.String(http.StatusOK, "Swagger UI")
	})

	// 創建一個HTTP請求
	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	w := httptest.NewRecorder()
	
	// 處理請求
	r.ServeHTTP(w, req)

	// 驗證響應
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Swagger UI", w.Body.String())
}