# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: [e.g., Python 3.11, Swift 5.9, Rust 1.75 or NEEDS CLARIFICATION]  
**Primary Dependencies**: [e.g., FastAPI, UIKit, LLVM or NEEDS CLARIFICATION]  
**Storage**: [if applicable, e.g., PostgreSQL, CoreData, files or N/A]  
**Testing**: [e.g., pytest, XCTest, cargo test or NEEDS CLARIFICATION]  
**Target Platform**: [e.g., Linux server, iOS 15+, WASM or NEEDS CLARIFICATION]
**Project Type**: [single/web/mobile - determines source structure]  
**Performance Goals**: [domain-specific, e.g., 1000 req/s, 10k lines/sec, 60 fps or NEEDS CLARIFICATION]  
**Constraints**: [domain-specific, e.g., <200ms p95, <100MB memory, offline-capable or NEEDS CLARIFICATION]  
**Scale/Scope**: [domain-specific, e.g., 10k users, 1M LOC, 50 screens or NEEDS CLARIFICATION]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Constitution Alignment Assessment:**

1. **模組化設計優先** (NON-NEGOTIABLE):
   - [ ] 驗證功能是否按照模組化方式設計
   - [ ] 確認程式碼分離到獨立包中（models、handlers、routes、middleware、utils）
   - [ ] 檢查是否避免將所有程式碼集中在單一文件中

2. **API 文件優先** (NON-NEGOTIABLE):
   - [ ] 驗證所有 API 端點是否包含完整的 Swagger 文件註解
   - [ ] 確認 API 設計遵循 RESTful 原則
   - [ ] 檢查錯誤處理和狀態碼是否適當
   - [ ] 驗證 API 文件是否與實現同步更新

3. **測試驅動開發** (NON-NEGOTIABLE):
   - [ ] 確認是否遵循 TDD 流程（先寫測試 → 測試通過驗證 → 測試失敗 → 實現功能）
   - [ ] 驗證是否執行 Red-Green-Refactor 循環
   - [ ] 檢查是否包含充足的單元和整合測試
   - [ ] 確認測試覆蓋率是否達到最低 80% 要求

4. **安全性和認證優先**:
   - [ ] 驗證所有端點是否實施適當的認證和授權機制
   - [ ] 確認 OAuth 2.0 與 JWT 是否作為主要認證方法
   - [ ] 檢查敏感操作是否需要額外的驗證層
   - [ ] 驗證所有安全漏洞是否在部署前修復

5. **可擴展性和性能**:
   - [ ] 確認系統架構是否支持未來功能擴展
   - [ ] 驗證資料庫查詢是否優化以避免效能瓶頸
   - [ ] 檢查 API 是否實施適當的分頁和快取機制
   - [ ] 確認響應時間是否保持在 500ms 以內

**Security & Compliance Requirements:**
- [ ] 所有使用者資料是否加密存儲
- [ ] OAuth 憑證和 JWT 密鑰是否安全管理
- [ ] API 是否實施速率限制以防止濫用
- [ ] 系統是否符合當地隱私權法規

**Development Workflow Compliance:**
- [ ] 所有程式碼變更是否通過代碼審查
- [ ] 功能分支是否包含完整的測試案例
- [ ] 合併前是否執行所有測試套件
- [ ] 是否遵循 Git 工作流程和有意義的提交訊息

**Gate Decision:**
- [ ] **PASS**: 所有憲法原則符合要求，可以繼續 Phase 0 研究
- [ ] **FAIL**: 存在憲法違反，必須在繼續前解決

**Notes:**
[記錄任何特殊的合規性考慮或偏離原因]

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |

