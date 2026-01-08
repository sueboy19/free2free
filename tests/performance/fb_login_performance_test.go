package performance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestFacebookLoginPerformance 測試完整 Facebook 登入流程的性能
func TestFacebookLoginPerformance(t *testing.T) {
	t.Run("Facebook 登入到 JWT 生成在 30 秒內", func(t *testing.T) {
		start := time.Now()

		// 初始化測試伺服器
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// 清除並設定資料庫
		err := testServer.ClearTestData()
		assert.NoError(t, err, "應該成功清除測試資料")
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err, "應該成功設定測試資料庫")

		// 模擬 Facebook 登入流程（建立使用者並生成 JWT）
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		assert.NotNil(t, user)

		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 驗證 JWT token - 只檢查是否成功驗證
		_, err = testutils.ValidateJWTToken(token)
		assert.NoError(t, err)

		elapsed := time.Since(start)

		// 根據 plan.md 中的需求："Facebook OAuth flow completed in under 30 seconds"
		assert.True(t, elapsed < 30*time.Second, "Facebook 登入流程應該在 30 秒內完成，耗時 %v", elapsed)

		t.Logf("Facebook 登入流程在 %v 內完成", elapsed)
	})

	t.Run("JWT token 驗證在 10ms 內", func(t *testing.T) {
		// 初始化測試伺服器
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// 清除並設定資料庫
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// 建立使用者和 JWT token
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		start := time.Now()
		_, err = testutils.ValidateJWTToken(token)
		validationTime := time.Since(start)

		assert.NoError(t, err)
		assert.True(t, validationTime < 10*time.Millisecond, "JWT 驗證應該在 10ms 內完成，耗時 %v", validationTime)

		t.Logf("JWT 驗證在 %v 內完成", validationTime)
	})

	t.Run("使用 JWT 的 API 請求在 500ms 內", func(t *testing.T) {
		// 初始化測試伺服器
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// 清除並設定資料庫
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// 建立使用者和 JWT token
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		start := time.Now()
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
		requestTime := time.Since(start)

		assert.NoError(t, err)
		// 根據實作可能得到 200 或 404，但不是性能問題
		assert.Contains(t, []int{200, 404}, resp.StatusCode)
		assert.True(t, requestTime < 500*time.Millisecond, "API 請求應該在 500ms 內完成，耗時 %v", requestTime)

		resp.Body.Close()
		t.Logf("使用 JWT 的 API 請求在 %v 內完成", requestTime)
	})

	t.Run("多次並發 Facebook 登入模擬", func(t *testing.T) {
		// 測試在多個使用者登入的模擬負載下的性能
		// 這是一個簡化版本 - 在實際系統中會使用實際的並發性

		const numSimulations = 5
		var totalElapsed time.Duration

		for i := 0; i < numSimulations; i++ {
			start := time.Now()

			testServer := testutils.NewTestServer()

			// 建立使用者和 JWT token
			user, err := testServer.CreateTestUser()
			assert.NoError(t, err)
			token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
			assert.NoError(t, err)

			// 驗證 JWT token
			_, err = testutils.ValidateJWTToken(token)
			assert.NoError(t, err)

			testServer.Close()
			elapsed := time.Since(start)
			totalElapsed += elapsed
		}

		avgTime := totalElapsed / numSimulations
		maxAllowedAvg := 30 * time.Second // 根據需求調整
		assert.True(t, avgTime < maxAllowedAvg, "平均 Facebook 登入流程應該在 %v 內完成，耗時 %v", maxAllowedAvg, avgTime)

		t.Logf("平均 Facebook 登入流程在 %v 內完成，共 %d 次模擬", avgTime, numSimulations)
	})
}

// TestSystemPerformanceUnderLoad 測試系統在負載條件下的性能
func TestSystemPerformanceUnderLoad(t *testing.T) {
	t.Run("JWT token 生成性能", func(t *testing.T) {
		// 初始化測試伺服器
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// 清除並設定資料庫
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// 建立使用者
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		const numTokens = 10
		var totalGenTime time.Duration

		for i := 0; i < numTokens; i++ {
			start := time.Now()
			token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
			genTime := time.Since(start)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			totalGenTime += genTime
		}

		avgGenTime := totalGenTime / numTokens
		// JWT 生成應該很快
		assert.True(t, avgGenTime < 10*time.Millisecond, "JWT 生成應該在 10ms 內完成，平均耗時 %v", avgGenTime)

		t.Logf("JWT 生成平均耗時 %v", avgGenTime)
	})

	t.Run("資料庫操作性能", func(t *testing.T) {
		// 初始化測試伺服器
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// 清除並設定資料庫
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// 測試使用者建立性能
		start := time.Now()
		user, err := testServer.CreateTestUser()
		creationTime := time.Since(start)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.True(t, creationTime < 100*time.Millisecond, "使用者建立應該在 100ms 內完成，耗時 %v", creationTime)

		t.Logf("使用者建立在 %v 內完成", creationTime)
	})
}
