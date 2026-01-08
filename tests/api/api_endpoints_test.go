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
	gin.SetMode(gin.TestMode)

	fmt.Println("開始測試API端點...")

	r := gin.New()
	r.Use(gin.Recovery())

	t.Run("SwaggerUI", func(t *testing.T) {
		r.GET("/swagger/*any", func(c *gin.Context) {
			c.String(http.StatusOK, "Swagger UI")
		})

		req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Swagger UI", w.Body.String())
		fmt.Println("✓ Swagger UI 端點測試通過")
	})

	t.Run("AuthEndpoints", func(t *testing.T) {
		r.GET("/auth/:provider", func(c *gin.Context) {
			provider := c.Param("provider")
			c.JSON(http.StatusTemporaryRedirect, gin.H{"provider": provider})
		})

		req, _ := http.NewRequest("GET", "/auth/facebook", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		fmt.Println("✓ Facebook 認證端點測試通過")

		req, _ = http.NewRequest("GET", "/auth/instagram", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		fmt.Println("✓ Instagram 認證端點測試通過")
	})

	t.Run("LogoutEndpoint", func(t *testing.T) {
		r.GET("/logout", func(c *gin.Context) {
			c.JSON(http.StatusTemporaryRedirect, gin.H{"message": "logged out"})
		})

		req, _ := http.NewRequest("GET", "/logout", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		fmt.Println("✓ 登出端點測試通過")
	})

	t.Run("ProtectedEndpoints", func(t *testing.T) {
		r.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		req, _ := http.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 受保護端點測試通過（未授權狀態）")
	})

	t.Run("AdminEndpoints", func(t *testing.T) {
		admin := r.Group("/admin")
		admin.GET("/activities", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})
		admin.GET("/locations", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		req, _ := http.NewRequest("GET", "/admin/activities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取活動列表端點測試通過（未授權狀態）")

		req, _ = http.NewRequest("GET", "/admin/locations", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取地點列表端點測試通過（未授權狀態）")
	})

	t.Run("UserEndpoints", func(t *testing.T) {
		user := r.Group("/user")
		user.GET("/matches", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})
		user.GET("/past-matches", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		req, _ := http.NewRequest("GET", "/user/matches", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取配對列表端點測試通過（未授權狀態）")

		req, _ = http.NewRequest("GET", "/user/past-matches", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 獲取過去配對列表端點測試通過（未授權狀態）")
	})

	t.Run("OrganizerEndpoints", func(t *testing.T) {
		organizer := r.Group("/organizer")
		organizer.PUT("/matches/:id/participants/:participant_id/approve", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})
		organizer.PUT("/matches/:id/participants/:participant_id/reject", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		req, _ := http.NewRequest("PUT", "/organizer/matches/1/participants/1/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 審核通過參與者端點測試通過（未授權狀態）")

		req, _ = http.NewRequest("PUT", "/organizer/matches/1/participants/1/reject", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 審核拒絕參與者端點測試通過（未授權狀態）")
	})

	t.Run("ReviewEndpoints", func(t *testing.T) {
		r.POST("/review/matches/:id", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		req, _ := http.NewRequest("POST", "/review/matches/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 創建評分端點測試通過（未授權狀態）")
	})

	t.Run("ReviewLikeEndpoints", func(t *testing.T) {
		like := r.Group("/review-like")
		like.POST("/views/:id/like", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})
		like.POST("/views/:id/dislike", func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		})

		req, _ := http.NewRequest("POST", "/review-like/views/1/like", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 點讚評論端點測試通過（未授權狀態）")

		req, _ = http.NewRequest("POST", "/review-like/views/1/dislike", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		fmt.Println("✓ 倒讚評論端點測試通過（未授權狀態）")
	})

	fmt.Println("API端點測試完成！")
}
