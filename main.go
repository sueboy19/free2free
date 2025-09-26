package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/instagram"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "free2free/docs" // 这里需要导入你项目的文档包
)

// 声明全局变量
var (
	db           *gorm.DB
	store        *sessions.CookieStore
	adminDB      *gorm.DB
	userDB       *gorm.DB
	organizerDB  *gorm.DB
	reviewDB     *gorm.DB
	reviewLikeDB *gorm.DB
)

// User 代表使用者資料結構
type User struct {
	ID             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SocialID       string `gorm:"uniqueIndex:social_provider" json:"social_id"`
	SocialProvider string `gorm:"uniqueIndex:social_provider" json:"social_provider"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
	IsAdmin        bool   `json:"is_admin"` // 添加管理員標記
	CreatedAt      int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// JwtClaims 自定义JWT声明结构
type JwtClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	IsAdmin  bool   `json:"is_admin"` // 添加管理員標記
	jwt.RegisteredClaims
}

func init() {
	// 載入 .env 檔案
	if err := godotenv.Load(); err != nil {
		log.Println("無法載入 .env 檔案，使用環境變數")
	}

	// 初始化資料庫連線
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("資料庫連線失敗:", err)
	}

	// 将db变量赋值给其他文件中的db变量
	adminDB = db
	userDB = db
	organizerDB = db
	reviewDB = db
	reviewLikeDB = db

	// 自動遷移所有資料表
	err = db.AutoMigrate(
		&User{},
		&Admin{},
		&Location{},
		&Activity{},
		&Match{},
		&MatchParticipant{},
		&Review{},
		&ReviewLike{},
	)
	if err != nil {
		log.Fatal("資料表遷移失敗:", err)
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

	gothic.Store = store
}

// generateJWT 生成JWT token
func generateJWT(user *User) (string, error) {
	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET 环境变量未设置")
	}

	// 创建声明
	claims := &JwtClaims{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  user.IsAdmin, // 包含管理員標記
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	return token.SignedString([]byte(jwtSecret))
}

// validateJWT 验证JWT token
func validateJWT(tokenString string) (*JwtClaims, error) {
	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET 环境变量未设置")
	}

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的token")
}

// sessionsMiddleware 将 session 存储在 context 中
func sessionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "free2free-session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "无法获取 session"})
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

	// 添加Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 設定 session middleware
	r.Use(sessionsMiddleware())

	// OAuth 認證路由
	r.GET("/auth/:provider", oauthBegin)
	r.GET("/auth/:provider/callback", oauthCallback)

	// 登出路由
	r.GET("/logout", logout)

	// JWT token 交換路由
	r.GET("/auth/token", exchangeToken)

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
// @Failure 500 {object} map[string]string "OAuth 回調錯誤"
// @Router /auth/{provider}/callback [get]
func oauthCallback(c *gin.Context) {
	// 使用 gothic 取得使用者資訊
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 儲存或更新使用者資訊到資料庫
	dbUser, err := saveOrUpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "儲存使用者資訊失敗"})
		return
	}

	// 將使用者資訊存入 session
	session := c.MustGet("session").(*sessions.Session)
	session.Values["user_id"] = dbUser.ID
	session.Values["user_name"] = dbUser.Name
	session.Save(c.Request, c.Writer)

	// 生成 JWT token
	tokenString, err := generateJWTToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 token 失敗"})
		return
	}

	// 返回使用者資訊和 token
	c.JSON(http.StatusOK, gin.H{
		"user":  dbUser,
		"token": tokenString,
	})
}

// logout 處理登出
// @Summary 處理登出
// @Description 登出使用者並清除 session
// @Tags 認證
// @Accept json
// @Produce json
// @Success 302 {string} string "重定向到首頁"
// @Router /logout [get]
func logout(c *gin.Context) {
	session := c.MustGet("session").(*sessions.Session)
	session.Options.MaxAge = -1 // 刪除 session
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
// @Failure 401 {object} map[string]string "未登入"
// @Router /auth/token [get]
func exchangeToken(c *gin.Context) {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	// 生成 JWT token
	tokenString, err := generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 token 失敗"})
		return
	}

	// 返回 token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// profile 受保護的路由範例
// @Summary 取得使用者資訊
// @Description 取得使用者資訊 (支援 session 和 JWT token 認證)
// @Tags 使用者
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 401 {object} map[string]string "未登入"
// @Failure 500 {object} map[string]string "無法取得使用者資訊"
// @Router /profile [get]
// @Security ApiKeyAuth
func profile(c *gin.Context) {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// saveOrUpdateUser 儲存或更新使用者資訊
func saveOrUpdateUser(gothUser goth.User) (*User, error) {
	var user User

	// 檢查使用者是否已存在
	err := db.Where("social_id = ? AND social_provider = ?", gothUser.UserID, gothUser.Provider).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		// 查詢出錯
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		// 使用者不存在，建立新使用者
		user = User{
			SocialID:       gothUser.UserID,
			SocialProvider: gothUser.Provider,
			Name:           gothUser.Name,
			Email:          gothUser.Email,
			AvatarURL:      gothUser.AvatarURL,
		}

		// 儲存新使用者
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// 使用者已存在，更新資訊
		user.Name = gothUser.Name
		user.Email = gothUser.Email
		user.AvatarURL = gothUser.AvatarURL

		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// JWT claims struct
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// generateJWTToken 生成 JWT token
func generateJWTToken(user *User) (string, error) {
	// 設定 JWT token 過期時間 (24小時)
	expirationTime := time.Now().Add(24 * time.Hour)

	// 建立 claims
	claims := &Claims{
		UserID:   user.ID,
		UserName: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 簽署 token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWTToken 驗證 JWT token
func validateJWTToken(tokenString string) (*Claims, error) {
	// 解析 token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 驗證簽署方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// 檢查解析錯誤
	if err != nil {
		return nil, err
	}

	// 驗證 token 有效性
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// getAuthenticatedUser 從 context 中取得已認證的使用者
func getAuthenticatedUser(c *gin.Context) (*User, error) {
	// 首先嘗試從 session 取得使用者
	session := c.MustGet("session").(*sessions.Session)
	if userID, ok := session.Values["user_id"]; ok {
		// 從資料庫取得使用者資訊
		var user User
		err := db.First(&user, userID).Error
		if err != nil {
			return nil, err
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
	var user User
	err = db.First(&user, claims.UserID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
