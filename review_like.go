package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReviewLike 代表評論點讚/倒讚
type ReviewLike struct {
	ID       int64 `db:"id" json:"id"`
	ReviewID int64 `db:"review_id" json:"review_id"`
	UserID   int64 `db:"user_id" json:"user_id"`
	IsLike   bool  `db:"is_like" json:"is_like"` // true: 點讚, false: 倒讚
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
	err = db.QueryRow("SELECT id, match_id, reviewer_id, reviewee_id, score, comment, created_at FROM reviews WHERE id = ?", reviewID).
		Scan(&review.ID, &review.MatchID, &review.ReviewerID, &review.RevieweeID, &review.Score, &review.Comment, &review.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的評論不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證評論"})
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike ReviewLike
	err = db.QueryRow("SELECT id, review_id, user_id, is_like FROM review_likes WHERE user_id = ? AND review_id = ?", userID, reviewID).
		Scan(&existingLike.ID, &existingLike.ReviewID, &existingLike.UserID, &existingLike.IsLike)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查點讚狀態"})
		return
	}

	// 如果已經點讚，返回錯誤
	if err != sql.ErrNoRows && existingLike.IsLike {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經點讚此評論"})
		return
	}

	// 如果已經倒讚，則更新為點讚
	if err != sql.ErrNoRows && !existingLike.IsLike {
		_, err = db.Exec("UPDATE review_likes SET is_like = true WHERE id = ?", existingLike.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新點讚狀態"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "已將倒讚改為點讚"})
		return
	}

	// 建立新的點讚記錄
	result, err := db.Exec("INSERT INTO review_likes (review_id, user_id, is_like) VALUES (?, ?, true)", reviewID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法點讚評論"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得點讚記錄 ID"})
		return
	}

	reviewLike := ReviewLike{
		ID:       id,
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   true,
	}

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
	err = db.QueryRow("SELECT id, match_id, reviewer_id, reviewee_id, score, comment, created_at FROM reviews WHERE id = ?", reviewID).
		Scan(&review.ID, &review.MatchID, &review.ReviewerID, &review.RevieweeID, &review.Score, &review.Comment, &review.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的評論不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證評論"})
		return
	}

	// 檢查使用者是否已經對此評論點讚或倒讚
	var existingLike ReviewLike
	err = db.QueryRow("SELECT id, review_id, user_id, is_like FROM review_likes WHERE user_id = ? AND review_id = ?", userID, reviewID).
		Scan(&existingLike.ID, &existingLike.ReviewID, &existingLike.UserID, &existingLike.IsLike)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢查點讚狀態"})
		return
	}

	// 如果已經倒讚，返回錯誤
	if err != sql.ErrNoRows && !existingLike.IsLike {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已經倒讚此評論"})
		return
	}

	// 如果已經點讚，則更新為倒讚
	if err != sql.ErrNoRows && existingLike.IsLike {
		_, err = db.Exec("UPDATE review_likes SET is_like = false WHERE id = ?", existingLike.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新點讚狀態"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "已將點讚改為倒讚"})
		return
	}

	// 建立新的倒讚記錄
	result, err := db.Exec("INSERT INTO review_likes (review_id, user_id, is_like) VALUES (?, ?, false)", reviewID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法倒讚評論"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得倒讚記錄 ID"})
		return
	}

	reviewLike := ReviewLike{
		ID:       id,
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   false,
	}

	c.JSON(http.StatusCreated, reviewLike)
}