import { Hono } from 'hono';
import { authMiddleware, optionalAuthMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();

router.get('/matches', optionalAuthMiddleware, async (c) => {
  const result = await c.env.DB.prepare(
    `SELECT * FROM matches WHERE status = ? AND match_time > datetime('now') ORDER BY match_time ASC`
  )
    .bind('open')
    .all();

  return c.json({ data: result.results || [] });
});

router.get('/user/matches', authMiddleware, async (c) => {
  const user = c.get('user' as never);

  if (!user) {
    throw new Error('User not found');
  }

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

router.get('/matches/:id', optionalAuthMiddleware, async (c) => {
  const id = c.req.param('id');
  const match = await c.env.DB.prepare('SELECT * FROM matches WHERE id = ?').bind(id).first();

  if (!match) {
    throw new Error('Match not found');
  }

  return c.json({ data: match });
});

router.get('/matches/:id/participants', async (c) => {
  const id = c.req.param('id');
  const result = await c.env.DB.prepare(
    'SELECT * FROM match_participants WHERE match_id = ? ORDER BY id DESC'
  )
    .bind(id)
    .all();

  return c.json({ data: result.results || [] });
});

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
