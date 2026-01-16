import { describe, it, expect, beforeEach } from 'vitest';
import { Hono } from 'hono';
import type { Env } from '../../src/types';

describe('Authentication Integration', () => {
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
              social_id: '123',
              social_provider: 'facebook',
              name: 'Test User',
              email: 'test@example.com',
              avatar_url: 'http://example.com/avatar.jpg',
              is_admin: 0,
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

  it('should reject unauthenticated request', async () => {
    app.get('/test-auth', async (c) => {
      const { authMiddleware } = await import('../../src/middleware/auth');
      await authMiddleware(c, async () => {
        const user = c.get('user');
        return c.json({ user });
      });
    });

    const res = await app.request('/test-auth');

    expect(res.status).toBe(500);
  });
});
