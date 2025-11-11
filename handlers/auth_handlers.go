package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"free2free/database"
	"free2free/models"
	"free2free/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	apperrors "free2free/errors"

	"github.com/go-playground/validator/v10"
)

// Session store
var store *sessions.CookieStore

// getDB returns the global database connection
func getDB() *gorm.DB {
	if database.GlobalDB == nil || database.GlobalDB.Conn == nil {
		panic("Database not initialized. Call database initialization first.")
	}
	return database.GlobalDB.Conn
}

func SetStore(s *sessions.CookieStore) {
	store = s
}

// oauthBegin 開始 OAuth 流程
// @Summary 開始 OAuth 流程
// @Description 開始 Facebook 或 Instagram OAuth 流程
// @Tags 認證
// @Accept json
// @Produce json
// @Param provider path string true "OAuth 提供者 (facebook 或 instagram)"
// @Success 302 {string} string "重定向到 OAuth 提供者"
// @Failure 500 {object} ErrorResponse "OAuth 開始失敗"
// @Router /auth/{provider} [get]
func OauthBegin(c *gin.Context) {
	// 使用 gothic 來處理 OAuth 流程
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// oauthCallback OAuth 回調處理
// @Summary OAuth 回調處理
// @Description 處理 OAuth 提供者的回調
// @Tags 認證
// @Accept json
// @Produce json
// @Param provider path string true "OAuth 提供者 (facebook 或 instagram)"
// @Success 200 {object} map[string]interface{} "使用者資訊和 JWT token"
// @Failure 400 {object} ErrorResponse "無效的提供者"
// @Failure 500 {object} ErrorResponse "OAuth 回調錯誤"
// @Router /auth/{provider}/callback [get]
func OauthCallback(c *gin.Context) {
	// 使用 gothic 取得使用者資訊
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.Error(apperrors.NewOAuthError("OAuth authentication failed"))
		return
	}

	// 儲存或更新使用者資訊到資料庫
	dbUser, err := saveOrUpdateUser(user)
	if err != nil {
		c.Error(apperrors.NewInternalError("failed to save user information"))
		return
	}

	// Delete existing refresh tokens for this user (revoke old ones)
	getDB().Where("user_id = ?", dbUser.ID).Delete(&models.RefreshToken{})

	// 將使用者資訊存入 session
	session := c.MustGet("session").(*sessions.Session)
	session.Values["user_id"] = dbUser.ID
	session.Values["user_name"] = dbUser.Name
	session.Save(c.Request, c.Writer)

	// 生成 JWT tokens
	accessToken, _, hashedRefresh, err := GenerateTokens(dbUser)
	if err != nil {
		c.Error(apperrors.NewInternalError("token generation failed"))
		return
	}

	// 創建 RefreshToken 記錄
	refreshRecord := &models.RefreshToken{
		UserID:    uint(dbUser.ID),
		Token:     string(hashedRefresh),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := getDB().Create(refreshRecord).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 返回 HTML 頁面來處理 OAuth 回調
	c.Header("Content-Type", "text/html")

	// 序列化用戶資料
	userJSON, err := json.Marshal(dbUser)
	if err != nil {
		c.Error(apperrors.NewInternalError("Failed to serialize user data"))
		return
	}

	// 前端期望的是字符串 token，不是對象
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	  <title>Facebook 登入成功</title>
	  <script>
	      (function() {
	          var response = {
	              type: 'auth_success',
	              payload: {
	                  user: %s,
	                  token: "%s"
	              }
	          };
	          
	          if (window.opener) {
	              window.opener.postMessage(response, '*');
	          }
	          
	          setTimeout(function() {
	              window.close();
	          }, 1000);
	      })();
	  </script>
</head>
<body>
	  <p>登入成功，正在返回...</p>
</body>
</html>`, string(userJSON), accessToken)

	c.String(http.StatusOK, html)
}

// logout 處理登出
// @Summary 處理登出
// @Description 登出使用者並清除 session
// @Tags 認證
// @Accept json
// @Produce json
// @Success 302 {string} string "重定向到首頁"
// @Failure 500 {object} ErrorResponse "登出失敗"
// @Router /logout [get]
func Logout(c *gin.Context) {
	// Safely get the session to avoid panic
	sessionVal, exists := c.Get("session")
	if !exists {
		c.Error(apperrors.NewInternalError("session not found in context"))
		return
	}

	s, ok := sessionVal.(*sessions.Session)
	if !ok {
		c.Error(apperrors.NewInternalError("session type assertion failed"))
		return
	}

	userID, ok := s.Values["user_id"].(int64)
	if ok {
		// Delete refresh tokens for this user
		if err := getDB().Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
			c.Error(apperrors.MapGORMError(err))
			// Don't abort, just log
		}
	}

	// Clear session
	s.Options.MaxAge = 0
	s.Options.Path = "/"
	s.Save(c.Request, c.Writer)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// exchangeToken 交換 session for JWT token
// @Summary 交換 session for JWT token
// @Description 將現有的 session 交換為 JWT token
// @Tags 認證
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "JWT token"
// @Failure 401 {object} ErrorResponse "未登入"
// @Router /auth/token [get]
func ExchangeToken(c *gin.Context) {
	// 取得已認證的使用者
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewAuthenticationError("not logged in"))
		return
	}

	// 生成 JWT tokens
	accessToken, refreshToken, hashedRefresh, err := GenerateTokens(user)
	if err != nil {
		c.Error(apperrors.NewInternalError("token generation failed"))
		return
	}

	// 創建或更新 RefreshToken 記錄 (for exchange, assume create new)
	// First, delete existing for this user to rotate
	getDB().Where("user_id = ?", user.ID).Delete(&models.RefreshToken{})

	refreshRecord := &models.RefreshToken{
		UserID:    uint(user.ID),
		Token:     string(hashedRefresh),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := getDB().Create(refreshRecord).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 返回 tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    15 * 60,
	})
}

// Profile 受保護的路由範例
// @Summary 取得使用者資訊
// @Description 取得使用者資訊 (支援 session 和 JWT token 認證)
// @Tags 使用者
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 401 {object} ErrorResponse "未登入"
// @Failure 500 {object} ErrorResponse "無法取得使用者資訊"
// @Router /profile [get]
// @Security ApiKeyAuth
func Profile(c *gin.Context) {
	// 取得已認證的使用者
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewAuthenticationError("not logged in"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// saveOrUpdateUser 儲存或更新使用者資訊
func saveOrUpdateUser(gothUser goth.User) (*models.User, error) {
	var user models.User

	// 檢查使用者是否已存在
	err := getDB().Where("social_id = ? AND social_provider = ?", gothUser.UserID, gothUser.Provider).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		// 查詢出錯
		return nil, apperrors.MapGORMError(err)
	}

	v := validator.New()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 使用者不存在，建立新使用者
		now := time.Now().UnixNano() / 1e6 // 毫秒時間戳
		user = models.User{
			SocialID:       gothUser.UserID,
			SocialProvider: gothUser.Provider,
			Name:           gothUser.Name,
			Email:          gothUser.Email,
			AvatarURL:      gothUser.AvatarURL,
			IsAdmin:        false,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// Validate new user
		if err := v.Struct(&user); err != nil {
			return nil, apperrors.NewValidationError("Invalid user data from OAuth: " + err.Error())
		}

		// 儲存新使用者
		if err := getDB().Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, apperrors.MapGORMError(err)
			}
			return nil, apperrors.MapGORMError(err)
		}
	} else {
		// 使用者已存在，更新資訊
		now := time.Now().UnixNano() / 1e6 // 毫秒時間戳

		// 只更新需要的欄位，避免意外修改其他欄位
		updates := map[string]interface{}{
			"name":       gothUser.Name,
			"email":      gothUser.Email,
			"avatar_url": gothUser.AvatarURL,
			"updated_at": now,
		}

		if err := getDB().Model(&user).Updates(updates).Error; err != nil {
			return nil, apperrors.MapGORMError(err)
		}
	}

	return &user, nil
}

// JWT claims struct
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	IsAdmin  bool   `json:"is_admin"` // 添加管理員標記
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	HashedRefresh string // internal
}

// GenerateTokens 生成 access 和 refresh tokens
func GenerateTokens(user *models.User) (string, string, string, error) {
	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", "", "", fmt.Errorf("JWT_SECRET 环境变量未设置")
	}
	if len(jwtSecret) < 32 {
		return "", "", "", fmt.Errorf("JWT_SECRET 長度不足 32 byte")
	}

	// Access token claims - 15 min expiry
	accessClaims := &Claims{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", "", err
	}

	// Generate random refresh token
	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", "", err
	}
	refreshToken := base64.StdEncoding.EncodeToString(refreshBytes)

	// Hash the refresh token
	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", err
	}

	return accessString, refreshToken, string(hashedRefresh), nil
}

// RefreshTokenHandler 處理 refresh token
// @Summary Refresh access token
// @Description 使用 refresh token 獲取新的 access token 和 refresh token
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token request"
// @Success 200 {object} map[string]interface{} "新 tokens"
// @Failure 400 {object} ErrorResponse "無效的請求"
// @Failure 401 {object} ErrorResponse "無效的 refresh token"
// @Router /auth/refresh [post]
func RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewBadRequestError("invalid request data"))
		return
	}

	if req.RefreshToken == "" {
		c.Error(apperrors.NewBadRequestError("refresh token required"))
		return
	}

	// Load all active refresh tokens
	var records []models.RefreshToken
	if err := getDB().Where("expires_at > ?", time.Now()).Find(&records).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	var validRecord *models.RefreshToken
	for i := range records {
		if err := bcrypt.CompareHashAndPassword([]byte(records[i].Token), []byte(req.RefreshToken)); err == nil {
			validRecord = &records[i]
			break
		}
	}

	if validRecord == nil {
		c.Error(apperrors.NewAuthenticationError("invalid refresh token"))
		return
	}

	// Get user
	var user models.User
	if err := getDB().First(&user, validRecord.UserID).Error; err != nil {
		c.Error(apperrors.NewInternalError("user not found"))
		return
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, newHashedRefresh, err := GenerateTokens(&user)
	if err != nil {
		c.Error(apperrors.NewInternalError("token generation failed"))
		return
	}

	// Rotate: delete old
	if err := getDB().Delete(validRecord).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// Create new
	newRecord := &models.RefreshToken{
		UserID:    validRecord.UserID,
		Token:     newHashedRefresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := getDB().Create(newRecord).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"expires_in":    15 * 60,
	})
}

// ValidateJWTToken 驗證 JWT token
func ValidateJWTToken(tokenString string) (*Claims, error) {
	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET 环境变量未设置")
	}

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, apperrors.NewUnauthorizedError(err.Error())
	}

	// 验证token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.NewUnauthorizedError("無效的 token")
}
