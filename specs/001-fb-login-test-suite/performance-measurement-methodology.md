# 效能測量方法論

## 概述
本文檔定義了規格說明中效能要求的具體測量方法和標準。

## 效能指標定義

### 1. Facebook OAuth 流程效能
**規格要求**: 30 秒內完成端到端流程

#### 測量方法
```mermaid
sequenceDiagram
    participant U as 使用者
    participant F as 前端
    participant B as 後端
    participant FB as Facebook

    Note over U,B: 開始計時 (T0)
    U->>F: 點擊 Facebook 登入
    F->>B: GET /auth/facebook
    B->>F: 重定向到 Facebook
    F->>FB: Facebook 登入流程
    FB->>B: GET /auth/facebook/callback
    Note over B: 使用者資料處理 (T1)
    B->>B: JWT Token 生成 (T2)
    B->>F: 返回 HTML + Token
    F->>U: 顯示登入成功
    Note over U,B: 結束計時 (T3)
```

#### 測量步驟
1. **準備階段**
   - 清空瀏覽器快取
   - 確保網路連線穩定
   - 使用乾淨的測試環境

2. **執行測量**
   ```bash
   # 使用瀏覽器開發工具測量
   console.time('facebook-oauth-flow');
   
   # 在登入按鈕點擊時開始計時
   document.getElementById('fb-login').addEventListener('click', function() {
     console.time('facebook-oauth-flow');
   });
   
   # 在收到成功回應時結束計時
   window.addEventListener('message', function(event) {
     if (event.data.type === 'auth_success') {
       console.timeEnd('facebook-oauth-flow');
     }
   });
   ```

3. **自動化測量**
   ```javascript
   // 使用 Playwright 測量
   const startTime = Date.now();
   await page.click('#fb-login');
   await page.waitForEvent('popup');
   // 等待 Facebook 登入完成
   await page.waitForFunction(() => window.authSuccess);
   const endTime = Date.now();
   const duration = endTime - startTime;
   ```

#### 成功標準
- **最佳情況**: < 15 秒
- **可接受**: < 30 秒 (規格要求)
- **警告**: 30-45 秒
- **失敗**: > 45 秒

### 2. API 回應時間效能
**規格要求**: < 500ms API 回應時間

#### 測量方法
```bash
# 使用 curl 測量單一請求
curl -w "@curl-format.txt" -o /dev/null -s "http://localhost:8080/profile"

# curl-format.txt 內容
     time_namelookup:  %{time_namelookup}\n
        time_connect:  %{time_connect}\n
     time_appconnect:  %{time_appconnect}\n
    time_pretransfer:  %{time_pretransfer}\n
       time_redirect:  %{time_redirect}\n
  time_starttransfer:  %{time_starttransfer}\n
                     ----------\n
          time_total:  %{time_total}\n
```

#### 自動化測量腳本
```bash
#!/bin/bash
# performance-test.sh

echo "=== API 效能測試報告 ==="
echo "測試時間: $(date)"
echo ""

# 測試端點列表
endpoints=(
  "/profile"
  "/user/matches"
  "/user/past-matches"
  "/admin/activities"
  "/auth/token"
)

# JWT Token (測試用)
TEST_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJUZXN0IFVzZXIiLCJpc19hZG1pbiI6ZmFsc2UsImlhdCI6MTY5Nzg5MDQwMCwiZXhwIjoxNjk3ODkxMzAwfQ.test-signature"

for endpoint in "${endpoints[@]}"; do
  echo "測試端點: $endpoint"
  
  # 執行 5 次測試取平均值
  total_time=0
  for i in {1..5}; do
    response_time=$(curl -w "%{time_total}" -o /dev/null -s \
      -H "Authorization: Bearer $TEST_TOKEN" \
      "http://localhost:8080$endpoint")
    
    response_time_ms=$(echo "$response_time * 1000" | bc)
    total_time=$(echo "$total_time + $response_time_ms" | bc)
    echo "  請求 $i: ${response_time_ms}ms"
  done
  
  average_time=$(echo "scale=2; $total_time / 5" | bc)
  echo "  平均回應時間: ${average_time}ms"
  
  if (( $(echo "$average_time > 500" | bc -l) )); then
    echo "  狀態: ❌ 超過 500ms"
  else
    echo "  狀態: ✅ 符合標準"
  fi
  echo ""
done
```

#### JWT Token 驗證效能
**規格要求**: < 10ms JWT token 驗證

```go
// 效能測試程式碼
func BenchmarkJWTValidation(b *testing.B) {
    token := "your-test-token"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := utils.ValidateJWTToken(token)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// 測量結果
func TestJWTValidationPerformance(t *testing.T) {
    token := "your-test-token"
    
    start := time.Now()
    _, err := utils.ValidateJWTToken(token)
    duration := time.Since(start)
    
    if err != nil {
        t.Fatal(err)
    }
    
    t.Logf("JWT 驗證時間: %v", duration)
    if duration > 10*time.Millisecond {
        t.Errorf("JWT 驗證時間超過 10ms: %v", duration)
    }
}
```

### 3. 資料庫查詢效能
**測量標準**: < 100ms 資料庫查詢

```go
// 資料庫效能測試
func TestDatabaseQueryPerformance(t *testing.T) {
    db := database.GlobalDB.Conn
    
    tests := []struct {
        name     string
        query    func() error
        expected time.Duration
    }{
        {
            name: "使用者查詢",
            query: func() error {
                var user models.User
                return db.First(&user, 1).Error
            },
            expected: 100 * time.Millisecond,
        },
        {
            name: "配對列表查詢",
            query: func() error {
                var matches []models.Match
                return db.Preload("Activity").Preload("Organizer").
                    Where("status = ? AND match_time > ?", "open", time.Now()).
                    Order("match_time ASC").
                    Find(&matches).Error
            },
            expected: 100 * time.Millisecond,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            start := time.Now()
            err := tt.query()
            duration := time.Since(start)
            
            if err != nil {
                t.Fatal(err)
            }
            
            t.Logf("%s 查詢時間: %v", tt.name, duration)
            if duration > tt.expected {
                t.Errorf("%s 查詢時間超過預期: %v", tt.name, duration)
            }
        })
    }
}
```

## 效能監控工具

### 1. 開發環境監控
```javascript
// 前端效能監控
function measureAPICall(endpoint, method = 'GET') {
    const startTime = performance.now();
    
    return fetch(endpoint, { method })
        .then(response => {
            const endTime = performance.now();
            const duration = endTime - startTime;
            
            console.log(`${method} ${endpoint}: ${duration.toFixed(2)}ms`);
            
            // 記錄到效能指標
            if (duration > 500) {
                console.warn(`⚠️ API 回應時間超過 500ms: ${duration.toFixed(2)}ms`);
            }
            
            return response;
        });
}
```

### 2. 伺服器端監控
```go
// Gin 中介層效能監控
func PerformanceMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        // 計算請求處理時間
        duration := param.TimeStamp.Sub(param.Request.Header.Get("X-Request-Start")).Seconds()
        
        // 記錄慢請求
        if duration > 0.5 {
            log.Printf("⚠️ 慢請求: %s %s - %.3f秒", 
                param.Method, param.Path, duration)
        }
        
        return fmt.Sprintf("%s - %s %s - %d - %.3f秒\n",
            param.ClientIP,
            param.Method,
            param.Path,
            param.StatusCode,
            duration)
    })
}
```

## 效能基準測試

### 壓力測試腳本
```bash
#!/bin/bash
# load-test.sh

echo "=== 效能基準測試 ==="
echo "測試時間: $(date)"
echo ""

# 使用 Apache Bench 進行壓力測試
echo "1. Profile 端點壓力測試 (100 請求, 10 並行)"
ab -n 100 -c 10 -H "Authorization: Bearer $TEST_TOKEN" \
   http://localhost:8080/profile

echo ""
echo "2. 使用者配對端點壓力測試 (50 請求, 5 並行)"
ab -n 50 -c 5 -H "Authorization: Bearer $TEST_TOKEN" \
   http://localhost:8080/user/matches

echo ""
echo "3. Facebook OAuth 端點壓力測試 (20 請求, 3 並行)"
ab -n 20 -c 3 http://localhost:8080/auth/facebook

echo ""
echo "=== 測試完成 ==="
```

## 效能報告範本

```markdown
# 效能測試報告

## 測試資訊
- **測試日期**: 2025-11-20
- **測試環境**: 本地開發環境
- **硬體規格**: [填寫硬體資訊]
- **網路環境**: [填寫網路資訊]

## 測試結果

### Facebook OAuth 流程
- **平均時間**: 18.5 秒
- **最佳時間**: 12.3 秒
- **最差時間**: 28.7 秒
- **符合標準**: ✅

### API 回應時間
| 端點 | 平均時間 | 符合標準 |
|------|----------|----------|
| /profile | 45ms | ✅ |
| /user/matches | 89ms | ✅ |
| /user/past-matches | 67ms | ✅ |
| /admin/activities | 123ms | ✅ |
| /auth/token | 23ms | ✅ |

### JWT 驗證效能
- **平均驗證時間**: 3.2ms
- **符合標準**: ✅

### 資料庫查詢效能
| 查詢類型 | 平均時間 | 符合標準 |
|----------|----------|----------|
| 使用者查詢 | 12ms | ✅ |
| 配對列表查詢 | 45ms | ✅ |
| 活動查詢 | 28ms | ✅ |

## 總結
- **整體符合率**: 100%
- **需要優化項目**: 無
- **下次測試時間**: [預計時間]
```

## 更新紀錄
- **2025-11-20**: 初次建立效能測量方法論
- **版本**: v1.0