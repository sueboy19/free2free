# Free2Free API 部署指南

## 前置要求

- Cloudflare 账户
- Node.js 18+
- Wrangler CLI

## 部署步骤

### 1. 安装 Wrangler CLI

```bash
npm install -g wrangler
```

### 2. 登入 Cloudflare

```bash
wrangler login
```

### 3. 配置环境变量

```bash
# 配置 JWT_SECRET（至少 32 字符随机字符串）
wrangler secret put JWT_SECRET

# 配置 SESSION_KEY（至少 32 字符随机字符串）
wrangler secret put SESSION_KEY

# 配置 Facebook OAuth
wrangler secret put FACEBOOK_KEY
wrangler secret put FACEBOOK_SECRET

# 配置 Instagram OAuth
wrangler secret put INSTAGRAM_KEY
wrangler secret put INSTAGRAM_SECRET
```

### 4. 创建 D1 数据库

```bash
# 创建 D1 数据库
wrangler d1 create free2free-db

# 记录 database_id 并更新 wrangler.toml
```

### 5. 执行数据库 Migration

```bash
# 执行初始 migration
wrangler d1 execute free2free-db --file=./migrations/0001_initial.sql

# 查看数据表
wrangler d1 execute free2free-db --command="SELECT name FROM sqlite_master WHERE type='table';"
```

### 6. 导入测试数据（可选）

```bash
# 导入测试数据到 D1
wrangler d1 execute free2free-db --file=./scripts/import-to-d1.sql
```

### 7. 部署应用

```bash
# 部署到 Cloudflare Workers
wrangler deploy
```

### 8. 验证部署

```bash
# 查看 worker 信息
wrangler deployments list

# 查看 worker 日志
wrangler tail
```

## 环境变量

| 变量名称 | 说明 | 示例 |
|-----------|--------|--------|
| JWT_SECRET | JWT 加密密钥 | 至少 32 字符随机字符串 |
| SESSION_KEY | Session 加密密钥 | 至少 32 字符随机字符串 |
| FACEBOOK_KEY | Facebook App ID | 从 Facebook Developer 获取 |
| FACEBOOK_SECRET | Facebook App Secret | 从 Facebook Developer 获取 |
| INSTAGRAM_KEY | Instagram App ID | 从 Instagram Developer 获取 |
| INSTAGRAM_SECRET | Instagram App Secret | 从 Instagram Developer 获取 |
| ENVIRONMENT | 环境标识 | production |
| CORS_ORIGINS | CORS 允许的来源 | https://free2free.example.com |

## 故障排除

### 部署失败

1. 检查 wrangler.toml 配置
2. 确认 Cloudflare 账户已登录
3. 检查数据库 ID 是否正确
4. 查看错误日志：`wrangler tail`

### 运行时错误

1. 使用 `wrangler tail` 查看实时日志
2. 检查环境变量是否正确设置
3. 检查数据库连接

## 回滚部署

```bash
# 查看部署历史
wrangler deployments list

# 回滚到上一个版本
wrangler rollback
```

## 监控

### 使用 Cloudflare Analytics

1. 登录 Cloudflare Dashboard
2. 导航到 Workers & Pages
3. 选择 free2free-api
4. 查看分析数据

### 使用日志

```bash
# 实时查看日志
wrangler tail
```

## 性能优化

### 数据库查询优化

1. 使用索引加速查询
2. 避免 N+1 查询
3. 限制查询结果数量

### 缓存策略

1. 考虑使用 Cloudflare KV 缓存常用数据
2. 实现查询结果缓存
3. 使用 CDN 缓存静态资源
