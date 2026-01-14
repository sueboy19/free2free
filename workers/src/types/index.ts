export interface Env {
  DB: D1Database;
  KV: KVNamespace;
  JWT_SECRET: string;
  SESSION_KEY: string;
  FACEBOOK_KEY: string;
  FACEBOOK_SECRET: string;
  INSTAGRAM_KEY: string;
  INSTAGRAM_SECRET: string;
  BASE_URL: string;
  FRONTEND_URL: string;
  CORS_ORIGINS: string;
}

export interface User {
  id: number;
  social_id: string;
  social_provider: 'facebook' | 'instagram';
  name: string;
  email: string;
  avatar_url?: string;
  is_admin: boolean;
  created_at: number;
  updated_at: number;
}

export interface Admin {
  id: number;
  username: string;
  email: string;
}

export interface Location {
  id: number;
  name: string;
  address: string;
  latitude: number;
  longitude: number;
}

export interface Activity {
  id: number;
  title: string;
  target_count: number;
  location_id: number;
  description?: string;
  created_by: number;
  location?: Location;
}

export interface Match {
  id: number;
  activity_id: number;
  organizer_id: number;
  match_time: string;
  status: 'open' | 'completed' | 'cancelled';
  activity?: Activity;
  organizer?: User;
}

export interface MatchParticipant {
  id: number;
  match_id: number;
  user_id: number;
  status: 'pending' | 'approved' | 'rejected';
  joined_at: string;
  match?: Match;
  user?: User;
}

export interface Review {
  id: number;
  match_id: number;
  reviewer_id: number;
  reviewee_id: number;
  score: number;
  comment?: string;
  created_at: string;
  match?: Match;
  reviewer?: User;
  reviewee?: User;
}

export interface ReviewLike {
  id: number;
  review_id: number;
  user_id: number;
  is_like: boolean;
  review?: Review;
  user?: User;
}

export interface RefreshToken {
  id: number;
  user_id: number;
  token: string;
  expires_at: string;
  created_at: string;
  user?: User;
}

export interface JWTPayload {
  user_id: number;
  user_name: string;
  is_admin: boolean;
  exp: number;
  iat: number;
}
