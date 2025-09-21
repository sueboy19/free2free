# 買一送一配對網站 - API 文件

## 1. 使用者認證

### 1.1 Facebook 登入
**請求:**
```
GET /auth/facebook
```

**回應:**
重新導向至 Facebook OAuth 頁面

### 1.2 Instagram 登入
**請求:**
```
GET /auth/instagram
```

**回應:**
重新導向至 Instagram OAuth 頁面

### 1.3 登出
**請求:**
```
GET /logout
```

**回應:**
重新導向至首頁

### 1.4 取得使用者資訊
**請求:**
```
GET /profile
Authorization: Bearer {token}
```

**回應:**
```json
{
  "id": 1,
  "social_id": "123456789",
  "social_provider": "facebook",
  "name": "張三",
  "email": "zhangsan@example.com",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

## 2. 管理後台

### 2.1 取得配對活動列表
**請求:**
```
GET /admin/activities
Authorization: Bearer {admin_token}
```

**回應:**
```json
[
  {
    "id": 1,
    "title": "全家咖啡買一送一",
    "target_count": 1,
    "location_id": 1,
    "description": "在全家便利商店購買咖啡，買一送一",
    "created_by": 1
  }
]
```

### 2.2 建立配對活動
**請求:**
```
POST /admin/activities
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "title": "7-11 咖啡買一送一",
  "target_count": 1,
  "location_id": 2,
  "description": "在7-11購買咖啡，買一送一"
}
```

**回應:**
```json
{
  "id": 2,
  "title": "7-11 咖啡買一送一",
  "target_count": 1,
  "location_id": 2,
  "description": "在7-11購買咖啡，買一送一",
  "created_by": 1
}
```

### 2.3 更新配對活動
**請求:**
```
PUT /admin/activities/{id}
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "title": "7-11 咖啡買一送一 (更新)",
  "target_count": 1,
  "location_id": 2,
  "description": "在7-11購買咖啡，買一送一"
}
```

**回應:**
```json
{
  "id": 2,
  "title": "7-11 咖啡買一送一 (更新)",
  "target_count": 1,
  "location_id": 2,
  "description": "在7-11購買咖啡，買一送一",
  "created_by": 1
}
```

### 2.4 刪除配對活動
**請求:**
```
DELETE /admin/activities/{id}
Authorization: Bearer {admin_token}
```

**回應:**
```json
{
  "message": "活動已刪除"
}
```

### 2.5 取得地點列表
**請求:**
```
GET /admin/locations
Authorization: Bearer {admin_token}
```

**回應:**
```json
[
  {
    "id": 1,
    "name": "全家便利商店",
    "address": "台北市信義區松山路123號",
    "latitude": 25.044094,
    "longitude": 121.568670
  }
]
```

### 2.6 建立地點
**請求:**
```
POST /admin/locations
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "name": "7-11便利商店",
  "address": "台北市大安區復興南路456號",
  "latitude": 25.043094,
  "longitude": 121.567670
}
```

**回應:**
```json
{
  "id": 2,
  "name": "7-11便利商店",
  "address": "台北市大安區復興南路456號",
  "latitude": 25.043094,
  "longitude": 121.567670
}
```

### 2.7 更新地點
**請求:**
```
PUT /admin/locations/{id}
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "name": "7-11便利商店 (更新)",
  "address": "台北市大安區復興南路456號",
  "latitude": 25.043094,
  "longitude": 121.567670
}
```

**回應:**
```json
{
  "id": 2,
  "name": "7-11便利商店 (更新)",
  "address": "台北市大安區復興南路456號",
  "latitude": 25.043094,
  "longitude": 121.567670
}
```

### 2.8 刪除地點
**請求:**
```
DELETE /admin/locations/{id}
Authorization: Bearer {admin_token}
```

**回應:**
```json
{
  "message": "地點已刪除"
}
```

## 3. 使用者功能

### 3.1 取得配對列表
**請求:**
```
GET /user/matches
Authorization: Bearer {token}
```

**回應:**
```json
[
  {
    "id": 1,
    "activity_id": 1,
    "organizer_id": 1,
    "match_time": "2023-06-15T14:00:00Z",
    "status": "open"
  }
]
```

### 3.2 建立配對局 (開局)
**請求:**
```
POST /user/matches
Authorization: Bearer {token}
Content-Type: application/json

{
  "activity_id": 1,
  "match_time": "2023-06-15T14:00:00Z"
}
```

**回應:**
```json
{
  "id": 1,
  "activity_id": 1,
  "organizer_id": 1,
  "match_time": "2023-06-15T14:00:00Z",
  "status": "open"
}
```

### 3.3 參與配對
**請求:**
```
POST /user/matches/{id}/join
Authorization: Bearer {token}
```

**回應:**
```json
{
  "id": 1,
  "match_id": 1,
  "user_id": 1,
  "status": "pending",
  "joined_at": "2023-06-10T10:00:00Z"
}
```

### 3.4 取得過去參與列表
**請求:**
```
GET /user/past-matches
Authorization: Bearer {token}
```

**回應:**
```json
[
  {
    "id": 1,
    "activity_id": 1,
    "organizer_id": 1,
    "match_time": "2023-06-01T14:00:00Z",
    "status": "completed"
  }
]
```

## 4. 開局者功能

### 4.1 審核通過參與者
**請求:**
```
PUT /organizer/matches/{id}/participants/{participant_id}/approve
Authorization: Bearer {token}
```

**回應:**
```json
{
  "id": 1,
  "match_id": 1,
  "user_id": 2,
  "status": "approved",
  "joined_at": "2023-06-10T10:00:00Z"
}
```

### 4.2 審核拒絕參與者
**請求:**
```
PUT /organizer/matches/{id}/participants/{participant_id}/reject
Authorization: Bearer {token}
```

**回應:**
```json
{
  "id": 1,
  "match_id": 1,
  "user_id": 2,
  "status": "rejected",
  "joined_at": "2023-06-10T10:00:00Z"
}
```

## 5. 評分與互動功能

### 5.1 建立評分與留言
**請求:**
```
POST /review/matches/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "reviewee_id": 2,
  "score": 5,
  "comment": "很好的夥伴，準時赴約"
}
```

**回應:**
```json
{
  "id": 1,
  "match_id": 1,
  "reviewer_id": 1,
  "reviewee_id": 2,
  "score": 5,
  "comment": "很好的夥伴，準時赴約",
  "created_at": "2023-06-01T18:00:00Z"
}
```

### 5.2 點讚評論
**請求:**
```
POST /review-like/reviews/{id}/like
Authorization: Bearer {token}
```

**回應:**
```json
{
  "id": 1,
  "review_id": 1,
  "user_id": 1,
  "is_like": true
}
```

### 5.3 倒讚評論
**請求:**
```
POST /review-like/reviews/{id}/dislike
Authorization: Bearer {token}
```

**回應:**
```json
{
  "id": 1,
  "review_id": 1,
  "user_id": 1,
  "is_like": false
}
```