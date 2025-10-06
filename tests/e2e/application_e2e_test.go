package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"free2free/handlers"
	"free2free/models"
	"free2free/routes"
)

// TestFullApplicationFlow 測試完整的應用程式流程
func TestFullApplicationFlow(t *testing.T) {
	// 設置測試模式
	gin.SetMode(gin.TestMode)

	fmt.Println("開始測試完整的應用程式流程...")

	// 1. 測試應用程式啟動
	t.Run("ApplicationStartup", func(t *testing.T) {
		// 創建路由器
		router := createTestRouter()

		// 創建測試服務器
		server, baseURL := createTestServer(router)
		defer func() {
			// 關閉服務器
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		}()

		// 測試服務器是否正常運行
		resp, err := http.Get(baseURL + "/swagger/index.html")
		if err != nil {
			t.Fatalf("無法連接到測試服務器: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		fmt.Println("✓ 應用程式啟動測試通過")
	})

	// 2. 測試API端點響應
	t.Run("APIEndpoints", func(t *testing.T) {
		// 創建路由器
		router := createTestRouter()

		// 創建測試服務器
		server, baseURL := createTestServer(router)
		defer func() {
			// 關閉服務器
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		}()

		// 測試認證端點
		testAuthEndpoints(t, baseURL)

		// 測試管理員端點（未授權）
		testAdminEndpointsUnauthorized(t, baseURL)

		// 測試使用者端點（未授權）
		testUserEndpointsUnauthorized(t, baseURL)

		// 測試開局者端點（未授權）
		testOrganizerEndpointsUnauthorized(t, baseURL)

		// 測試評分端點（未授權）
		testReviewEndpointsUnauthorized(t, baseURL)

		// 測試評論點讚/倒讚端點（未授權）
		testReviewLikeEndpointsUnauthorized(t, baseURL)

		fmt.Println("✓ API端點響應測試通過")
	})

	fmt.Println("完整的應用程式流程測試完成！")
}

// createTestRouter 創建測試路由器
func createTestRouter() *gin.Engine {
	// 創建路由器
	router := gin.New()
	router.Use(gin.Recovery())

	// 添加Swagger路由
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.String(http.StatusOK, "Swagger UI")
	})

	// 設定 session middleware
	router.Use(testSessionsMiddlewareE2E())

	// OAuth 認證路由
	router.GET("/auth/:provider", handlers.OauthBegin)
	router.GET("/auth/:provider/callback", handlers.OauthCallback)

	// 登出路由
	router.GET("/logout", handlers.Logout)

	// 受保護的路由範例
	router.GET("/profile", handlers.Profile)

	// 設定管理後台路由
	routes.SetupAdminRoutes(router)

	// 設定使用者路由
	routes.SetupUserRoutes(router)

	// 設定開局者路由
	routes.SetupOrganizerRoutes(router)

	// 設定評分路由
	routes.SetupReviewRoutes(router)

	// 設定評論點讚/倒讚路由
	routes.SetupReviewLikeRoutes(router)

	return router
}

// createTestServer 創建測試服務器
func createTestServer(router *gin.Engine) (*http.Server, string) {
	// 監聽隨機端口
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("無法創建監聽器: %v", err))
	}

	// 創建服務器
	server := &http.Server{Addr: listener.Addr().String(), Handler: router}

	// 在goroutine中啟動服務器
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("服務器錯誤: %v", err))
		}
	}()

	// 等待服務器啟動
	time.Sleep(100 * time.Millisecond)

	// 返回服務器和基礎URL
	return server, "http://" + listener.Addr().String()
}

// testAuthEndpoints 測試認證端點
func testAuthEndpoints(t *testing.T, baseURL string) {
	// 測試Facebook OAuth端點
	resp, err := http.Get(baseURL + "/auth/facebook")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	// 測試Instagram OAuth端點
	resp, err = http.Get(baseURL + "/auth/instagram")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	// 測試登出端點
	resp, err = http.Get(baseURL + "/logout")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	fmt.Println("  ✓ 認證端點測試通過")
}

// testAdminEndpointsUnauthorized 測試管理員端點（未授權）
func testAdminEndpointsUnauthorized(t *testing.T, baseURL string) {
	// 測試獲取活動列表
	resp, err := http.Get(baseURL + "/admin/activities")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// 測試獲取地點列表
	resp, err = http.Get(baseURL + "/admin/locations")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	fmt.Println("  ✓ 管理員端點（未授權）測試通過")
}

// testUserEndpointsUnauthorized 測試使用者端點（未授權）
func testUserEndpointsUnauthorized(t *testing.T, baseURL string) {
	// 測試獲取配對列表
	resp, err := http.Get(baseURL + "/user/matches")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// 測試獲取過去的配對列表
	resp, err = http.Get(baseURL + "/user/past-matches")
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	fmt.Println("  ✓ 使用者端點（未授權）測試通過")
}

// testOrganizerEndpointsUnauthorized 測試開局者端點（未授權）
func testOrganizerEndpointsUnauthorized(t *testing.T, baseURL string) {
	// 測試審核通過參與者
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", baseURL+"/organizer/matches/1/participants/1/approve", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// 測試審核拒絕參與者
	req, _ = http.NewRequest("PUT", baseURL+"/organizer/matches/1/participants/1/reject", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	fmt.Println("  ✓ 開局者端點（未授權）測試通過")
}

// testReviewEndpointsUnauthorized 測試評分端點（未授權）
func testReviewEndpointsUnauthorized(t *testing.T, baseURL string) {
	// 測試創建評分
	review := models.Review{
		RevieweeID: 2,
		Score:      5,
		Comment:    "測試評分",
	}
	jsonValue, _ := json.Marshal(review)
	resp, err := http.Post(baseURL+"/review/matches/1", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	fmt.Println("  ✓ 評分端點（未授權）測試通過")
}

// testReviewLikeEndpointsUnauthorized 測試評論點讚/倒讚端點（未授權）
func testReviewLikeEndpointsUnauthorized(t *testing.T, baseURL string) {
	// 測試點讚評論
	client := &http.Client{}
	req, _ := http.NewRequest("POST", baseURL+"/review-like/reviews/1/like", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// 測試倒讚評論
	req, _ = http.NewRequest("POST", baseURL+"/review-like/reviews/1/dislike", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("請求失敗: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	fmt.Println("  ✓ 評論點讚/倒讚端點（未授權）測試通過")
}

// testSessionsMiddlewareE2E 模擬 session 中介層，用於E2E測試
func testSessionsMiddlewareE2E() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 為測試目的，模擬一個空的 session
		c.Set("session", nil)
		c.Next()
	}
}