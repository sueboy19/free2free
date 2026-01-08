package performance

import (
	"net/http"
	"testing"
	"time"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOAuthFlowPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 設定測試環境
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("OAuth begin endpoint 回應時間", func(t *testing.T) {
		// 建立測試伺服器
		ts := testutils.NewTestServer()
		defer ts.Close()

		// 測量回應時間
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// 檢查回應時間在 500ms 以下，符合需求
		assert.Less(t, duration.Milliseconds(), int64(500),
			"OAuth begin endpoint 應該在 500ms 內回應，耗時 %dms", duration.Milliseconds())
	})

	t.Run("OAuth callback endpoint 回應時間", func(t *testing.T) {
		// 建立測試伺服器
		ts := testutils.NewTestServer()
		defer ts.Close()

		// 測量回應時間，使用有效的 code 參數
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/auth/facebook/callback?code=test-code", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// 檢查回應時間在 500ms 以下，符合需求
		assert.Less(t, duration.Milliseconds(), int64(500),
			"OAuth callback endpoint 應該在 500ms 內回應，耗時 %dms", duration.Milliseconds())
	})

	t.Run("Token 交換 endpoint 回應時間", func(t *testing.T) {
		// 建立測試伺服器
		ts := testutils.NewTestServer()
		defer ts.Close()

		// 測量回應時間
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/auth/token", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		// 接受 200（成功）和 401（未授權）進行性能測試
		assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp.StatusCode)

		// 檢查回應時間在 500ms 以下，符合需求
		assert.Less(t, duration.Milliseconds(), int64(500),
			"Token 交換 endpoint 應該在 500ms 內回應，耗時 %dms", duration.Milliseconds())
	})

	t.Run("Logout endpoint 回應時間", func(t *testing.T) {
		// 建立測試伺服器
		ts := testutils.NewTestServer()
		defer ts.Close()

		// 測量回應時間
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/logout", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		// Logout 在 mock 中返回 200，也接受重導向狀態
		assert.Condition(t, func() bool {
			return resp.StatusCode == http.StatusOK ||
				resp.StatusCode == http.StatusTemporaryRedirect ||
				resp.StatusCode == http.StatusFound
		}, "Logout 應該返回成功或重導向狀態")

		// 檢查回應時間在 500ms 以下，符合需求
		assert.Less(t, duration.Milliseconds(), int64(500),
			"Logout endpoint 應該在 500ms 內回應，耗時 %dms", duration.Milliseconds())
	})

	t.Run("OAuth 流程完成時間在 10 秒內", func(t *testing.T) {
		// 建立測試伺服器
		ts := testutils.NewTestServer()
		defer ts.Close()

		// 測量完整 OAuth 流程模擬的時間
		start := time.Now()

		// 呼叫 OAuth begin endpoint
		resp1, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		// 呼叫 OAuth callback endpoint，使用 code 參數
		resp2, err := ts.DoRequest("GET", "/auth/facebook/callback?code=test-code", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		// 呼叫 token 交換 endpoint（在沒有 session 的情況下會返回 401，這對於性能測試來說是可以接受的）
		resp3, err := ts.DoRequest("GET", "/auth/token", nil, nil)
		assert.NoError(t, err)
		assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp3.StatusCode)

		duration := time.Since(start)

		// 完整的 OAuth 流程應該在 10 秒內完成，符合需求
		assert.Less(t, duration.Seconds(), float64(10),
			"完整的 OAuth 流程應該在 10 秒內完成，耗時 %f 秒", duration.Seconds())
	})
}
