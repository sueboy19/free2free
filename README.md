# 買一送一配對網站

這是一個使用 Go 語言開發的買一送一配對網站。

## 功能特色
- 使用者可以透過 Facebook 或 Instagram 登入
- 管理者可以建立配對活動與地點
- 使用者可以建立配對局或加入他人建立的配對局
- 開局者可以審核參與者
- 配對完成後可互相評分與留言
- 評論可點讚或倒讚

## 技術架構
- 後端：Go 1.25 + Gin 框架
- 資料庫：MariaDB (透過 Docker) + GORM
- OAuth 認證：Goth 套件
- Session 管理：Gorilla Sessions

## 安裝與設定

### 環境變數
需要設定以下環境變數：
- `SESSION_KEY` - Session 加密金鑰
- `DB_USER` - 資料庫使用者名稱
- `DB_PASSWORD` - 資料庫密碼
- `DB_NAME` - 資料庫名稱
- `FACEBOOK_KEY` - Facebook OAuth 應用程式金鑰
- `FACEBOOK_SECRET` - Facebook OAuth 應用程式密鑰
- `INSTAGRAM_KEY` - Instagram OAuth 應用程式金鑰
- `INSTAGRAM_SECRET` - Instagram OAuth 應用程式密鑰
- `BASE_URL` - 應用程式基礎 URL (例如: http://localhost:8080)

可以複製 `.env.example` 檔案為 `.env` 並填入相應的值：
```bash
cp .env.example .env
```

### 資料庫設定
本專案使用 Docker Compose 來建立 MariaDB 資料庫環境，並使用 GORM 進行自動遷移。

#### 使用 Docker Compose 建立資料庫
1. 確保已安裝 Docker 和 Docker Compose
2. 在專案根目錄執行以下命令啟動資料庫：
   ```bash
   docker-compose up -d
   ```
3. 資料庫將在埠 3306 上運行
4. 資料庫名稱: `free2free`
5. 使用者名稱: `free2free_user`
6. 密碼: `free2free_password`

GORM 會在應用程式啟動時自動建立所需的資料表結構，無需手動執行 SQL 語句。

#### 手動建立資料庫 (可選)
如果您不想使用 Docker，也可以手動建立 MariaDB 資料庫，應用程式啟動時會自動建立資料表結構。

### 安裝相依套件
```bash
go mod tidy
```

### 執行應用程式
```bash
# 普通運行
go run .

# 使用air進行開發（支持熱重載）
air

# 或者直接運行編譯後的程序
go build
./free2free.exe  # Windows
# 或者
./free2free      # Linux/Mac
```

## API 端點

### OAuth 認證
- `GET /auth/:provider` - 開始 OAuth 認證流程
- `GET /auth/:provider/callback` - OAuth 認證回調
- `GET /logout` - 登出

### 使用者相關
- `GET /profile` - 取得使用者資訊 (需登入)

## 專案結構
- `main.go` - 應用程式入口點
- `main_test.go` - 測試設定
- `go.mod`, `go.sum` - 相依套件管理
- `database_design.md` - 資料庫設計文件
- `security_design.md` - 資訊安全設計文件
- `analysis_worklist.md` - 需求分析與工作列表

## API 文件
應用程式集成了 Swagger UI，可以在以下地址訪問 API 文檔：
```
http://localhost:8080/swagger/index.html
```

## 開發指南
1. 確保所有環境變數都已正確設定
2. 建立資料庫並執行 `docker-compose up -d` 啟動資料庫
3. 使用 `air` 啟動開發伺服器（支持熱重載）
4. 開發新功能時請參考 `analysis_worklist.md` 的工作列表
5. 遵循 `security_design.md` 中的安全規範
6. 如需進行 OAuth 測試，請參考 `VSCODE.md` 中的反向代理設定說明