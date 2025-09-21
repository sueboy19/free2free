package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrganizerAuthMiddleware 開局者認證中介層
func OrganizerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作開局者認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者是否為指定配對局的開局者
		// 為了簡化，這裡假設有一個 isMatchOrganizer 函數
		matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的配對局 ID"})
			c.Abort()
			return
		}

		if !isMatchOrganizer(c, matchID) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要開局者權限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// isMatchOrganizer 檢查是否為指定配對局的開局者
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isMatchOrganizer(c *gin.Context, matchID int64) bool {
	// 這裡應該實作實際的認證邏輯
	// 例如檢查 session 中的 user_id 是否與配對局的 organizer_id 相同
	// 為了示範，這裡暫時回傳 true
	return true
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

// approveParticipant 審核通過參與者
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