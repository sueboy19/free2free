package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Review 代表評分與留言
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

// ReviewAuthMiddleware 評分認證中介層
func ReviewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作評分認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者是否參與了指定的配對局
		// 且配對局已結束但在評分時間範圍內 (結束後4小時)
		// 為了簡化，這裡假設有一個 canReviewMatch 函數
		matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
			c.Abort()
			return
		}

		if !canReviewMatch(c, matchID) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無評分權限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// canReviewMatch 檢查是否可以對指定配對局進行評分
// 這是一個簡化的實作，實際應用中需要檢查使用者是否參與了配對局
// 且配對局已結束但在評分時間範圍內
func canReviewMatch(c *gin.Context, matchID int64) bool {
	// 這裡應該實作實際的認證邏輯
	// 為了示範，這裡暫時回傳 true
	return true
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

// createReview 建立評分與留言
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

	// 這裡應該從 session 或 token 取得評分者 ID
	// 為了簡化，這裡暫時設為 1
	review.ReviewerID = 1
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