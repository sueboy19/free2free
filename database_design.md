# 資料庫設計

## 系統架構
採用六角架構 (Hexagonal Architecture) 來分離核心業務邏輯與外部依賴 (如資料庫、第三方API)。
- **核心層**: 業務邏輯與領域模型
- **適配器層**: 資料庫存取、第三方API整合
- **入口層**: HTTP API、CLI

## 資料表設計

### 1. users (使用者)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    social_id VARCHAR(255) NOT NULL UNIQUE, -- Facebook/Instagram ID
    social_provider ENUM('facebook', 'instagram') NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_social_id_provider (social_id, social_provider)
);
```

### 2. admins (管理員)
```sql
CREATE TABLE admins (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL, -- 使用 bcrypt 加密
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 3. locations (地點)
```sql
CREATE TABLE locations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL, -- 例: "全家便利商店"
    address TEXT NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
);
```

### 4. activities (配對活動)
```sql
CREATE TABLE activities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL, -- 例: "全家咖啡買一送一"
    target_count INT NOT NULL, -- 需求人數
    location_id BIGINT NOT NULL,
    description TEXT,
    created_by BIGINT NOT NULL, -- 管理員 ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES admins(id) ON DELETE CASCADE,
    INDEX idx_location_id (location_id)
);
```

### 5. matches (配對局)
```sql
CREATE TABLE matches (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL,
    organizer_id BIGINT NOT NULL, -- 開局者 ID
    match_time DATETIME NOT NULL,
    status ENUM('open', 'closed', 'completed') DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE,
    FOREIGN KEY (organizer_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_activity_id (activity_id),
    INDEX idx_organizer_id (organizer_id),
    INDEX idx_match_time_status (match_time, status)
);
```

### 6. match_participants (配對參與者)
```sql
CREATE TABLE match_participants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    match_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending', -- 審核狀態
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_match_user (match_id, user_id),
    INDEX idx_match_id (match_id),
    INDEX idx_user_id (user_id)
);
```

### 7. reviews (評分與留言)
```sql
CREATE TABLE reviews (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    match_id BIGINT NOT NULL,
    reviewer_id BIGINT NOT NULL, -- 評分者
    reviewee_id BIGINT NOT NULL, -- 被評分者
    score TINYINT NOT NULL CHECK (score >= 3 AND score <= 5), -- 3-5分
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewee_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_reviewer_reviewee_match (reviewer_id, reviewee_id, match_id),
    INDEX idx_match_id (match_id)
);
```

### 8. review_likes (評論點讚/倒讚)
```sql
CREATE TABLE review_likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    review_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    is_like BOOLEAN NOT NULL, -- true: 點讚, false: 倒讚
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_review (user_id, review_id),
    INDEX idx_review_id (review_id)
);
```

## 索引策略
1. 在經常查詢的欄位上建立索引 (如 foreign keys, status)
2. 在時間相關查詢上建立複合索引 (如 match_time + status)
3. 在唯一性約束上建立唯一索引

## 資料完整性
1. 使用 foreign key constraints 確保關聯資料一致性
2. 使用 enum 限制欄位值範圍
3. 使用 CHECK constraint 驗證資料有效性 (如 score 範圍)