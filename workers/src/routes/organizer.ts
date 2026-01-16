import { Hono } from 'hono';
import { authMiddleware, organizerAuthMiddleware } from '../middleware/auth';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();

router.post('/matches', authMiddleware, async (c) => {
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

router.put('/matches/:id/status', authMiddleware, async (c) => {
  const id = c.req.param('id');
  const body = await c.req.json();
  const { status } = body;

  if (!['open', 'completed', 'cancelled'].includes(status)) {
    throw new Error('Invalid status');
  }

  await c.env.DB.prepare('UPDATE matches SET status = ? WHERE id = ?').bind(status, id).run();

  const match = await c.env.DB.prepare('SELECT * FROM matches WHERE id = ?').bind(id).first();

  return c.json({ data: match });
});

router.post('/matches/:id/join', authMiddleware, async (c) => {
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

router.put('/matches/:matchId/participants/:participantId', organizerAuthMiddleware, async (c) => {
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

router.delete('/matches/:id/join', authMiddleware, async (c) => {
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
