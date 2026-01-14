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

請參考 ../API.md

## 遷移進度

- ✅ 階段 1：基礎架構設置
- 🚧 階段 2：資料層遷移
- ⬜ 階段 3：認證系統遷移
- ⬜ 階段 4：路由處理器遷移
- ⬜ 階段 5：測試遷移
- ⬜ 階段 6：部署與驗證

## 授權

MIT
