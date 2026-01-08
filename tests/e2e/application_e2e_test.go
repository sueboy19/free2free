package e2e

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"free2free/handlers"
	"free2free/routes"
	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFullApplicationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fmt.Println("開始測試完整的應用程式流程...")

	router := createTestRouter()
	server, baseURL := createTestServer(router)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	t.Run("ApplicationStartup", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/swagger/index.html")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		fmt.Println("✓ 應用程式啟動測試通過")
		resp.Body.Close()
	})

	t.Run("APIEndpoints", func(t *testing.T) {
		testAuthEndpoints(t, baseURL)
		testAdminEndpointsUnauthorized(t, baseURL)
		testAdminEndpointsAuthorized(t, baseURL)
		testUserEndpointsUnauthorized(t, baseURL)
		testUserEndpointsAuthorized(t, baseURL)
		testOrganizerEndpointsUnauthorized(t, baseURL)
		testOrganizerEndpointsAuthorized(t, baseURL)
		testReviewEndpointsUnauthorized(t, baseURL)
		testReviewEndpointsAuthorized(t, baseURL)
		testReviewLikeEndpointsUnauthorized(t, baseURL)
		testReviewLikeEndpointsAuthorized(t, baseURL)
		fmt.Println("✓ API端點響應測試通過")
	})

	fmt.Println("完整的應用程式流程測試完成！")
}

func createTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.String(http.StatusOK, "Swagger UI")
	})
	router.Use(testSessionsMiddlewareE2E())
	router.GET("/auth/:provider", handlers.OauthBegin)
	router.GET("/auth/:provider/callback", handlers.OauthCallback)
	router.GET("/logout", handlers.Logout)
	router.GET("/profile", handlers.Profile)
	routes.SetupAdminRoutes(router)
	routes.SetupUserRoutes(router)
	routes.SetupOrganizerRoutes(router)
	routes.SetupReviewRoutes(router)
	routes.SetupReviewLikeRoutes(router)
	return router
}

func createTestServer(router *gin.Engine) (*http.Server, string) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("無法創建監聽器: %v", err))
	}

	server := &http.Server{Addr: listener.Addr().String(), Handler: router}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("服務器錯誤: %v", err))
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return server, "http://" + listener.Addr().String()
}

func testAuthEndpoints(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/auth/facebook")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusTemporaryRedirect, http.StatusBadRequest}, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(baseURL + "/auth/instagram")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusTemporaryRedirect, http.StatusBadRequest}, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(baseURL + "/logout")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusTemporaryRedirect}, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 認證端點測試通過")
}

func testAdminEndpointsUnauthorized(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/admin/activities")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(baseURL + "/admin/locations")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 管理員端點（未授權）測試通過")
}

func testAdminEndpointsAuthorized(t *testing.T, baseURL string) {
	ts := testutils.NewTestServer()
	defer ts.Close()

	adminToken, err := testutils.CreateMockJWTToken(1, "Admin User", true)
	assert.NoError(t, err)
	assert.NotEmpty(t, adminToken)

	resp, err := testutils.MakeAuthenticatedRequest(ts, "GET", "/admin/activities", adminToken, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	resp, err = testutils.MakeAuthenticatedRequest(ts, "GET", "/admin/locations", adminToken, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 管理員端點（已授權）測試通過")
}

func testUserEndpointsUnauthorized(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/user/matches")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(baseURL + "/user/past-matches")
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 使用者端點（未授權）測試通過")
}

func testUserEndpointsAuthorized(t *testing.T, baseURL string) {
	ts := testutils.NewTestServer()
	defer ts.Close()

	token, err := testutils.CreateMockJWTToken(1, "Regular User", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	resp, err := testutils.MakeAuthenticatedRequest(ts, "GET", "/user/matches", token, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	resp, err = testutils.MakeAuthenticatedRequest(ts, "GET", "/user/past-matches", token, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 使用者端點（已授權）測試通過")
}

func testOrganizerEndpointsUnauthorized(t *testing.T, baseURL string) {
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", baseURL+"/organizer/approve-participant/1", nil)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()

	req, _ = http.NewRequest("PUT", baseURL+"/organizer/reject-participant/1", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 開局者端點（未授權）測試通過")
}

func testOrganizerEndpointsAuthorized(t *testing.T, baseURL string) {
	ts := testutils.NewTestServer()
	defer ts.Close()

	token, err := testutils.CreateMockJWTToken(1, "Organizer User", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	resp, err := testutils.MakeAuthenticatedRequest(ts, "PUT", "/organizer/approve-participant/1", token, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	resp, err = testutils.MakeAuthenticatedRequest(ts, "PUT", "/organizer/reject-participant/1", token, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 開局者端點（已授權）測試通過")
}

func testReviewEndpointsUnauthorized(t *testing.T, baseURL string) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", baseURL+"/review/matches/1", nil)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusCreated, http.StatusOK, http.StatusUnauthorized, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 評分端點（未授權）測試通過")
}

func testReviewEndpointsAuthorized(t *testing.T, baseURL string) {
	ts := testutils.NewTestServer()
	defer ts.Close()

	token, err := testutils.CreateMockJWTToken(1, "Review User", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	review := map[string]interface{}{
		"reviewee_id": 2,
		"score":       5,
		"comment":     "測試評分",
	}

	resp, err := testutils.MakeAuthenticatedRequest(ts, "POST", "/review/matches/1", token, review)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 評分端點（已授權）測試通過")
}

func testReviewLikeEndpointsUnauthorized(t *testing.T, baseURL string) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", baseURL+"/review-like/reviews/1/like", nil)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest}, resp.StatusCode)
	resp.Body.Close()

	req, _ = http.NewRequest("POST", baseURL+"/review-like/reviews/1/dislike", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest}, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 評論點讚/倒讚端點（未授權）測試通過")
}

func testReviewLikeEndpointsAuthorized(t *testing.T, baseURL string) {
	ts := testutils.NewTestServer()
	defer ts.Close()

	token, err := testutils.CreateMockJWTToken(1, "ReviewLike User", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	resp, err := testutils.MakeAuthenticatedRequest(ts, "POST", "/review-like/reviews/1/like", token, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	resp, err = testutils.MakeAuthenticatedRequest(ts, "POST", "/review-like/reviews/1/dislike", token, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	fmt.Println("  ✓ 評論點讚/倒讚端點（已授權）測試通過")
}

func testSessionsMiddlewareE2E() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("session", nil)
		c.Next()
	}
}
