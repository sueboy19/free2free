package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Activity 代表配對活動
type Activity struct {
	ID          int64  `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	TargetCount int    `db:"target_count" json:"target_count"`
	LocationID  int64  `db:"location_id" json:"location_id"`
	Description string `db:"description" json:"description"`
	CreatedBy   int64  `db:"created_by" json:"created_by"`
}

// Location 代表地點
type Location struct {
	ID        int64   `db:"id" json:"id"`
	Name      string  `db:"name" json:"name"`
	Address   string  `db:"address" json:"address"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
}

// Admin 代表管理員
type Admin struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
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
	rows, err := db.Query("SELECT id, title, target_count, location_id, description, created_by FROM activities ORDER BY id DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得活動列表"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var activity Activity
		if err := rows.Scan(&activity.ID, &activity.Title, &activity.TargetCount, &activity.LocationID, &activity.Description, &activity.CreatedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析活動資料"})
			return
		}
		activities = append(activities, activity)
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
	err := db.QueryRow("SELECT id, name, address, latitude, longitude FROM locations WHERE id = ?", activity.LocationID).
		Scan(&location.ID, &location.Name, &location.Address, &location.Latitude, &location.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的地點不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證地點"})
		return
	}

	// 這裡應該從 session 或 token 取得管理員 ID
	// 為了簡化，這裡暫時設為 1
	activity.CreatedBy = 1

	result, err := db.Exec("INSERT INTO activities (title, target_count, location_id, description, created_by) VALUES (?, ?, ?, ?, ?)",
		activity.Title, activity.TargetCount, activity.LocationID, activity.Description, activity.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立活動"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得新活動 ID"})
		return
	}

	activity.ID = id
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
	err = db.QueryRow("SELECT id, name, address, latitude, longitude FROM locations WHERE id = ?", activity.LocationID).
		Scan(&location.ID, &location.Name, &location.Address, &location.Latitude, &location.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的地點不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法驗證地點"})
		return
	}

	// 更新活動
	_, err = db.Exec("UPDATE activities SET title = ?, target_count = ?, location_id = ?, description = ? WHERE id = ?",
		activity.Title, activity.TargetCount, activity.LocationID, activity.Description, id)
	if err != nil {
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

	_, err = db.Exec("DELETE FROM activities WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除活動"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活動已刪除"})
}

// listLocations 取得地點列表
func listLocations(c *gin.Context) {
	var locations []Location
	rows, err := db.Query("SELECT id, name, address, latitude, longitude FROM locations ORDER BY id DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得地點列表"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var location Location
		if err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.Latitude, &location.Longitude); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析地點資料"})
			return
		}
		locations = append(locations, location)
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

	result, err := db.Exec("INSERT INTO locations (name, address, latitude, longitude) VALUES (?, ?, ?, ?)",
		location.Name, location.Address, location.Latitude, location.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立地點"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得新地點 ID"})
		return
	}

	location.ID = id
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
	_, err = db.Exec("UPDATE locations SET name = ?, address = ?, latitude = ?, longitude = ? WHERE id = ?",
		location.Name, location.Address, location.Latitude, location.Longitude, id)
	if err != nil {
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

	_, err = db.Exec("DELETE FROM locations WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法刪除地點"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "地點已刪除"})
}