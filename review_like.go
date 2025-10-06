package main

import (
	"net/http"
	"strconv"

	"free2free/models"
	"free2free/database"
	"free2free/utils"

	apperrors "free2free/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
	idStr := c.Param("id")
	reviewID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || reviewID <= 0 {
		c.Error(apperrors.NewValidationError("無效的評論 ID"))
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	userID := user.ID

	// 檢查評論是否存在
	var review models.Review
	if err := database.GlobalDB.Conn.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(apperrors.NewValidationError("指定的評論不存在"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike models.ReviewLike
	err = database.GlobalDB.Conn.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 如果已經點讚，返回錯誤
	if err == nil && existingLike.IsLike {
		c.Error(apperrors.NewValidationError("您已經點讚此評論"))
		return
	}

	// 如果已經倒讚，則更新為點讚
	if err == nil && !existingLike.IsLike {
		if err := database.GlobalDB.Conn.Model(&existingLike).Update("is_like", true).Error; err != nil {
			c.Error(apperrors.MapGORMError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "已將倒讚改為點讚"})
		return
	}

	// 建立新的點讚記錄
	reviewLike := models.ReviewLike{
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   true,
	}

	if err := database.GlobalDB.Conn.Create(&reviewLike).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 預加載關聯資料
	database.GlobalDB.Conn.Preload("Review").Preload("User").First(&reviewLike, reviewLike.ID)
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
	idStr := c.Param("id")
	reviewID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || reviewID <= 0 {
		c.Error(apperrors.NewValidationError("無效的評論 ID"))
		return
	}

	// 從認證資訊取得使用者 ID
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		c.Error(apperrors.NewUnauthorizedError("未登入"))
		return
	}

	userID := user.ID

	// 檢查評論是否存在
	var review models.Review
	if err := database.GlobalDB.Conn.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(apperrors.NewValidationError("指定的評論不存在"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike models.ReviewLike
	err = database.GlobalDB.Conn.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 如果已經倒讚，返回錯誤
	if err == nil && !existingLike.IsLike {
		c.Error(apperrors.NewValidationError("您已經倒讚此評論"))
		return
	}

	// 如果已經點讚，則更新為倒讚
	if err == nil && existingLike.IsLike {
		if err := database.GlobalDB.Conn.Model(&existingLike).Update("is_like", false).Error; err != nil {
			c.Error(apperrors.MapGORMError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "已將點讚改為倒讚"})
		return
	}

	// 建立新的倒讚記錄
	reviewLike := models.ReviewLike{
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   false,
	}

	if err := database.GlobalDB.Conn.Create(&reviewLike).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 預加載關聯資料
	database.GlobalDB.Conn.Preload("Review").Preload("User").First(&reviewLike, reviewLike.ID)
	c.JSON(http.StatusCreated, reviewLike)
}
