package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"free2free/handlers"
	"free2free/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"free2free/routes"
)

// E2ETestSuite 代表端到端測試套件
type E2ETestSuite struct {
	router *gin.Engine
}

// NewE2ETestSuite 創建新的測試套件
func NewE2ETestSuite() *E2ETestSuite {
	// 設置測試模式
	gin.SetMode(gin.TestMode)

	// 創建路由器
	router := gin.New()
	router.Use(gin.Recovery())

	return &E2ETestSuite{
		router: router,
	}
}

// SetupRoutes 設置所有路由
func (suite *E2ETestSuite) SetupRoutes() {
	// 添加Swagger路由
	suite.router.GET("/swagger/*any", func(c *gin.Context) {
		c.String(http.StatusOK, "Swagger UI")
	})

	// 設定 session middleware
	suite.router.Use(testSessionsMiddlewareSuite())

	// OAuth 認證路由
	suite.router.GET("/auth/:provider", handlers.OauthBegin)
	suite.router.GET("/auth/:provider/callback", handlers.OauthCallback)

	// 登出路由
	suite.router.GET("/logout", handlers.Logout)

	// 受保護的路由範例
	suite.router.GET("/profile", handlers.Profile)

	// 設定管理後台路由
	routes.SetupAdminRoutes(suite.router)

	// 設定使用者路由
	routes.SetupUserRoutes(suite.router)

	// 設定開局者路由
	routes.SetupOrganizerRoutes(suite.router)

	// 設定評分路由
	routes.SetupReviewRoutes(suite.router)

	// 設定評論點讚/倒讚路由
	routes.SetupReviewLikeRoutes(suite.router)
}

// TestAPIEndpoint 測試單個API端點
func (suite *E2ETestSuite) TestAPIEndpoint(t *testing.T, method, url string, body interface{}, expectedCode int) {
	var req *http.Request
	var err error

	if body != nil {
		jsonValue, _ := json.Marshal(body)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Fatalf("創建請求失敗: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			t.Fatalf("創建請求失敗: %v", err)
		}
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, expectedCode, w.Code, fmt.Sprintf("請求 %s %s 應返回狀態碼 %d，但實際返回 %d", method, url, expectedCode, w.Code))
}

// TestE2EFlow 進行端到端測試
func TestE2EFlow(t *testing.T) {
	// 創建測試套件
	suite := NewE2ETestSuite()
	suite.SetupRoutes()

	fmt.Println("開始進行端到端測試...")

	// 1. 測試Swagger UI
	t.Run("SwaggerUI", func(t *testing.T) {
		suite.TestAPIEndpoint(t, "GET", "/swagger/index.html", nil, http.StatusOK)
		fmt.Println("✓ Swagger UI 測試通過")
	})

	// 2. 測試OAuth認證端點
	t.Run("OAuthEndpoints", func(t *testing.T) {
		suite.TestAPIEndpoint(t, "GET", "/auth/facebook", nil, http.StatusTemporaryRedirect)
		fmt.Println("✓ Facebook OAuth端點測試通過")

		suite.TestAPIEndpoint(t, "GET", "/auth/instagram", nil, http.StatusTemporaryRedirect)
		fmt.Println("✓ Instagram OAuth端點測試通過")
	})

	// 3. 測試登出端點
	t.Run("LogoutEndpoint", func(t *testing.T) {
		suite.TestAPIEndpoint(t, "GET", "/logout", nil, http.StatusTemporaryRedirect)
		fmt.Println("✓ 登出端點測試通過")
	})

	// 4. 測試受保護的個人檔案端點（應該返回未授權）
	t.Run("ProfileEndpoint", func(t *testing.T) {
		suite.TestAPIEndpoint(t, "GET", "/profile", nil, http.StatusUnauthorized)
		fmt.Println("✓ 個人檔案端點測試通過（未授權狀態）")
	})

	// 5. 測試管理員端點（應該返回未授權）
	t.Run("AdminEndpoints", func(t *testing.T) {
		// 測試獲取活動列表
		suite.TestAPIEndpoint(t, "GET", "/admin/activities", nil, http.StatusUnauthorized)
		fmt.Println("✓ 獲取活動列表端點測試通過（未授權狀態）")

		// 測試獲取地點列表
		suite.TestAPIEndpoint(t, "GET", "/admin/locations", nil, http.StatusUnauthorized)
		fmt.Println("✓ 獲取地點列表端點測試通過（未授權狀態）")
	})

	// 6. 測試使用者端點（應該返回未授權）
	t.Run("UserEndpoints", func(t *testing.T) {
		// 測試獲取配對列表
		suite.TestAPIEndpoint(t, "GET", "/user/matches", nil, http.StatusUnauthorized)
		fmt.Println("✓ 獲取配對列表端點測試通過（未授權狀態）")

		// 測試獲取過去的配對列表
		suite.TestAPIEndpoint(t, "GET", "/user/past-matches", nil, http.StatusUnauthorized)
		fmt.Println("✓ 獲取過去配對列表端點測試通過（未授權狀態）")
	})

	// 7. 測試開局者端點（應該返回未授權）
	t.Run("OrganizerEndpoints", func(t *testing.T) {
		// 測試審核通過參與者
		suite.TestAPIEndpoint(t, "PUT", "/organizer/matches/1/participants/1/approve", nil, http.StatusUnauthorized)
		fmt.Println("✓ 審核通過參與者端點測試通過（未授權狀態）")

		// 測試審核拒絕參與者
		suite.TestAPIEndpoint(t, "PUT", "/organizer/matches/1/participants/1/reject", nil, http.StatusUnauthorized)
		fmt.Println("✓ 審核拒絕參與者端點測試通過（未授權狀態）")
	})

	// 8. 測試評分端點（應該返回未授權）
	t.Run("ReviewEndpoints", func(t *testing.T) {
		// 測試創建評分
		review := models.Review{
			RevieweeID: 2,
			Score:      5,
			Comment:    "測試評分",
		}
		suite.TestAPIEndpoint(t, "POST", "/review/matches/1", review, http.StatusUnauthorized)
		fmt.Println("✓ 創建評分端點測試通過（未授權狀態）")
	})

	// 9. 測試評論點讚/倒讚端點（應該返回未授權）
	t.Run("ReviewLikeEndpoints", func(t *testing.T) {
		// 測試點讚評論
		suite.TestAPIEndpoint(t, "POST", "/review-like/reviews/1/like", nil, http.StatusUnauthorized)
		fmt.Println("✓ 點讚評論端點測試通過（未授權狀態）")

		// 測試倒讚評論
		suite.TestAPIEndpoint(t, "POST", "/review-like/reviews/1/dislike", nil, http.StatusUnauthorized)
		fmt.Println("✓ 倒讚評論端點測試通過（未授權狀態）")
	})

	fmt.Println("端到端測試完成！")
}

// testSessionsMiddlewareSuite 模擬 session 中介層，用於測試套件
func testSessionsMiddlewareSuite() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 為測試目的，模擬一個空的 session
		c.Set("session", nil)
		c.Next()
	}
}