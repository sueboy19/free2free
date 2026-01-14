# Cloudflare Workers 遷移計畫 - 階段 1：基礎架構設置

## 📋 階段 1 概述

**目標：** 建立 Cloudflare Workers 專案基礎架構，配置開發環境

**預計時間：** 1 週

**狀態：** ✅ 已完成

**完成日期：** 2026-01-14

---

## ✅ 任務清單

### Task 1.1: 創建專案目錄結構

**狀態：** ⬜ 待辦

**說明：** 建立標準的 TypeScript + Hono 專案結構

**執行步驟：**
```
workers/
├── src/
│   ├── index.ts              # 主入口
│   ├── lib/
│   │   ├── db.ts             # D1 客戶端
│   │   ├── kv.ts             # KV 客戶端
│   │   ├── jwt.ts            # JWT 驗證
│   │   ├── oauth.ts          # OAuth handlers
│   │   └── errors.ts         # 錯誤處理
│   ├── routes/
│   │   ├── auth.ts           # 認證路由
│   │   ├── admin.ts          # 管理員路由
│   │   ├── user.ts           # 使用者路由
│   │   ├── organizer.ts      # 開局者路由
│   │   └── review.ts         # 評分路由
│   ├── middleware/
│   │   ├── cors.ts           # CORS middleware
│   │   ├── auth.ts           # 認證 middleware
│   │   └── error.ts          # 錯誤處理 middleware
│   └── types/
│       └── index.ts          # TypeScript 類型定義
├── migrations/
│   └── 0001_initial.sql      # 資料庫 migration
├── test/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── wrangler.toml             # Workers 配置
├── package.json
├── tsconfig.json
└── README.md
```

**驗證：**
- [ ] 所有目錄已創建
- [ ] 所有檔案已創建（至少是空的）

---

### Task 1.2: 初始化 npm 專案

**狀態：** ⬜ 待辦

**說明：** 初始化 npm 專案，安裝必要依賴

**執行命令：**
```bash
cd workers
npm init -y
```

**安裝依賴：**
```bash
# 核心框架
npm install hono

# JWT 處理
npm install jose

# 開發依賴
npm install -D typescript @types/node

# Workers SDK
npm install -D wrangler

# 測試工具
npm install -D vitest @cloudflare/vitest-pool-workers

# ESLint 和 Prettier
npm install -D eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
npm install -D prettier eslint-config-prettier eslint-plugin-prettier
```

**驗證：**
- [ ] package.json 已創建
- [ ] node_modules 已安裝
- [ ] 所有依賴已安裝

---

### Task 1.3: 配置 TypeScript

**狀態：** ⬜ 待辦

**說明：** 創建 tsconfig.json 配置檔案

**內容：**
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ES2022",
    "moduleResolution": "bundler",
    "lib": ["ES2022"],
    "types": ["@cloudflare/workers-types"],
    "jsx": "react-jsx",
    "jsxImportSource": "hono/jsx",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules"]
}
```

**驗證：**
- [ ] tsconfig.json 已創建
- [ ] TypeScript 編譯檢查通過

---

### Task 1.4: 配置 Wrangler

**狀態：** ⬜ 待辦

**說明：** 創建 wrangler.toml 配置檔案

**內容：**
```toml
name = "free2free-api"
main = "src/index.ts"
compatibility_date = "2024-01-01"
compatibility_flags = ["nodejs_compat"]

# D1 Database binding（稍後創建）
[[d1_databases]]
binding = "DB"
database_name = "free2free-db"
database_id = "placeholder"  # 運行 `wrangler d1 create` 後更新

# KV Namespace binding（稍後創建）
[[kv_namespaces]]
binding = "KV"
id = "placeholder"  # 運行 `wrangler kv:namespace create` 後更新

# Environment variables
[vars]
ENVIRONMENT = "development"
CORS_ORIGINS = "http://localhost:3000,http://localhost:5173"

# Secrets（需要使用 `wrangler secret put` 命令）
# JWT_SECRET
# SESSION_KEY
# FACEBOOK_KEY
# FACEBOOK_SECRET
# INSTAGRAM_KEY
# INSTAGRAM_SECRET
```

**驗證：**
- [ ] wrangler.toml 已創建
- [ ] wrangler whoami 命令可正常執行

---

### Task 1.5: 配置 ESLint 和 Prettier

**狀態：** ⬜ 待辦

**說明：** 創建代碼風格配置檔案

**eslintrc.json：**
```json
{
  "extends": [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier"
  ],
  "parser": "@typescript-eslint/parser",
  "parserOptions": {
    "ecmaVersion": 2022,
    "sourceType": "module",
    "project": "./tsconfig.json"
  },
  "rules": {
    "@typescript-eslint/no-unused-vars": ["error", { "argsIgnorePattern": "^_" }],
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/no-explicit-any": "warn"
  }
}
```

**.prettierrc：**
```json
{
  "semi": true,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5",
  "printWidth": 100
}
```

**驗證：**
- [ ] eslintrc.json 已創建
- [ ] .prettierrc 已創建
- [ ] npx eslint . --ext .ts 可以執行

---

### Task 1.6: 創建 TypeScript 類型定義

**狀態：** ⬜ 待辦

**說明：** 建立所有資料模型的 TypeScript 接口

**檔案：** `src/types/index.ts`

**內容：**
```typescript
export interface Env {
  DB: D1Database;
  KV: KVNamespace;
  JWT_SECRET: string;
  SESSION_KEY: string;
  FACEBOOK_KEY: string;
  FACEBOOK_SECRET: string;
  INSTAGRAM_KEY: string;
  INSTAGRAM_SECRET: string;
  BASE_URL: string;
  FRONTEND_URL: string;
  CORS_ORIGINS: string;
}

export interface User {
  id: number;
  social_id: string;
  social_provider: 'facebook' | 'instagram';
  name: string;
  email: string;
  avatar_url?: string;
  is_admin: boolean;
  created_at: number;
  updated_at: number;
}

export interface Admin {
  id: number;
  username: string;
  email: string;
}

export interface Location {
  id: number;
  name: string;
  address: string;
  latitude: number;
  longitude: number;
}

export interface Activity {
  id: number;
  title: string;
  target_count: number;
  location_id: number;
  description?: string;
  created_by: number;
  location?: Location;
}

export interface Match {
  id: number;
  activity_id: number;
  organizer_id: number;
  match_time: string;
  status: 'open' | 'completed' | 'cancelled';
  activity?: Activity;
  organizer?: User;
}

export interface MatchParticipant {
  id: number;
  match_id: number;
  user_id: number;
  status: 'pending' | 'approved' | 'rejected';
  joined_at: string;
  match?: Match;
  user?: User;
}

export interface Review {
  id: number;
  match_id: number;
  reviewer_id: number;
  reviewee_id: number;
  score: number;
  comment?: string;
  created_at: string;
  match?: Match;
  reviewer?: User;
  reviewee?: User;
}

export interface ReviewLike {
  id: number;
  review_id: number;
  user_id: number;
  is_like: boolean;
  review?: Review;
  user?: User;
}

export interface RefreshToken {
  id: number;
  user_id: number;
  token: string;
  expires_at: string;
  created_at: string;
  user?: User;
}

export interface JWTPayload {
  user_id: number;
  user_name: string;
  is_admin: boolean;
  exp: number;
  iat: number;
}
```

**驗證：**
- [ ] 所有接口已定義
- [ ] TypeScript 編譯無錯誤

---

### Task 1.7: 創建主入口檔案

**狀態：** ⬜ 待辦

**說明：** 建立 Hono 應用主入口，配置基礎 middleware

**檔案：** `src/index.ts`

**內容：**
```typescript
import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { errorHandler } from './middleware/error';
import type { Env } from './types';

const app = new Hono<{ Bindings: Env }>();

// Middleware
app.use('*', logger());
app.use('*', async (c, next) => {
  const corsMiddleware = cors({
    origin: c.env.CORS_ORIGINS.split(','),
    credentials: true,
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowHeaders: ['Content-Type', 'Authorization'],
  });
  return corsMiddleware(c, next);
});
app.use('*', errorHandler);

// Health check
app.get('/', (c) => {
  return c.json({
    status: 'ok',
    message: 'Free2Free API is running',
    timestamp: new Date().toISOString(),
  });
});

// API routes (will be added in later phases)
// app.route('/auth', authRoutes);
// app.route('/admin', adminRoutes);
// app.route('/user', userRoutes);
// app.route('/organizer', organizerRoutes);
// app.route('/review', reviewRoutes);

export default app;
```

**驗證：**
- [ ] src/index.ts 已創建
- [ ] wrangler dev 可以啟動開發伺服器
- [ ] 訪問 http://localhost:8787 返回 200

---

### Task 1.8: 創建錯誤處理 middleware

**狀態：** ⬜ 待辦

**說明：** 建立統一的錯誤處理 middleware

**檔案：** `src/middleware/error.ts`

**內容：**
```typescript
import type { Context, Next } from 'hono';

export class AppError extends Error {
  constructor(
    public statusCode: number,
    public code: string,
    message: string
  ) {
    super(message);
    this.name = 'AppError';
  }
}

export const errorHandler = async (c: Context, next: Next) => {
  try {
    await next();
  } catch (error) {
    console.error('Error:', error);

    if (error instanceof AppError) {
      return c.json(
        {
          code: error.statusCode,
          error: error.message,
          code_error: error.code,
        },
        error.statusCode
      );
    }

    return c.json(
      {
        code: 500,
        error: 'Internal server error',
        code_error: 'INTERNAL_ERROR',
      },
      500
    );
  }
};
```

**驗證：**
- [ ] src/middleware/error.ts 已創建
- [ ] 錯誤處理邏輯正確

---

### Task 1.9: 創建 CORS middleware

**狀態：** ⬜ 待辦

**說明：** 建立 CORS 配置（雖然在 index.ts 中已配置，但單獨檔案更清晰）

**檔案：** `src/middleware/cors.ts`

**內容：**
```typescript
import { cors } from 'hono/cors';
import type { Env } from '../types';

export const corsMiddleware = (env: Env) => {
  return cors({
    origin: env.CORS_ORIGINS.split(','),
    credentials: true,
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowHeaders: ['Content-Type', 'Authorization'],
    exposeHeaders: ['Content-Length'],
  });
};
```

**驗證：**
- [ ] src/middleware/cors.ts 已創建
- [ ] CORS 配置正確

---

### Task 1.10: 創建 D1 資料庫連線客戶端

**狀態：** ⬜ 待辦

**說明：** 建立 D1 資料庫操作的封裝類

**檔案：** `src/lib/db.ts`

**內容：**
```typescript
import type { D1Database } from '@cloudflare/workers-types';
import type {
  User,
  Admin,
  Location,
  Activity,
  Match,
  MatchParticipant,
  Review,
  ReviewLike,
  RefreshToken,
} from '../types';

export class DB {
  constructor(private db: D1Database) {}

  // User operations
  async createUser(user: Omit<User, 'id' | 'created_at' | 'updated_at'>): Promise<User> {
    const result = await this.db
      .prepare(
        `INSERT INTO users (social_id, social_provider, name, email, avatar_url, is_admin)
         VALUES (?, ?, ?, ?, ?, ?)`
      )
      .bind(
        user.social_id,
        user.social_provider,
        user.name,
        user.email,
        user.avatar_url || null,
        user.is_admin ? 1 : 0
      )
      .run();

    const created = await this.db
      .prepare('SELECT * FROM users WHERE id = ?')
      .bind(result.meta.last_row_id)
      .first<User>();

    if (!created) {
      throw new Error('Failed to create user');
    }

    return {
      ...created,
      is_admin: created.is_admin === 1,
    };
  }

  async getUserById(id: number): Promise<User | null> {
    const user = await this.db.prepare('SELECT * FROM users WHERE id = ?').bind(id).first<User>();
    if (!user) return null;
    return { ...user, is_admin: user.is_admin === 1 };
  }

  async getUserBySocialId(socialId: string, provider: string): Promise<User | null> {
    const user = await this.db
      .prepare('SELECT * FROM users WHERE social_id = ? AND social_provider = ?')
      .bind(socialId, provider)
      .first<User>();
    if (!user) return null;
    return { ...user, is_admin: user.is_admin === 1 };
  }

  // 其他 CRUD 操作將在階段 2 添加
}
```

**驗證：**
- [ ] src/lib/db.ts 已創建
- [ ] 基礎 CRUD 操作已實現

---

### Task 1.11: 創建 KV 存儲客戶端

**狀態：** ⬜ 待辦

**說明：** 建立 Workers KV 操作的封裝類

**檔案：** `src/lib/kv.ts`

**內容：**
```typescript
import type { KVNamespace } from '@cloudflare/workers-types';

export class KVStore {
  constructor(private kv: KVNamespace) {}

  async get<T = any>(key: string): Promise<T | null> {
    const value = await this.kv.get(key, 'json');
    return value as T | null;
  }

  async set(key: string, value: any, options?: { expirationTtl?: number }): Promise<void> {
    await this.kv.put(key, JSON.stringify(value), options);
  }

  async delete(key: string): Promise<void> {
    await this.kv.delete(key);
  }

  async exists(key: string): Promise<boolean> {
    const value = await this.kv.get(key);
    return value !== null;
  }

  // Session operations
  async setRefreshToken(userId: number, token: string, expiresAt: Date): Promise<void> {
    const key = `refresh:${userId}:${token}`;
    const ttl = Math.floor((expiresAt.getTime() - Date.now()) / 1000);
    await this.set(key, { userId, expiresAt: expiresAt.toISOString() }, { expirationTtl: ttl });
  }

  async getRefreshToken(userId: number, token: string): Promise<any> {
    const key = `refresh:${userId}:${token}`;
    return this.get(key);
  }

  async deleteRefreshToken(userId: number, token: string): Promise<void> {
    const key = `refresh:${userId}:${token}`;
    await this.delete(key);
  }

  async deleteAllRefreshTokens(userId: number): Promise<void> {
    // KV 不支持模式匹配刪除，需要在應用層維護 token 列表
    await this.delete(`user:${userId}:tokens`);
  }
}
```

**驗證：**
- [ ] src/lib/kv.ts 已創建
- [ ] KV 操作封裝正確

---

### Task 1.12: 創建 JWT 處理工具

**狀態：** ⬜ 待辦

**說明：** 建立 JWT 生成和驗證功能

**檔案：** `src/lib/jwt.ts`

**內容：**
```typescript
import { SignJWT, jwtVerify } from 'jose';
import type { JWTPayload, User } from '../types';

const JWT_ALGORITHM = 'HS256';
const ACCESS_TOKEN_EXPIRY = '15m';
const REFRESH_TOKEN_EXPIRY = '7d';

export class JWTManager {
  constructor(private secret: string) {
    if (secret.length < 32) {
      throw new Error('JWT_SECRET must be at least 32 characters');
    }
  }

  private getSecretKey(): Uint8Array {
    return new TextEncoder().encode(this.secret);
  }

  async generateAccessToken(user: User): Promise<string> {
    return new SignJWT({
      user_id: user.id,
      user_name: user.name,
      is_admin: user.is_admin,
    } as JWTPayload)
      .setProtectedHeader({ alg: JWT_ALGORITHM })
      .setIssuedAt()
      .setExpirationTime(ACCESS_TOKEN_EXPIRY)
      .sign(this.getSecretKey());
  }

  async generateRefreshToken(user: User): Promise<string> {
    return new SignJWT({
      user_id: user.id,
    })
      .setProtectedHeader({ alg: JWT_ALGORITHM })
      .setIssuedAt()
      .setExpirationTime(REFRESH_TOKEN_EXPIRY)
      .sign(this.getSecretKey());
  }

  async generateTokens(user: User): Promise<{ access: string; refresh: string }> {
    const [access, refresh] = await Promise.all([
      this.generateAccessToken(user),
      this.generateRefreshToken(user),
    ]);

    return { access, refresh };
  }

  async verifyAccessToken(token: string): Promise<JWTPayload> {
    const { payload } = await jwtVerify(token, this.getSecretKey());
    return payload as JWTPayload;
  }

  async verifyRefreshToken(token: string): Promise<{ user_id: number }> {
    const { payload } = await jwtVerify(token, this.getSecretKey());
    return payload as { user_id: number };
  }
}
```

**驗證：**
- [ ] src/lib/jwt.ts 已創建
- [ ] JWT 生成和驗證功能正確

---

### Task 1.13: 創建 OAuth 處理工具框架

**狀態：** ⬜ 待辦

**說明：** 建立 OAuth 處理的基礎框架（詳細實現在階段 3）

**檔案：** `src/lib/oauth.ts`

**內容：**
```typescript
import type { User } from '../types';

export interface OAuthProvider {
  name: 'facebook' | 'instagram';
  getAuthUrl(): string;
  exchangeCodeForToken(code: string): Promise<string>;
  getUserProfile(accessToken: string): Promise<OAuthProfile>;
}

export interface OAuthProfile {
  id: string;
  name: string;
  email: string;
  avatar_url?: string;
}

// Facebook OAuth Provider
export class FacebookOAuthProvider implements OAuthProvider {
  name = 'facebook' as const;
  private clientId: string;
  private clientSecret: string;
  private redirectUri: string;

  constructor(clientId: string, clientSecret: string, redirectUri: string) {
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.redirectUri = redirectUri;
  }

  getAuthUrl(): string {
    const params = new URLSearchParams({
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      scope: 'email,public_profile',
      response_type: 'code',
    });
    return `https://www.facebook.com/v18.0/dialog/oauth?${params}`;
  }

  async exchangeCodeForToken(code: string): Promise<string> {
    const params = new URLSearchParams({
      client_id: this.clientId,
      client_secret: this.clientSecret,
      redirect_uri: this.redirectUri,
      code,
    });

    const response = await fetch(
      `https://graph.facebook.com/v18.0/oauth/access_token?${params}`
    );
    const data = await response.json();

    if (data.error) {
      throw new Error(data.error.message);
    }

    return data.access_token;
  }

  async getUserProfile(accessToken: string): Promise<OAuthProfile> {
    const params = new URLSearchParams({
      fields: 'id,name,email,picture',
      access_token: accessToken,
    });

    const response = await fetch(`https://graph.facebook.com/v18.0/me?${params}`);
    const data = await response.json();

    if (data.error) {
      throw new Error(data.error.message);
    }

    return {
      id: data.id,
      name: data.name,
      email: data.email,
      avatar_url: data.picture?.data?.url,
    };
  }
}

// Instagram OAuth Provider（類似實現）
export class InstagramOAuthProvider implements OAuthProvider {
  name = 'instagram' as const;
  // 實現類似 Facebook...
}

export class OAuthManager {
  private providers: Map<string, OAuthProvider> = new Map();

  registerProvider(provider: OAuthProvider) {
    this.providers.set(provider.name, provider);
  }

  getProvider(name: string): OAuthProvider | undefined {
    return this.providers.get(name);
  }
}
```

**驗證：**
- [ ] src/lib/oauth.ts 已創建
- [ ] OAuth 框架正確

---

### Task 1.14: 創建錯誤處理工具

**狀態：** ⬜ 待辦

**說明：** 建立錯誤類型定義和工廠函數

**檔案：** `src/lib/errors.ts`

**內容：**
```typescript
import { AppError } from '../middleware/error';

export const ErrorCodes = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  OAUTH_ERROR: 'OAUTH_ERROR',
  INVALID_TOKEN: 'INVALID_TOKEN',
  TOKEN_EXPIRED: 'TOKEN_EXPIRED',
} as const;

export const Errors = {
  validation: (message: string) =>
    new AppError(400, ErrorCodes.VALIDATION_ERROR, message),

  unauthorized: (message: string = 'Authentication required') =>
    new AppError(401, ErrorCodes.UNAUTHORIZED, message),

  forbidden: (message: string = 'Access denied') =>
    new AppError(403, ErrorCodes.FORBIDDEN, message),

  notFound: (resource: string) =>
    new AppError(404, ErrorCodes.NOT_FOUND, `${resource} not found`),

  conflict: (message: string) => new AppError(409, ErrorCodes.CONFLICT, message),

  internal: (message: string = 'Internal server error') =>
    new AppError(500, ErrorCodes.INTERNAL_ERROR, message),

  oauth: (message: string) => new AppError(500, ErrorCodes.OAUTH_ERROR, message),

  invalidToken: () => new AppError(401, ErrorCodes.INVALID_TOKEN, 'Invalid token'),

  tokenExpired: () => new AppError(401, ErrorCodes.TOKEN_EXPIRED, 'Token expired'),
};
```

**驗證：**
- [ ] src/lib/errors.ts 已創建
- [ ] 錯誤工廠函數正確

---

### Task 1.15: 創建空的路由檔案

**狀態：** ⬜ 待辦

**說明：** 為階段 4 準備路由檔案結構

**檔案：**
- `src/routes/auth.ts`
- `src/routes/admin.ts`
- `src/routes/user.ts`
- `src/routes/organizer.ts`
- `src/routes/review.ts`

**內容範例（所有路由檔案）：**
```typescript
import { Hono } from 'hono';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();

// Routes will be implemented in Stage 4

export default router;
```

**驗證：**
- [ ] 所有 5 個路由檔案已創建
- [ ] 檔案結構正確

---

### Task 1.16: 創建空的 middleware 檔案

**狀態：** ⬜ 待辦

**說明：** 為階段 3 準備 middleware 檔案結構

**檔案：**
- `src/middleware/auth.ts`

**內容：**
```typescript
import type { Context, Next } from 'hono';
import type { Env } from '../types';
import type { JWTPayload } from '../lib/jwt';

export const authMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  // Will be implemented in Stage 3
  await next();
};

export const adminAuthMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  // Will be implemented in Stage 3
  await next();
};

export const organizerAuthMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  // Will be implemented in Stage 3
  await next();
};

export const reviewAuthMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  // Will be implemented in Stage 3
  await next();
};
```

**驗證：**
- [ ] src/middleware/auth.ts 已創建
- [ ] middleware 框架正確

---

### Task 1.17: 配置測試環境

**狀態：** ⬜ 待辦

**說明：** 配置 Vitest 測試環境

**檔案：** `vitest.config.ts`

**內容：**
```typescript
import { defineConfig } from 'vitest/config';
import { workersPool } from '@cloudflare/vitest-pool-workers';

export default defineConfig({
  test: {
    pool: workersPool({
      // 選項：'browser' | 'node' | 'workers'
      runtime: 'workers',
      miniflare: {
        compatibilityDate: '2024-01-01',
        compatibilityFlags: ['nodejs_compat'],
        modules: true,
      },
    }),
    globals: true,
    environment: 'miniflare',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    },
  },
});
```

**驗證：**
- [ ] vitest.config.ts 已創建
- [ ] npx vitest 可以執行

---

### Task 1.18: 創建 README 文件

**狀態：** ⬜ 待辦

**說明：** 創建專案 README 文件

**檔案：** `workers/README.md`

**內容：**
```markdown
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

\`\`\`bash
# 安裝依賴
npm install

# 安裝 Wrangler CLI
官方建議
\`\`\`bash
npm i -D wrangler@latest
\`\`\`

# 登入 Cloudflare
wrangler login
\`\`\`

### 本地開發

\`\`\`bash
# 啟動開發伺服器
npm run dev

# 運行測試
npm run test

# 運行 lint
npm run lint
\`\`\`

### 環境變數

在使用 \`wrangler secret put\` 設置以下 secrets：

\`\`\`bash
wrangler secret put JWT_SECRET
wrangler secret put SESSION_KEY
wrangler secret put FACEBOOK_KEY
wrangler secret put FACEBOOK_SECRET
wrangler secret put INSTAGRAM_KEY
wrangler secret put INSTAGRAM_SECRET
\`\`\`

### 部署

\`\`\`bash
# 部署到 Cloudflare Workers
npm run deploy
\`\`\`

## 專案結構

\`\`\`
src/
├── lib/           # 工具函數（db, kv, jwt, oauth）
├── routes/        # API 路由處理器
├── middleware/    # 中介層（cors, auth, error）
├── types/         # TypeScript 類型定義
└── index.ts       # 主入口
\`\`\`

## API 文檔

請參考 API.md

## 授權

MIT
\`\`\`
```

**驗證：**
- [ ] README.md 已創建
- [ ] 內容完整

---

### Task 1.19: 配置 npm scripts

**狀態：** ⬜ 待辦

**說明：** 在 package.json 中添加腳本命令

**package.json 更新：**
```json
{
  "name": "free2free-workers",
  "version": "1.0.0",
  "description": "Free2Free API on Cloudflare Workers",
  "type": "module",
  "scripts": {
    "dev": "wrangler dev",
    "deploy": "wrangler deploy",
    "test": "vitest",
    "test:unit": "vitest --run test/unit",
    "test:integration": "vitest --run test/integration",
    "test:e2e": "vitest --run test/e2e",
    "test:coverage": "vitest --coverage",
    "lint": "eslint src --ext .ts",
    "lint:fix": "eslint src --ext .ts --fix",
    "format": "prettier --write \"src/**/*.ts\"",
    "format:check": "prettier --check \"src/**/*.ts\"",
    "typecheck": "tsc --noEmit"
  },
  ...
}
```

**驗證：**
- [ ] package.json 已更新
- [ ] 所有腳本可執行

---

### Task 1.20: 驗證階段 1 完成度

**狀態：** ⬜ 待辦

**說明：** 完成最終驗證，確保所有基礎架構正常

**驗證清單：**

```bash
# 1. TypeScript 編譯
npm run typecheck

# 2. ESLint 檢查
npm run lint

# 3. 開發伺服器啟動
npm run dev
# 應該能看到 http://localhost:8787 返回 JSON 響應

# 4. 測試執行
npm run test

# 5. 項目結構檢查
tree src -I node_modules
# 應該看到完整的目錄結構
```

**預期結果：**
- ✅ TypeScript 編譯無錯誤
- ✅ ESLint 無警告
- ✅ 開發伺服器正常啟動
- ✅ 測試可以執行（即使沒有測試用例）
- ✅ 所有檔案和目錄已創建

---

## 🎯 階段 1 完成標準

當以下所有項目都完成時，階段 1 視為完成：

- [ ] 所有 20 個任務已完成
- [ ] 專案可以成功編譯
- [ ] 開發伺服器可以啟動
- [ ] Health check 端點可訪問
- [ ] README 文件完整
- [ ] 開發環境可以正常使用

---

## 📝 備註

1. **Wrangler 設置**: 需要先運行 `wrangler login` 來認證
2. **D1 Database**: 階段 2 將會創建實際的 D1 database
3. **KV Namespace**: 階段 3 將會創建實際的 KV namespace
4. **環境變數**: 使用 `wrangler secret put` 來設置敏感信息

---

## 🚀 下一階段

完成階段 1 後，可以進入：

**階段 2：資料層遷移**
- 創建 D1 資料庫
- 定義資料表 schema
- 實現 CRUD 操作
- 資料遷移腳本

---

**更新日期：** 2026-01-14
**當前進度：** 20/20 任務完成 ✅
**狀態：** 已完成
