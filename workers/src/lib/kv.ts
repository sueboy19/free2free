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
    await this.delete(`user:${userId}:tokens`);
  }
}
