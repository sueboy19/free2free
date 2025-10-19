# Quickstart: Facebook 登入與 API 測試套件

## 環境設置

1. **設置 Facebook 應用**：
   - 在 Facebook 開發者中心創建應用
   - 設置有效的重定向 URI: `http://localhost:8080/auth/facebook/callback`
   - 獲取應用 ID 和密鑰

2. **配置環境變量**：
   ```bash
   cp .env.example .env
   # 編輯 .env 文件並添加 Facebook 應用憑證
   FACEBOOK_KEY=your_facebook_app_id
   FACEBOOK_SECRET=your_facebook_app_secret
   ```

3. **確保數據庫正在運行**：
   ```bash
   docker-compose up -d
   ```

## 運行測試

1. **運行完整的 Facebook 登入測試**：
   ```bash
   go test ./tests/e2e/fb_login_e2e_test.go -v
   ```

2. **運行 OAuth 整合測試**：
   ```bash
   go test ./tests/integration/fb_auth_integration_test.go -v
   ```

3. **運行 API 功能測試**：
   ```bash
   go test ./tests/integration/api_integration_test.go -v
   ```

## 主要測試流程

1. **Facebook 登入測試**：
   - 啟動本地服務器
   - 訪問 Facebook 登入端點
   - 模擬 OAuth 重定向和回調流程
   - 驗證 JWT token 生成

2. **API 訪問測試**：
   - 使用獲得的 JWT token
   - 訪問所有受保護的 API 端點
   - 驗證響應和權限

3. **端到端測試**：
   - 完整的用戶流程測試
   - 從登入到執行所有 API 操作

## 預期結果

- Facebook 登入流程在本地環境中 100% 成功
- JWT token 正確生成且有效
- 所有 API 端點可使用 JWT token 成功訪問
- 測試套件在 5 分鐘內完成