package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Match 代表配對局
type Match struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID int64     `json:"activity_id"`
	OrganizerID int64    `json:"organizer_id"`
	MatchTime  time.Time `json:"match_time"`
	Status     string    `json:"status"` // open, closed, completed
	Activity   Activity  `gorm:"foreignKey:ActivityID" json:"activity"`
	Organizer  User      `gorm:"foreignKey:OrganizerID" json:"organizer"`
}

// MatchParticipant 代表配對參與者
type MatchParticipant struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID   int64     `json:"match_id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"` // pending, approved, rejected
	JoinedAt  time.Time `json:"joined_at"`
	Match     Match     `gorm:"foreignKey:MatchID" json:"match"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// UserAuthMiddleware 使用者認證中介層
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作使用者認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者身份
		// 為了簡化，這裡假設有一個 isAuthenticatedUser 函數
		if !isAuthenticatedUser(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要使用者權限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// isAuthenticatedUser 檢查是否為已認證的使用者
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isAuthenticatedUser(c *gin.Context) bool {
	// 這裡應該實作實際的認證邏輯
	// 例如檢查 session 中是否有 user_id
	// 或解析 JWT token 驗證使用者身份
	// 為了示範，這裡暫時回傳 true
	return true
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

// listMatches 取得時間未到的配對列表
func listMatches(c *gin.Context) {
	var matches []Match
	// 只顯示狀態為 open 且時間未到的配對
	if err := db.Preload("Activity").Preload("Organizer").Where("status = ? AND match_time > ?", "open", time.Now()).Order("match_time ASC").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得配對列表"})
		return
	}

	c.JSON(http.StatusOK, matches)
}

// createMatch 建立新的配對局 (開局)
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
	if err := db.First(&activity, match.ActivityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的活動不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證活動"})
		return
	}

	// 這裡應該從 session 或 token 取得使用者 ID
	// 為了簡化，這裡暫時設為 1
	match.OrganizerID = 1
	match.Status = "open"

	if err := db.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立配對局"})
		return
	}

	// 預加載關聯資料
	db.Preload("Activity").Preload("Organizer").First(&match, match.ID)
	c.JSON(http.StatusCreated, match)
}

// joinMatch 參與配對
func joinMatch(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
		return
	}

	// 檢查配對局是否存在且可參與
	var match Match
	if err := db.Where("id = ? AND status = ? AND match_time > ?", matchID, "open", time.Now()).First(&match).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的配對局不存在或已關閉"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證配對局"})
		return
	}

	// 這裡應該從 session 或 token 取得使用者 ID
	// 為了簡化，這裡暫時設為 1
	userID := int64(1)

	// 檢查使用者是否已經參與此配對局
	var existingParticipant MatchParticipant
	err = db.Where("match_id = ? AND user_id = ?", matchID, userID).First(&existingParticipant).Error
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

	if err := db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法參與配對局"})
		return
	}

	// 預加載關聯資料
	db.Preload("Match").Preload("User").First(&participant, participant.ID)
	c.JSON(http.StatusCreated, participant)
}

// listPastMatches 取得過去參與的配對列表
func listPastMatches(c *gin.Context) {
	// 這裡應該從 session 或 token 取得使用者 ID
	// 為了簡化，這裡暫時設為 1
	userID := int64(1)

	var matches []Match
	// 取得該使用者參與過的已完成的配對局
	if err := db.Joins("JOIN match_participants mp ON matches.id = mp.match_id").
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