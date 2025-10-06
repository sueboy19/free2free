package routes

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"free2free/models"
	"free2free/database"
	"free2free/utils"

	apperrors "free2free/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// UserAuthMiddleware 使用者認證中介層
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作使用者認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者身份
		// 為了簡化，這裡假設有一個 isAuthenticatedUser 函數
		if !isAuthenticatedUser(c) {
			c.Error(apperrors.NewUnauthorizedError("需要使用者權限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// isAuthenticatedUser 檢查是否為已認證的使用者
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isAuthenticatedUser(c *gin.Context) bool {
	// 取得已認證的使用者
	_, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return false
	}
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
	var matches []models.Match
	// 只顯示狀態為 open 且時間未到的配對
	if err := database.GlobalDB.Conn.Preload("Activity").Preload("Organizer").Where("status = ? AND match_time > ?", "open", time.Now()).Order("match_time ASC").Find(&matches).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
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
	var match models.Match
	if err := c.ShouldBindJSON(&match); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	// Validate struct
	v := validator.New()
	if err := v.Struct(&match); err != nil {
		c.Error(apperrors.NewValidationError(err.Error()))
		return
	}

	// 檢查活動是否存在
	var activity models.Activity
	if err := database.GlobalDB.Conn.First(&activity, match.ActivityID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(apperrors.NewValidationError("指定的活動不存在"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	// 設定開局者為當前使用者
	match.OrganizerID = user.ID
	match.Status = "open"

	if err := database.GlobalDB.Conn.Create(&match).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 預加載關聯資料
	database.GlobalDB.Conn.Preload("Activity").Preload("Organizer").First(&match, match.ID)
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
	idStr := c.Param("id")
	matchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || matchID <= 0 {
		c.Error(apperrors.NewValidationError("無效的配對局 ID"))
		return
	}

	// 檢查配對局是否存在且可參與
	var match models.Match
	if err := database.GlobalDB.Conn.Where("id = ? AND status = ? AND match_time > ?", matchID, "open", time.Now()).First(&match).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(apperrors.NewValidationError("指定的配對局不存在或已關閉"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	userID := user.ID

	// 檢查使用者是否已經參與此配對局
	var existingParticipant models.MatchParticipant
	err = database.GlobalDB.Conn.Where("match_id = ? AND user_id = ?", matchID, userID).First(&existingParticipant).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 如果已經參與，返回錯誤
	if err == nil {
		c.Error(apperrors.NewValidationError("您已經參與此配對局"))
		return
	}

	// 建立新的參與記錄
	participant := models.MatchParticipant{
		MatchID:  matchID,
		UserID:   userID,
		Status:   "pending",
		JoinedAt: time.Now(),
	}

	if err := database.GlobalDB.Conn.Create(&participant).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 預加載關聯資料
	database.GlobalDB.Conn.Preload("Match").Preload("User").First(&participant, participant.ID)
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
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	userID := user.ID

	var matches []models.Match
	// 取得該使用者參與過的已完成的配對局
	if err := database.GlobalDB.Conn.Joins("JOIN match_participants mp ON matches.id = mp.match_id").
		Where("mp.user_id = ? AND matches.status = ?", userID, "completed").
		Order("matches.match_time DESC").
		Preload("Activity").
		Preload("Organizer").
		Find(&matches).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, matches)
}
