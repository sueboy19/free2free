<!-- 
Sync Impact Report:
Version change: 1.0.0 → 1.0.1
Added sections: None
Removed sections: None
Templates requiring updates: 
- ✅ plan-template.md: Constitution Check section aligns with existing principles
- ✅ spec-template.md: No mandatory changes needed for current principles
- ✅ tasks-template.md: No mandatory changes needed for current principles
- ⚠ README.md: Contains project-specific references but no constitution violations
Follow-up TODOs: None
-->
# 買一送一配對網站 Constitution

## Core Principles

### I. 模組化設計優先 (NON-NEGOTIABLE)
所有功能必須以模組化方式設計和實現；程式碼應分離到獨立的包中（如 models、handlers、routes、middleware、utils），每個包應具有明確職責且可獨立測試；禁止將所有程式碼集中在 main.go 中。

### II. API 文件優先 (NON-NEGOTIABLE)
所有 API 端點必須包含完整的 Swagger 文件註解；API 設計需遵循 RESTful 原則；所有端點需提供適當的錯誤處理和狀態碼；API 文件必須保持與實現同步更新。

### III. 測試驅動開發 (NON-NEGOTIABLE)
TDD 強制執行：先寫測試 → 測試通過驗證 → 測試失敗 → 然後實現功能；嚴格執行 Red-Green-Refactor 循環；所有功能代碼必須伴隨充足的單元和整合測試；覆蓋率最低需達到 80%。

### IV. 安全性和認證優先
所有端點必須實施適當的認證和授權機制；OAuth 2.0 與 JWT 應作為主要認證方法；敏感操作需要額外的驗證層；所有安全漏洞必須在部署前修復。

### V. 可擴展性和性能
系統架構必須支持未來功能擴展；資料庫查詢應優化以避免效能瓶頸；API 需處理適當的分頁和快取機制；響應時間需保持在 500ms 以內。

## 安全與合規性要求

所有使用者資料必須加密存儲；OAuth 憑證和 JWT 密鑰必須安全管理；API 需實施速率限制以防止濫用；系統需符合當地隱私權法規。

## 開發工作流程

所有程式碼變更必須通過代碼審查；功能分支必須包含完整的測試案例；合併前必須執行所有測試套件；遵循 Git 工作流程，使用有意義的提交訊息。

## Governance

本憲法優先於所有其他開發實踐；任何與此憲法衝突的實踐必須修改以符合原則；修改憲法需要明確的文件、批准和遷移計劃；所有 PR/審查必須驗證合規性；複雜變更必須有適當的理由說明。

**Version**: 1.0.1 | **Ratified**: 2025-10-19 | **Last Amended**: 2025-11-09
