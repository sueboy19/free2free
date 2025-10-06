package routes

import (
	"errors"
	"net/http"
	"strconv"

	"free2free/models"
	"free2free/database"
	"free2free/utils"

	apperrors "free2free/errors"

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
			c.Error(apperrors.NewValidationError("無效的配對局 ID"))
			c.Abort()
			return
		}

		if !isMatchOrganizer(c, matchID) {
			c.Error(apperrors.NewForbiddenError("需要開局者權限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// isMatchOrganizer 檢查是否為指定配對局的開局者
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isMatchOrganizer(c *gin.Context, matchID int64) bool {
	// 取得已認證的使用者
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return false
	}

	// 檢查配對局是否存在且開局者為當前使用者
	var match models.Match
	err = database.GlobalDB.Conn.Where("id = ? AND organizer_id = ?", matchID, user.ID).First(&match).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		return false
	}
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
	idStr := c.Param("id")
	matchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || matchID <= 0 {
		c.Error(apperrors.NewValidationError("無效的配對局 ID"))
		return
	}

	participantIdStr := c.Param("participant_id")
	participantID, err := strconv.ParseInt(participantIdStr, 10, 64)
	if err != nil || participantID <= 0 {
		c.Error(apperrors.NewValidationError("無效的參與者 ID"))
		return
	}

	// 檢查參與者是否屬於此配對局
	var participant models.MatchParticipant
	if err := database.GlobalDB.Conn.Where("id = ? AND match_id = ?", participantID, matchID).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(apperrors.NewValidationError("指定的參與者不存在或不屬於此配對局"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 更新參與者狀態為 approved
	if err := database.GlobalDB.Conn.Model(&participant).Update("status", "approved").Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
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
	idStr := c.Param("id")
	matchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || matchID <= 0 {
		c.Error(apperrors.NewValidationError("無效的配對局 ID"))
		return
	}

	participantIdStr := c.Param("participant_id")
	participantID, err := strconv.ParseInt(participantIdStr, 10, 64)
	if err != nil || participantID <= 0 {
		c.Error(apperrors.NewValidationError("無效的參與者 ID"))
		return
	}

	// 檢查參與者是否屬於此配對局
	var participant models.MatchParticipant
	if err := database.GlobalDB.Conn.Where("id = ? AND match_id = ?", participantID, matchID).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(apperrors.NewValidationError("指定的參與者不存在或不屬於此配對局"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 更新參與者狀態為 rejected
	if err := database.GlobalDB.Conn.Model(&participant).Update("status", "rejected").Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	participant.Status = "rejected"
	c.JSON(http.StatusOK, participant)
}
