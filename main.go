package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/instagram"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	apperrors "free2free/errors"
	middlewarepkg "free2free/middleware"
	"free2free/models"

	"github.com/go-playground/validator/v10"

	_ "free2free/docs" // 这里需要导入你项目的文档包
)

// 声明全局变量
var (
	store *sessions.CookieStore
)

var (
	db           DB
	adminDB      DB
	userDB       DB
	organizerDB  DB
	reviewDB     DB
	reviewLikeDB DB
)

func init() {
	// 載入 .env 檔案
	if err := godotenv.Load(); err != nil {
		log.Println("無法載入 .env 檔案，使用環境變數")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET 未設定")
	}

	// 初始化資料庫連線
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("資料庫連線失敗:", err)
	}

	// 将db变量赋值给其他文件中的db变量
	dbImpl := &dbImpl{conn: gormDB} // dbConn 是你原本 gorm.Open 回傳的 *gorm.DB
	db = dbImpl
	adminDB = dbImpl
	userDB = dbImpl
	organizerDB = dbImpl
	reviewDB = dbImpl
	reviewLikeDB = dbImpl

	var migrateOn, _ = strconv.ParseBool(os.Getenv("AUTO_MIGRATE"))
	if migrateOn {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// 自動遷移所有資料表
		if err := db.WithContext(ctx).AutoMigrate(
			&models.User{},
			&models.Admin{},
			&models.Location{},
			&models.Activity{},
			&models.Match{},
			&models.MatchParticipant{},
			&models.Review{},
			&models.ReviewLike{},
			&models.RefreshToken{},
		); err != nil {
			log.Fatal("資料表遷移失敗:", err)
		}
	}

	// 設定 OAuth 提供者
	goth.UseProviders(
		facebook.New(
			os.Getenv("FACEBOOK_KEY"),
			os.Getenv("FACEBOOK_SECRET"),
			fmt.Sprintf("%s/auth/facebook/callback", os.Getenv("BASE_URL")),
		),
		instagram.New(
			os.Getenv("INSTAGRAM_KEY"),
			os.Getenv("INSTAGRAM_SECRET"),
			fmt.Sprintf("%s/auth/instagram/callback", os.Getenv("BASE_URL")),
		),
	)

	// 初始化 session store，需要提供 auth key 和 encryption key
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY 环境变量未设置")
	}

	// 将 sessionKey 分为 auth key 和 encryption key
	var authKey, encryptionKey []byte
	if len(sessionKey) >= 32 {
		authKey = []byte(sessionKey[:32])
		if len(sessionKey) >= 64 {
			encryptionKey = []byte(sessionKey[32:64])
		} else {
			encryptionKey = []byte(sessionKey)
		}
	} else {
		// 如果 key 太短，重复以达到所需长度
		authKey = make([]byte, 32)
		encryptionKey = make([]byte, 32)
		for i := 0; i < 32; i++ {
			authKey[i] = sessionKey[i%len(sessionKey)]
			encryptionKey[i] = sessionKey[i%len(sessionKey)]
		}
	}

	store = sessions.NewCookieStore(authKey, encryptionKey)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   os.Getenv("SECURE_COOKIE") == "true",
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store
}

// sessionsMiddleware 将 session 存储在 context 中
func sessionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "free2free-session")
		if err != nil {
			c.Error(apperrors.NewAppError(http.StatusInternalServerError, "无法获取 session"))
			c.Abort()
			return
		}
		c.Set("session", session)
		c.Next()
	}
}

// @title 買一送一配對網站 API
// @version 1.0
// @description 這是一個買一送一配對網站的API文檔
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description 輸入 'Bearer &lt;JWT token&gt;' 進行認證。可以先透過 Facebook 登入取得 token。
func main() {
	// 設定 session 名稱
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		provider := req.URL.Query().Get("provider")
		if provider == "" {
			// 如果查詢參數中沒有提供者，嘗試從路徑中獲取
			// 這對於處理 /auth/facebook 這樣的路由很有用
			path := req.URL.Path
			parts := strings.Split(path, "/")
			if len(parts) >= 3 {
				provider = parts[2]
			}
		}
		return provider, nil
	}

	r := gin.Default()
	r.Use(cors.Default())
	// 生產環境請鎖域：
	// config := cors.Config{
	// 	AllowOrigins: []string{"https://yourdomain.com"},
	// 	AllowCredentials: true,
	// }
	// r.Use(cors.New(config))

	// 添加Swagger路由
	if os.Getenv("GIN_MODE") != "release" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 設定 session middleware
	r.Use(sessionsMiddleware())

	// 統一錯誤處理中間件
	r.Use(middlewarepkg.CustomRecovery())
	r.Use(middlewarepkg.ErrorHandler())

	// OAuth 認證路由
	r.GET("/auth/:provider", oauthBegin)
	r.GET("/auth/:provider/callback", oauthCallback)

	// 登出路由
	r.GET("/logout", logout)

	// JWT token 交換路由
	r.GET("/auth/token", exchangeToken)

	// Refresh token 路由
	r.POST("/auth/refresh", refreshTokenHandler)

	// 受保護的路由範例
	r.GET("/profile", profile)

	// 設定管理後台路由
	SetupAdminRoutes(r)

	// 設定使用者路由
	SetupUserRoutes(r)

	// 設定開局者路由
	SetupOrganizerRoutes(r)

	// 設定評分路由
	SetupReviewRoutes(r)

	// 設定評論點讚/倒讚路由
	SetupReviewLikeRoutes(r)

	// 啟動伺服器
	r.Run(":8080")
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
func oauthBegin(c *gin.Context) {
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
func oauthCallback(c *gin.Context) {
	// 使用 gothic 取得使用者資訊
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.Error(apperrors.NewAppError(http.StatusInternalServerError, err.Error()))
		return
	}

	// 儲存或更新使用者資訊到資料庫
	dbUser, err := saveOrUpdateUser(user)
	if err != nil {
		c.Error(apperrors.NewAppError(http.StatusInternalServerError, "儲存使用者資訊失敗"))
		return
	}

	// Delete existing refresh tokens for this user (revoke old ones)
	db.Where("user_id = ?", dbUser.ID).Delete(&models.RefreshToken{})

	// 將使用者資訊存入 session
	session := c.MustGet("session").(*sessions.Session)
	session.Values["user_id"] = dbUser.ID
	session.Values["user_name"] = dbUser.Name
	session.Save(c.Request, c.Writer)

	// 生成 JWT tokens
	accessToken, refreshToken, hashedRefresh, err := generateTokens(dbUser)
	if err != nil {
		c.Error(apperrors.NewAppError(http.StatusInternalServerError, "生成 token 失敗"))
		return
	}

	// 創建 RefreshToken 記錄
	refreshRecord := &models.RefreshToken{
		UserID:    uint(dbUser.ID),
		Token:     string(hashedRefresh),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := db.Create(refreshRecord).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 返回使用者資訊和 tokens
	c.JSON(http.StatusOK, gin.H{
		"user":          dbUser,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    15 * 60, // 15 minutes in seconds
	})
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
func logout(c *gin.Context) {
	session := c.MustGet("session").(*sessions.Session)
	userID, ok := session.Values["user_id"].(int64)
	if ok {
		// Delete refresh tokens for this user
		if err := db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
			c.Error(apperrors.MapGORMError(err))
			// Don't abort, just log
		}
	}

	// Clear session
	session.Options.MaxAge = 0
	session.Options.Path = "/"
	session.Save(c.Request, c.Writer)
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
func exchangeToken(c *gin.Context) {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	// 生成 JWT tokens
	accessToken, refreshToken, hashedRefresh, err := generateTokens(user)
	if err != nil {
		c.Error(apperrors.NewAppError(http.StatusInternalServerError, "生成 token 失敗"))
		return
	}

	// 創建或更新 RefreshToken 記錄 (for exchange, assume create new)
	// First, delete existing for this user to rotate
	db.Where("user_id = ?", user.ID).Delete(&models.RefreshToken{})

	refreshRecord := &models.RefreshToken{
		UserID:    uint(user.ID),
		Token:     string(hashedRefresh),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := db.Create(refreshRecord).Error; err != nil {
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

// profile 受保護的路由範例
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
func profile(c *gin.Context) {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// saveOrUpdateUser 儲存或更新使用者資訊
func saveOrUpdateUser(gothUser goth.User) (*models.User, error) {
	var user models.User

	// 檢查使用者是否已存在
	err := db.Where("social_id = ? AND social_provider = ?", gothUser.UserID, gothUser.Provider).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		// 查詢出錯
		return nil, apperrors.MapGORMError(err)
	}

	v := validator.New()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 使用者不存在，建立新使用者
		user = models.User{
			SocialID:       gothUser.UserID,
			SocialProvider: gothUser.Provider,
			Name:           gothUser.Name,
			Email:          gothUser.Email,
			AvatarURL:      gothUser.AvatarURL,
		}

		// Validate new user
		if err := v.Struct(&user); err != nil {
			return nil, apperrors.NewValidationError("Invalid user data from OAuth: " + err.Error())
		}

		// 儲存新使用者
		if err := db.Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, apperrors.MapGORMError(err)
			}
			return nil, apperrors.MapGORMError(err)
		}
	} else {
		// 使用者已存在，更新資訊
		updatedUser := models.User{
			ID:             user.ID,
			SocialID:       user.SocialID,
			SocialProvider: user.SocialProvider,
			Name:           gothUser.Name,
			Email:          gothUser.Email,
			AvatarURL:      gothUser.AvatarURL,
			IsAdmin:        user.IsAdmin,
		}

		// Validate updated user
		if err := v.Struct(&updatedUser); err != nil {
			return nil, apperrors.NewValidationError("Invalid update data from OAuth: " + err.Error())
		}

		user.Name = gothUser.Name
		user.Email = gothUser.Email
		user.AvatarURL = gothUser.AvatarURL

		if err := db.Save(&user).Error; err != nil {
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

// generateTokens 生成 access 和 refresh tokens
func generateTokens(user *models.User) (string, string, string, error) {
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

// refreshTokenHandler 處理 refresh token
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
func refreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	if req.RefreshToken == "" {
		c.Error(apperrors.NewUnauthorizedError("缺少 refresh token"))
		return
	}

	// Load all active refresh tokens
	var records []models.RefreshToken
	if err := db.Where("expires_at > ?", time.Now()).Find(&records).Error; err != nil {
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
		c.Error(apperrors.NewUnauthorizedError("無效的 refresh token"))
		return
	}

	// Get user
	var user models.User
	if err := db.First(&user, validRecord.UserID).Error; err != nil {
		c.Error(apperrors.NewAppError(http.StatusInternalServerError, "無法取得使用者"))
		return
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, newHashedRefresh, err := generateTokens(&user)
	if err != nil {
		c.Error(apperrors.NewAppError(http.StatusInternalServerError, "生成新 token 失敗"))
		return
	}

	// Rotate: delete old
	if err := db.Delete(validRecord).Error; err != nil {
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
	if err := db.Create(newRecord).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"expires_in":    15 * 60,
	})
}

// validateJWTToken 驗證 JWT token
func validateJWTToken(tokenString string) (*Claims, error) {
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

// getAuthenticatedUser 從 context 中取得已認證的使用者
var getAuthenticatedUser func(*gin.Context) (*models.User, error) = func(c *gin.Context) (*models.User, error) {
	// 首先嘗試從 session 取得使用者
	session := c.MustGet("session").(*sessions.Session)
	if userID, ok := session.Values["user_id"]; ok {
		// 從資料庫取得使用者資訊
		var user models.User
		err := db.First(&user, userID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewAppError(http.StatusNotFound, "user not found")
		}
		if err != nil {
			return nil, apperrors.MapGORMError(err)
		}
		return &user, nil
	}

	// 如果 session 中沒有使用者，嘗試從 JWT token 取得
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	// 檢查 Bearer token 格式
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	// 取得 token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 驗證 token
	claims, err := validateJWTToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 從資料庫取得使用者資訊
	var user models.User
	err = db.First(&user, claims.UserID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.NewAppError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return nil, apperrors.MapGORMError(err)
	}

	return &user, nil
}
