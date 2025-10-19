# 測試套件：Facebook 登入與 API 測試

本目錄包含用於測試 Facebook 登入功能和所有 API 端點的完整測試套件。

## 目錄結構

```
tests/
├── contract/              # API 合約測試
│   ├── test_fb_oauth_contract.go
│   └── test_protected_endpoints_contract.go
├── e2e/                 # 端到端測試
│   ├── fb_login_e2e_test.go
│   ├── env_setup_test.go
│   ├── test_suite_validation_test.go
│   └── complete_flow_test.go
├── integration/          # 整合測試
│   ├── fb_auth_integration_test.go
│   ├── api_integration_test.go
│   ├── user_api_integration_test.go
│   ├── admin_api_integration_test.go
│   ├── organizer_api_integration_test.go
│   └── review_api_integration_test.go
├── unit/                # 單元測試
│   └── jwt_token_test.go
├── performance/         # 效能測試
│   └── fb_login_performance_test.go
├── testutils/           # 測試工具
│   ├── config.go
│   ├── test_server.go
│   ├── mock_fb_provider.go
│   ├── jwt_validator.go
│   ├── api_helpers.go
│   ├── fb_test_helpers.go
│   ├── test_data.go
│   ├── result_reporter.go
│   └── test_cleanup.go
└── README.md            # 本文件
```

## 測試類型

### 1. 單元測試 (unit/)
- 測試 JWT 令牌生成和驗證功能
- 驗證單一函數和方法的正確性

### 2. 整合測試 (integration/)
- 測試 Facebook OAuth 流程
- 測試與資料庫的整合
- 測試不同 API 端點的整合

### 3. 端到端測試 (e2e/)
- 測試完整的 Facebook 登入流程
- 測試從登入到使用 API 的完整用戶旅程

### 4. 合約測試 (contract/)
- 驗證 API 端點的合約和響應結構
- 確保端點符合預期的行為

### 5. 效能測試 (performance/)
- 測試 Facebook 登入流程的效能 (應在 30 秒內完成)
- 測試 JWT 驗證效能 (應在 10 毫秒內完成)
- 測試 API 響應時間 (應在 500 毫秒內完成)

## 環境設置

### 使用測試設置腳本
```bash
# Windows
scripts/test_setup.bat

# 或者手動設置環境變數
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3306
export TEST_DB_USER=root
export TEST_DB_PASSWORD=password
export TEST_DB_NAME=free2free_test
export TEST_JWT_SECRET=test-jwt-secret-key-32-chars-long-enough!!
export TEST_FACEBOOK_KEY=your_test_facebook_app_id
export TEST_FACEBOOK_SECRET=your_test_facebook_app_secret
```

## 運行測試

### 運行所有測試
```bash
go test ./tests/...
```

### 運行特定類型的測試
```bash
# 單元測試
go test ./tests/unit/...

# 整合測試
go test ./tests/integration/...

# 端到端測試
go test ./tests/e2e/...

# 合約測試
go test ./tests/contract/...

# 效能測試
go test ./tests/performance/...
```

### 帶有詳細輸出和覆蓋率的測試
```bash
go test -v -race ./tests/...
go test -cover ./tests/...
go test -coverprofile=coverage.out ./tests/... && go tool cover -html=coverage.out
```

## 測試報告

測試結果報告器會生成 JSON 和文本格式的報告，包含：
- 測試套件摘要
- 個別測試結果
- 執行時間統計
- 失敗測試的詳細資訊

## 測試配置

測試配置由 `tests/testutils/config.go` 管理，支援環境變數覆蓋。主要配置包括：
- 資料庫連接
- JWT 密鑰
- Facebook OAuth 憑證
- 伺服器端口

## 清理測試數據

測試清理機制在 `tests/testutils/test_cleanup.go` 中實現，確保每次測試後數據庫和資源都被恰當清理。

## 本地測試執行

1. 確保 MariaDB 服務正在運行
2. 設置測試環境變數
3. 執行測試套件
4. 檢查測試報告

## 故障排除

如果遇到測試失敗：
1. 確保測試數據庫存在且可訪問
2. 檢查環境變數設置
3. 確認 JWT 密鑰長度至少為 32 字符
4. 查看測試日誌獲取詳細錯誤訊息

## 測試標準

- 測試應遵循快速、獨立、可重複的原則
- 每個測試應測試單一功能
- 應覆蓋正常情況、邊界情況和錯誤情況
- 測試命名應清楚說明測試內容