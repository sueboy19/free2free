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
npm i -D wrangler@latest

# 登入 Cloudflare
wrangler login
```

### 本地開發

```bash
# 啟動本地開發伺服器（會自動建立本地 D1 資料庫）
wrangler dev

# 或者使用 npm scripts
npm run dev

# 運行測試
npm run test

# 運行 lint
npm run lint
```

### 環境設定

本專案支援三個環境：**本地開發**、**Staging**、**Production**。

#### 1. 本地開發（Development）

**資料庫**：
- 使用**純本地 D1 資料庫**（在 `.wrangler` 目錄中）
- `wrangler dev` 首次執行時會自動建立，無需手動建立
- 使用 `--local` 參數操作本地資料庫

**Secrets 管理**：
- 使用 `.dev.vars` 檔案（已在 `.gitignore` 中，不會上 git）
- 建立 `.dev.vars`（參考 `.dev.vars.example`）：

```bash
JWT_SECRET=your_local_jwt_secret_at_least_32_characters_long
SESSION_KEY=your_local_session_secret_key_at_least_32_characters
FACEBOOK_KEY=your_local_facebook_app_key
FACEBOOK_SECRET=your_local_facebook_app_secret
INSTAGRAM_KEY=your_local_instagram_app_key
INSTAGRAM_SECRET=your_local_instagram_app_secret
```

**重要配置**：

⚠️ **本地開發必須設定 `database_id`**：
- 在 `wrangler.toml` 的頂層 `[[d1_databases]]` 中必須設定 `database_id`
- 本地開發時 `database_id` 可以是任意值（如隨機生成的 UUID 或自定義標識）
- `--local` 模式下 `database_id` 只作識別用，實際使用本地 `.wrangler` 目錄中的 SQLite 檔案
- **不設定 `database_id` 會導致 `wrangler d1 execute --local` 無法找到資料庫**

```toml
[[d1_databases]]
binding = "DB"
database_name = "free2free-db"
database_id = "local-dev-placeholder"  # 任意值均可
```

**常用指令**：

```bash
# 啟動本地開發伺服器（自動建立本地 D1 資料庫）
wrangler dev

# 執行本地資料庫遷移
wrangler d1 execute DB --local --file=./migrations/0001_initial.sql

# 查詢本地資料庫
wrangler d1 execute DB --local --command="SELECT * FROM users LIMIT 10;"
```

#### 2. Staging 環境

**資料庫**：
- 需要先建立雲端 D1 資料庫

```bash
# 建立 staging D1 資料庫
wrangler d1 create free2free-db-staging

# 將 database_id 填入 wrangler.toml 的 [env.staging.d1_databases].database_id
```

**Secrets 管理**：
- 使用 `wrangler secret put` 將 secrets 加密存儲到 Cloudflare 雲端

```bash
# 設定 staging 環境 secrets
wrangler secret put JWT_SECRET --env staging
wrangler secret put SESSION_KEY --env staging
wrangler secret put FACEBOOK_KEY --env staging
wrangler secret put FACEBOOK_SECRET --env staging
wrangler secret put INSTAGRAM_KEY --env staging
wrangler secret put INSTAGRAM_SECRET --env staging
```

**常用指令**：

```bash
# 執行 staging 資料庫遷移
wrangler d1 execute free2free-db-staging --remote --file=./migrations/0001_initial.sql

# 部署到 staging（重要：必須加 --env staging）
wrangler deploy --env staging
```

#### 3. Production 環境

**資料庫**：
- 已建立雲端 D1 資料庫（`7e2b5c27-c755-42e5-b349-b035b5df3534`）
- 已在 `wrangler.toml` 的 `[env.production.d1_databases].database_id` 中設定

**Secrets 管理**：
- 使用 `wrangler secret put` 將 secrets 加密存儲到 Cloudflare 雲端

```bash
# 設定 production 環境 secrets
wrangler secret put JWT_SECRET --env production
wrangler secret put SESSION_KEY --env production
wrangler secret put FACEBOOK_KEY --env production
wrangler secret put FACEBOOK_SECRET --env production
wrangler secret put INSTAGRAM_KEY --env production
wrangler secret put INSTAGRAM_SECRET --env production
```

**常用指令**：

```bash
# 執行 production 資料庫遷移
wrangler d1 execute free2free-db --remote --file=./migrations/0001_initial.sql

# 部署到 production（重要：必須加 --env production）
wrangler deploy --env production
```

### 環境對比表

| 項目 | 本地開發 | Staging | Production |
|------|---------|---------|------------|
| **啟動指令** | `wrangler dev` | `wrangler deploy --env staging` | `wrangler deploy --env production` |
| **資料庫位置** | 本地檔案（`.wrangler`） | 雲端 D1 | 雲端 D1 |
| **資料庫建立** | 自動建立 | `wrangler d1 create` | 已建立 |
| **資料庫遷移** | `--local` | `--remote` | `--remote` |
| **Secrets 管理** | `.dev.vars` | `wrangler secret put --env staging` | `wrangler secret put --env production` |
| **Secrets 存放** | 本地檔案 | Cloudflare 雲端（加密） | Cloudflare 雲端（加密） |
| **重開機影響** | 無（本地檔案） | 無（雲端） | 無（雲端） |
| **上 Git** | ❌ 已在 .gitignore | ✅ 安全（雲端） | ✅ 安全（雲端） |

### 重要提醒

⚠️ **`--env` 參數很重要**：
- 部署時**必須指定** `--env staging` 或 `--env production`
- 不指定 `--env` 則使用預設開發環境

⚠️ **Secrets 環境隔離**：
- 每個環境的 secrets 是**完全獨立**的
- staging 和 production 的 OAuth keys、資料庫等應該分開
- 避免使用 production 的 secrets 在 staging 或本地開發

⚠️ **資料庫操作注意**：
- 本地開發用 `--local`
- 雲端環境用 `--remote`
- 避免誤操作 production 資料庫

### 部署

#### 部署到 Staging

```bash
# 部署到 staging 環境（重要：必須加 --env staging）
wrangler deploy --env staging
```

#### 部署到 Production

```bash
# 部署到 production 環境（重要：必須加 --env production）
wrangler deploy --env production
```

#### 部署前檢查清單

在部署到 production 之前，請確認：
- [ ] 所有測試已通過（`npm run test`）
- [ ] 資料庫遷移已在 staging 測試過
- [ ] production 的 secrets 已設定（`wrangler secret put --env production`）
- [ ] CORS_ORIGINS 已更新為生產域名
- [ ] 檢查 wrangler.toml 中的生產環境設定

## 資料庫

### 本地開發

使用 Miniflare 本地模擬 D1 資料庫：

```bash
wrangler dev
```

### 創建資料庫

```bash
# 創建 D1 資料庫
wrangler d1 create free2free-db-staging

# 記錄 database_id 並更新 wrangler.toml 的 [env.staging.d1_databases].database_id
```

### 執行 Migration

#### 本地開發（Local）

```bash
# 執行本地資料庫 schema migration（使用 --local）
wrangler d1 execute free2free-db --local --file=./migrations/0001_initial.sql

# 查看本地資料表
wrangler d1 execute free2free-db --local --command="SELECT name FROM sqlite_master WHERE type='table';"
```

#### Staging 環境

```bash
# 執行 staging 資料庫 schema migration（使用 --remote）
wrangler d1 execute free2free-db-staging --remote --file=./migrations/0001_initial.sql

# 查看資料表
wrangler d1 execute free2free-db-staging --remote --command="SELECT name FROM sqlite_master WHERE type='table';"
```

#### Production 環境

```bash
# 執行 production 資料庫 schema migration（使用 --remote）
wrangler d1 execute free2free-db --remote --file=./migrations/0001_initial.sql

# 查看資料表
wrangler d1 execute free2free-db --remote --command="SELECT name FROM sqlite_master WHERE type='table';"
```

### 匯入測試資料

#### 本地開發

```bash
# 匯入測試資料到本地 D1
wrangler d1 execute free2free-db --local --file=./scripts/import-to-d1.sql
```

#### Staging 環境

```bash
# 匯入測試資料到 staging D1
wrangler d1 execute free2free-db-staging --remote --file=./scripts/import-to-d1.sql
```

⚠️ **Production 環境通常不匯入測試資料**，應使用真實資料。

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

請參考 ../API.md

## 遷移進度

- ✅ 階段 1：基礎架構設置
- ✅ 階段 2：資料層遷移
- ✅ 階段 3：認證系統遷移
- ✅ 階段 4：API 路由實現
- ✅ 階段 5：測試與部署

## 授權

MIT
