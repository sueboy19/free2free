package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReviewLike 代表評論點讚/倒讚
type ReviewLike struct {
	ID       int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID int64 `json:"review_id"`
	UserID   int64 `json:"user_id"`
	IsLike   bool  `json:"is_like"` // true: 點讚, false: 倒讚
	Review   Review `gorm:"foreignKey:ReviewID" json:"review"`
	User     User   `gorm:"foreignKey:UserID" json:"user"`
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

// likeReview 點讚評論
func likeReview(c *gin.Context) {
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的評論 ID"})
		return
	}

	// 這裡應該從 session 或 token 取得使用者 ID
	// 為了簡化，這裡暫時設為 1
	userID := int64(1)

	// 檢查評論是否存在
	var review Review
	if err := db.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的評論不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證評論"})
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike ReviewLike
	err = db.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike).Error
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
		if err := db.Model(&existingLike).Update("is_like", true).Error; err != nil {
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

	if err := db.Create(&reviewLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法點讚評論"})
		return
	}

	// 預加載關聯資料
	db.Preload("Review").Preload("User").First(&reviewLike, reviewLike.ID)
	c.JSON(http.StatusCreated, reviewLike)
}

// dislikeReview 倒讚評論
func dislikeReview(c *gin.Context) {
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的評論 ID"})
		return
	}

	// 這裡應該從 session 或 token 取得使用者 ID
	// 為了簡化，這裡暫時設為 1
	userID := int64(1)

	// 檢查評論是否存在
	var review Review
	if err := db.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的評論不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證評論"})
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike ReviewLike
	err = db.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike).Error
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
		if err := db.Model(&existingLike).Update("is_like", false).Error; err != nil {
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

	if err := db.Create(&reviewLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法倒讚評論"})
		return
	}

	// 預加載關聯資料
	db.Preload("Review").Preload("User").First(&reviewLike, reviewLike.ID)
	c.JSON(http.StatusCreated, reviewLike)
}