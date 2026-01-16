import { describe, it, expect, vi } from 'vitest';
import { FacebookOAuthProvider, InstagramOAuthProvider, OAuthManager } from '../../src/lib/oauth';

global.fetch = vi.fn();

describe('OAuthProvider', () => {
  const redirectUri = 'http://localhost:8787/auth/facebook/callback';

  describe('FacebookOAuthProvider', () => {
    const provider = new FacebookOAuthProvider('test-id', 'test-secret', redirectUri);

    it('should generate auth URL', () => {
      const url = provider.getAuthUrl();

      expect(url).toContain('facebook.com');
      expect(url).toContain('client_id=test-id');
      expect(url).toContain('redirect_uri');
    });

    it('should exchange code for token', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        json: async () => ({ access_token: 'test-access-token' }),
      } as Response);

      const token = await provider.exchangeCodeForToken('test-code');

      expect(token).toBe('test-access-token');
    });

    it('should get user profile', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        json: async () => ({
          id: '123',
          name: 'Test User',
          email: 'test@example.com',
          picture: { data: { url: 'http://example.com/avatar.jpg' } },
        }),
      } as Response);

      const profile = await provider.getUserProfile('test-access-token');

      expect(profile.id).toBe('123');
      expect(profile.name).toBe('Test User');
      expect(profile.email).toBe('test@example.com');
    });

    it('should throw error on OAuth failure', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        json: async () => ({ error: { message: 'Invalid code' } }),
      } as Response);

      await expect(provider.exchangeCodeForToken('invalid-code')).rejects.toThrow();
    });
  });

  describe('InstagramOAuthProvider', () => {
    const provider = new InstagramOAuthProvider('test-id', 'test-secret', redirectUri);

    it('should generate auth URL', () => {
      const url = provider.getAuthUrl();

      expect(url).toContain('instagram.com');
      expect(url).toContain('client_id=test-id');
    });

    it('should exchange code for token', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        json: async () => ({ access_token: 'test-access-token' }),
      } as Response);

      const token = await provider.exchangeCodeForToken('test-code');

      expect(token).toBe('test-access-token');
    });

    it('should get user profile', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        json: async () => ({
          id: '123',
          username: 'testuser',
          profile_picture_url: 'http://example.com/avatar.jpg',
        }),
      } as Response);

      const profile = await provider.getUserProfile('test-access-token');

      expect(profile.id).toBe('123');
      expect(profile.name).toBe('testuser');
    });
  });

  describe('OAuthManager', () => {
    const manager = new OAuthManager();

    it('should register and get provider', () => {
      const provider = new FacebookOAuthProvider('test-id', 'test-secret', redirectUri);
      manager.registerProvider(provider);

      const retrieved = manager.getProvider('facebook');
      expect(retrieved).toBe(provider);
    });

    it('should return undefined for unknown provider', () => {
      const retrieved = manager.getProvider('unknown');
      expect(retrieved).toBeUndefined();
    });
  });
});
