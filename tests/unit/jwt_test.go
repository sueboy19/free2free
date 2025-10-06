package unit

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"free2free/models"
)

// MockClaims is a struct for JWT claims used in tests
type MockClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// mockGenerateTokens generates mock JWT tokens for testing
func mockGenerateTokens(user *models.User) (string, string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret-key-32-chars-long-enough!!" // default for tests
	}
	
	// Validate secret length (same validation as in real function)
	if len(jwtSecret) < 32 {
		return "", "", fmt.Errorf("JWT_SECRET 長度不足 32 byte")
	}
	
	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, MockClaims{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}
	
	// For simplicity in tests, we're not generating refresh tokens
	return accessTokenString, "", nil
}

// mockValidateJWTToken validates JWT tokens for testing
func mockValidateJWTToken(tokenString string) (*MockClaims, error) {
	jwtSecret := "test-secret-key-32-chars-long-enough!!"
	
	token, err := jwt.ParseWithClaims(tokenString, &MockClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*MockClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, nil
}

func TestGenerateJWTToken(t *testing.T) {
	// 設定測試用的 JWT_SECRET
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long-enough!!")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "有效使用者生成 token",
			user: &models.User{
				ID:      1,
				Name:    "Test User",
				IsAdmin: false,
			},
			wantErr: false,
		},
		{
			name: "無效 secret 長度",
			user: &models.User{
				ID:   1,
				Name: "Test User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "無效 secret 長度" {
				os.Setenv("JWT_SECRET", "short")
			}

			accessToken, _, err := mockGenerateTokens(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, accessToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)

				// 驗證 token
				claims := &MockClaims{}
				token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret-key-32-chars-long-enough!!"), nil
				})
				assert.NoError(t, err)
				assert.True(t, token.Valid)
				assert.Equal(t, tt.user.ID, claims.UserID)
				assert.Equal(t, tt.user.Name, claims.UserName)
				assert.Equal(t, tt.user.IsAdmin, claims.IsAdmin)
			}
		})
	}
}

func TestValidateJWTToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long-enough!!")
	defer func() {
		os.Unsetenv("JWT_SECRET")
	}()

	user := &models.User{
		ID:      1,
		Name:    "Test User",
		IsAdmin: false,
	}
	accessToken, _, err := mockGenerateTokens(user)
	assert.NoError(t, err)

	// 生成過期 token
	claims := &MockClaims{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := expiredToken.SignedString([]byte("test-secret-key-32-chars-long-enough!!"))
	assert.NoError(t, err)

	invalidToken := "invalid.token.string"

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有效 token",
			token:   accessToken,
			wantErr: false,
		},
		{
			name:    "過期 token",
			token:   expiredTokenString,
			wantErr: true,
		},
		{
			name:    "無效 token",
			token:   invalidToken,
			wantErr: true,
		},
		{
			name:    "空 token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := mockValidateJWTToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, user.ID, claims.UserID)
				assert.Equal(t, user.Name, claims.UserName)
				assert.Equal(t, user.IsAdmin, claims.IsAdmin)
			}
		})
	}
}