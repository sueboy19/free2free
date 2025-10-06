# 專案摘要
請用中文，每次工作請細分成一份工作明細清單，每次執行完後都要更這件工作清單的狀態。把不必要的檔案刪除或歸放到相關目錄內。工作完成把相關結果濃縮在這件文件，變成記憶工作文件。將檔案歸檔，建立必要的目錄，減少每次讀取檔案，避免讀取太多。

## 整體目標
目標是開發一個「買一送一」配對網站，透過 Facebook/Instagram 進行使用者認證、管理活動/地點的管理員面板、使用者配對功能，以及帶有 Swagger API 文件的評論系統。

## 關鍵知識
- **技術堆疊**：Go 1.25 + Gin framework + GORM + MariaDB + Goth OAuth library
- **開發工具**：Air 用於熱重載，Swagger 用於 API 文件
- **資料庫**：透過 Docker Compose 使用 MariaDB，並使用 GORM 進行自動 schema 遷移
- **認證**：OAuth 2.0 與 Facebook 和 Instagram 提供者
- **環境配置**：使用 .env 檔案，包含 DB 連線、session 金鑰和 OAuth 憑證的變數
- **專案結構**：模組化設計，包含 admin、user、organizer、review 和 review-like 功能的獨立檔案
- **Windows 相容性**：避免使用 Makefile，改用 batch 腳本，使用 air 而非 make 進行開發
- **API 文件**：涵蓋所有端點和資料模型的全面 Swagger/OpenAPI 文件
- 專案曾有冗餘的 JWT 相關函數（`generateJWT`/`validateJWT`）和結構（`JwtClaims`），已移除並改用 `generateJWTToken`/`validateJWTToken` 和 `Claims` 結構，以消除程式碼重複
- 模型定義散佈於各檔案（建議提取至獨立 models 包）；認證機制結合 Goth OAuth 與 JWT（24小時過期，包含 IsAdmin 旗標）

## 最近動作
- 成功從 MySQL 遷移資料庫至 MariaDB，並設定 Docker Compose
- 實作基於 GORM 的全面資料模型和關係
- 為所有 API 端點新增 Swagger/OpenAPI 文件註解
- 配置 Air 熱重載開發環境
- 建立 Windows 相容的 batch 腳本，用於建置和運行應用程式
- 更新環境變數處理，包含 DB_HOST 配置
- 修復多項編譯問題和依賴衝突
- 識別並移除冗餘的 JWT 函數（`generateJWT`、`validateJWT`）和結構（`JwtClaims`）
- 標準化單一 JWT 實作，使用 `generateJWTToken`、`validateJWTToken` 和 `Claims` 結構，包含 IsAdmin 欄位
- 更新 JWT token 生成，包含 claims 中的 IsAdmin 欄位
- 進行全面專案結構和代碼分析，確認模組化設計良好、依賴完整、運行狀態穩定（Air 在 :8080 運行，Docker MariaDB 設定正常）

## 目前計劃
1.  [DONE] 設定 MariaDB 資料庫與 Docker Compose
2.  [DONE] 實作 GORM 模型和自動遷移
3.  [DONE] 為所有端點新增 Swagger API 文件
4.  [DONE] 配置 Air 熱重載開發環境
5.  [DONE] 建立 Windows 相容的開發腳本
6.  [DONE] 分析程式碼邏輯並移除冗餘程式碼
7.  [DONE] 運行測試以確保移除後的程式碼完整性
8.  [DONE] 驗證所有 API 端點在 Swagger UI 中正常運作
9.  [TODO] 實作全面錯誤處理和驗證
10. [TODO] 新增單元和整合測試
11. [TODO] 部署並在 staging 環境中測試
12. [TODO] 提取模型至獨立包
13. [TODO] 統一錯誤處理中間件
14. [TODO] 增強 JWT 安全（添加 refresh token）
15. [TODO] 優化 Docker 生產環境

關於用戶在 Swagger 中使用 Facebook 登入的問題，這需要實作一個特殊的認證機制，因為 Swagger UI 本身無法直接處理 OAuth 重定向。通常的做法是：
1. 在 Swagger 中新增一個 API 金鑰認證選項
2. 用戶先透過網站前端完成 Facebook 登入，獲取 JWT token 或 session
3. 將 token/session ID 手動輸入至 Swagger UI 的認證欄位中
4. Swagger 會在後續請求中將該 token 作為 Authorization header 發送

這需要在後端實作相應的 JWT token 生成和驗證機制，或允許 Swagger 直接使用 session ID 進行認證。

---

## 系統分析
**強項**：模組化設計、Swagger 文件完整、Air 熱重載開發環境高效。
**弱點**：模型定義散佈於多檔案、管理員檢查邏輯過於簡化、重複驗證邏輯存在。
**建議**：提取 models 至獨立包、統一錯誤處理中間件、強化 OAuth 驗證流程、增加更多單元和整合測試、審核 JWT 安全實作。

## 摘要元數據
**Update time**: 2025-10-01T15:42:35.780Z
