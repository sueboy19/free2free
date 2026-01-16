import { describe, it, expect, beforeEach } from 'vitest';
import { DB } from '../../src/lib/db';
import type { Env } from '../../src/types';

describe('DB Operations', () => {
  let db: DB;
  let env: Env;

  beforeEach(async () => {
    env = {
      DB: {
        prepare: (query: string) => {
          return {
            bind: (...args: any[]) => {
              const lowerQuery = query.toLowerCase();
              return {
                run: async () => ({ meta: { last_row_id: 1, changes: 1 } }),
                first: async <T>() => {
                  if (lowerQuery.includes('locations')) {
                    return { id: 1, name: 'Test Location', address: '123 Test St', latitude: 25.0479, longitude: 121.5170 } as T;
                  }
                  if (lowerQuery.includes('activities')) {
                    return { id: 1, title: 'Test Activity', target_count: 4, location_id: 1, description: 'Test Description', created_by: 1, location: { id: 1, name: 'Test Location' } } as T;
                  }
                  if (lowerQuery.includes('matches')) {
                    return { id: 1, activity_id: 1, organizer_id: 1, match_time: '2024-01-01T10:00:00Z', status: 'open', activity: { id: 1, title: 'Test Activity' }, organizer: { id: 1, name: 'Test User' } } as T;
                  }
                  if (lowerQuery.includes('match_participants')) {
                    return { id: 1, match_id: 1, user_id: 1, status: 'pending', joined_at: '2024-01-01T10:00:00Z', match: { id: 1, activity_id: 1, organizer_id: 1 }, user: { id: 1, name: 'Test User' } } as T;
                  }
                  if (lowerQuery.includes('reviews')) {
                    return { id: 1, match_id: 1, reviewer_id: 1, reviewee_id: 2, score: 5, comment: 'Great!', created_at: '2024-01-01T10:00:00Z', match: { id: 1, activity_id: 1 }, reviewer: { id: 1, name: 'Test User' }, reviewee: { id: 2, name: 'Test User 2' } } as T;
                  }
                  if (lowerQuery.includes('review_likes')) {
                    return { id: 1, review_id: 1, user_id: 1, is_like: 1, review: { id: 1, score: 5 }, user: { id: 1, name: 'Test User' } } as T;
                  }
                  if (lowerQuery.includes('refresh_tokens')) {
                    return { id: 1, user_id: 1, token: 'test-token', expires_at: '2024-01-01T10:00:00Z', created_at: '2024-01-01T09:00:00Z', user: { id: 1, name: 'Test User' } } as T;
                  }
                  if (lowerQuery.includes('admins')) {
                    return { id: 1, username: 'admin', email: 'admin@example.com' } as T;
                  }
                  return { id: 1, name: 'Test User', email: 'test@example.com', is_admin: 0 } as T;
                },
                all: async <T>() => ({ results: [] as T[] }),
              };
            },
          };
        },
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

    db = new DB(env.DB);
  });

  describe('User Operations', () => {
    it('should create a user', async () => {
      const user = await db.createUser({
        social_id: '123',
        social_provider: 'facebook',
        name: 'Test User',
        email: 'test@example.com',
        avatar_url: 'http://example.com/avatar.jpg',
        is_admin: false,
      });

      expect(user.id).toBeGreaterThan(0);
      expect(user.name).toBe('Test User');
      expect(user.email).toBe('test@example.com');
      expect(user.is_admin).toBe(false);
    });
  });

  describe('Location Operations', () => {
    it('should create a location', async () => {
      const location = await db.createLocation({
        name: 'Test Location',
        address: '123 Test St',
        latitude: 25.0479,
        longitude: 121.5170,
      });

      expect(location.id).toBeGreaterThan(0);
      expect(location.name).toBe('Test Location');
    });
  });

  describe('Activity Operations', () => {
    it('should create an activity', async () => {
      const activity = await db.createActivity({
        title: 'Test Activity',
        target_count: 4,
        location_id: 1,
        description: 'Test Description',
        created_by: 1,
      });

      expect(activity.id).toBeGreaterThan(0);
      expect(activity.title).toBe('Test Activity');
    });
  });

  describe('Match Operations', () => {
    it('should create a match', async () => {
      const match = await db.createMatch({
        activity_id: 1,
        organizer_id: 1,
        match_time: '2024-01-01T10:00:00Z',
        status: 'open',
      });

      expect(match.id).toBeGreaterThan(0);
      expect(match.status).toBe('open');
    });
  });

  describe('MatchParticipant Operations', () => {
    it('should join a match', async () => {
      const participant = await db.joinMatch(1, 1);

      expect(participant.id).toBeGreaterThan(0);
      expect(participant.match_id).toBe(1);
      expect(participant.user_id).toBe(1);
      expect(participant.status).toBe('pending');
    });
  });

  describe('Review Operations', () => {
    it('should create a review', async () => {
      const review = await db.createReview({
        match_id: 1,
        reviewer_id: 1,
        reviewee_id: 2,
        score: 5,
        comment: 'Great!',
      });

      expect(review.id).toBeGreaterThan(0);
      expect(review.score).toBe(5);
      expect(review.comment).toBe('Great!');
    });
  });

  describe('ReviewLike Operations', () => {
    it('should like a review', async () => {
      const like = await db.likeReview(1, 1, true);

      expect(like.id).toBeGreaterThan(0);
      expect(like.review_id).toBe(1);
      expect(like.user_id).toBe(1);
      expect(like.is_like).toBe(true);
    });
  });

  describe('RefreshToken Operations', () => {
    it('should create a refresh token', async () => {
      const token = await db.createRefreshToken(1, 'test-token', '2024-01-01T10:00:00Z');

      expect(token.id).toBeGreaterThan(0);
      expect(token.user_id).toBe(1);
      expect(token.token).toBe('test-token');
    });
  });

  describe('Admin Operations', () => {
    it('should create an admin', async () => {
      const admin = await db.createAdmin({
        username: 'admin',
        email: 'admin@example.com',
      });

      expect(admin.id).toBeGreaterThan(0);
      expect(admin.username).toBe('admin');
      expect(admin.email).toBe('admin@example.com');
    });
  });
});
