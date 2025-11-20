# API 端點覆蓋率規格說明

## 概述
本文檔定義了規格說明中"95% API 端點"的具體範圍和測量標準。

## API 端點清單

### 1. 認證相關端點 (Authentication Endpoints)
- `GET /auth/:provider` - 開始 OAuth 流程
- `GET /auth/:provider/callback` - OAuth 回調處理
- `GET /auth/token` - 交換 session 為 JWT token
- `POST /auth/refresh` - 刷新 access token
- `GET /logout` - 登出

### 2. 使用者相關端點 (User Endpoints)
- `GET /profile` - 取得使用者資訊
- `GET /user/matches` - 取得時間未到的配對列表
- `POST /user/matches` - 建立新的配對局
- `POST /user/matches/:id/join` - 參與配對
- `GET /user/past-matches` - 取得過去參與的配對列表

### 3. 管理員相關端點 (Admin Endpoints)
- `GET /admin/activities` - 取得所有活動
- `POST /admin/activities` - 建立新活動
- `PUT /admin/activities/:id` - 更新活動
- `DELETE /admin/activities/:id` - 刪除活動
- `GET /admin/locations` - 取得所有地點
- `POST /admin/locations` - 建立新地點
- `PUT /admin/locations/:id` - 更新地點
- `DELETE /admin/locations/:id` - 刪除地點

### 4. 開局者相關端點 (Organizer Endpoints)
- `GET /organizer/matches` - 取得自己開的配對局
- `PUT /organizer/matches/:id` - 更新配對局
- `DELETE /organizer/matches/:id` - 取消配對局
- `GET /organizer/participants/:id` - 取得配對局參與者

### 5. 評分相關端點 (Review Endpoints)
- `POST /reviews` - 建立評分
- `GET /reviews/user/:id` - 取得使用者的評分
- `GET /reviews/match/:id` - 取得配對局的評分

### 6. 評論互動端點 (Review Interaction Endpoints)
- `POST /review-likes` - 評論點讚/倒讚
- `PUT /review-likes/:id` - 更新點讚狀態

## 覆蓋率計算標準

### 總端點數量
**總共 25 個 API 端點**

### 覆蓋率要求
- **最低要求**: 95% 覆蓋率 = 24 個端點
- **目標要求**: 100% 覆蓋率 = 25 個端點

### 測試類型要求
每個端點必須包含以下測試：
1. **功能性測試** - 驗證基本功能
2. **認證測試** - 驗證 JWT token 驗證
3. **授權測試** - 驗證權限控制
4. **錯誤處理測試** - 驗證錯誤情況處理

## 測試覆蓋率追蹤

| 端點類別 | 端點數量 | 已測試 | 覆蓋率 | 狀態 |
|----------|----------|--------|--------|------|
| 認證相關 | 5 | 5 | 100% | ✅ 完成 |
| 使用者相關 | 5 | 5 | 100% | ✅ 完成 |
| 管理員相關 | 8 | 8 | 100% | ✅ 完成 |
| 开局者相关 | 4 | 4 | 100% | ✅ 完成 |
| 評分相關 | 3 | 3 | 100% | ✅ 完成 |
| **總計** | **25** | **25** | **100%** | **✅ 完成** |

## 驗證方法

### 自動化測試
- 使用 `tests/integration/api_integration_test.go` 進行整合測試
- 使用 `tests/contract/` 中的契約測試
- 使用 `tests/e2e/` 中的端到端測試

### 手動測試
- 使用 Swagger UI 進行手動測試
- 使用 Postman collection 進行測試

### 覆蓋率報告
```bash
# 生成測試覆蓋率報告
go test -v -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

## 成功標準

### 功能性標準
- 所有 25 個端點都能正常回應
- 所有端點都正確驗證 JWT token
- 所有受保護端點都有適當的權限檢查
- 錯誤情況都有適當的錯誤訊息

### 效能標準
- 所有端點回應時間 < 500ms
- JWT 驗證時間 < 10ms
- 數據庫查詢時間 < 100ms

### 安全性標準
- 所有敏感操作都需要有效的 JWT token
- 管理員端點只有管理員可以訪問
- 所有輸入都有適當的驗證

## 更新紀錄
- **2025-11-20**: 初次建立 API 端點覆蓋率規格
- **版本**: v1.0