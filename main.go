package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/instagram"
)

// User 代表使用者資料結構
type User struct {
	ID             int64  `db:"id"`
	SocialID       string `db:"social_id"`
	SocialProvider string `db:"social_provider"`
	Name           string `db:"name"`
	Email          string `db:"email"`
	AvatarURL      string `db:"avatar_url"`
}

var (
	db *sql.DB
	// 使用 Gorilla Sessions 管理 session
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

func init() {
	// 載入 .env 檔案
	if err := godotenv.Load(); err != nil {
		log.Println("無法載入 .env 檔案，使用環境變數")
	}

	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)

	// 初始化資料庫連線
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal("資料庫連線失敗:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("无法连接到数据库:", err)
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

	gothic.Store = store
}

func main() {
	// 設定 session 名稱
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return req.URL.Query().Get("provider"), nil
	}

	r := gin.Default()

	// 設定 session middleware
	r.Use(sessionsMiddleware())

	// OAuth 認證路由
	r.GET("/auth/:provider", oauthBegin)
	r.GET("/auth/:provider/callback", oauthCallback)

	// 登出路由
	r.GET("/logout", logout)

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

// sessionsMiddleware 處理 session
func sessionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c.Request, "free2free-session")
		c.Set("session", session)
		c.Next()
	}
}

// oauthBegin 開始 OAuth 流程
func oauthBegin(c *gin.Context) {
	// 使用 gothic 來處理 OAuth 流程
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// oauthCallback OAuth 回調處理
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

	// 重新導向到首頁或其他頁面
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// logout 處理登出
func logout(c *gin.Context) {
	session := c.MustGet("session").(*sessions.Session)
	session.Options.MaxAge = -1 // 刪除 session
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// profile 受保護的路由範例
func profile(c *gin.Context) {
	session := c.MustGet("session").(*sessions.Session)
	userID, ok := session.Values["user_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	// 從資料庫取得使用者資訊
	var user User
	err := db.QueryRow("SELECT id, social_id, social_provider, name, email, avatar_url FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.SocialID, &user.SocialProvider, &user.Name, &user.Email, &user.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得使用者資訊"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// saveOrUpdateUser 儲存或更新使用者資訊
func saveOrUpdateUser(gothUser goth.User) (*User, error) {
	var user User

	// 檢查使用者是否已存在
	err := db.QueryRow("SELECT id, social_id, social_provider, name, email, avatar_url FROM users WHERE social_id = ? AND social_provider = ?",
		gothUser.UserID, gothUser.Provider).Scan(&user.ID, &user.SocialID, &user.SocialProvider, &user.Name, &user.Email, &user.AvatarURL)

	if err != nil && err != sql.ErrNoRows {
		// 查詢出錯
		return nil, err
	}

	if err == sql.ErrNoRows {
		// 使用者不存在，建立新使用者
		result, err := db.Exec("INSERT INTO users (social_id, social_provider, name, email, avatar_url) VALUES (?, ?, ?, ?, ?)",
			gothUser.UserID, gothUser.Provider, gothUser.Name, gothUser.Email, gothUser.AvatarURL)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		user = User{
			ID:             id,
			SocialID:       gothUser.UserID,
			SocialProvider: gothUser.Provider,
			Name:           gothUser.Name,
			Email:          gothUser.Email,
			AvatarURL:      gothUser.AvatarURL,
		}
	} else {
		// 使用者已存在，更新資訊
		_, err := db.Exec("UPDATE users SET name = ?, email = ?, avatar_url = ? WHERE id = ?",
			gothUser.Name, gothUser.Email, gothUser.AvatarURL, user.ID)
		if err != nil {
			return nil, err
		}

		user.Name = gothUser.Name
		user.Email = gothUser.Email
		user.AvatarURL = gothUser.AvatarURL
	}

	return &user, nil
}
