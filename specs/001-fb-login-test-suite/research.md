# Research Summary: Facebook 登入與 API 測試套件

## Decision: Facebook OAuth 2.0 測試方法
**Rationale**: 基於現有的 Goth 庫實現，我們將創建端到端測試來驗證 Facebook OAuth 2.0 流程，包括重定向、回調處理和 JWT token 生成。這確保了從用戶點擊 Facebook 登入到獲得可用 JWT token 的完整流程都能正常工作。

**Alternatives considered**: 
- 模擬 OAuth 提供者：雖然更快，但無法測試真實的 Facebook 集成
- 僅單元測試：無法驗證完整的 OAuth 流程

## Decision: JWT Token 使用方式
**Rationale**: 在測試中，我們將使用 Facebook 登入後獲得的 JWT token 來訪問受保護的 API 端點，這模擬了真實用戶的操作流程，確保從認證到授權的整個鏈條都正常工作。

**Alternatives considered**:
- 使用預生成的固定 token：無法驗證 JWT 生成和驗證邏輯
- 測試會話機制：不符合現有架構（主要使用 JWT）

## Decision: 測試環境配置
**Rationale**: 將創建一個本地測試環境配置，包含測試用的 Facebook 應用憑證，並使用環境變量來管理這些測試憑證，確保與生產環境分離。

**Alternatives considered**:
- 使用生產 Facebook 應用：可能會污染生產數據
- 不使用真實 Facebook 應用：無法測試真實的 OAuth 流程

## Decision: API 測試覆蓋範圍
**Rationale**: 將測試所有受保護的 API 端點，確保 Facebook 登入後的 JWT token 可以訪問所有預期的功能，包括管理員、開局者和一般用戶的端點。

**Alternatives considered**:
- 只測試部分端點：無法確保完整的功能可用性
- 測試未受保護的端點：不在本功能範圍內