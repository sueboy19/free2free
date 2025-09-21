package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Match 代表配對局
type Match struct {
	ID         int64     `db:"id" json:"id"`
	ActivityID int64     `db:"activity_id" json:"activity_id"`
	OrganizerID int64    `db:"organizer_id" json:"organizer_id"`
	MatchTime  time.Time `db:"match_time" json:"match_time"`
	Status     string    `db:"status" json:"status"` // open, closed, completed
}

// MatchParticipant 代表配對參與者
type MatchParticipant struct {
	ID     int64     `db:"id" json:"id"`
	MatchID int64    `db:"match_id" json:"match_id"`
	UserID  int64    `db:"user_id" json:"user_id"`
	Status  string   `db:"status" json:"status"` // pending, approved, rejected
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
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
	rows, err := db.Query("SELECT id, activity_id, organizer_id, match_time, status FROM matches WHERE status = 'open' AND match_time > NOW() ORDER BY match_time ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得配對列表"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.ActivityID, &match.OrganizerID, &match.MatchTime, &match.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析配對資料"})
			return
		}
		matches = append(matches, match)
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
	err := db.QueryRow("SELECT id, title, target_count, location_id, description, created_by FROM activities WHERE id = ?", match.ActivityID).
		Scan(&activity.ID, &activity.Title, &activity.TargetCount, &activity.LocationID, &activity.Description, &activity.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
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

	result, err := db.Exec("INSERT INTO matches (activity_id, organizer_id, match_time, status) VALUES (?, ?, ?, ?)",
		match.ActivityID, match.OrganizerID, match.MatchTime, match.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立配對局"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得新配對局 ID"})
		return
	}

	match.ID = id
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
	err = db.QueryRow("SELECT id, activity_id, organizer_id, match_time, status FROM matches WHERE id = ? AND status = 'open' AND match_time > NOW()", matchID).
		Scan(&match.ID, &match.ActivityID, &match.OrganizerID, &match.MatchTime, &match.Status)
	if err != nil {
		if err == sql.ErrNoRows {
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
	err = db.QueryRow("SELECT id, match_id, user_id, status, joined_at FROM match_participants WHERE match_id = ? AND user_id = ?", matchID, userID).
		Scan(&existingParticipant.ID, &existingParticipant.MatchID, &existingParticipant.UserID, &existingParticipant.Status, &existingParticipant.JoinedAt)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查參與狀態"})
		return
	}

	// 如果已經參與，返回錯誤
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經參與此配對局"})
		return
	}

	// 建立新的參與記錄
	result, err := db.Exec("INSERT INTO match_participants (match_id, user_id, status) VALUES (?, ?, 'pending')",
		matchID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法參與配對局"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得參與記錄 ID"})
		return
	}

	participant := MatchParticipant{
		ID:     id,
		MatchID: matchID,
		UserID:  userID,
		Status:  "pending",
		JoinedAt: time.Now(),
	}

	c.JSON(http.StatusCreated, participant)
}

// listPastMatches 取得過去參與的配對列表
func listPastMatches(c *gin.Context) {
	// 這裡應該從 session 或 token 取得使用者 ID
	// 為了簡化，這裡暫時設為 1
	userID := int64(1)

	var matches []Match
	// 取得該使用者參與過的已完成的配對局
	rows, err := db.Query(`
		SELECT m.id, m.activity_id, m.organizer_id, m.match_time, m.status 
		FROM matches m 
		JOIN match_participants mp ON m.id = mp.match_id 
		WHERE mp.user_id = ? AND m.status = 'completed' 
		ORDER BY m.match_time DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得過去參與的配對列表"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.ActivityID, &match.OrganizerID, &match.MatchTime, &match.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析配對資料"})
			return
		}
		matches = append(matches, match)
	}

	c.JSON(http.StatusOK, matches)
}