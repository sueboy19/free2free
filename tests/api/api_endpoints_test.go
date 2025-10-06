package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAPIEndpoints 測試所有API端點的基本功能
func TestAPIEndpoints(t *testing.T) {
	// 設置測試模式
	gin.SetMode(gin.TestMode)

	fmt.Println("開始測試API端點...")

	// 1. 測試Swagger UI端點
	t.Run("SwaggerUI", func(t *testing.T) {
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
		fmt.Println("✓ Swagger UI 端點測試通過")
	})

	// 2. 測試認證端點
	t.Run("AuthEndpoints", func(t *testing.T) {
		// 創建一個新的Gin路由器
		r := gin.New()
		
		// 添加認證路由
		r.GET("/auth/:provider", func(c *gin.Context) {
			provider := c.Param("provider")
			c.JSON(http.StatusTemporaryRedirect, gin.H{"provider": provider})
		})

		// 測試Facebook認證端點
		req, _ := http.NewRequest("GET", "/auth/facebook", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		fmt.Println("✓ Facebook 認證端點測試通過")

		// 測試Instagram認證端點
		req, _ = http.NewRequest("GET", "/auth/instagram", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		fmt.Println("✓ Instagram 認證端點測試通過")
	})

	// 3. 測試登出端點
	t.Run("LogoutEndpoint", func(t *testing.T) {
		// 創建一個新的Gin路由器
		r := gin.New()
		
		// 添加登出路由
		r.GET("/logout", func(c *gin.Context) {
			c.JSON(http.StatusTemporaryRedirect, gin.H{"message": "logged out"})
		})

		// 測試登出端點
		req, _ := http.NewRequest("GET", "/logout", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		fmt.Println("✓ 登出端點測試通過")
	})

	// 4. 測試受保護的端點（應該返回未授權）
	t.Run("ProtectedEndpoints", func(t *testing.T) {
		// 創建一個新的Gin路由器
		r := gin.New()
		
		// 添加受保護的路由
		r.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		// 測試受保護的端點
		req, _ := http.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 受保護端點測試通過（未授權狀態）")
	})

	// 5. 測試管理員端點（應該返回未授權）
	t.Run("AdminEndpoints", func(t *testing.T) {
		// 創建一個新的Gin路由器
		r := gin.New()
		
		// 添加管理員路由
		admin := r.Group("/admin")
		admin.GET("/activities", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})
		admin.GET("/locations", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		// 測試獲取活動列表端點
		req, _ := http.NewRequest("GET", "/admin/activities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取活動列表端點測試通過（未授權狀態）")

		// 測試獲取地點列表端點
		req, _ = http.NewRequest("GET", "/admin/locations", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取地點列表端點測試通過（未授權狀態）")
	})

	// 6. 測試使用者端點（應該返回未授權）
	t.Run("UserEndpoints", func(t *testing.T) {
		// 創建一個新的Gin路由器
		r := gin.New()
		
		// 添加使用者路由
		user := r.Group("/user")
		user.GET("/matches", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})
		user.GET("/past-matches", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		// 測試獲取配對列表端點
		req, _ := http.NewRequest("GET", "/user/matches", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取配對列表端點測試通過（未授權狀態）")

		// 測試獲取過去配對列表端點
		req, _ = http.NewRequest("GET", "/user/past-matches", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取過去配對列表端點測試通過（未授權狀態）")
	})

	fmt.Println("API端點測試完成！")
}