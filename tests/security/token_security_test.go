package security

import (
	"testing"

	"free2free/tests/testutils"
	"github.com/stretchr/testify/assert"
)

// TestJWTSecurity 驗證 JWT 實作的安全性方面
func TestJWTSecurity(t *testing.T) {
	authHelper := testutils.NewAuthTestHelper()

	t.Run("Token 過期驗證", func(t *testing.T) {
		// 建立一個過期的 token
		expiredToken, err := authHelper.CreateExpiredUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// 驗證 token 被正確拒絕
		_, err = testutils.ValidateToken(expiredToken, authHelper.Secret)
		assert.Error(t, err)
	})

	t.Run("Token 簽章驗證", func(t *testing.T) {
		// 建立一個有效的 token
		token, err := authHelper.CreateValidUserToken(2, "user2@example.com", "Test User 2", "facebook")
		assert.NoError(t, err)

		// 嘗試用錯誤的 secret 驗證
		_, err = testutils.ValidateToken(token, "wrong-secret")
		assert.Error(t, err)
	})

	t.Run("Token 篡改檢測", func(t *testing.T) {
		// 建立一個有效的 token
		token, err := authHelper.CreateValidUserToken(5, "user5@example.com", "Test User 5", "facebook")
		assert.NoError(t, err)

		// 步驟 1：驗證原始 token 是有效的
		_, err = testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err, "原始 token 應該是有效的")

		// 步驟 2：驗證原始 token 簽章是有效的
		isValid, err := testutils.ValidateTokenSignature(token, authHelper.Secret)
		assert.NoError(t, err)
		assert.True(t, isValid, "原始 token 簽章應該是有效的")

		// 步驟 3：篡改 token（修改 payload 中的 user_id）
		tamperedToken, err := testutils.TamperWithJWTToken(token, map[string]interface{}{
			"user_id": 999, // 改成不同的使用者 ID
		})
		assert.NoError(t, err, "應該成功篡改 token")

		// 步驟 4：驗證原始 token 仍然有效（我們沒有修改它）
		_, err = testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err, "原始 token 應該仍然有效")

		// 步驟 5：驗證篡改後的 token 簽章是無效的（簽章不匹配）
		isValid, err = testutils.ValidateTokenSignature(tamperedToken, authHelper.Secret)
		assert.NoError(t, err)
		assert.False(t, isValid, "篡改後的 token 簽章應該是無效的")
	})

	t.Run("基於角色的存取驗證", func(t *testing.T) {
		// 建立不同角色的 token
		userToken, err := authHelper.CreateValidUserToken(3, "user3@example.com", "Regular User", "facebook")
		assert.NoError(t, err)

		adminToken, err := authHelper.CreateValidAdminToken(4, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// 驗證角色提取
		userClaims, err := testutils.ValidateToken(userToken, authHelper.Secret)
		assert.NoError(t, err)
		assert.Equal(t, "user", userClaims["role"])

		adminClaims, err := testutils.ValidateToken(adminToken, authHelper.Secret)
		assert.NoError(t, err)
		assert.Equal(t, "admin", adminClaims["role"])
	})
}

// TestOAuthSecurity 驗證 OAuth 實作的安全性方面
func TestOAuthSecurity(t *testing.T) {
	// 設定 mock OAuth provider
	mockProvider := testutils.NewMockAuthProvider()

	t.Run("拒絕無效的 OAuth 代碼", func(t *testing.T) {
		// 嘗試驗證無效的 OAuth 代碼
		_, valid := mockProvider.ValidateAuthCode("invalid-code")
		assert.False(t, valid)
	})

	t.Run("防止 OAuth 代碼重複使用", func(t *testing.T) {
		// 建立一個 mock 使用者
		mockUser := testutils.MockUser{
			ID:       "123456",
			Email:    "test@example.com",
			Name:     "Test User",
			Provider: "facebook",
		}

		// 新增一個有效的 auth code
		authCode := "test-auth-code"
		mockProvider.AddValidAuthCode(authCode, mockUser)

		// 第一次驗證應該成功
		returnedUser, valid := mockProvider.ValidateAuthCode(authCode)
		assert.True(t, valid)
		assert.Equal(t, mockUser, returnedUser)

		// MockAuthProvider 應該在第一次驗證時就刪除代碼
		// 驗證它不再存在於 ValidAuthCodes 中
		_, exists := mockProvider.ValidAuthCodes[authCode]
		assert.False(t, exists, "Auth code 應該在第一次驗證後被移除")
	})

	t.Run("State 參數驗證", func(t *testing.T) {
		mockProvider := testutils.NewMockAuthProvider()

		// 步驟 1：生成並儲存 state 參數
		state := mockProvider.StateManager.GenerateState()
		assert.NotEmpty(t, state, "State 應該被生成")
		assert.True(t, mockProvider.StateManager.ValidStates[state], "State 應該儲存在有效狀態中")

		// 步驟 2：驗證正確的 state（應該成功）
		isValid := mockProvider.ValidateState(state)
		assert.True(t, isValid, "有效的 state 應該被接受")
		_, exists := mockProvider.StateManager.ValidStates[state]
		assert.False(t, exists, "State 應該在驗證後被消耗（一次性使用）")

		// 步驟 3：生成另一個 state 進行 CSRF 保護測試
		state2 := mockProvider.StateManager.GenerateState()
		assert.NotEqual(t, state, state2, "每個 state 應該是唯一的")

		// 步驟 4：驗證第二個 state（第一次使用應該成功）
		isValid = mockProvider.ValidateState(state2)
		assert.True(t, isValid, "第二個 state 應該在第一次使用時被接受")

		// 步驟 5：嘗試重複使用第二個 state（應該失敗 - CSRF 保護）
		isValid = mockProvider.ValidateState(state2)
		assert.False(t, isValid, "State 不應該可重複使用（防止 CSRF 攻擊）")

		// 步驟 6：嘗試驗證無效的 state（應該失敗）
		invalidState := "invalid-state-12345"
		isValid = mockProvider.ValidateState(invalidState)
		assert.False(t, isValid, "無效的 state 應該被拒絕")
	})

	t.Run("PKCE 驗證", func(t *testing.T) {
		// 設定帶有 PKCE 的 mock OAuth provider
		mockProvider := testutils.NewMockAuthProvider()

		// 步驟 1：生成 code verifier 和 challenge
		verifier, challenge, err := mockProvider.PKCEManager.GenerateCodeChallenge()
		assert.NoError(t, err, "應該成功生成 code challenge")
		assert.NotEmpty(t, verifier, "Code verifier 不應該為空")
		assert.NotEmpty(t, challenge, "Code challenge 不應該為空")
		assert.NotEqual(t, verifier, challenge, "Verifier 和 challenge 應該不同")

		// 步驟 2：新增一個帶有 PKCE verifier 的有效 auth code
		authCode := "pkce-auth-code-123"
		mockProvider.PKCEManager.AddValidAuthCodeWithPKCE(authCode, verifier)
		assert.NotEmpty(t, mockProvider.PKCEManager.Codes[authCode], "帶有 verifier 的 auth code 應該被儲存")

		// 步驟 3：驗證 verifier 在驗證前已經被儲存
		_, exists := mockProvider.PKCEManager.Codes[authCode]
		assert.True(t, exists, "Code verifier 應該在驗證前存在於儲存中")

		// 步驟 4：驗證正確的 verifier（應該成功並消耗它）
		isValid := mockProvider.PKCEManager.ValidateCodeVerifier(verifier)
		assert.True(t, isValid, "有效的 code verifier 應該被接受")

		// 步驟 5：驗證 verifier 現在被消耗了（一次性使用）
		_, exists = mockProvider.PKCEManager.Codes[authCode]
		assert.False(t, exists, "Code verifier 應該在驗證後被消耗")

		// 步驟 6：生成另一個 verifier 進行重放攻擊測試
		verifier2, _, err := mockProvider.PKCEManager.GenerateCodeChallenge()
		assert.NoError(t, err)

		authCode2 := "pkce-auth-code-456"
		mockProvider.PKCEManager.AddValidAuthCodeWithPKCE(authCode2, verifier2)

		// 步驟 7：驗證正確的 verifier2 仍然有效（尚未使用）
		isValid = mockProvider.PKCEManager.ValidateCodeVerifier(verifier2)
		assert.True(t, isValid, "第二個有效的 verifier 應該被接受")

		// 步驟 8：嘗試再次驗證 verifier2（應該失敗 - 一次性使用）
		isValid = mockProvider.PKCEManager.ValidateCodeVerifier(verifier2)
		assert.False(t, isValid, "Code verifier2 不應該可重複使用")

		// 步驟 9：嘗試驗證不正確的 verifier（應該失敗）
		invalidVerifier := "invalid-verifier-xyz"
		isValid = mockProvider.PKCEManager.ValidateCodeVerifier(invalidVerifier)
		assert.False(t, isValid, "無效的 code verifier 應該被拒絕")

		// 步驟 10：驗證 verifier2 仍然被消耗
		_, exists = mockProvider.PKCEManager.Codes[authCode2]
		assert.False(t, exists, "Code verifier2 應該在使用後仍然被消耗")
	})
}

// TestTokenLeakagePrevention 確保敏感資訊不會洩漏
func TestTokenLeakagePrevention(t *testing.T) {
	authHelper := testutils.NewAuthTestHelper()

	t.Run("Token 中不包含敏感資料", func(t *testing.T) {
		token, err := authHelper.CreateValidUserToken(5, "user5@example.com", "Test User 5", "facebook")
		assert.NoError(t, err)

		claims, err := testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err)

		// 驗證敏感資料不在 token 中
		assert.NotContains(t, claims, "password")
		assert.NotContains(t, claims, "credit_card")
		assert.NotContains(t, claims, "secret_key")

		// 驗證預期的資料存在
		assert.Contains(t, claims, "user_id")
		assert.Contains(t, claims, "email")
		assert.Contains(t, claims, "role")
		assert.Contains(t, claims, "exp")
		assert.Contains(t, claims, "iat")
	})

	t.Run("錯誤訊息不洩露敏感資訊", func(t *testing.T) {
		// 使用無效的 secret 建立 token
		_, err := testutils.ValidateToken("invalid.token.format", "some-secret")
		assert.Error(t, err)

		// 驗證錯誤訊息不洩露內部細節
		errMsg := err.Error()
		assert.NotContains(t, errMsg, "secret")
		assert.NotContains(t, errMsg, "internal")
		assert.NotContains(t, errMsg, "database")
	})
}
