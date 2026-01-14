/**
 * 資料匯入腳本：用於 D1 測試
 *
 * 使用方式：
 *   wrangler d1 execute free2free-db --file=./scripts/import-to-d1.sql
 */

-- 清空所有資料表
DELETE FROM review_likes;
DELETE FROM reviews;
DELETE FROM match_participants;
DELETE FROM matches;
DELETE FROM activities;
DELETE FROM locations;
DELETE FROM refresh_tokens;
DELETE FROM users;
DELETE FROM admins;

-- 插入測試資料

-- Admins
INSERT INTO admins (username, email) VALUES ('admin', 'admin@free2free.com');

-- Users
INSERT INTO users (social_id, social_provider, name, email, avatar_url, is_admin) VALUES
  ('1234567890', 'facebook', 'Test User 1', 'user1@example.com', 'https://example.com/avatar1.jpg', 0),
  ('0987654321', 'facebook', 'Test User 2', 'user2@example.com', 'https://example.com/avatar2.jpg', 0),
  ('111222333', 'facebook', 'Admin User', 'admin@free2free.com', 'https://example.com/avatar3.jpg', 1);

-- Locations
INSERT INTO locations (name, address, latitude, longitude) VALUES
  ('台北車站', '台北市中正區北平西路3號', 25.0479, 121.5170),
  ('新北板橋', '新北市板橋區縣民大道二段7號', 25.0124, 121.4635),
  ('台中車站', '台中市東區台灣大道一段1號', 24.1477, 120.6736);

-- Activities
INSERT INTO activities (title, target_count, location_id, description, created_by) VALUES
  ('羽毛球雙打', 4, 1, '週末羽毛球雙打，歡迎初學者', 3),
  ('跑步團', 10, 2, '週末晨跑，配速 5:30-6:00', 3),
  ('桌派對', 6, 3, '歡迎所有桌派對愛好者參加', 1);

-- Matches
INSERT INTO matches (activity_id, organizer_id, match_time, status) VALUES
  (1, 1, datetime('now', '+1 day'), 'open'),
  (2, 1, datetime('now', '+2 days'), 'open'),
  (1, 2, datetime('now', '+3 days'), 'completed');

-- Match Participants
INSERT INTO match_participants (match_id, user_id, status, joined_at) VALUES
  (1, 2, 'approved', datetime('now')),
  (2, 2, 'pending', datetime('now')),
  (3, 1, 'approved', datetime('now'));

-- Reviews
INSERT INTO reviews (match_id, reviewer_id, reviewee_id, score, comment, created_at) VALUES
  (3, 1, 2, 5, '很好的搭檔！', datetime('now', '-1 day')),
  (3, 2, 1, 4, '愉快的體驗', datetime('now', '-1 day'));

-- Review Likes
INSERT INTO review_likes (review_id, user_id, is_like) VALUES
  (1, 2, 1),
  (2, 1, 1);

-- Refresh Tokens (測試用）
INSERT INTO refresh_tokens (user_id, token, expires_at, created_at) VALUES
  (1, 'test-refresh-token-1', datetime('now', '+7 days'), datetime('now'));
