# Free2Free API - Cloudflare Workers

買一送一配對網站的 Cloudflare Workers 後端 API。

## 技術棧

- **框架**: Hono
- **語言**: TypeScript
- **資料庫**: Cloudflare D1 (SQLite)
- **存儲**: Cloudflare KV
- **部署**: Cloudflare Workers

## 開發環境設置

### 前置要求

- Node.js 18+
- npm 或 yarn
- Wrangler CLI

### 安裝

```bash
# 安裝依賴
npm install

# 安裝 Wrangler CLI
官方建議
npm i -D wrangler@latest

# 登入 Cloudflare
wrangler login
```

### 本地開發

```bash
# 啟動開發伺服器
npm run dev

# 運行測試
npm run test

# 運行 lint
npm run lint
```

### 環境變數

在使用 `wrangler secret put` 設置以下 secrets：

```bash
wrangler secret put JWT_SECRET
wrangler secret put SESSION_KEY
wrangler secret put FACEBOOK_KEY
wrangler secret put FACEBOOK_SECRET
wrangler secret put INSTAGRAM_KEY
wrangler secret put INSTAGRAM_SECRET
```

### 部署

```bash
# 部署到 Cloudflare Workers
npm run deploy
```

## 資料庫

### 本地開發

使用 Miniflare 本地模擬 D1 資料庫：

```bash
wrangler dev
```

### 創建資料庫

```bash
# 創建 D1 資料庫
wrangler d1 create free2free-db

# 記錄 database_id 並更新 wrangler.toml
```

### 執行 Migration

```bash
# 執行資料庫 schema migration
wrangler d1 execute free2free-db --file=./migrations/0001_initial.sql

# 查看資料表
wrangler d1 execute free2free-db --command="SELECT name FROM sqlite_master WHERE type='table';"
```

### 匯入測試資料

```bash
# 匯入測試資料到 D1
wrangler d1 execute free2free-db --file=./scripts/import-to-d1.sql
```

### 資料表結構

- `users` - 使用者資料
- `admins` - 管理員資料
- `locations` - 地點資料
- `activities` - 活動資料
- `matches` - 配對局資料
- `match_participants` - 參與者資料
- `reviews` - 評分資料
- `review_likes` - 評分點讚資料
- `refresh_tokens` - 重新整理 token 資料
- `sessions` - Session 資料

## 認證

### OAuth 登入流程

1. 獲取 OAuth 授權 URL
   ```
   GET /auth/:provider
   ```

2. 用戶授權後，系統回調
   ```
   GET /auth/:provider/callback?code=...
   ```

3. 返回 JWT token 和 session

### JWT Token

- **Access Token**: 15 分鐘過期
- **Refresh Token**: 7 天過期

### 使用 Token

在請求頭中添加 Authorization：
```
Authorization: Bearer <access_token>
```

### 刷新 Token

```
POST /auth/refresh
{
  "refresh_token": "<refresh_token>"
}
```

### 登出

```
POST /auth/logout
{
  "refresh_token": "<refresh_token>",
  "session_id": "<session_id>"
}
```

## 專案結構

```
src/
├── lib/           # 工具函數（db, kv, jwt, oauth）
├── routes/        # API 路由處理器
├── middleware/    # 中介層（cors, auth, error）
├── types/         # TypeScript 類型定義
└── index.ts       # 主入口
migrations/        # 資料庫 migration 檔案
scripts/          # 腳本（資料遷移、匯入）
test/             # 測試檔案
```

## API 文檔

### Auth 路由

#### OAuth 登入
- `GET /auth/:provider` - 獲取 OAuth 授權 URL
- `GET /auth/:provider/callback?code=...` - OAuth 回調
- `POST /auth/refresh` - 刷新 token
- `POST /auth/logout` - 登出
- `GET /auth/me` - 獲取當前用戶

### Admin 路由

#### Locations
- `POST /admin/locations` - 創建地點
- `GET /admin/locations` - 列出所有地點
- `GET /admin/locations/:id` - 獲取地點詳情
- `PUT /admin/locations/:id` - 更新地點
- `DELETE /admin/locations/:id` - 刪除地點

#### Activities
- `POST /admin/activities` - 創建活動
- `GET /admin/activities` - 列出所有活動
- `GET /admin/activities/:id` - 獲取活動詳情
- `PUT /admin/activities/:id` - 更新活動
- `DELETE /admin/activities/:id` - 刪除活動

### User 路由

#### Matches
- `GET /matches` - 獲取公開的開放配對局
- `GET /user/matches` - 獲取用戶的配對局
- `GET /matches/:id` - 獲取配對局詳情
- `GET /matches/:id/participants` - 獲取參與者列表
- `GET /matches/:id/reviews` - 獲取評分列表

### Organizer 路由

#### Matches
- `POST /matches` - 創建配對局
- `PUT /matches/:id/status` - 更新配對局狀態
- `POST /matches/:id/join` - 加入配對局
- `PUT /matches/:matchId/participants/:participantId` - 審核參與者
- `DELETE /matches/:id/join` - 離開配對局

### Review 路由

#### Reviews
- `POST /reviews` - 創建評分
- `GET /reviews` - 獲取評分列表
- `POST /reviews/:id/like` - 點讚/倒讚評分
- `DELETE /reviews/:id` - 刪除評分

請參考 ../API.md 詳細的 API 文檔

## 遷移進度

- ✅ 階段 1：基礎架構設置
- ✅ 階段 2：資料層遷移
- ✅ 階段 3：認證系統遷移
- ✅ 階段 4：API 路由實現
- ⬜ 階段 5：測試與部署

## 授權

MIT
