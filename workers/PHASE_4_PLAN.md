# Cloudflare Workers 遷移計畫 - 階段 4：API 路由實現

## 📋 階段 4 概述

**目標：** 實現完整的 API 路由，包括 Admin、User、Organizer 和 Review 功能

**預計時間：** 1.5 週

**狀態：** ✅ 已完成

**完成日期：** 2026-01-16

---

## ✅ 任務清單

### Task 4.1: 實現 Admin 路由 - Locations 管理

**狀態：** ⬜ 待辦

**說明：** 實現管理員管理地點的路由

**檔案：** `src/routes/admin.ts`

**實現內容：**
```typescript
import { Hono } from 'hono';
import { adminAuthMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();
router.use('/*', adminAuthMiddleware);

// Create location
router.post('/admin/locations', async (c) => {
  const body = await c.req.json();
  const { name, address, latitude, longitude } = body;

  if (!name || !address || latitude === undefined || longitude === undefined) {
    throw new Error('Missing required fields');
  }

  const result = await c.env.DB.prepare(
    `INSERT INTO locations (name, address, latitude, longitude)
     VALUES (?, ?, ?, ?)`
  )
    .bind(name, address, latitude, longitude)
    .run();

  const location = await c.env.DB.prepare('SELECT * FROM locations WHERE id = ?')
    .bind(result.meta.last_row_id)
    .first();

  return c.json({ data: location });
});

// List locations
router.get('/admin/locations', async (c) => {
  const result = await c.env.DB.prepare('SELECT * FROM locations ORDER BY id DESC').all();
  return c.json({ data: result.results || [], total: result.results?.length || 0 });
});

// Get location by id
router.get('/admin/locations/:id', async (c) => {
  const id = c.req.param('id');
  const location = await c.env.DB.prepare('SELECT * FROM locations WHERE id = ?')
    .bind(id)
    .first();

  if (!location) {
    throw new Error('Location not found');
  }

  return c.json({ data: location });
});

// Update location
router.put('/admin/locations/:id', async (c) => {
  const id = c.req.param('id');
  const body = await c.req.json();

  const updates: string[] = [];
  const values: any[] = [];

  for (const [key, value] of Object.entries(body)) {
    updates.push(`${key} = ?`);
    values.push(value);
  }

  if (updates.length > 0) {
    await c.env.DB.prepare(`UPDATE locations SET ${updates.join(', ')} WHERE id = ?`)
      .bind(...values, id)
      .run();
  }

  const location = await c.env.DB.prepare('SELECT * FROM locations WHERE id = ?')
    .bind(id)
    .first();

  return c.json({ data: location });
});

// Delete location
router.delete('/admin/locations/:id', async (c) => {
  const id = c.req.param('id');
  const result = await c.env.DB.prepare('DELETE FROM locations WHERE id = ?')
    .bind(id)
    .run();

  return c.json({ success: (result.meta.changes || 0) > 0 });
});

export default router;
```

**驗證：**
- [ ] 所有 location CRUD 路由已實現
- [ ] Admin 認證正常工作

---

### Task 4.2: 實現 Admin 路由 - Activities 管理

**狀態：** ⬜ 待辦

**說明：** 實現管理員管理活動的路由

**實現內容：**
```typescript
// Create activity
router.post('/admin/activities', async (c) => {
  const body = await c.req.json();
  const { title, target_count, location_id, description } = body;

  if (!title || !target_count || !location_id) {
    throw new Error('Missing required fields');
  }

  const user = c.get('user' as never);

  const result = await c.env.DB.prepare(
    `INSERT INTO activities (title, target_count, location_id, description, created_by)
     VALUES (?, ?, ?, ?, ?)`
  )
    .bind(title, target_count, location_id, description || null, (user as any).id)
    .run();

  const activity = await c.env.DB.prepare('SELECT * FROM activities WHERE id = ?')
    .bind(result.meta.last_row_id)
    .first();

  return c.json({ data: activity });
});

// List activities
router.get('/admin/activities', async (c) => {
  const result = await c.env.DB.prepare('SELECT * FROM activities ORDER BY id DESC').all();
  return c.json({ data: result.results || [], total: result.results?.length || 0 });
});

// Get activity by id
router.get('/admin/activities/:id', async (c) => {
  const id = c.req.param('id');
  const activity = await c.env.DB.prepare('SELECT * FROM activities WHERE id = ?')
    .bind(id)
    .first();

  if (!activity) {
    throw new Error('Activity not found');
  }

  return c.json({ data: activity });
});

// Update activity
router.put('/admin/activities/:id', async (c) => {
  const id = c.req.param('id');
  const body = await c.req.json();

  const updates: string[] = [];
  const values: any[] = [];

  for (const [key, value] of Object.entries(body)) {
    updates.push(`${key} = ?`);
    values.push(value);
  }

  if (updates.length > 0) {
    await c.env.DB.prepare(`UPDATE activities SET ${updates.join(', ')} WHERE id = ?`)
      .bind(...values, id)
      .run();
  }

  const activity = await c.env.DB.prepare('SELECT * FROM activities WHERE id = ?')
    .bind(id)
    .first();

  return c.json({ data: activity });
});

// Delete activity
router.delete('/admin/activities/:id', async (c) => {
  const id = c.req.param('id');
  const result = await c.env.DB.prepare('DELETE FROM activities WHERE id = ?')
    .bind(id)
    .run();

  return c.json({ success: (result.meta.changes || 0) > 0 });
});
```

**驗證：**
- [ ] 所有 activity CRUD 路由已實現

---

### Task 4.3: 實現 User 路由 - Matches 查詢

**狀態：** ⬜ 待辦

**說明：** 實現用戶查看配對局的路由

**檔案：** `src/routes/user.ts`

**實現內容：**
```typescript
import { Hono } from 'hono';
import { authMiddleware, optionalAuthMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();

// Get open matches (public)
router.get('/matches', async (c) => {
  const result = await c.env.DB.prepare(
    `SELECT * FROM matches WHERE status = ? AND match_time > datetime('now') ORDER BY match_time ASC`
  )
    .bind('open')
    .all();

  return c.json({ data: result.results || [] });
});

// Get user's matches (authenticated)
router.get('/user/matches', authMiddleware, async (c) => {
  const user = c.get('user' as never);
  const status = c.req.query('status') || 'completed';

  const result = await c.env.DB.prepare(
    `SELECT DISTINCT m.* FROM matches m
     JOIN match_participants mp ON m.id = mp.match_id
     WHERE mp.user_id = ? AND m.status = ?
     ORDER BY m.match_time DESC`
  )
    .bind((user as any).id, status)
    .all();

  return c.json({ data: result.results || [] });
});

// Get match details
router.get('/matches/:id', optionalAuthMiddleware, async (c) => {
  const id = c.req.param('id');
  const match = await c.env.DB.prepare('SELECT * FROM matches WHERE id = ?')
    .bind(id)
    .first();

  if (!match) {
    throw new Error('Match not found');
  }

  return c.json({ data: match });
});

// Get match participants
router.get('/matches/:id/participants', async (c) => {
  const id = c.req.param('id');
  const result = await c.env.DB.prepare(
    'SELECT * FROM match_participants WHERE match_id = ? ORDER BY id DESC'
  )
    .bind(id)
    .all();

  return c.json({ data: result.results || [] });
});

// Get match reviews
router.get('/matches/:id/reviews', async (c) => {
  const id = c.req.param('id');
  const result = await c.env.DB.prepare(
    'SELECT * FROM reviews WHERE match_id = ? ORDER BY created_at DESC'
  )
    .bind(id)
    .all();

  return c.json({ data: result.results || [] });
});

export default router;
```

**驗證：**
- [ ] 所有用戶查詢路由已實現

---

### Task 4.4: 實現 Organizer 路由 - Matches 管理

**狀態：** ⬜ 待辦

**說明：** 實現開局者管理配對局的路由

**檔案：** `src/routes/organizer.ts`

**實現內容：**
```typescript
import { Hono } from 'hono';
import { authMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();
router.use('/*', authMiddleware);

// Create match
router.post('/matches', async (c) => {
  const user = c.get('user' as never);
  const body = await c.req.json();
  const { activity_id, match_time } = body;

  if (!activity_id || !match_time) {
    throw new Error('Missing required fields');
  }

  const result = await c.env.DB.prepare(
    `INSERT INTO matches (activity_id, organizer_id, match_time, status)
     VALUES (?, ?, ?, 'open')`
  )
    .bind(activity_id, (user as any).id, match_time)
    .run();

  const match = await c.env.DB.prepare('SELECT * FROM matches WHERE id = ?')
    .bind(result.meta.last_row_id)
    .first();

  return c.json({ data: match });
});

// Update match status
router.put('/matches/:id/status', async (c) => {
  const id = c.req.param('id');
  const body = await c.req.json();
  const { status } = body;

  if (!['open', 'completed', 'cancelled'].includes(status)) {
    throw new Error('Invalid status');
  }

  await c.env.DB.prepare('UPDATE matches SET status = ? WHERE id = ?')
    .bind(status, id)
    .run();

  const match = await c.env.DB.prepare('SELECT * FROM matches WHERE id = ?')
    .bind(id)
    .first();

  return c.json({ data: match });
});

// Join a match
router.post('/matches/:id/join', async (c) => {
  const user = c.get('user' as never);
  const matchId = c.req.param('id');

  const result = await c.env.DB.prepare(
    `INSERT INTO match_participants (match_id, user_id, status, joined_at)
     VALUES (?, ?, 'pending', datetime('now'))`
  )
    .bind(matchId, (user as any).id)
    .run();

  const participant = await c.env.DB.prepare('SELECT * FROM match_participants WHERE id = ?')
    .bind(result.meta.last_row_id)
    .first();

  return c.json({ data: participant });
});

// Approve participant
router.put('/matches/:matchId/participants/:participantId', async (c) => {
  const matchId = c.req.param('matchId');
  const participantId = c.req.param('participantId');
  const body = await c.req.json();
  const { status } = body;

  if (!['approved', 'rejected'].includes(status)) {
    throw new Error('Invalid status');
  }

  await c.env.DB.prepare('UPDATE match_participants SET status = ? WHERE id = ?')
    .bind(status, participantId)
    .run();

  const participant = await c.env.DB.prepare('SELECT * FROM match_participants WHERE id = ?')
    .bind(participantId)
    .first();

  return c.json({ data: participant });
});

// Leave a match
router.delete('/matches/:id/join', async (c) => {
  const user = c.get('user' as never);
  const matchId = c.req.param('id');

  const result = await c.env.DB.prepare(
    'DELETE FROM match_participants WHERE match_id = ? AND user_id = ?'
  )
    .bind(matchId, (user as any).id)
    .run();

  return c.json({ success: (result.meta.changes || 0) > 0 });
});

export default router;
```

**驗證：**
- [ ] 所有開局者路由已實現

---

### Task 4.5: 實現 Review 路由

**狀態：** ⬜ 待辦

**說明：** 實現評分和點讚的路由

**檔案：** `src/routes/review.ts`

**實現內容：**
```typescript
import { Hono } from 'hono';
import { authMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();
router.use('/*', authMiddleware);

// Create review
router.post('/reviews', async (c) => {
  const user = c.get('user' as never);
  const body = await c.req.json();
  const { match_id, reviewee_id, score, comment } = body;

  if (!match_id || !reviewee_id || score === undefined) {
    throw new Error('Missing required fields');
  }

  if (score < 1 || score > 5) {
    throw new Error('Score must be between 1 and 5');
  }

  const result = await c.env.DB.prepare(
    `INSERT INTO reviews (match_id, reviewer_id, reviewee_id, score, comment, created_at)
     VALUES (?, ?, ?, ?, ?, datetime('now'))`
  )
    .bind(match_id, (user as any).id, reviewee_id, score, comment || null)
    .run();

  const review = await c.env.DB.prepare('SELECT * FROM reviews WHERE id = ?')
    .bind(result.meta.last_row_id)
    .first();

  return c.json({ data: review });
});

// Get reviews by match
router.get('/reviews', async (c) => {
  const matchId = c.req.query('match_id');
  const reviewerId = c.req.query('reviewer_id');

  let query = 'SELECT * FROM reviews WHERE 1=1';
  const params: any[] = [];

  if (matchId) {
    query += ' AND match_id = ?';
    params.push(matchId);
  }

  if (reviewerId) {
    query += ' AND reviewer_id = ?';
    params.push(reviewerId);
  }

  query += ' ORDER BY created_at DESC';

  const stmt = c.env.DB.prepare(query);
  for (const param of params) {
    stmt.bind(param);
  }

  const result = await stmt.all();

  return c.json({ data: result.results || [] });
});

// Like/Unlike review
router.post('/reviews/:id/like', async (c) => {
  const user = c.get('user' as never);
  const reviewId = c.req.param('id');
  const body = await c.req.json();
  const { is_like } = body;

  if (typeof is_like !== 'boolean') {
    throw new Error('is_like must be a boolean');
  }

  const existing = await c.env.DB.prepare(
    'SELECT * FROM review_likes WHERE review_id = ? AND user_id = ?'
  )
    .bind(reviewId, (user as any).id)
    .first();

  if (existing) {
    await c.env.DB.prepare(
      'UPDATE review_likes SET is_like = ? WHERE review_id = ? AND user_id = ?'
    )
      .bind(is_like ? 1 : 0, reviewId, (user as any).id)
      .run();
  } else {
    await c.env.DB.prepare(
      'INSERT INTO review_likes (review_id, user_id, is_like) VALUES (?, ?, ?)'
    )
      .bind(reviewId, (user as any).id, is_like ? 1 : 0)
      .run();
  }

  return c.json({ success: true });
});

// Delete review
router.delete('/reviews/:id', async (c) => {
  const user = c.get('user' as never);
  const reviewId = c.req.param('id');

  const result = await c.env.DB.prepare(
    'DELETE FROM reviews WHERE id = ? AND reviewer_id = ?'
  )
    .bind(reviewId, (user as any).id)
    .run();

  return c.json({ success: (result.meta.changes || 0) > 0 });
});

export default router;
```

**驗證：**
- [ ] 所有評分路由已實現

---

### Task 4.6: 更新主入口註冊所有路由

**狀態：** ⬜ 待辦

**說明：** 在 src/index.ts 中註冊所有路由

**實現內容：**
```typescript
import adminRoutes from './routes/admin';
import userRoutes from './routes/user';
import organizerRoutes from './routes/organizer';
import reviewRoutes from './routes/review';

// ... existing code

// Admin routes (requires admin auth)
app.route('/', adminRoutes);

// User routes
app.route('/', userRoutes);

// Organizer routes
app.route('/', organizerRoutes);

// Review routes
app.route('/', reviewRoutes);
```

**驗證：**
- [ ] 所有路由已註冊

---

### Task 4.7: 寫入 API 測試

**狀態：** ⬜ 待辦

**檔案：** `test/integration/api.test.ts`

**實現內容：**
```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { Hono } from 'hono';
import type { Env } from '../../src/types';

describe('API Integration', () => {
  let app: Hono<{ Bindings: Env }>;
  let env: Env;

  beforeEach(async () => {
    env = {
      DB: {
        prepare: (query: string) => ({
          bind: (...args: any[]) => ({
            run: async () => ({ meta: { last_row_id: 1, changes: 1 } }),
            first: async () => ({
              id: 1,
              name: 'Test Location',
              address: '123 Test St',
              latitude: 25.0479,
              longitude: 121.5170,
            }),
          }),
        }),
      } as any,
      KV: {} as any,
      JWT_SECRET: 'test-secret-key-at-least-32-characters-long',
      SESSION_KEY: 'test-session-key-at-least-32-characters-long',
      FACEBOOK_KEY: 'test',
      FACEBOOK_SECRET: 'test',
      INSTAGRAM_KEY: 'test',
      INSTAGRAM_SECRET: 'test',
      BASE_URL: 'http://localhost',
      FRONTEND_URL: 'http://localhost:3000',
      CORS_ORIGINS: 'http://localhost:3000',
    };

    app = new Hono<{ Bindings: Env }>();
  });

  it('should create a location', async () => {
    const res = await app.request('/admin/locations', {
      method: 'POST',
      env,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: 'Test Location',
        address: '123 Test St',
        latitude: 25.0479,
        longitude: 121.5170,
      }),
    });

    expect(res.status).toBe(200);
  });

  // More tests...
});
```

**驗證：**
- [ ] 測試已寫入

---

### Task 4.8: 寫入路由單元測試

**狀態：** ⬜ 待辦

**檔案：**
- `test/unit/admin.test.ts`
- `test/unit/user.test.ts`
- `test/unit/organizer.test.ts`
- `test/unit/review.test.ts`

**驗證：**
- [ ] 所有單元測試已寫入

---

### Task 4.9: 更新 README 文檔

**狀態：** ⬜ 待辦

**說明：** 在 README.md 中添加 API 端點說明

**實現內容：**
```markdown
## API 端點

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
```

**驗證：**
- [ ] 文檔已更新

---

### Task 4.10: 最終驗證

**狀態：** ⬜ 待辦

**說明：** 完成所有任務的最終驗證

**驗證清單：**

```bash
# 1. TypeScript 編譯
npm run typecheck

# 2. 執行測試
npm run test

# 3. Lint 檢查
npm run lint

# 4. 本地開發服務器
npm run dev
```

**預期結果：**
- ✅ TypeScript 編譯無錯誤
- ✅ 所有測試通過
- ✅ Lint 無警告
- ✅ 所有 API 端點可訪問

---

## 🎯 階段 4 完成標準

當以下所有項目都完成時，階段 4 視為完成：

- [ ] 所有 10 個任務已完成
- [ ] Admin 路由完整實現
- [ ] User 路由完整實現
- [ ] Organizer 路由完整實現
- [ ] Review 路由完整實現
- [ ] 所有測試通過
- [ ] 測試覆蓋率 > 80%

---

## 📝 備註

1. **權限控制**：
   - Admin 路由需要管理員權限
   - User 路由需要登入（部分可選）
   - Organizer 路由需要登入
   - Review 路由需要登入

2. **輸入驗證**：
   - 必需字段檢查
   - 數據格式驗證
   - 範圍檢查（如評分 1-5）

3. **錯誤處理**：
   - 統一的錯誤響應格式
   - 清晰的錯誤消息
   - 適當的 HTTP 狀態碼

---

## 🚀 下一階段

完成階段 4 後，可以進入：

**階段 5：測試與部署**
- 完整的端到端測試
- 性能測試
- 部署到 Cloudflare Workers
- 生產環境驗證

---

**更新日期：** 2026-01-16
**當前進度：** 10/10 任務完成 ✅
**狀態：** 已完成
