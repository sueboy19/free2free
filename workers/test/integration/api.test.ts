import { describe, it, expect, beforeEach } from 'vitest';
import { Hono } from 'hono';
import type { Env } from '../../src/types';

describe('API Integration', () => {
  let app: Hono<{ Bindings: Env }>;
  let env: Env;
  let authToken: string;

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

  it('should return health check', async () => {
    app.get('/', (c) => c.json({ status: 'ok', message: 'API is running' }));
    const res = await app.request('/', { env });
    expect(res.status).toBe(200);
  });

  it('should create a location', async () => {
    app.post('/admin/locations', async (c) => {
      return c.json({ data: { id: 1, name: 'Test Location' } });
    });

    const res = await app.request('/admin/locations', {
      method: 'POST',
      env,
      headers: { 'Content-Type': 'application/json', Authorization: 'Bearer admin-token' },
      body: JSON.stringify({
        name: 'Test Location',
        address: '123 Test St',
        latitude: 25.0479,
        longitude: 121.5170,
      }),
    });

    expect(res.status).toBe(200);
  });

  it('should list open matches', async () => {
    app.get('/matches', async (c) => {
      return c.json({ data: [] });
    });

    const res = await app.request('/matches', { env });
    expect(res.status).toBe(200);
  });
});
