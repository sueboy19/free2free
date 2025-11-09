package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	"free2free/database"
	"free2free/handlers"
	"free2free/models"
	"free2free/routes"

	middlewarepkg "free2free/middleware"

	_ "free2free/docs" // 这里需要导入你项目的文档包
	
	// Use modernc.org/sqlite as the underlying driver (no CGO required)
	_ "modernc.org/sqlite"
)

// 声明全局变量
var (
	store *sessions.CookieStore
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

	// 設定全局 DB instance
	database.GlobalDB = &database.ActualGormDB{Conn: gormDB}

	var migrateOn, _ = strconv.ParseBool(os.Getenv("AUTO_MIGRATE"))
	if migrateOn {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// 自動遷移所有資料表
		if err := database.GlobalDB.Conn.WithContext(ctx).AutoMigrate(
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
	
	// Set the store in handlers package
	handlers.SetStore(store)
}

// sessionsMiddleware 将 session 存储在 context 中
func sessionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or create a session
		session, err := store.Get(c.Request, "free2free-session")
		if err != nil {
			// If there's an error getting the session, log it but don't panic
			// Just create a new empty session
			session, _ = store.New(c.Request, "free2free-session")
		}

		// Make sure session is never nil
		if session == nil {
			session, _ = store.New(c.Request, "free2free-session")
		}

		// Set the session in the context
		c.Set("session", session)
		
		// Continue with the request
		c.Next()
		
		// Save the session if it was modified
		if session != nil && session.Options != nil {
			err := store.Save(c.Request, c.Writer, session)
			if err != nil {
				// Log the error but don't fail the request
				fmt.Printf("Error saving session: %v\n", err)
			}
		}
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
// @description 輸入 'Bearer <JWT token>' 進行認證。可以先透過 Facebook 登入取得 token。
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
	r.GET("/auth/:provider", handlers.OauthBegin)
	r.GET("/auth/:provider/callback", handlers.OauthCallback)
	r.GET("/logout", handlers.Logout)

	// JWT token 交換路由
	r.GET("/auth/token", handlers.ExchangeToken)

	// Refresh token 路由
	r.POST("/auth/refresh", handlers.RefreshTokenHandler)

	// 受保護的路由範例
	r.GET("/profile", handlers.Profile)

	// 設定管理後台路由
	routes.SetupAdminRoutes(r)

	// 設定使用者路由
	routes.SetupUserRoutes(r)

	// 設定開局者路由
	routes.SetupOrganizerRoutes(r)

	// 設定評分路由
	routes.SetupReviewRoutes(r)

	// 設定評論點讚/倒讚路由
	routes.SetupReviewLikeRoutes(r)

	// 啟動伺服器
	r.Run(":8080")
}


