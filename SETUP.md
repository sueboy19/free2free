# Free2Free 系統設定指南

本文說明如何從頭開始設定 Free2Free 系統，包括 Facebook Login、HTTPS 處理和資料庫設定。

## 前置需求

- Go 1.25 或更高版本
- Docker 和 Docker Compose（用於 MariaDB）
- Node.js 和 npm（用於前端）
- VS Code（用於 Port Forwarding）
- Facebook 開發者帳號

---

## 步驟一：產生安全金鑰

系統需要兩個安全金鑰：`SESSION_KEY` 和 `JWT_SECRET`。

### 1.1 使用金鑰產生腳本

在專案根目錄執行：

```bash
go run generate_keys.go
```

### 1.2 產生的金鑰範例

```
SESSION_KEY=JIcrKHinkYwsRoKFmQNIqO5a78Q8sRFyV8JJQJAeuXnSbcyhk9R29dIbmrRjNmWKuuApQnZJHIGa39KFRNT0wA==
JWT_SECRET=9gz_YSoUOPGli_QvOG3g2XS-5O6DtAspr9LSsY5w6sKDz0grrNQPXNb7lLgVJljxmA0CeaYP0yrCxCNsnyGThg==
```

**重要**：請妥善保管這些金鑰，不要提交到版本控制系統。

---

## 步驟二：設定環境變數

### 2.1 建立環境變數檔案

```bash
# 在專案根目錄
cp .env.example .env

# 前端環境變數
cd frontend
cp .env.example .env
cd ..
```

### 2.2 編輯 `.env` 檔案

編輯專案根目錄的 `.env` 檔案：

```bash
# 資料庫設定
DB_HOST=localhost
DB_USER=free2free_user
DB_PASSWORD=free2free_password
DB_NAME=free2free

# Session 金鑰（從步驟一取得）
SESSION_KEY=JIcrKHinkYwsRoKFmQNIqO5a78Q8sRFyV8JJQJAeuXnSbcyhk9R29dIbmrRjNmWKuuApQnZJHIGa39KFRNT0wA==

# JWT 金鑰（從步驟一取得）
JWT_SECRET=9gz_YSoUOPGli_QvOG3g2XS-5O6DtAspr9LSsY5w6sKDz0grrNQPXNb7lLgVJljxmA0CeaYP0yrCxCNsnyGThg==

# Facebook OAuth 設定
FACEBOOK_KEY=your_facebook_app_id
FACEBOOK_SECRET=your_facebook_app_secret

# Instagram OAuth 設定（可選）
INSTAGRAM_KEY=your_instagram_app_key
INSTAGRAM_SECRET=your_instagram_app_secret

# 應用程式基礎 URL（從步驟三取得）
BASE_URL=http://localhost:8080

# Cookie 安全性（使用 HTTPS 時設為 true）
SECURE_COOKIE=false

# 資料庫自動遷移（開發時建議設為 true）
AUTO_MIGRATE=true

# 前端 URL（CORS 設定用）
FRONTEND_URL=http://localhost:3000

# Gin 執行模式（debug/release）
#GIN_MODE=release
```

**注意**：此時 `BASE_URL` 設為 `http://localhost:8080`，稍後在取得 HTTPS URL 後需要更新。

### 2.3 編輯前端 `.env` 檔案

編輯 `frontend/.env` 檔案：

```bash
# API 基礎URL（稍後會更新為 HTTPS URL）
VITE_API_BASE_URL=http://localhost:8080

VITE_APP_TITLE=買一送一配對網站
VITE_APP_VERSION=1.0.0
```

---

## 步驟三：設定 Facebook OAuth

### 3.1 取得 Facebook 憑證

1. 前往 [Facebook Developers](https://developers.facebook.com/apps/)
2. 選擇或建立應用程式
3. 啟用「Facebook 登入」產品
4. 複製應用程式編號（`FACEBOOK_KEY`）和應用程式金鑰（`FACEBOOK_SECRET`）

### 3.2 設定 OAuth 回調 URL

**重要**：Facebook 要求回調 URL 必須是 HTTPS。我們使用 VS Code Port Forwarding 來取得 HTTPS URL。

#### 在 【使用案例】>【自訂】> 【Facebook 登入】> 設定 
【有效的 Oauth 重新導向 URI】 裡面設定  
```
https://36pdmllw-8080.asse.devtunnels.ms/auth/facebook/callback
```
設定完後，記得最底下存檔  

##### **用戶端 OAuth 設定:**
- 確認「有效的 OAuth 重新導向 URI」啟用
- 確認「用戶端 OAuth 設定」啟用

##### **用戶端 OAuth 進階設定:**
- **Use Strict Mode for Redirect URIs**: 設為 **OFF**
- **Force Web OAuth Reauthentication**: 設為 **OFF**
- **Login URIs**: 包含前端網址

#### 選項 A：使用 VS Code Port Forwarding（推薦）

**步驟 1：啟動後端**
```bash
# 終端機 1
go run main.go
```

**步驟 2：設定 VS Code Port Forwarding**
1. 開啟「出口」面板（Ctrl+Shift+E）
2. 切換到「PORTS」標籤頁
3. 找到 8080 連接埠
4. 點選「Visibility」欄位
5. 將「localhost」改成「Public」
6. 複製顯示的公開 HTTPS URL，例如：
   ```
   https://abc123-8080.asse.devtunnels.ms
   ```

**步驟 3：更新 Facebook 應用程式**
1. 在 Facebook Developers，進入您的應用程式
2. 前往【使用案例】>【自訂】> 【Facebook 登入】→「Facebook 登入」→「設定」
3. 在「允許的重新導向 URL」中加入：
   ```
   https://abc123-8080.asse.devtunnels.ms/auth/facebook/callback
   ```
4. 點選「儲存變更」

**步驟 4：更新環境變數**
更新 `.env` 和 `frontend/.env` 中的 URL：
```bash
# .env
BASE_URL=https://abc123-8080.asse.devtunnels.ms
SECURE_COOKIE=true

# frontend/.env
VITE_API_BASE_URL=https://abc123-8080.asse.devtunnels.ms
```

#### 選項 B：使用 Azure Devtunnel

```bash
# 安裝 devtunnel CLI
# 從 https://aka.ms/devtunnel/cli 下載

# 啟動後端
go run main.go

# 在新終端機執行
devtunnel host --protocol=https 8080

# 複製顯示的 URL
```

然後將 URL 設定到 Facebook 應用程式和環境變數中。

#### 選項 C：使用 ngrok

```bash
# 安裝 ngrok
# 從 https://ngrok.com/download 下載

# 執行
ngrok http 8080

# 複製顯示的 HTTPS URL
```

---

## 步驟四：啟動資料庫

### 4.1 使用 Docker Compose 啟動 MariaDB

```bash
docker-compose up -d
```

**注意**：
- 服務名稱為 `mariadb`，容器名稱為 `free2free-mariadb`
- 開發環境使用 `localhost` 連線（透過連接埠轉發 `3306:3306`）
- Staging 環境使用 `mariadb` 連線（Docker 內部網路）

### 4.2 驗證資料庫

```bash
# 檢查容器狀態
docker-compose ps

# 進入資料庫（選用）
docker exec -it free2free-mariadb mariadb -u free2free_user -p
# 密碼：free2free_password

# 在 MariaDB 中
SHOW DATABASES;
USE free2free;
SHOW TABLES;
EXIT;
```

---

## 步驟五：啟動應用程式

### 5.1 啟動後端

```bash
# 方式一：使用 air（熱重載，開發推薦）
air

# 方式二：直接執行
go run main.go

# 方式三：編譯後執行
go build -o free2free .
./free2free  # Windows: free2free.exe
```

### 5.2 啟動前端

```bash
cd frontend
npm install  # 首次執行需要
npm run dev
```

前端會運行在 `http://localhost:3000`（或 Vite 分配的連接埠）。

#### 如果有 build
```
npx serve dist/
```

### 5.3 設定 VS Code Port Forwarding

如果還沒設定，參考步驟三中的說明。

---

## 步驟六：測試系統

### 6.1 測試後端 API

1. 訪問 Swagger UI：
   ```
   https://your-https-url/swagger/index.html
   ```

2. 測試健康檢查：
   ```bash
   curl https://your-https-url/health
   ```

### 6.2 測試 Facebook Login

1. 在瀏覽器訪問：
   ```
   https://your-https-url/auth/facebook
   ```

2. 會重新導向到 Facebook 登入頁面

3. 登入並授權應用程式

4. 成功後，訪問取得 JWT token：
   ```
   https://your-https-url/auth/token
   ```

5. 在 Swagger UI 使用 token：
   - 點擊「Authorize」按鈕
   - 輸入：`Bearer <your_jwt_token>`
   - 測試 `/profile` 端點

### 6.3 測試前端

1. 開啟 `http://localhost:3000`
2. 點擊「使用 Facebook 登入」按鈕
3. 確認能順利完成登入流程

---

## 常見問題排除

### 問題 1：Facebook 回調錯誤「URL 已封鎖」

**原因**：重新導向 URL 未在 Facebook 應用程式中設定

**解決**：
1. 確認在 Facebook 應用程式中加入了完整的 HTTPS 回調 URL
2. URL 格式必須完全正確：`https://your-url/auth/facebook/callback`

### 問題 2：CORS 錯誤

**原因**：前端 URL 未在 CORS 允許清單中

**解決**：
- 在 `main.go` 中的 CORS 中介層加入前端 URL
- 或設定環境變數 `FRONTEND_URL=http://localhost:3000`

### 問題 3：Session 錯誤

**原因**：Session 金鑰太短或未設定

**解決**：
- 確認 `SESSION_KEY` 至少 32 個字符
- 重新啟動應用程式

### 問題 4：資料庫連線失敗

**原因**：Docker 容器未啟動或連線設定錯誤

**解決**：
```bash
# 檢查容器狀態
docker-compose ps

# 重啟容器
docker-compose down
docker-compose up -d

# 檢查容器日誌
docker-compose logs mysql
```

### 問題 5：VS Code Port Forwarding URL 失效

**原因**：重啟 VS Code 或關閉 Port Forwarding

**解決**：
- 重新設定 Port Forwarding
- 更新 Facebook 應用程式中的回調 URL（如果 URL 改變）
- 更新 `.env` 中的 `BASE_URL`

---

## 快速啟動指令（設定完成後）

```bash
# 終端機 1：啟動資料庫
docker-compose up -d

# 終端機 2：啟動後端
go run main.go

# 終端機 3：啟動前端
cd frontend
npm run dev

# 在 VS Code 中設定 Port Forwarding：
# 1. Ctrl+Shift+E 開啟「出口」面板
# 2. 切換到「PORTS」標籤頁
# 3. 找到 8080 連接埠
# 4. 點選「Visibility」→ 改為「Public」
```

---

## 環境變數對應表

| 變數 | 來源 | 範例值 |
|------|------|--------|
| `SESSION_KEY` | `generate_keys.go` 產生 | `JIcrKHinkYwsRoKFmQNIqO5a78Q8sRFyV8JJQJAeuXnSbcyhk9R29dIbmrRjNmWKuuApQnZJHIGa39KFRNT0wA==` |
| `JWT_SECRET` | `generate_keys.go` 產生 | `9gz_YSoUOPGli_QvOG3g2XS-5O6DtAspr9LSsY5w6sKDz0grrNQPXNb7lLgVJljxmA0CeaYP0yrCxCNsnyGThg==` |
| `FACEBOOK_KEY` | Facebook Developers | `0000000000000000` |
| `FACEBOOK_SECRET` | Facebook Developers | `000000000000000000000` |
| `BASE_URL` | VS Code Port Forwarding | `https://abc123-8080.asse.devtunnels.ms` |
| `SECURE_COOKIE` | 手動設定 | `true`（使用 HTTPS 時） |
| `VITE_API_BASE_URL` | 與 `BASE_URL` 相同 | `https://abc123-8080.asse.devtunnels.ms` |

---

## HTTPS 解決方案比較

| 方案 | 優點 | 缺點 | 推薦度 |
|------|------|------|--------|
| **VS Code Port Forwarding** | 整合在 VS Code、最方便、免費 | 需要保持 VS Code 開啟 | ⭐⭐⭐ 推薦 |
| **Azure Devtunnel** | 免費、穩定、由 Microsoft 提供 | 需要安裝額外 CLI | ⭐⭐ 備選 |
| **ngrok** | 功能強大、有額外功能 | 免費版有限制、需註冊 | ⭐ 備選 |

---

## 安全性建議

1. **金鑰管理**：
   - 不要將 `.env` 檔案提交到版本控制系統
   - 使用不同的金鑰於開發、測試、生產環境
   - 定期更換金鑰

2. **Facebook 憑證**：
   - 妥善保管 `FACEBOOK_SECRET`
   - 定期檢查 Facebook 應用程式的安全性設定
   - 啟用 Facebook 應用程式的「應用程式密鑰」驗證

3. **Cookie 安全性**：
   - 生產環境必須設定 `SECURE_COOKIE=true`
   - 確保使用 HTTPS
   - 設定適當的 `SameSite` 屬性

4. **資料庫**：
   - 使用強壯的密碼
   - 限制資料庫只接受本地連線
   - 定期備份資料

---

## 相關文件

- [CLAUDE.md](CLAUDE.md) - 專案架構和開發指南
- [VSCODE.md](VSCODE.md) - VS Code 開發設定
- [docs/deployment.md](docs/deployment.md) - 部署指南

---

## 下一步

設定完成後，您可以：

1. 開始開發新功能
2. 查看 [CLAUDE.md](CLAUDE.md) 了解專案架構
3. 執行測試：`go test ./tests/... -v`
4. 查看Swagger UI 文件測試 API

---

## 聯絡與支援

如有問題，請檢查：
1. 專案 Issues
2. 程式碼註解
3. API 文件（Swagger UI）
