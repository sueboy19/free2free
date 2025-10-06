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
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// AdminAuthMiddleware 管理員認證中介層
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作管理員認證邏輯
		// 例如檢查 session 或 JWT token 中的管理員身份
		// 為了簡化，這裡假設有一個 isAuthenticatedAdmin 函數
		if !isAuthenticatedAdmin(c) {
			c.Error(apperrors.NewForbiddenError("需要管理員權限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// isAuthenticatedAdmin 檢查是否為已認證的管理員
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isAuthenticatedAdmin(c *gin.Context) bool {
	// 在這個簡化的實作中，我們假設 userID 為 1 的使用者是管理員
	// 實際應用中應該有一個管理員表或管理員標記欄位

	// 取得已認證的使用者
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return false
	}

	return user.IsAdmin
}

// SetupAdminRoutes 設定管理後台路由
func SetupAdminRoutes(r *gin.Engine) {
	// 管理員認證路由組
	admin := r.Group("/admin")
	admin.Use(AdminAuthMiddleware())
	{
		// 配對活動管理
		admin.GET("/activities", listActivities)
		admin.POST("/activities", createActivity)
		admin.PUT("/activities/:id", updateActivity)
		admin.DELETE("/activities/:id", deleteActivity)

		// 地點管理
		admin.GET("/locations", listLocations)
		admin.POST("/locations", createLocation)
		admin.PUT("/locations/:id", updateLocation)
		admin.DELETE("/locations/:id", deleteLocation)
	}
}

// listActivities 取得配對活動列表
// @Summary 取得配對活動列表
// @Description 取得所有配對活動的列表
// @Tags 管理員
// @Accept json
// @Produce json
// @Success 200 {array} Activity
// @Failure 500 {object} map[string]string "無法取得活動列表"
// @Router /admin/activities [get]
// @Security ApiKeyAuth
func listActivities(c *gin.Context) {
	var activities []models.Activity
	if err := database.GlobalDB.Conn.Preload("Location").Order("id DESC").Find(&activities).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, activities)
}

// createActivity 建立新的配對活動
// @Summary 建立新的配對活動
// @Description 建立新的配對活動
// @Tags 管理員
// @Accept json
// @Produce json
// @Param activity body Activity true "配對活動資訊"
// @Success 201 {object} Activity
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法建立活動"
// @Router /admin/activities [post]
// @Security ApiKeyAuth
func createActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	// Validate struct
	v := validator.New()
	if err := v.Struct(&activity); err != nil {
		c.Error(apperrors.NewValidationError(err.Error()))
		return
	}

	// 檢查地點是否存在
	var location models.Location
	if err := database.GlobalDB.Conn.First(&location, activity.LocationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(apperrors.NewValidationError("指定的地點不存在"))
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

	// 設定活動建立者為當前使用者
	activity.CreatedBy = user.ID

	if err := database.GlobalDB.Conn.Create(&activity).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// updateActivity 更新配對活動
// @Summary 更新配對活動
// @Description 更新指定ID的配對活動
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "活動ID"
// @Param activity body Activity true "配對活動資訊"
// @Success 200 {object} Activity
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法更新活動"
// @Router /admin/activities/{id} [put]
// @Security ApiKeyAuth
func updateActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Error(apperrors.NewValidationError("無效的活動 ID"))
		return
	}

	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	// Validate struct
	v := validator.New()
	if err := v.Struct(&activity); err != nil {
		c.Error(apperrors.NewValidationError(err.Error()))
		return
	}

	// 檢查地點是否存在
	var location models.Location
	if err := database.GlobalDB.Conn.First(&location, activity.LocationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(apperrors.NewValidationError("指定的地點不存在"))
			return
		}
		c.Error(apperrors.MapGORMError(err))
		return
	}

	// 更新活動
	if err := database.GlobalDB.Conn.Model(&models.Activity{}).Where("id = ?", id).Updates(activity).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	activity.ID = id
	c.JSON(http.StatusOK, activity)
}

// deleteActivity 刪除配對活動
// @Summary 刪除配對活動
// @Description 刪除指定ID的配對活動
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "活動ID"
// @Success 200 {object} map[string]string "活動已刪除"
// @Failure 400 {object} map[string]string "無效的活動 ID"
// @Failure 500 {object} map[string]string "無法刪除活動"
// @Router /admin/activities/{id} [delete]
// @Security ApiKeyAuth
func deleteActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Error(apperrors.NewValidationError("無效的活動 ID"))
		return
	}

	if err := database.GlobalDB.Conn.Delete(&models.Activity{}, id).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活動已刪除"})
}

// listLocations 取得地點列表
// @Summary 取得地點列表
// @Description 取得所有地點的列表
// @Tags 管理員
// @Accept json
// @Produce json
// @Success 200 {array} Location
// @Failure 500 {object} map[string]string "無法取得地點列表"
// @Router /admin/locations [get]
// @Security ApiKeyAuth
func listLocations(c *gin.Context) {
	var locations []models.Location
	if err := database.GlobalDB.Conn.Order("id DESC").Find(&locations).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, locations)
}

// createLocation 建立新的地點
// @Summary 建立新的地點
// @Description 建立新的地點
// @Tags 管理員
// @Accept json
// @Produce json
// @Param location body Location true "地點資訊"
// @Success 201 {object} Location
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法建立地點"
// @Router /admin/locations [post]
// @Security ApiKeyAuth
func createLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	// Validate struct
	v := validator.New()
	if err := v.Struct(&location); err != nil {
		c.Error(apperrors.NewValidationError(err.Error()))
		return
	}

	if err := database.GlobalDB.Conn.Create(&location).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusCreated, location)
}

// updateLocation 更新地點
// @Summary 更新地點
// @Description 更新指定ID的地點
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "地點ID"
// @Param location body Location true "地點資訊"
// @Success 200 {object} Location
// @Failure 400 {object} map[string]string "無效的請求資料"
// @Failure 500 {object} map[string]string "無法更新地點"
// @Router /admin/locations/{id} [put]
// @Security ApiKeyAuth
func updateLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Error(apperrors.NewValidationError("無效的地點 ID"))
		return
	}

	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.Error(apperrors.NewValidationError("無效的請求資料"))
		return
	}

	// Validate struct
	v := validator.New()
	if err := v.Struct(&location); err != nil {
		c.Error(apperrors.NewValidationError(err.Error()))
		return
	}

	// 更新地點
	if err := database.GlobalDB.Conn.Model(&models.Location{}).Where("id = ?", id).Updates(location).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	location.ID = id
	c.JSON(http.StatusOK, location)
}

// deleteLocation 刪除地點
// @Summary 刪除地點
// @Description 刪除指定ID的地點
// @Tags 管理員
// @Accept json
// @Produce json
// @Param id path int true "地點ID"
// @Success 200 {object} map[string]string "地點已刪除"
// @Failure 400 {object} map[string]string "無效的地點 ID"
// @Failure 500 {object} map[string]string "無法刪除地點"
// @Router /admin/locations/{id} [delete]
// @Security ApiKeyAuth
func deleteLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.Error(apperrors.NewValidationError("無效的地點 ID"))
		return
	}

	if err := database.GlobalDB.Conn.Delete(&models.Location{}, id).Error; err != nil {
		c.Error(apperrors.MapGORMError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "地點已刪除"})
}
