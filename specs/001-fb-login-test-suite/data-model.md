# Data Model: Facebook 登入與 API 測試套件

## JWT Token
- **Type**: String (JWT格式)
- **Fields**: 
  - user_id (int64): 用戶唯一標識符
  - user_name (string): 用戶名稱
  - is_admin (bool): 管理員權限標記
  - exp (int64): 過期時間戳
  - iat (int64): 簽發時間戳
- **Relationships**: 關聯到 User 實體
- **Validation**: 必須是有效的 JWT 格式，未過期，簽名有效

## Facebook OAuth Session
- **Type**: Session 實體
- **Fields**:
  - session_id (string): 會話唯一標識符
  - user_id (int64): 關聯的用戶ID
  - facebook_user_id (string): Facebook用戶唯一標識符
  - facebook_access_token (string): Facebook訪問令牌
  - created_at (time.Time): 會話創建時間
  - expires_at (time.Time): 會話過期時間
- **Relationships**: 關聯到 User 實體
- **State Transitions**: pending → authenticated → expired

## User (現有實體)
- **Type**: 資料庫模型
- **Fields**:
  - id (int64): 用戶唯一標識符
  - social_id (string): 社交媒體用戶ID (Facebook/Instagram)
  - social_provider (string): 社交媒體提供者 (facebook, instagram)
  - name (string): 用戶名稱
  - email (string): 電子郵箱
  - avatar_url (string): 頭像URL
  - is_admin (bool): 管理員權限
- **Relationships**: 與 JWT Token 和 OAuth Session 關聯

## API Request (測試中使用的虛擬實體)
- **Type**: 測試請求結構
- **Fields**:
  - endpoint (string): API端點路徑
  - method (string): HTTP方法 (GET, POST, PUT, DELETE)
  - headers (map[string]string): 請求頭，包含Authorization
  - body (interface{}): 請求主體
  - expected_status (int): 預期響應狀態碼
- **Validation**: 必須包含有效的 JWT token 在 Authorization 頭中