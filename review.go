package main

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

// ReviewAuthMiddleware 評分認證中介層
func ReviewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作評分認證邏輯
		// 例如檢查 session 或 JWT token 中的使用者是否參與了指定的配對局
		// 且配對局已結束但在評分時間範圍內 (結束後4小時)
		// 為了簡化，這裡假設有一個 canReviewMatch 函數
		matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.Error(apperrors.NewValidationError("無效的配對局 ID"))
			c.Abort()
			return
		}

		if !canReviewMatch(c, matchID) {
			c.Error(apperrors.NewForbiddenError("無評分權限"))
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
	// 取得已認證的使用者
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return false
	}

	// 檢查使用者是否參與了指定的配對局
	var participant models.MatchParticipant
	err = database.GlobalDB.Conn.Where("match_id = ? AND user_id = ? AND status = ?", matchID, user.ID, "approved").First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		return false
	}

	// 檢查配對局是否已完成
	var match models.Match
	err = database.GlobalDB.Conn.Where("id = ? AND status = ?", matchID, "completed").First(&match).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		return false
	}

	// 檢查是否在評分時間範圍內（結束後4小時內）
	// 假設配對局結束時間存儲在 match.MatchTime 中
	// 實際應用中可能需要一個額外的字段來記錄配對局的結束時間
	endTime := match.MatchTime
	reviewDeadline := endTime.Add(4 * time.Hour)

	return time.Now().Before(reviewDeadline)
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
	idStr := c.Param("id")
	matchID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || matchID <= 0 {
		c.Error(apperrors.NewValidationError("無效的配對局 ID"))
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	// Validate struct
	v := validator.New()
	if err := v.Struct(&review); err != nil {
		c.Error(apperrors.NewValidationError(err.Error()))
		return
	}

	// 從認證資訊取得評分者 ID
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	review.ReviewerID = user.ID
	review.MatchID = matchID
	review.CreatedAt = time.Now()

	// 檢查是否已經對此人在此配對局評分過
	var existingReview models.Review
	err = database.GlobalDB.Conn.Where("reviewer_id = ? AND reviewee_id = ? AND match_id = ?",
		review.ReviewerID, review.RevieweeID, review.MatchID).First(&existingReview).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 如果已經評分過，返回錯誤
	if err == nil {
		c.Error(apperrors.NewValidationError("您已經對此人評分過"))
		return
	}

	// 建立新的評分記錄
	if err := database.GlobalDB.Conn.Create(&review).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 預加載關聯資料
	database.GlobalDB.Conn.Preload("Match").Preload("Reviewer").Preload("Reviewee").First(&review, review.ID)
	c.JSON(http.StatusCreated, review)
}
