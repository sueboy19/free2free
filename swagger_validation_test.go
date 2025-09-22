package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/docs"
)

// TestSwaggerGeneration 測試Swagger文檔是否正確生成
func TestSwaggerGeneration(t *testing.T) {
	fmt.Println("開始測試Swagger文檔生成...")

	// 測試Swagger信息是否正確設置
	assert.Equal(t, "買一送一配對網站 API", docs.SwaggerInfo.Title)
	assert.Equal(t, "這是一個買一送一配對網站的API文檔", docs.SwaggerInfo.Description)
	assert.Equal(t, "1.0", docs.SwaggerInfo.Version)
	assert.Equal(t, "localhost:8080", docs.SwaggerInfo.Host)
	assert.Equal(t, "/", docs.SwaggerInfo.BasePath)
	
	fmt.Println("✓ Swagger基本信息驗證通過")

	// 測試Swagger模板是否存在
	assert.NotEmpty(t, docs.SwaggerInfo.SwaggerTemplate)
	fmt.Println("✓ Swagger模板存在")

	// 測試必要的端點是否在模板中定義
	swaggerTemplate := docs.SwaggerInfo.SwaggerTemplate
	
	// 測試認證相關端點
	assert.Contains(t, swaggerTemplate, "/auth/{provider}")
	assert.Contains(t, swaggerTemplate, "/auth/{provider}/callback")
	fmt.Println("✓ 認證端點定義驗證通過")
	
	// 測試管理員相關端點
	assert.Contains(t, swaggerTemplate, "/admin/activities")
	assert.Contains(t, swaggerTemplate, "/admin/locations")
	fmt.Println("✓ 管理員端點定義驗證通過")
	
	// 測試使用者相關端點
	assert.Contains(t, swaggerTemplate, "/user/matches")
	assert.Contains(t, swaggerTemplate, "/user/past-matches")
	fmt.Println("✓ 使用者端點定義驗證通過")
	
	// 測試開局者相關端點
	assert.Contains(t, swaggerTemplate, "/organizer/matches")
	fmt.Println("✓ 開局者端點定義驗證通過")
	
	// 測試評分相關端點
	assert.Contains(t, swaggerTemplate, "/review/matches")
	fmt.Println("✓ 評分端點定義驗證通過")
	
	// 測試評論點讚/倒讚相關端點
	assert.Contains(t, swaggerTemplate, "/review-like/reviews")
	fmt.Println("✓ 評論點讚/倒讚端點定義驗證通過")

	// 測試必要的數據模型是否定義
	assert.Contains(t, swaggerTemplate, "main.User")
	assert.Contains(t, swaggerTemplate, "main.Activity")
	assert.Contains(t, swaggerTemplate, "main.Location")
	assert.Contains(t, swaggerTemplate, "main.Match")
	assert.Contains(t, swaggerTemplate, "main.MatchParticipant")
	assert.Contains(t, swaggerTemplate, "main.Review")
	assert.Contains(t, swaggerTemplate, "main.ReviewLike")
	fmt.Println("✓ 數據模型定義驗證通過")

	fmt.Println("Swagger文檔生成測試完成！")
}

// TestAPIEndpointRegistration 測試API端點是否正確註冊
func TestAPIEndpointRegistration(t *testing.T) {
	fmt.Println("開始測試API端點註冊...")

	// 測試Swagger文檔中是否包含必要的HTTP方法
	swaggerTemplate := docs.SwaggerInfo.SwaggerTemplate
	
	// 測試GET方法
	assert.Contains(t, swaggerTemplate, "\"get\":")
	fmt.Println("✓ GET方法定義驗證通過")
	
	// 測試POST方法
	assert.Contains(t, swaggerTemplate, "\"post\":")
	fmt.Println("✓ POST方法定義驗證通過")
	
	// 測試PUT方法
	assert.Contains(t, swaggerTemplate, "\"put\":")
	fmt.Println("✓ PUT方法定義驗證通過")
	
	// 測試DELETE方法
	assert.Contains(t, swaggerTemplate, "\"delete\":")
	fmt.Println("✓ DELETE方法定義驗證通過")

	fmt.Println("API端點註冊測試完成！")
}