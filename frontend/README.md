# 買一送一配對網站前端

這是一個基於 Vue 3 + TypeScript + Tailwind CSS 的響應式前端應用，為買一送一配對網站提供完整的用戶介面。

## 功能特色

### 🔐 用戶認證
- Facebook 登入
- Instagram 登入
- JWT Token 管理
- 自動會話恢復

### 🎯 配對功能
- 瀏覽可用的配對機會
- 創建新的配對活動
- 參與他人開局的配對
- 配對狀態管理

### 👥 用戶管理
- 個人資料管理
- 配對歷史記錄
- 開局與參與統計

### ⚙️ 管理後台
- 配對活動管理 (CRUD)
- 地點管理 (CRUD)
- 數據統計

### ⭐ 評分系統
- 配對完成後評分
- 評論功能
- 互動反饋

### 📱 響應式設計
- 支援手機端
- 平板電腦
- 桌面端

## 技術棧

- **Vue 3** - 現代化前端框架
- **TypeScript** - 類型安全
- **Vite** - 快速建置工具
- **Vue Router** - 路由管理
- **Pinia** - 狀態管理
- **Axios** - HTTP 客戶端
- **Tailwind CSS** - 實用優先的 CSS 框架
- **Vue Toastification** - 通知組件
- **Date-fns** - 日期處理

## 專案結構

```
frontend/
├── public/                 # 靜態資源
├── src/
│   ├── components/        # 可重用組件
│   ├── views/            # 頁面視圖
│   │   ├── Home.vue         # 首頁
│   │   ├── Login.vue        # 登入頁
│   │   ├── Matches.vue      # 配對列表
│   │   ├── CreateMatch.vue  # 創建配對
│   │   ├── MyMatches.vue    # 我的配對
│   │   ├── MatchDetails.vue # 配對詳情
│   │   ├── Profile.vue      # 個人資料
│   │   └── Admin.vue        # 管理後台
│   ├── stores/           # Pinia 狀態管理
│   │   └── auth.ts          # 認證狀態
│   ├── services/         # API 服務
│   │   └── api.ts           # API 客戶端
│   ├── router/           # 路由配置
│   │   └── index.ts         # 路由定義
│   ├── style.css         # 全域樣式
│   ├── main.ts           # 應用入口
│   └── App.vue           # 根組件
├── index.html            # HTML 模板
├── package.json          # 專案依賴
├── tsconfig.json         # TypeScript 配置
├── tailwind.config.js    # Tailwind 配置
├── vite.config.ts        # Vite 配置
└── README.md             # 專案說明
```

## 安裝和運行

### 環境要求
- Node.js 16+ 
- npm 或 yarn

### 安裝依賴
```bash
cd frontend
npm install
```

### 環境配置
創建 `.env` 檔案：
```env
VITE_API_BASE_URL=http://localhost:8080
```

### 開發模式
```bash
npm run dev
```
應用將在 http://localhost:3000 運行

### 建置生產版本
```bash
npm run build
```

### 預覽生產版本
```bash
npm run preview
```

## 頁面路由

| 路徑 | 頁面 | 需要登入 | 需要管理員 |
|------|------|----------|------------|
| `/` | 首頁 | ❌ | ❌ |
| `/login` | 登入頁 | ❌ | ❌ |
| `/matches` | 配對列表 | ✅ | ❌ |
| `/matches/create` | 創建配對 | ✅ | ❌ |
| `/matches/:id` | 配對詳情 | ✅ | ❌ |
| `/my-matches` | 我的配對 | ✅ | ❌ |
| `/profile` | 個人資料 | ✅ | ❌ |
| `/admin` | 管理後台 | ✅ | ✅ |

## API 整合

應用已整合以下 API 端點：

### 認證相關
- `GET /auth/facebook` - Facebook 登入
- `GET /auth/instagram` - Instagram 登入  
- `GET /profile` - 獲取用戶資料
- `GET /auth/token` - 交換 JWT Token
- `GET /logout` - 登出

### 配對活動
- `GET /admin/activities` - 獲取活動列表
- `POST /admin/activities` - 創建活動
- `PUT /admin/activities/:id` - 更新活動
- `DELETE /admin/activities/:id` - 刪除活動

### 地點管理
- `GET /admin/locations` - 獲取地點列表
- `POST /admin/locations` - 創建地點
- `PUT /admin/locations/:id` - 更新地點
- `DELETE /admin/locations/:id` - 刪除地點

### 配對功能
- `GET /user/matches` - 獲取配對列表
- `POST /user/matches` - 創建配對
- `POST /user/matches/:id/join` - 參與配對
- `GET /user/past-matches` - 獲取歷史配對

### 審核功能
- `PUT /organizer/matches/:id/participants/:participant_id/approve` - 審核通過
- `PUT /organizer/matches/:id/participants/:participant_id/reject` - 審核拒絕

### 評分功能
- `POST /review/matches/:id` - 建立評分
- `POST /review-like/reviews/:id/like` - 點讚
- `POST /review-like/reviews/:id/dislike` - 倒讚

## 響應式設計

應用採用移動優先的設計原則：

### 斷點
- `sm`: 640px+
- `md`: 768px+  
- `lg`: 1024px+
- `xl`: 1280px+

### 特性
- 適配手機、平板、桌面
- 觸控友好的按鈕尺寸
- 簡化的手機端導航
- 優化的表單體驗

## 開發說明

### 狀態管理
使用 Pinia 進行狀態管理，主要包含：
- `useAuthStore` - 用戶認證狀態
- API 調用和錯誤處理
- 本地存儲管理

### 樣式系統
- Tailwind CSS 實用類
- 自訂設計系統
- 響應式組件
- 深色模式支援 (可擴展)

### 類型安全
- TypeScript 嚴格模式
- 完整類型定義
- 編譯時檢查

## 部署建議

### 生產環境
1. 設定正確的 `VITE_API_BASE_URL`
2. 啟用 HTTPS
3. 設定 CORS 政策
4. 啟用 Gzip 壓縮
5. 設定 CDN

### Docker 部署
```dockerfile
FROM node:16-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 故障排除

### 常見問題
1. **模組找不到錯誤**: 確保所有依賴已安裝
2. **API 調用失敗**: 檢查後端服務和 CORS 設定
3. **樣式問題**: 確保 Tailwind CSS 正確配置
4. **路由問題**: 檢查 Vue Router 配置

### 調試技巧
- 使用 Vue DevTools
- 檢查瀏覽器控制台
- 查看網路請求
- 驗證環境變數

## 貢獻指南

1. Fork 專案
2. 創建功能分支
3. 提交變更
4. 發起 Pull Request

## 授權

MIT License