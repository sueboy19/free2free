# Facebook OAuth 設定指南

## 概述
本文件記錄了 Facebook OAuth 登入功能設置的完整流程和常見問題解決方案。

## Facebook App 設定

### 1. 基本設定 (Settings > Basic)

**App Domains:**
```
89z5cmtz-3000.asse.devtunnels.ms
89z5cmtz-8080.asse.devtunnels.ms
```

**Site URL:**
```
https://89z5cmtz-3000.asse.devtunnels.ms
```

**App Contact Email:**
- 設置有效的聯繫郵箱

### 2. Facebook Login 設定 (Facebook Login > Settings)
使用案例 ->自訂 -> 設定

**有效的 OAuth 重新導向 URI:**
```
https://89z5cmtz-8080.asse.devtunnels.ms/auth/facebook/callback
```

**登入 URI:**
```
https://89z5cmtz-3000.asse.devtunnels.ms
```

**用戶端 OAuth 設定:**
- 確認「有效的 OAuth 重新導向 URI」啟用
- 確認「用戶端 OAuth 設定」啟用

**用戶端 OAuth 進階設定:**
- **Use Strict Mode for Redirect URIs**: 設為 **OFF**
- **Force Web OAuth Reauthentication**: 設為 **OFF**
- **Login URIs**: 包含前端網址

### 3. 應用程式狀態
- **App Mode**: Development (開發模式)
- **App Status**: 確保沒有審核中的變更

## 環境變數設定

### 後端 (.env)
```bash
# 應用程式基礎 URL (必須使用外部可訪問的 URL)
BASE_URL=https://89z5cmtz-8080.asse.devtunnels.ms

# Facebook OAuth 設定
FACEBOOK_KEY=1455559675710916
FACEBOOK_SECRET=98d4b91f07d77c750fbcdf975a51add9
```

### 前端 (frontend/.env)
```bash
# API 基礎URL (不包含結尾斜線)
VITE_API_BASE_URL=https://89z5cmtz-8080.asse.devtunnels.ms
```

## 數據庫設定

### 必需的表格

**users 表結構:**
```sql
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  social_id VARCHAR(255) UNIQUE NOT NULL,
  social_provider VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  avatar_url VARCHAR(500),
  is_admin BOOLEAN DEFAULT FALSE,
  created_at BIGINT,
  updated_at BIGINT,
  UNIQUE KEY social_provider (social_id, social_provider)
);
```

**refresh_tokens 表結構:**
```sql
CREATE TABLE refresh_tokens (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  token VARCHAR(255) NOT NULL,
  expires_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_expires_at (expires_at)
);
```

## 常見問題與解決方案

### 問題 1: "重定向 URI 並未列入應用程式用戶端 OAuth 設定的許可名單中"

**原因分析:**
1. BASE_URL 設定不匹配
2. 協議不一致 (http vs https)
3. Facebook 設定不完整

**解決方案:**
1. 確保 BASE_URL 設定為外部可訪問的 URL
2. 統一使用 https 協議
3. 重新檢查 Facebook 設定，等待 5-10 分鐘生效

### 問題 2: "Unknown column 'is_admin' in 'SET'"

**原因分析:**
數據庫結構與代碼模型不匹配

**解決方案:**
```sql
ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;
```

### 問題 3: "Table 'refresh_tokens' doesn't exist"

**解決方案:**
執行 `refresh_tokens` 表創建 SQL 語句

### 問題 4: 前端顯示 JSON 響應而非正常跳轉

**原因分析:**
後端返回 JSON 響應到彈窗，但前端期望接收 postMessage 事件

**解決方案:**
修改 OAuth 回調函數，返回 HTML 頁面使用 `window.opener.postMessage()` 通信

## 網址架構說明

- **前端應用**: `https://89z5cmtz-3000.asse.devtunnels.ms/`
- **後端 API**: `https://89z5cmtz-8080.asse.devtunnels.ms/`
- **Facebook 重定向流程**:
  1. 前端 → 開啟後端 `/auth/facebook` 彈窗
  2. 後端 → 重定向到 Facebook 登入
  3. Facebook → 重定向到後端 `/auth/facebook/callback`
  4. 後端 → 返回 HTML 頁面，使用 postMessage 通知前端
  5. 前端 → 收到消息後關閉彈窗並更新狀態

## 測試清單

- [ ] Facebook App 設定完整
- [ ] 環境變數正確配置
- [ ] 數據庫表格結構完整
- [ ] 後端服務正常運行
- [ ] 前端服務正常運行
- [ ] Facebook 登入流程測試
- [ ] 登入後用戶信息顯示正常
- [ ] 登出功能正常

## 設定檢查要點

### 1. URL 一致性
- 前端 .env 中的 VITE_API_BASE_URL
- 後端 .env 中的 BASE_URL
- Facebook 設定中的所有 URL

### 2. 協議一致性
- 所有 URL 必須使用相同協議 (https)
- Site URL 必須使用 https

### 3. 域名設定
- App Domains 必須包含所有使用的域名
- Facebook Login 設定中的 URI 必須精確匹配

### 4. 設定生效時間
- Facebook 設定修改後需要 5-10 分鐘生效
- 環境變數修改後需要重啟服務

## 部署注意事項

1. **生產環境 URL 替換**
   - 將所有 devtunnels.ms URL 替換為實際域名
   - 更新 Facebook 設定中的 URL
   - 更新環境變數

2. **安全考量**
   - 確保所有通信使用 HTTPS
   - 定期更新 Facebook App Secret
   - 監控登入失敗次數

3. **備份設定**
   - 記錄 Facebook App 設定
   - 備份環境變數
   - 保存數據庫結構

---

**最後更新**: 2025-11-10
**測試狀態**: ✅ 已通過完整測試
**版本**: v1.0