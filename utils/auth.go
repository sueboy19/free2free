package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"free2free/models"
	"free2free/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	apperrors "free2free/errors"
)

// getDB returns the global database connection for utils package
func getDBForUtils() *gorm.DB {
	if database.GlobalDB == nil || database.GlobalDB.Conn == nil {
		panic("Database not initialized. Call database initialization first.")
	}
	return database.GlobalDB.Conn
}

// GetAuthenticatedUser 從 context 中取得已認證的使用者
func GetAuthenticatedUser(c *gin.Context) (*models.User, error) {
	// 首先嘗試從 session 取得使用者
	session := c.MustGet("session").(*sessions.Session)
	if userID, ok := session.Values["user_id"]; ok {
		// 從資料庫取得使用者資訊
		var user models.User
		err := getDBForUtils().First(&user, userID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewAppError(http.StatusNotFound, "user not found")
		}
		if err != nil {
			return nil, apperrors.MapGORMError(err)
		}
		return &user, nil
	}

	// 如果 session 中沒有使用者，嘗試從 JWT token 取得
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("no authorization header")
	}

	// 檢查 Bearer token 格式
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid authorization header format")
	}

	// 取得 token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 驗證 token
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 從資料庫取得使用者資訊
	var user models.User
	err = getDBForUtils().First(&user, claims.UserID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.NewAppError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return nil, apperrors.MapGORMError(err)
	}

	return &user, nil
}

// ValidateJWTToken 驗證 JWT token
func ValidateJWTToken(tokenString string) (*Claims, error) {
	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET 环境变量未设置")
	}

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, apperrors.NewUnauthorizedError(err.Error())
	}

	// 验证token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.NewUnauthorizedError("無效的 token")
}

// JWT claims struct
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	IsAdmin  bool   `json:"is_admin"` // 添加管理員標記
	jwt.RegisteredClaims
}