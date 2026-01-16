import type { D1Database } from '@cloudflare/workers-types';
import type {
  User,
  Admin,
  Location,
  Activity,
  Match,
  MatchParticipant,
  Review,
  ReviewLike,
  RefreshToken,
} from '../types';

type DBUser = Omit<User, 'is_admin'> & { is_admin: number };
type DBReviewLike = Omit<ReviewLike, 'is_like'> & { is_like: number };

export class DB {
  constructor(private db: D1Database) {}

  async createUser(user: Omit<User, 'id' | 'created_at' | 'updated_at'>): Promise<User> {
    const result = await this.db
      .prepare(
        `INSERT INTO users (social_id, social_provider, name, email, avatar_url, is_admin)
         VALUES (?, ?, ?, ?, ?, ?)`
      )
      .bind(
        user.social_id,
        user.social_provider,
        user.name,
        user.email,
        user.avatar_url || null,
        user.is_admin ? 1 : 0
      )
      .run();

    const created = await this.db
      .prepare('SELECT * FROM users WHERE id = ?')
      .bind(result.meta.last_row_id)
      .first<DBUser>();

    if (!created) {
      throw new Error('Failed to create user');
    }

    return {
      ...created,
      is_admin: created.is_admin === 1,
    };
  }

  async getUserById(id: number): Promise<User | null> {
    const user = await this.db.prepare('SELECT * FROM users WHERE id = ?').bind(id).first<DBUser>();
    if (!user) return null;
    return { ...user, is_admin: user.is_admin === 1 };
  }

  async getUserBySocialId(socialId: string, provider: string): Promise<User | null> {
    const user = await this.db
      .prepare('SELECT * FROM users WHERE social_id = ? AND social_provider = ?')
      .bind(socialId, provider)
      .first<DBUser>();
    if (!user) return null;
    return { ...user, is_admin: user.is_admin === 1 };
  }

  async updateUser(id: number, user: Partial<User>): Promise<User | null> {
    const updates: string[] = [];
    const values: any[] = [];

    for (const [key, value] of Object.entries(user)) {
      if (key === 'id') continue;

      updates.push(`${key} = ?`);

      if (key === 'is_admin') {
        values.push(value ? 1 : 0);
      } else {
        values.push(value);
      }
    }

    if (updates.length > 0) {
      await this.db
        .prepare(`UPDATE users SET ${updates.join(', ')} WHERE id = ?`)
        .bind(...values, id)
        .run();
    }

    return this.getUserById(id);
  }

  async deleteUser(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM users WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  async listUsers(): Promise<User[]> {
    const result = await this.db.prepare('SELECT * FROM users ORDER BY id DESC').all<DBUser>();
    return (result.results || []).map((u) => ({ ...u, is_admin: u.is_admin === 1 }));
  }

  // Location operations
  async createLocation(location: Omit<Location, 'id'>): Promise<Location> {
    const result = await this.db
      .prepare(`INSERT INTO locations (name, address, latitude, longitude) VALUES (?, ?, ?, ?)`)
      .bind(location.name, location.address, location.latitude, location.longitude)
      .run();

    return this.getLocationById(result.meta.last_row_id) as Promise<Location>;
  }

  async getLocationById(id: number): Promise<Location | null> {
    return await this.db.prepare('SELECT * FROM locations WHERE id = ?').bind(id).first<Location>();
  }

  async listLocations(): Promise<Location[]> {
    const result = await this.db
      .prepare('SELECT * FROM locations ORDER BY id DESC')
      .all<Location>();
    return result.results || [];
  }

  async updateLocation(
    id: number,
    location: Partial<Omit<Location, 'id'>>
  ): Promise<Location | null> {
    const updates: string[] = [];
    const values: any[] = [];

    for (const [key, value] of Object.entries(location)) {
      updates.push(`${key} = ?`);
      values.push(value);
    }

    if (updates.length > 0) {
      await this.db
        .prepare(`UPDATE locations SET ${updates.join(', ')} WHERE id = ?`)
        .bind(...values, id)
        .run();
    }

    return this.getLocationById(id);
  }

  async deleteLocation(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM locations WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  // Activity operations
  async createActivity(activity: Omit<Activity, 'id' | 'location'>): Promise<Activity> {
    const result = await this.db
      .prepare(
        `INSERT INTO activities (title, target_count, location_id, description, created_by)
         VALUES (?, ?, ?, ?, ?)`
      )
      .bind(
        activity.title,
        activity.target_count,
        activity.location_id,
        activity.description || null,
        activity.created_by
      )
      .run();

    return this.getActivityById(result.meta.last_row_id) as Promise<Activity>;
  }

  async getActivityById(id: number): Promise<Activity | null> {
    const activity = await this.db
      .prepare('SELECT * FROM activities WHERE id = ?')
      .bind(id)
      .first<Activity>();

    if (!activity) return null;

    const location = await this.getLocationById(activity.location_id);
    return { ...activity, location: location || undefined };
  }

  async listActivities(): Promise<Activity[]> {
    const result = await this.db
      .prepare('SELECT * FROM activities ORDER BY id DESC')
      .all<Activity>();

    const activities = result.results || [];

    for (const activity of activities) {
      const location = await this.getLocationById(activity.location_id);
      (activity as any).location = location || undefined;
    }

    return activities;
  }

  async updateActivity(
    id: number,
    activity: Partial<Omit<Activity, 'id' | 'location'>>
  ): Promise<Activity | null> {
    const updates: string[] = [];
    const values: any[] = [];

    for (const [key, value] of Object.entries(activity)) {
      updates.push(`${key} = ?`);
      values.push(value);
    }

    if (updates.length > 0) {
      await this.db
        .prepare(`UPDATE activities SET ${updates.join(', ')} WHERE id = ?`)
        .bind(...values, id)
        .run();
    }

    return this.getActivityById(id);
  }

  async deleteActivity(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM activities WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  // Match operations
  async createMatch(match: Omit<Match, 'id' | 'activity' | 'organizer'>): Promise<Match> {
    const result = await this.db
      .prepare(
        `INSERT INTO matches (activity_id, organizer_id, match_time, status) VALUES (?, ?, ?, ?)`
      )
      .bind(match.activity_id, match.organizer_id, match.match_time, match.status || 'open')
      .run();

    return this.getMatchById(result.meta.last_row_id) as Promise<Match>;
  }

  async getMatchById(id: number): Promise<Match | null> {
    const match = await this.db
      .prepare('SELECT * FROM matches WHERE id = ?')
      .bind(id)
      .first<Match>();

    if (!match) return null;

    const activity = await this.getActivityById(match.activity_id);
    const organizer = await this.getUserById(match.organizer_id);

    return { ...match, activity: activity || undefined, organizer: organizer || undefined };
  }

  async listOpenMatches(): Promise<Match[]> {
    const result = await this.db
      .prepare(
        `SELECT * FROM matches WHERE status = ? AND match_time > datetime('now') ORDER BY match_time ASC`
      )
      .bind('open')
      .all<Match>();

    const matches = result.results || [];

    for (const match of matches) {
      const activity = await this.getActivityById(match.activity_id);
      const organizer = await this.getUserById(match.organizer_id);
      (match as any).activity = activity || undefined;
      (match as any).organizer = organizer || undefined;
    }

    return matches;
  }

  async listMatchesByUser(userId: number, status: string = 'completed'): Promise<Match[]> {
    const result = await this.db
      .prepare(
        `SELECT DISTINCT m.* FROM matches m
         JOIN match_participants mp ON m.id = mp.match_id
         WHERE mp.user_id = ? AND m.status = ?
         ORDER BY m.match_time DESC`
      )
      .bind(userId, status)
      .all<Match>();

    const matches = result.results || [];

    for (const match of matches) {
      const activity = await this.getActivityById(match.activity_id);
      const organizer = await this.getUserById(match.organizer_id);
      (match as any).activity = activity || undefined;
      (match as any).organizer = organizer || undefined;
    }

    return matches;
  }

  async updateMatchStatus(
    id: number,
    status: 'open' | 'completed' | 'cancelled'
  ): Promise<Match | null> {
    await this.db.prepare('UPDATE matches SET status = ? WHERE id = ?').bind(status, id).run();
    return this.getMatchById(id);
  }

  async deleteMatch(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM matches WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  // MatchParticipant operations
  async joinMatch(matchId: number, userId: number): Promise<MatchParticipant> {
    const result = await this.db
      .prepare(
        `INSERT INTO match_participants (match_id, user_id, status, joined_at)
         VALUES (?, ?, 'pending', datetime('now'))`
      )
      .bind(matchId, userId)
      .run();

    return this.getMatchParticipantById(result.meta.last_row_id) as Promise<MatchParticipant>;
  }

  async getMatchParticipantById(id: number): Promise<MatchParticipant | null> {
    const participant = await this.db
      .prepare('SELECT * FROM match_participants WHERE id = ?')
      .bind(id)
      .first<MatchParticipant>();

    if (!participant) return null;

    const match = await this.getMatchById(participant.match_id);
    const user = await this.getUserById(participant.user_id);

    return { ...participant, match: match || undefined, user: user || undefined };
  }

  async getMatchParticipant(matchId: number, userId: number): Promise<MatchParticipant | null> {
    const participant = await this.db
      .prepare('SELECT * FROM match_participants WHERE match_id = ? AND user_id = ?')
      .bind(matchId, userId)
      .first<MatchParticipant>();

    if (!participant) return null;

    const match = await this.getMatchById(participant.match_id);
    const user = await this.getUserById(participant.user_id);

    return { ...participant, match: match || undefined, user: user || undefined };
  }

  async listMatchParticipants(matchId: number): Promise<MatchParticipant[]> {
    const result = await this.db
      .prepare('SELECT * FROM match_participants WHERE match_id = ? ORDER BY id DESC')
      .bind(matchId)
      .all<MatchParticipant>();

    const participants = result.results || [];

    for (const participant of participants) {
      const match = await this.getMatchById(participant.match_id);
      const user = await this.getUserById(participant.user_id);
      (participant as any).match = match || undefined;
      (participant as any).user = user || undefined;
    }

    return participants;
  }

  async updateParticipantStatus(
    id: number,
    status: 'pending' | 'approved' | 'rejected'
  ): Promise<MatchParticipant | null> {
    await this.db
      .prepare('UPDATE match_participants SET status = ? WHERE id = ?')
      .bind(status, id)
      .run();
    return this.getMatchParticipantById(id);
  }

  async deleteMatchParticipant(id: number): Promise<boolean> {
    const result = await this.db
      .prepare('DELETE FROM match_participants WHERE id = ?')
      .bind(id)
      .run();
    return (result.meta.changes || 0) > 0;
  }

  // Review operations
  async createReview(
    review: Omit<Review, 'id' | 'match' | 'reviewer' | 'reviewee' | 'created_at'>
  ): Promise<Review> {
    const result = await this.db
      .prepare(
        `INSERT INTO reviews (match_id, reviewer_id, reviewee_id, score, comment, created_at)
         VALUES (?, ?, ?, ?, ?, datetime('now'))`
      )
      .bind(
        review.match_id,
        review.reviewer_id,
        review.reviewee_id,
        review.score,
        review.comment || null
      )
      .run();

    return this.getReviewById(result.meta.last_row_id) as Promise<Review>;
  }

  async getReviewById(id: number): Promise<Review | null> {
    const review = await this.db
      .prepare('SELECT * FROM reviews WHERE id = ?')
      .bind(id)
      .first<Review>();

    if (!review) return null;

    const match = await this.getMatchById(review.match_id);
    const reviewer = await this.getUserById(review.reviewer_id);
    const reviewee = await this.getUserById(review.reviewee_id);

    return {
      ...review,
      match: match || undefined,
      reviewer: reviewer || undefined,
      reviewee: reviewee || undefined,
    };
  }

  async listReviewsByMatch(matchId: number): Promise<Review[]> {
    const result = await this.db
      .prepare('SELECT * FROM reviews WHERE match_id = ? ORDER BY created_at DESC')
      .bind(matchId)
      .all<Review>();

    const reviews = result.results || [];

    for (const review of reviews) {
      const match = await this.getMatchById(review.match_id);
      const reviewer = await this.getUserById(review.reviewer_id);
      const reviewee = await this.getUserById(review.reviewee_id);
      (review as any).match = match || undefined;
      (review as any).reviewer = reviewer || undefined;
      (review as any).reviewee = reviewee || undefined;
    }

    return reviews;
  }

  async listReviewsByReviewer(reviewerId: number): Promise<Review[]> {
    const result = await this.db
      .prepare('SELECT * FROM reviews WHERE reviewer_id = ? ORDER BY created_at DESC')
      .bind(reviewerId)
      .all<Review>();

    const reviews = result.results || [];

    for (const review of reviews) {
      const match = await this.getMatchById(review.match_id);
      const reviewer = await this.getUserById(review.reviewer_id);
      const reviewee = await this.getUserById(review.reviewee_id);
      (review as any).match = match || undefined;
      (review as any).reviewer = reviewer || undefined;
      (review as any).reviewee = reviewee || undefined;
    }

    return reviews;
  }

  async updateReview(id: number, review: Partial<Review>): Promise<Review | null> {
    const updates: string[] = [];
    const values: any[] = [];

    for (const [key, value] of Object.entries(review)) {
      updates.push(`${key} = ?`);
      values.push(value);
    }

    if (updates.length > 0) {
      await this.db
        .prepare(`UPDATE reviews SET ${updates.join(', ')} WHERE id = ?`)
        .bind(...values, id)
        .run();
    }

    return this.getReviewById(id);
  }

  async deleteReview(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM reviews WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  async hasReviewed(reviewerId: number, revieweeId: number, matchId: number): Promise<boolean> {
    const review = await this.db
      .prepare('SELECT 1 FROM reviews WHERE reviewer_id = ? AND reviewee_id = ? AND match_id = ?')
      .bind(reviewerId, revieweeId, matchId)
      .first();

    return review !== undefined;
  }

  // ReviewLike operations
  async likeReview(reviewId: number, userId: number, isLike: boolean): Promise<ReviewLike> {
    const existing = await this.getReviewLike(reviewId, userId);

    if (existing) {
      await this.db
        .prepare('UPDATE review_likes SET is_like = ? WHERE review_id = ? AND user_id = ?')
        .bind(isLike ? 1 : 0, reviewId, userId)
        .run();

      return { ...existing, is_like: isLike };
    }

    const result = await this.db
      .prepare(`INSERT INTO review_likes (review_id, user_id, is_like) VALUES (?, ?, ?)`)
      .bind(reviewId, userId, isLike ? 1 : 0)
      .run();

    return this.getReviewLikeById(result.meta.last_row_id) as Promise<ReviewLike>;
  }

  async getReviewLikeById(id: number): Promise<ReviewLike | null> {
    const like = await this.db
      .prepare('SELECT * FROM review_likes WHERE id = ?')
      .bind(id)
      .first<DBReviewLike>();

    if (!like) return null;

    const review = await this.getReviewById(like.review_id);
    const user = await this.getUserById(like.user_id);

    return {
      ...like,
      review: review || undefined,
      user: user || undefined,
      is_like: like.is_like === 1,
    };
  }

  async getReviewLike(reviewId: number, userId: number): Promise<ReviewLike | null> {
    const like = await this.db
      .prepare('SELECT * FROM review_likes WHERE review_id = ? AND user_id = ?')
      .bind(reviewId, userId)
      .first<DBReviewLike>();

    if (!like) return null;

    const review = await this.getReviewById(like.review_id);
    const user = await this.getUserById(like.user_id);

    return {
      ...like,
      review: review || undefined,
      user: user || undefined,
      is_like: like.is_like === 1,
    };
  }

  async deleteReviewLike(reviewId: number, userId: number): Promise<boolean> {
    const result = await this.db
      .prepare('DELETE FROM review_likes WHERE review_id = ? AND user_id = ?')
      .bind(reviewId, userId)
      .run();

    return (result.meta.changes || 0) > 0;
  }

  // RefreshToken operations
  async createRefreshToken(
    userId: number,
    token: string,
    expiresAt: string
  ): Promise<RefreshToken> {
    const result = await this.db
      .prepare(
        `INSERT INTO refresh_tokens (user_id, token, expires_at, created_at) VALUES (?, ?, ?, datetime('now'))`
      )
      .bind(userId, token, expiresAt)
      .run();

    return this.getRefreshTokenById(result.meta.last_row_id) as Promise<RefreshToken>;
  }

  async getRefreshTokenById(id: number): Promise<RefreshToken | null> {
    const token = await this.db
      .prepare('SELECT * FROM refresh_tokens WHERE id = ?')
      .bind(id)
      .first<RefreshToken>();

    if (!token) return null;

    const user = await this.getUserById(token.user_id);

    return { ...token, user: user || undefined };
  }

  async getRefreshTokenByToken(token: string): Promise<RefreshToken | null> {
    const refreshToken = await this.db
      .prepare('SELECT * FROM refresh_tokens WHERE token = ?')
      .bind(token)
      .first<RefreshToken>();

    if (!refreshToken) return null;

    const user = await this.getUserById(refreshToken.user_id);

    return { ...refreshToken, user: user || undefined };
  }

  async deleteRefreshToken(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM refresh_tokens WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  async deleteRefreshTokensByUserId(userId: number): Promise<number> {
    const result = await this.db
      .prepare('DELETE FROM refresh_tokens WHERE user_id = ?')
      .bind(userId)
      .run();
    return result.meta.changes || 0;
  }

  async cleanupExpiredRefreshTokens(): Promise<number> {
    const result = await this.db
      .prepare("DELETE FROM refresh_tokens WHERE expires_at < datetime('now')")
      .run();

    return result.meta.changes || 0;
  }

  // Admin operations
  async createAdmin(admin: Omit<Admin, 'id'>): Promise<Admin> {
    const result = await this.db
      .prepare(`INSERT INTO admins (username, email) VALUES (?, ?)`)
      .bind(admin.username, admin.email)
      .run();

    return this.getAdminById(result.meta.last_row_id) as Promise<Admin>;
  }

  async getAdminById(id: number): Promise<Admin | null> {
    return await this.db.prepare('SELECT * FROM admins WHERE id = ?').bind(id).first<Admin>();
  }

  async getAdminByUsername(username: string): Promise<Admin | null> {
    return await this.db
      .prepare('SELECT * FROM admins WHERE username = ?')
      .bind(username)
      .first<Admin>();
  }

  async listAdmins(): Promise<Admin[]> {
    const result = await this.db.prepare('SELECT * FROM admins ORDER BY id DESC').all<Admin>();
    return result.results || [];
  }

  async updateAdmin(id: number, admin: Partial<Admin>): Promise<Admin | null> {
    const updates: string[] = [];
    const values: any[] = [];

    for (const [key, value] of Object.entries(admin)) {
      updates.push(`${key} = ?`);
      values.push(value);
    }

    if (updates.length > 0) {
      await this.db
        .prepare(`UPDATE admins SET ${updates.join(', ')} WHERE id = ?`)
        .bind(...values, id)
        .run();
    }

    return this.getAdminById(id);
  }

  async deleteAdmin(id: number): Promise<boolean> {
    const result = await this.db.prepare('DELETE FROM admins WHERE id = ?').bind(id).run();
    return (result.meta.changes || 0) > 0;
  }

  // Helper methods
  async clearAllTables(): Promise<void> {
    await this.db.prepare('DELETE FROM review_likes').run();
    await this.db.prepare('DELETE FROM reviews').run();
    await this.db.prepare('DELETE FROM match_participants').run();
    await this.db.prepare('DELETE FROM matches').run();
    await this.db.prepare('DELETE FROM activities').run();
    await this.db.prepare('DELETE FROM locations').run();
    await this.db.prepare('DELETE FROM refresh_tokens').run();
    await this.db.prepare('DELETE FROM users').run();
    await this.db.prepare('DELETE FROM admins').run();
  }

  async getTableCount(tableName: string): Promise<number> {
    const result = await this.db
      .prepare(`SELECT COUNT(*) as count FROM ${tableName}`)
      .first<{ count: number }>();
    return result?.count || 0;
  }

  async getMatchCountByStatus(status: string): Promise<number> {
    const result = await this.db
      .prepare(`SELECT COUNT(*) as count FROM matches WHERE status = ?`)
      .bind(status)
      .first<{ count: number }>();

    return result?.count || 0;
  }

  async getUserMatchCount(userId: number): Promise<number> {
    const result = await this.db
      .prepare(`SELECT COUNT(*) as count FROM match_participants WHERE user_id = ?`)
      .bind(userId)
      .first<{ count: number }>();

    return result?.count || 0;
  }
}
