import { describe, it, expect } from 'vitest';
import { JWTManager } from '../../src/lib/jwt';
import type { User } from '../../src/types';

describe('JWTManager', () => {
  const secret = 'test-secret-key-at-least-32-characters-long';
  const jwtManager = new JWTManager(secret);

  const mockUser: User = {
    id: 1,
    social_id: '123',
    social_provider: 'facebook',
    name: 'Test User',
    email: 'test@example.com',
    avatar_url: 'http://example.com/avatar.jpg',
    is_admin: false,
    created_at: Date.now(),
    updated_at: Date.now(),
  };

  it('should generate access token', async () => {
    const token = await jwtManager.generateAccessToken(mockUser);
    expect(token).toBeTruthy();
    expect(typeof token).toBe('string');
  });

  it('should verify access token', async () => {
    const token = await jwtManager.generateAccessToken(mockUser);
    const payload = await jwtManager.verifyAccessToken(token);

    expect(payload.user_id).toBe(mockUser.id);
    expect(payload.user_name).toBe(mockUser.name);
    expect(payload.is_admin).toBe(mockUser.is_admin);
    expect(payload.exp).toBeTruthy();
  });

  it('should generate refresh token', async () => {
    const token = await jwtManager.generateRefreshToken(mockUser);
    expect(token).toBeTruthy();
    expect(typeof token).toBe('string');
  });

  it('should generate both tokens', async () => {
    const tokens = await jwtManager.generateTokens(mockUser);

    expect(tokens.access).toBeTruthy();
    expect(tokens.refresh).toBeTruthy();
    expect(tokens.access).not.toBe(tokens.refresh);
  });

  it('should throw error for invalid token', async () => {
    await expect(jwtManager.verifyAccessToken('invalid-token')).rejects.toThrow();
  });

  it('should throw error for secret less than 32 characters', () => {
    expect(() => new JWTManager('short-secret')).toThrow();
  });
});
