package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Activity 代表配對活動
type Activity struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `json:"title"`
	TargetCount int    `json:"target_count"`
	LocationID  int64  `json:"location_id"`
	Description string `json:"description"`
	CreatedBy   int64  `json:"created_by"`
	Location    Location `gorm:"foreignKey:LocationID" json:"location"`
}

// Location 代表地點
type Location struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Admin 代表管理員
type Admin struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
}

// AdminAuthMiddleware 管理員認證中介層
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 這裡應該實作管理員認證邏輯
		// 例如檢查 session 或 JWT token 中的管理員身份
		// 為了簡化，這裡假設有一個 isAuthenticatedAdmin 函數
		if !isAuthenticatedAdmin(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要管理員權限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// isAuthenticatedAdmin 檢查是否為已認證的管理員
// 這是一個簡化的實作，實際應用中需要檢查 session 或 token
func isAuthenticatedAdmin(c *gin.Context) bool {
	// 這裡應該實作實際的認證邏輯
	// 例如檢查 session 中是否有 admin_id
	// 或解析 JWT token 驗證管理員身份
	// 為了示範，這裡暫時回傳 true
	return true
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
func listActivities(c *gin.Context) {
	var activities []Activity
	if err := adminDB.Preload("Location").Order("id DESC").Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得活動列表"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// createActivity 建立新的配對活動
func createActivity(c *gin.Context) {
	var activity Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if activity.Title == "" || activity.TargetCount <= 0 || activity.LocationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "標題、目標人數和地點為必填欄位"})
		return
	}

	// 檢查地點是否存在
	var location Location
	if err := adminDB.First(&location, activity.LocationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的地點不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證地點"})
		return
	}

	// 這裡應該從 session 或 token 取得管理員 ID
	// 為了簡化，這裡暫時設為 1
	activity.CreatedBy = 1

	if err := adminDB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立活動"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// updateActivity 更新配對活動
func updateActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的活動 ID"})
		return
	}

	var activity Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if activity.Title == "" || activity.TargetCount <= 0 || activity.LocationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "標題、目標人數和地點為必填欄位"})
		return
	}

	// 檢查地點是否存在
	var location Location
	if err := adminDB.First(&location, activity.LocationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的地點不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證地點"})
		return
	}

	// 更新活動
	if err := adminDB.Model(&Activity{}).Where("id = ?", id).Updates(activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新活動"})
		return
	}

	activity.ID = id
	c.JSON(http.StatusOK, activity)
}

// deleteActivity 刪除配對活動
func deleteActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的活動 ID"})
		return
	}

	if err := adminDB.Delete(&Activity{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除活動"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活動已刪除"})
}

// listLocations 取得地點列表
func listLocations(c *gin.Context) {
	var locations []Location
	if err := adminDB.Order("id DESC").Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得地點列表"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

// createLocation 建立新的地點
func createLocation(c *gin.Context) {
	var location Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if location.Name == "" || location.Address == "" || location.Latitude == 0 || location.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "名稱、地址和座標為必填欄位"})
		return
	}

	if err := adminDB.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立地點"})
		return
	}

	c.JSON(http.StatusCreated, location)
}

// updateLocation 更新地點
func updateLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的地點 ID"})
		return
	}

	var location Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料"})
		return
	}

	// 驗證必要欄位
	if location.Name == "" || location.Address == "" || location.Latitude == 0 || location.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "名稱、地址和座標為必填欄位"})
		return
	}

	// 更新地點
	if err := adminDB.Model(&Location{}).Where("id = ?", id).Updates(location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法更新地點"})
		return
	}

	location.ID = id
	c.JSON(http.StatusOK, location)
}

// deleteLocation 刪除地點
func deleteLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的地點 ID"})
		return
	}

	if err := adminDB.Delete(&Location{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除地點"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "地點已刪除"})
}
