import { describe, it, expect, beforeEach } from 'vitest';
import { SessionManager } from '../../src/lib/session';

describe('SessionManager', () => {
  let db: any;
  let sessionManager: SessionManager;

  beforeEach(() => {
    db = {
      prepare: () => ({
        bind: () => ({
          run: async () => ({ meta: { changes: 1 } }),
          first: async () => ({
            id: 'session-id',
            user_id: 1,
            data: JSON.stringify({ test: 'data' }),
            expires_at: new Date(Date.now() + 1000000).toISOString(),
            created_at: new Date().toISOString(),
          }),
        }),
      }),
    };
    sessionManager = new SessionManager(db);
  });

  it('should create a session', async () => {
    const session = await sessionManager.createSession(1, { test: 'data' });

    expect(session.id).toBeTruthy();
    expect(session.user_id).toBe(1);
    expect(session.data).toEqual({ test: 'data' });
  });

  it('should get a session', async () => {
    const session = await sessionManager.getSession('session-id');

    expect(session).toBeTruthy();
    expect(session?.id).toBe('session-id');
  });

  it('should return null for expired session', async () => {
    db.prepare = () => ({
      bind: () => ({
        first: async () => null,
      }),
    });

    const session = await sessionManager.getSession('session-id');
    expect(session).toBeNull();
  });

  it('should delete a session', async () => {
    const result = await sessionManager.deleteSession('session-id');
    expect(result).toBe(true);
  });

  it('should delete sessions by user id', async () => {
    const count = await sessionManager.deleteSessionsByUserId(1);
    expect(typeof count).toBe('number');
  });
});
