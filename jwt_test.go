package main

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTToken(t *testing.T) {
	// 設定測試用的 JWT_SECRET
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long!!")

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "有效使用者生成 token",
			user: &User{
				ID:      1,
				Name:    "Test User",
				IsAdmin: false,
			},
			wantErr: false,
		},
		{
			name: "無效 secret 長度",
			user: &User{
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
				defer os.Unsetenv("JWT_SECRET")
			}

			tokenString, err := generateJWTToken(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)

				// 驗證 token
				claims := &Claims{}
				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret-key-32-chars-long!!"), nil
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
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long!!")
	defer os.Unsetenv("JWT_SECRET")

	user := &User{
		ID:      1,
		Name:    "Test User",
		IsAdmin: false,
	}
	tokenString, err := generateJWTToken(user)
	assert.NoError(t, err)

	// 生成過期 token
	claims := &Claims{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := expiredToken.SignedString([]byte("test-secret-key-32-chars-long!!"))
	assert.NoError(t, err)

	invalidToken := "invalid.token.string"

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有效 token",
			token:   tokenString,
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
			claims, err := validateJWTToken(tt.token)

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
