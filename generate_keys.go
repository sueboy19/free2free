//go:build ignore

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// generateKey generates a cryptographically secure random key
func generateKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func main() {
	// Generate 64 bytes random keys (encoded to ~86 characters)
	sessionKey := generateKey(64)
	jwtSecret := generateKey(64)

	fmt.Println("=== 安全金鑰產生完成 ===")
	fmt.Printf("SESSION_KEY=%s\n", sessionKey)
	fmt.Println()
	fmt.Printf("JWT_SECRET=%s", jwtSecret)
	fmt.Println()
	fmt.Println()
	fmt.Println("請將以上金鑰複製到 .env 檔案中")
	fmt.Println()
	fmt.Println("使用方式：")
	fmt.Println("1. 複製 SESSION_KEY 到 .env 檔案的 SESSION_KEY 變數")
	fmt.Println("2. 複製 JWT_SECRET 到 .env 檔案的 JWT_SECRET 變數")
	fmt.Println("3. 儲存 .env 檔案並重新啟動應用程式")
}
