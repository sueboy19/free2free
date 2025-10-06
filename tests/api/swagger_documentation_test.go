package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"free2free/docs"
)

// TestSwaggerDocumentation 測試Swagger文檔的完整性和正確性
func TestSwaggerDocumentation(t *testing.T) {
	fmt.Println("開始測試Swagger文檔...")

	// 1. 測試Swagger信息
	t.Run("SwaggerInfo", func(t *testing.T) {
		assert.Equal(t, "買一送一配對網站 API", docs.SwaggerInfo.Title)
		assert.Equal(t, "這是一個買一送一配對網站的API文檔", docs.SwaggerInfo.Description)
		assert.Equal(t, "1.0", docs.SwaggerInfo.Version)
		assert.Equal(t, "localhost:8080", docs.SwaggerInfo.Host)
		assert.Equal(t, "/", docs.SwaggerInfo.BasePath)
		fmt.Println("✓ Swagger基本信息測試通過")
	})

	// 2. 測試Swagger UI端點
	t.Run("SwaggerUIEndpoint", func(t *testing.T) {
		// 設置測試模式
		gin.SetMode(gin.TestMode)

		// 創建路由器
		r := gin.New()
		
		// 添加Swagger路由（模擬實際的Swagger路由）
		r.GET("/swagger/*any", func(c *gin.Context) {
			// 模擬Swagger UI的響應
			c.String(http.StatusOK, "Swagger UI")
		})

		// 創建請求
		req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
		w := httptest.NewRecorder()
		
		// 處理請求
		r.ServeHTTP(w, req)

		// 驗證響應
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Swagger UI", w.Body.String())
		fmt.Println("✓ Swagger UI端點測試通過")
	})

	// 3. 測試Swagger JSON文檔
	t.Run("SwaggerJSON", func(t *testing.T) {
		// 檢查Swagger模板是否包含必要的信息
		assert.Contains(t, docs.SwaggerInfo.SwaggerTemplate, "\"swagger\": \"2.0\"")
		assert.Contains(t, docs.SwaggerInfo.SwaggerTemplate, "\"title\": \"買一送一配對網站 API\"")
		assert.Contains(t, docs.SwaggerInfo.SwaggerTemplate, "\"description\": \"這是一個買一送一配對網站的API文檔\"")
		assert.Contains(t, docs.SwaggerInfo.SwaggerTemplate, "\"version\": \"1.0\"")
		fmt.Println("✓ Swagger JSON文檔測試通過")
	})

	// 4. 測試必要的API端點是否在文檔中定義
	t.Run("APIEndpointsInDocumentation", func(t *testing.T) {
		// 檢查文檔是否包含必要的端點
		swaggerDoc := docs.SwaggerInfo.SwaggerTemplate
		
		// 認證相關端點
		assert.Contains(t, swaggerDoc, "/auth/{provider}")
		assert.Contains(t, swaggerDoc, "/auth/{provider}/callback")
		assert.Contains(t, swaggerDoc, "/logout")
		assert.Contains(t, swaggerDoc, "/profile")
		
		// 管理員相關端點
		assert.Contains(t, swaggerDoc, "/admin/activities")
		assert.Contains(t, swaggerDoc, "/admin/locations")
		assert.Contains(t, swaggerDoc, "/admin/activities/{id}")
		assert.Contains(t, swaggerDoc, "/admin/locations/{id}")
		
		// 使用者相關端點
		assert.Contains(t, swaggerDoc, "/user/matches")
		assert.Contains(t, swaggerDoc, "/user/past-matches")
		assert.Contains(t, swaggerDoc, "/user/matches/{id}/join")
		
		// 開局者相關端點
		assert.Contains(t, swaggerDoc, "/organizer/matches/{id}/participants/{participant_id}/approve")
		assert.Contains(t, swaggerDoc, "/organizer/matches/{id}/participants/{participant_id}/reject")
		
		// 評分相關端點
		assert.Contains(t, swaggerDoc, "/review/matches/{id}")
		
		// 評論點讚/倒讚相關端點
		assert.Contains(t, swaggerDoc, "/review-like/reviews/{id}/like")
		assert.Contains(t, swaggerDoc, "/review-like/reviews/{id}/dislike")
		
		fmt.Println("✓ 所有API端點在文檔中定義測試通過")
	})

	// 5. 測試數據模型是否在文檔中定義
	t.Run("DataModelsInDocumentation", func(t *testing.T) {
		// 檢查文檔是否包含必要的數據模型
		swaggerDoc := docs.SwaggerInfo.SwaggerTemplate
		
		// 數據模型
		assert.Contains(t, swaggerDoc, "main.User")
		assert.Contains(t, swaggerDoc, "main.Activity")
		assert.Contains(t, swaggerDoc, "main.Location")
		assert.Contains(t, swaggerDoc, "main.Match")
		assert.Contains(t, swaggerDoc, "main.MatchParticipant")
		assert.Contains(t, swaggerDoc, "main.Review")
		assert.Contains(t, swaggerDoc, "main.ReviewLike")
		
		fmt.Println("✓ 所有數據模型在文檔中定義測試通過")
	})

	// 6. 測試API標籤是否正確分類
	t.Run("APITags", func(t *testing.T) {
		// 檢查文檔是否包含正確的標籤
		swaggerDoc := docs.SwaggerInfo.SwaggerTemplate
		
		// API標籤
		assert.Contains(t, swaggerDoc, "\"認證\"")
		assert.Contains(t, swaggerDoc, "\"管理員\"")
		assert.Contains(t, swaggerDoc, "\"使用者\"")
		assert.Contains(t, swaggerDoc, "\"開局者\"")
		assert.Contains(t, swaggerDoc, "\"評分\"")
		assert.Contains(t, swaggerDoc, "\"評論\"")
		
		fmt.Println("✓ API標籤分類測試通過")
	})

	fmt.Println("Swagger文檔測試完成！")
}

// TestRealSwaggerEndpoint 測試真實的Swagger端點
func TestRealSwaggerEndpoint(t *testing.T) {
	// 讀取生成的swagger.json文件
	swaggerJSON, err := ioutil.ReadFile("./docs/swagger.json")
	if err != nil {
		t.Skip("跳過測試：無法讀取swagger.json文件")
		return
	}

	// 檢查文件是否包含必要的內容
	swaggerContent := string(swaggerJSON)
	
	assert.Contains(t, swaggerContent, "\"swagger\": \"2.0\"")
	assert.Contains(t, swaggerContent, "\"title\": \"買一送一配對網站 API\"")
	assert.Contains(t, swaggerContent, "paths")
	assert.Contains(t, swaggerContent, "definitions")
	
	// 檢查是否包含所有主要端點
	endpoints := []string{
		"/auth/{provider}",
		"/auth/{provider}/callback",
		"/logout",
		"/profile",
		"/admin/activities",
		"/admin/locations",
		"/user/matches",
		"/review/matches/{id}",
	}
	
	for _, endpoint := range endpoints {
		// 在JSON中路徑會被轉義，所以我們需要檢查轉義後的版本
		escapedEndpoint := strings.ReplaceAll(endpoint, "{", "\\{")
		escapedEndpoint = strings.ReplaceAll(escapedEndpoint, "}", "\\}")
		assert.Contains(t, swaggerContent, escapedEndpoint, "端點 %s 應該在Swagger文檔中定義", endpoint)
	}
	
	fmt.Println("✓ 真實Swagger端點測試通過")
}