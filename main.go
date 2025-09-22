package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	CreatedAt      int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// JwtClaims 自定义JWT声明结构
type JwtClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// Admin 代表管理員
type Admin struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
}

// Location 代表地點
// @Description 地點資訊
type Location struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Activity 代表配對活動
// @Description 配對活動資訊
type Activity struct {
	ID          int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string   `json:"title"`
	TargetCount int      `json:"target_count"`
	LocationID  int64    `json:"location_id"`
	Description string   `json:"description"`
	CreatedBy   int64    `json:"created_by"`
	Location    Location `gorm:"foreignKey:LocationID" json:"location"`
}

// Match 代表配對局
// @Description 配對局資訊
type Match struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID  int64     `json:"activity_id"`
	OrganizerID int64     `json:"organizer_id"`
	MatchTime   time.Time `json:"match_time"`
	Status      string    `json:"status"` // open, closed, completed
	Activity    Activity  `gorm:"foreignKey:ActivityID" json:"activity"`
	Organizer   User      `gorm:"foreignKey:OrganizerID" json:"organizer"`
}

// MatchParticipant 代表配對參與者
// @Description 配對參與者資訊
type MatchParticipant struct {
	ID       int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID  int64     `json:"match_id"`
	UserID   int64     `json:"user_id"`
	Status   string    `json:"status"` // pending, approved, rejected
	JoinedAt time.Time `json:"joined_at"`
	Match    Match     `gorm:"foreignKey:MatchID" json:"match"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
}

// Review 代表評分與留言
// @Description 評分與留言資訊
type Review struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID    int64     `json:"match_id"`
	ReviewerID int64     `json:"reviewer_id"`
	RevieweeID int64     `json:"reviewee_id"`
	Score      int       `json:"score"` // 3-5分
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	Match      Match     `gorm:"foreignKey:MatchID" json:"match"`
	Reviewer   User      `gorm:"foreignKey:ReviewerID" json:"reviewer"`
	Reviewee   User      `gorm:"foreignKey:RevieweeID" json:"reviewee"`
}

// ReviewLike 代表評論點讚/倒讚
// @Description 評論點讚/倒讚資訊
type ReviewLike struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID int64  `json:"review_id"`
	UserID   int64  `json:"user_id"`
	IsLike   bool   `json:"is_like"` // true: 點讚, false: 倒讚
	Review   Review `gorm:"foreignKey:ReviewID" json:"review"`
	User     User   `gorm:"foreignKey:UserID" json:"user"`
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

// isAuthenticatedAdmin 檢查是否為已認證的管理員
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isAuthenticatedAdmin(c *gin.Context) bool {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		return false
	}

	// 這裡簡化處理，假設 ID 為 1 的使用者是管理員
	// 實際應用中應該檢查管理員表或使用者的管理員標記
	return user.ID == 1
}

// isAuthenticatedUser 檢查是否為已認證的使用者
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isAuthenticatedUser(c *gin.Context) bool {
	// 取得已認證的使用者
	_, err := getAuthenticatedUser(c)
	return err == nil
}

// isMatchOrganizer 檢查是否為指定配對局的開局者
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isMatchOrganizer(c *gin.Context, matchID int64) bool {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		return false
	}

	// 檢查配對局是否存在且開局者為當前使用者
	var match Match
	err = organizerDB.Where("id = ? AND organizer_id = ?", matchID, user.ID).First(&match).Error
	return err == nil
}

// canReviewMatch 檢查是否可以對指定配對局進行評分
// 這是一個簡化的實作，實際應用中需要檢查使用者是否參與了配對局
// 且配對局已結束但在評分時間範圍內
func canReviewMatch(c *gin.Context, matchID int64) bool {
	// 取得已認證的使用者
	user, err := getAuthenticatedUser(c)
	if err != nil {
		return false
	}

	// 檢查使用者是否參與了指定的配對局
	var participant MatchParticipant
	err = reviewDB.Where("match_id = ? AND user_id = ?", matchID, user.ID).First(&participant).Error
	if err != nil {
		return false
	}

	// 檢查配對局是否已完成且在評分時間範圍內
	var match Match
	err = reviewDB.Where("id = ? AND status = ?", matchID, "completed").First(&match).Error
	if err != nil {
		return false
	}

	// 這裡簡化處理，實際應用中應該檢查是否在評分時間範圍內 (結束後4小時)
	// 例如: match.EndTime.Add(4 * time.Hour).After(time.Now())
	return true
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

// AdminAuthMiddleware 管理員認證中介層
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作管理員認證邏輯
		// 例如檢查 session 或 JWT token 中的管理員身份
		// 為了簡化，這裡假設有一個 isAuthenticatedAdmin 函數
		if !isAuthenticatedAdmin(c) {
			c.JSON(401, gin.H{"error": "需要管理員權限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupAdminRoutes 設定管理後台路由
func SetupAdminRoutes(r *gin.Engine) {
	// 管理員認證路由組
	admin := r.Group("/admin")
	admin.Use(AdminAuthMiddleware())
	{
		// 配對活動管理
		admin.GET("/activities", listActivities)
		admin.POST("/activities", createActivity)
		admin.PUT("/activities/:id", updateActivity)
		admin.DELETE("/activities/:id", deleteActivity)

		// 地點管理
		admin.GET("/locations", listLocations)
		admin.POST("/locations", createLocation)
		admin.PUT("/locations/:id", updateLocation)
		admin.DELETE("/locations/:id", deleteLocation)
	}
}

// UserAuthMiddleware 使用者認證中介層
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作使用者認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者身份
		// 為了簡化，這裡假設有一個 isAuthenticatedUser 函數
		if !isAuthenticatedUser(c) {
			c.JSON(401, gin.H{"error": "需要使用者權限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupUserRoutes 設定使用者路由
func SetupUserRoutes(r *gin.Engine) {
	// 使用者認證路由組
	user := r.Group("/user")
	user.Use(UserAuthMiddleware())
	{
		// 配對列表
		user.GET("/matches", listMatches)

		// 開局功能
		user.POST("/matches", createMatch)

		// 參與配對
		user.POST("/matches/:id/join", joinMatch)

		// 過去參與列表
		user.GET("/past-matches", listPastMatches)
	}
}

// OrganizerAuthMiddleware 開局者認證中介層
func OrganizerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作開局者認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者是否為指定配對局的開局者
		// 為了簡化，這裡假設有一個 isMatchOrganizer 函數
		// matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		// if err != nil {
		// 	c.JSON(400, gin.H{"error": "無效的配對局 ID"})
		// 	c.Abort()
		// 	return
		// }

		// if !isMatchOrganizer(c, matchID) {
		// 	c.JSON(401, gin.H{"error": "需要開局者權限"})
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}

// SetupOrganizerRoutes 設定開局者路由
func SetupOrganizerRoutes(r *gin.Engine) {
	// 開局者認證路由組
	organizer := r.Group("/organizer")
	organizer.Use(UserAuthMiddleware())
	{
		// 審核參與者
		organizer.PUT("/matches/:id/participants/:participant_id/approve", OrganizerAuthMiddleware(), approveParticipant)
		organizer.PUT("/matches/:id/participants/:participant_id/reject", OrganizerAuthMiddleware(), rejectParticipant)
	}
}

// ReviewAuthMiddleware 評分認證中介層
func ReviewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作評分認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者是否參與了指定的配對局
		// 且配對局已結束但在評分時間範圍內 (結束後4小時)
		// 為了簡化，這裡假設有一個 canReviewMatch 函數
		// matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		// if err != nil {
		// 	c.JSON(400, gin.H{"error": "無效的配對局 ID"})
		// 	c.Abort()
		// 	return
		// }

		// if !canReviewMatch(c, matchID) {
		// 	c.JSON(401, gin.H{"error": "無評分權限"})
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}

// SetupReviewRoutes 設定評分路由
func SetupReviewRoutes(r *gin.Engine) {
	// 評分認證路由組
	review := r.Group("/review")
	review.Use(UserAuthMiddleware())
	{
		// 互相評分與留言
		review.POST("/matches/:id", ReviewAuthMiddleware(), createReview)
	}
}

// SetupReviewLikeRoutes 設定評論點讚/倒讚路由
func SetupReviewLikeRoutes(r *gin.Engine) {
	// 評論點讚/倒讚路由組
	reviewLike := r.Group("/review-like")
	reviewLike.Use(UserAuthMiddleware())
	{
		// 點讚評論
		reviewLike.POST("/reviews/:id/like", likeReview)
		// 倒讚評論
		reviewLike.POST("/reviews/:id/dislike", dislikeReview)
	}
}

// listActivities 取得配對活動列表
// @Summary 取得配對活動列表
// @Description 取得所有配對活動的列表
// @Tags 管理員
// @Accept json
// @Produce json
// @Success 200 {array} Activity
// @Failure 500 {object} map[string]string "無法取得活動列表"
// @Router /admin/activities [get]
// @Security ApiKeyAuth
func listActivities(c *gin.Context) {
	var activities []Activity
	if err := adminDB.Preload("Location").Order("id DESC").Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得活動列表"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// createActivity 建立新的配對活動
// @Summary 建立新的配對活動
// @Description 建立新的配對活動
// @Tags 管理員
// @Accept json
// @Produce json
// @Param activity body Activity true "配對活動資訊"
// @Success 201 {object} Activity
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法建立活動"
// @Router /admin/activities [post]
// @Security ApiKeyAuth
func createActivity(c *gin.Context) {
	var activity Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if activity.Title == "" || activity.TargetCount <= 0 || activity.LocationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "標題、目標人數和地點為必填欄位"})
		return
	}

	// 檢查地點是否存在
	var location Location
	if err := adminDB.First(&location, activity.LocationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的地點不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證地點"})
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	// 設定活動建立者為當前使用者
	activity.CreatedBy = user.ID

	if err := adminDB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立活動"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// updateActivity 更新配對活動
// @Summary 更新配對活動
// @Description 更新指定ID的配對活動
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "活動ID"
// @Param activity body Activity true "配對活動資訊"
// @Success 200 {object} Activity
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法更新活動"
// @Router /admin/activities/{id} [put]
// @Security ApiKeyAuth
func updateActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的活動 ID"})
		return
	}

	var activity Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if activity.Title == "" || activity.TargetCount <= 0 || activity.LocationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "標題、目標人數和地點為必填欄位"})
		return
	}

	// 檢查地點是否存在
	var location Location
	if err := adminDB.First(&location, activity.LocationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的地點不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證地點"})
		return
	}

	// 更新活動
	if err := adminDB.Model(&Activity{}).Where("id = ?", id).Updates(activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新活動"})
		return
	}

	activity.ID = id
	c.JSON(http.StatusOK, activity)
}

// deleteActivity 刪除配對活動
// @Summary 刪除配對活動
// @Description 刪除指定ID的配對活動
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "活動ID"
// @Success 200 {object} map[string]string "活動已刪除"
// @Failure 400 {object} map[string]string "無效的活動 ID"
// @Failure 500 {object} map[string]string "無法刪除活動"
// @Router /admin/activities/{id} [delete]
// @Security ApiKeyAuth
func deleteActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的活動 ID"})
		return
	}

	if err := adminDB.Delete(&Activity{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除活動"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活動已刪除"})
}

// listLocations 取得地點列表
// @Summary 取得地點列表
// @Description 取得所有地點的列表
// @Tags 管理員
// @Accept json
// @Produce json
// @Success 200 {array} Location
// @Failure 500 {object} map[string]string "無法取得地點列表"
// @Router /admin/locations [get]
// @Security ApiKeyAuth
func listLocations(c *gin.Context) {
	var locations []Location
	if err := adminDB.Order("id DESC").Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得地點列表"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

// createLocation 建立新的地點
// @Summary 建立新的地點
// @Description 建立新的地點
// @Tags 管理員
// @Accept json
// @Produce json
// @Param location body Location true "地點資訊"
// @Success 201 {object} Location
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法建立地點"
// @Router /admin/locations [post]
// @Security ApiKeyAuth
func createLocation(c *gin.Context) {
	var location Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if location.Name == "" || location.Address == "" || location.Latitude == 0 || location.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "名稱、地址和座標為必填欄位"})
		return
	}

	if err := adminDB.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立地點"})
		return
	}

	c.JSON(http.StatusCreated, location)
}

// updateLocation 更新地點
// @Summary 更新地點
// @Description 更新指定ID的地點
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "地點ID"
// @Param location body Location true "地點資訊"
// @Success 200 {object} Location
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法更新地點"
// @Router /admin/locations/{id} [put]
// @Security ApiKeyAuth
func updateLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的地點 ID"})
		return
	}

	var location Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if location.Name == "" || location.Address == "" || location.Latitude == 0 || location.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "名稱、地址和座標為必填欄位"})
		return
	}

	// 更新地點
	if err := adminDB.Model(&Location{}).Where("id = ?", id).Updates(location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新地點"})
		return
	}

	location.ID = id
	c.JSON(http.StatusOK, location)
}

// deleteLocation 刪除地點
// @Summary 刪除地點
// @Description 刪除指定ID的地點
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "地點ID"
// @Success 200 {object} map[string]string "地點已刪除"
// @Failure 400 {object} map[string]string "無效的地點 ID"
// @Failure 500 {object} map[string]string "無法刪除地點"
// @Router /admin/locations/{id} [delete]
// @Security ApiKeyAuth
func deleteLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的地點 ID"})
		return
	}

	if err := adminDB.Delete(&Location{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除地點"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "地點已刪除"})
}

// listMatches 取得時間未到的配對列表
// @Summary 取得時間未到的配對列表
// @Description 取得所有時間未到且狀態為open的配對列表
// @Tags 使用者
// @Accept json
// @Produce json
// @Success 200 {array} Match
// @Failure 500 {object} map[string]string "無法取得配對列表"
// @Router /user/matches [get]
// @Security ApiKeyAuth
func listMatches(c *gin.Context) {
	var matches []Match
	// 只顯示狀態為 open 且時間未到的配對
	if err := userDB.Preload("Activity").Preload("Organizer").Where("status = ? AND match_time > ?", "open", time.Now()).Order("match_time ASC").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得配對列表"})
		return
	}

	c.JSON(http.StatusOK, matches)
}

// createMatch 建立新的配對局 (開局)
// @Summary 建立新的配對局
// @Description 建立新的配對局 (開局)
// @Tags 使用者
// @Accept json
// @Produce json
// @Param match body Match true "配對局資訊"
// @Success 201 {object} Match
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法建立配對局"
// @Router /user/matches [post]
// @Security ApiKeyAuth
func createMatch(c *gin.Context) {
	var match Match
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if match.ActivityID <= 0 || match.MatchTime.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "活動 ID 和配對時間為必填欄位"})
		return
	}

	// 檢查活動是否存在
	var activity Activity
	if err := userDB.First(&activity, match.ActivityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的活動不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證活動"})
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	// 設定開局者為當前使用者
	match.OrganizerID = user.ID
	match.Status = "open"

	if err := userDB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立配對局"})
		return
	}

	// 預加載關聯資料
	userDB.Preload("Activity").Preload("Organizer").First(&match, match.ID)
	c.JSON(http.StatusCreated, match)
}

// joinMatch 參與配對
// @Summary 參與配對
// @Description 參與指定ID的配對局
// @Tags 使用者
// @Accept json
// @Produce json
// @Param id path int true "配對局ID"
// @Success 201 {object} MatchParticipant
// @Failure 400 {object} map[string]string "無效的配對局 ID 或已參與"
// @Failure 500 {object} map[string]string "無法參與配對局"
// @Router /user/matches/{id}/join [post]
// @Security ApiKeyAuth
func joinMatch(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
		return
	}

	// 檢查配對局是否存在且可參與
	var match Match
	if err := userDB.Where("id = ? AND status = ? AND match_time > ?", matchID, "open", time.Now()).First(&match).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的配對局不存在或已關閉"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證配對局"})
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	userID := user.ID

	// 檢查使用者是否已經參與此配對局
	var existingParticipant MatchParticipant
	err = userDB.Where("match_id = ? AND user_id = ?", matchID, userID).First(&existingParticipant).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查參與狀態"})
		return
	}

	// 如果已經參與，返回錯誤
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經參與此配對局"})
		return
	}

	// 建立新的參與記錄
	participant := MatchParticipant{
		MatchID:  matchID,
		UserID:   userID,
		Status:   "pending",
		JoinedAt: time.Now(),
	}

	if err := userDB.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法參與配對局"})
		return
	}

	// 預加載關聯資料
	userDB.Preload("Match").Preload("User").First(&participant, participant.ID)
	c.JSON(http.StatusCreated, participant)
}

// listPastMatches 取得過去參與的配對列表
// @Summary 取得過去參與的配對列表
// @Description 取得該使用者參與過的已完成的配對局列表
// @Tags 使用者
// @Accept json
// @Produce json
// @Success 200 {array} Match
// @Failure 500 {object} map[string]string "無法取得過去參與的配對列表"
// @Router /user/past-matches [get]
// @Security ApiKeyAuth
func listPastMatches(c *gin.Context) {
	// 從認證資訊取得使用者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	userID := user.ID

	var matches []Match
	// 取得該使用者參與過的已完成的配對局
	if err := userDB.Joins("JOIN match_participants mp ON matches.id = mp.match_id").
		Where("mp.user_id = ? AND matches.status = ?", userID, "completed").
		Order("matches.match_time DESC").
		Preload("Activity").
		Preload("Organizer").
		Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得過去參與的配對列表"})
		return
	}

	c.JSON(http.StatusOK, matches)
}

// approveParticipant 審核通過參與者
// @Summary 審核通過參與者
// @Description 開局者審核通過指定配對局的參與者
// @Tags 開局者
// @Accept json
// @Produce json
// @Param id path int true "配對局ID"
// @Param participant_id path int true "參與者ID"
// @Success 200 {object} MatchParticipant
// @Failure 400 {object} map[string]string "無效的配對局 ID 或參與者 ID"
// @Failure 500 {object} map[string]string "無法審核通過參與者"
// @Router /organizer/matches/{id}/participants/{participant_id}/approve [put]
// @Security ApiKeyAuth
func approveParticipant(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
		return
	}

	participantID, err := strconv.ParseInt(c.Param("participant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的參與者 ID"})
		return
	}

	// 檢查參與者是否屬於此配對局
	var participant MatchParticipant
	if err := organizerDB.Where("id = ? AND match_id = ?", participantID, matchID).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的參與者不存在或不屬於此配對局"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證參與者"})
		return
	}

	// 更新參與者狀態為 approved
	if err := organizerDB.Model(&participant).Update("status", "approved").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法審核通過參與者"})
		return
	}

	participant.Status = "approved"
	c.JSON(http.StatusOK, participant)
}

// rejectParticipant 審核拒絕參與者
// @Summary 審核拒絕參與者
// @Description 開局者審核拒絕指定配對局的參與者
// @Tags 開局者
// @Accept json
// @Produce json
// @Param id path int true "配對局ID"
// @Param participant_id path int true "參與者ID"
// @Success 200 {object} MatchParticipant
// @Failure 400 {object} map[string]string "無效的配對局 ID 或參與者 ID"
// @Failure 500 {object} map[string]string "無法審核拒絕參與者"
// @Router /organizer/matches/{id}/participants/{participant_id}/reject [put]
// @Security ApiKeyAuth
func rejectParticipant(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
		return
	}

	participantID, err := strconv.ParseInt(c.Param("participant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的參與者 ID"})
		return
	}

	// 檢查參與者是否屬於此配對局
	var participant MatchParticipant
	if err := organizerDB.Where("id = ? AND match_id = ?", participantID, matchID).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的參與者不存在或不屬於此配對局"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證參與者"})
		return
	}

	// 更新參與者狀態為 rejected
	if err := organizerDB.Model(&participant).Update("status", "rejected").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法審核拒絕參與者"})
		return
	}

	participant.Status = "rejected"
	c.JSON(http.StatusOK, participant)
}

// createReview 建立評分與留言
// @Summary 建立評分與留言
// @Description 建立對指定配對局中其他參與者的評分與留言
// @Tags 評分
// @Accept json
// @Produce json
// @Param id path int true "配對局ID"
// @Param review body Review true "評分與留言資訊"
// @Success 201 {object} Review
// @Failure 400 {object} map[string]string "無效的請求資料或已評分過"
// @Failure 500 {object} map[string]string "無法建立評分記錄"
// @Router /review/matches/{id} [post]
// @Security ApiKeyAuth
func createReview(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
		return
	}

	var review Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if review.RevieweeID <= 0 || review.Score < 3 || review.Score > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "被評分者和評分為必填欄位，評分範圍為3-5分"})
		return
	}

	// 從認證資訊取得評分者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	review.ReviewerID = user.ID
	review.MatchID = matchID
	review.CreatedAt = time.Now()

	// 檢查是否已經對此人在此配對局評分過
	var existingReview Review
	err = reviewDB.Where("reviewer_id = ? AND reviewee_id = ? AND match_id = ?",
		review.ReviewerID, review.RevieweeID, review.MatchID).First(&existingReview).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查評分記錄"})
		return
	}

	// 如果已經評分過，返回錯誤
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經對此人評分過"})
		return
	}

	// 建立新的評分記錄
	if err := reviewDB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立評分記錄"})
		return
	}

	// 預加載關聯資料
	reviewDB.Preload("Match").Preload("Reviewer").Preload("Reviewee").First(&review, review.ID)
	c.JSON(http.StatusCreated, review)
}

// likeReview 點讚評論
// @Summary 點讚評論
// @Description 對指定ID的評論進行點讚
// @Tags 評論
// @Accept json
// @Produce json
// @Param id path int true "評論ID"
// @Success 201 {object} ReviewLike
// @Failure 400 {object} map[string]string "無效的評論 ID 或已點讚"
// @Failure 500 {object} map[string]string "無法點讚評論"
// @Router /review-like/reviews/{id}/like [post]
// @Security ApiKeyAuth
func likeReview(c *gin.Context) {
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的評論 ID"})
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	userID := user.ID

	// 檢查評論是否存在
	var review Review
	if err := reviewLikeDB.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的評論不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證評論"})
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike ReviewLike
	err = reviewLikeDB.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查點讚狀態"})
		return
	}

	// 如果已經點讚，返回錯誤
	if err == nil && existingLike.IsLike {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經點讚此評論"})
		return
	}

	// 如果已經倒讚，則更新為點讚
	if err == nil && !existingLike.IsLike {
		if err := reviewLikeDB.Model(&existingLike).Update("is_like", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新點讚狀態"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "已將倒讚改為點讚"})
		return
	}

	// 建立新的點讚記錄
	reviewLike := ReviewLike{
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   true,
	}

	if err := reviewLikeDB.Create(&reviewLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法點讚評論"})
		return
	}

	// 預加載關聯資料
	reviewLikeDB.Preload("Review").Preload("User").First(&reviewLike, reviewLike.ID)
	c.JSON(http.StatusCreated, reviewLike)
}

// dislikeReview 倒讚評論
// @Summary 倒讚評論
// @Description 對指定ID的評論進行倒讚
// @Tags 評論
// @Accept json
// @Produce json
// @Param id path int true "評論ID"
// @Success 201 {object} ReviewLike
// @Failure 400 {object} map[string]string "無效的評論 ID 或已倒讚"
// @Failure 500 {object} map[string]string "無法倒讚評論"
// @Router /review-like/reviews/{id}/dislike [post]
// @Security ApiKeyAuth
func dislikeReview(c *gin.Context) {
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的評論 ID"})
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := getAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	userID := user.ID

	// 檢查評論是否存在
	var review Review
	if err := reviewLikeDB.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的評論不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證評論"})
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike ReviewLike
	err = reviewLikeDB.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查點讚狀態"})
		return
	}

	// 如果已經倒讚，返回錯誤
	if err == nil && !existingLike.IsLike {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經倒讚此評論"})
		return
	}

	// 如果已經點讚，則更新為倒讚
	if err == nil && existingLike.IsLike {
		if err := reviewLikeDB.Model(&existingLike).Update("is_like", false).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新點讚狀態"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "已將點讚改為倒讚"})
		return
	}

	// 建立新的倒讚記錄
	reviewLike := ReviewLike{
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   false,
	}

	if err := reviewLikeDB.Create(&reviewLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法倒讚評論"})
		return
	}

	// 預加載關聯資料
	reviewLikeDB.Preload("Review").Preload("User").First(&reviewLike, reviewLike.ID)
	c.JSON(http.StatusCreated, reviewLike)
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
