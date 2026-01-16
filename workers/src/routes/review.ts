import { Hono } from 'hono';
import { authMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();

router.post('/reviews', authMiddleware, async (c) => {
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

  let stmt = c.env.DB.prepare(query);
  for (const param of params) {
    stmt = stmt.bind(param);
  }

  const result = await stmt.all();

  return c.json({ data: result.results || [] });
});

router.post('/reviews/:id/like', authMiddleware, async (c) => {
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

router.delete('/reviews/:id', authMiddleware, async (c) => {
  const user = c.get('user' as never);
  const reviewId = c.req.param('id');

  const result = await c.env.DB.prepare('DELETE FROM reviews WHERE id = ? AND reviewer_id = ?')
    .bind(reviewId, (user as any).id)
    .run();

  return c.json({ success: (result.meta.changes || 0) > 0 });
});

export default router;
