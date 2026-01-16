import type { D1Database } from '@cloudflare/workers-types';
import type { User } from '../types';

export interface Session {
  id: string;
  user_id: number;
  data: Record<string, unknown>;
  expires_at: string;
  created_at: string;
}

export class SessionManager {
  constructor(private db: D1Database) {}

  private generateSessionId(): string {
    return crypto.randomUUID();
  }

  private getExpiryTime(minutes: number): string {
    const expiry = new Date(Date.now() + minutes * 60 * 1000);
    return expiry.toISOString();
  }

  async createSession(
    userId: number,
    data: Record<string, unknown> = {},
    expiryMinutes: number = 1440
  ): Promise<Session> {
    const sessionId = this.generateSessionId();
    const expiresAt = this.getExpiryTime(expiryMinutes);

    await this.db
      .prepare(
        `INSERT INTO sessions (id, user_id, data, expires_at, created_at)
         VALUES (?, ?, ?, ?, datetime('now'))`
      )
      .bind(sessionId, userId, JSON.stringify(data), expiresAt)
      .run();

    return this.getSession(sessionId) as Promise<Session>;
  }

  async getSession(sessionId: string): Promise<Session | null> {
    const row = await this.db
      .prepare('SELECT * FROM sessions WHERE id = ? AND expires_at > datetime("now")')
      .bind(sessionId)
      .first<Session>();

    if (!row) return null;

    return {
      ...row,
      data: JSON.parse(row.data as unknown as string),
    };
  }

  async updateSession(sessionId: string, data: Record<string, unknown>): Promise<Session | null> {
    await this.db
      .prepare(
        'UPDATE sessions SET data = ?, expires_at = datetime("now", "+1440 minutes") WHERE id = ?'
      )
      .bind(JSON.stringify(data), sessionId)
      .run();

    return this.getSession(sessionId);
  }

  async deleteSession(sessionId: string): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM sessions WHERE id = ?').bind(sessionId).run();
    return (result.meta.changes || 0) > 0;
  }

  async deleteSessionsByUserId(userId: number): Promise<number> {
    const result = await this.db
      .prepare('DELETE FROM sessions WHERE user_id = ?')
      .bind(userId)
      .run();
    return result.meta.changes || 0;
  }

  async cleanupExpiredSessions(): Promise<number> {
    const result = await this.db
      .prepare('DELETE FROM sessions WHERE expires_at <= datetime("now")')
      .run();
    return result.meta.changes || 0;
  }

  async getUserFromSession(sessionId: string): Promise<User | null> {
    const session = await this.getSession(sessionId);
    if (!session) return null;

    const user = await this.db
      .prepare('SELECT * FROM users WHERE id = ?')
      .bind(session.user_id)
      .first<User>();

    if (!user) return null;

    return { ...user, is_admin: (user.is_admin as unknown as number) === 1 };
  }

  async refreshSession(sessionId: string): Promise<Session | null> {
    const expiresAt = this.getExpiryTime(1440);

    await this.db
      .prepare('UPDATE sessions SET expires_at = ? WHERE id = ?')
      .bind(expiresAt, sessionId)
      .run();

    return this.getSession(sessionId);
  }
}

export function createSessionManager(db: D1Database): SessionManager {
  return new SessionManager(db);
}
