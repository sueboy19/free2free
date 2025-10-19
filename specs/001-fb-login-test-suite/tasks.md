---
description: "Task list for Facebook login and API test suite implementation"
---

# Tasks: Facebook 登入與 API 測試套件

**Input**: Design documents from `/specs/001-fb-login-test-suite/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `tests/` at repository root
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project structure per implementation plan in tests/
- [X] T002 Set up test environment configuration files for local testing
- [X] T003 [P] Create test utilities directory structure in tests/testutils/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Create base test configuration in tests/testutils/config.go
- [X] T005 [P] Implement test server setup for Facebook OAuth testing in tests/testutils/test_server.go
- [X] T006 Create mock Facebook OAuth provider for testing in tests/testutils/mock_fb_provider.go
- [X] T007 Set up JWT token validation utility for tests in tests/testutils/jwt_validator.go
- [X] T008 Create helper functions for API testing in tests/testutils/api_helpers.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Facebook 登入功能測試 (Priority: P1) 🎯 MVP

**Goal**: 在本地開發環境中，開發者需要測試 Facebook 登入功能，確保能正確完成 OAuth 流程，取得必要的 JWT token，並驗證登入狀態的有效性。

**Independent Test**: 可以獨立測試 Facebook OAuth 流程的完整性和正確性，確保使用者能夠成功登入並取得適當的認證 token。

### Tests for User Story 1 (REQUIRED - based on spec) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T009 [P] [US1] Contract test for Facebook OAuth endpoints in tests/contract/test_fb_oauth_contract.go
- [X] T010 [P] [US1] Integration test for Facebook OAuth flow in tests/integration/fb_auth_integration_test.go
- [X] T011 [P] [US1] Unit test for JWT token generation in tests/unit/jwt_token_test.go

### Implementation for User Story 1

- [X] T012 [P] [US1] Create Facebook OAuth test helpers in tests/testutils/fb_test_helpers.go
- [X] T013 [US1] Implement Facebook login flow test in tests/e2e/fb_login_e2e_test.go
- [X] T014 [US1] Implement JWT token validation test after Facebook login in tests/e2e/fb_login_e2e_test.go
- [X] T015 [US1] Add test cases for Facebook OAuth callback handling in tests/e2e/fb_login_e2e_test.go
- [X] T016 [US1] Add test for Facebook login failure scenarios in tests/e2e/fb_login_e2e_test.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - 完整 API 功能測試 (Priority: P2)

**Goal**: 在本地環境中，已透過 Facebook 登入的使用者需要能夠正確執行所有 API 功能，包括配對活動、評論、管理員功能等，確保系統功能完整。

**Independent Test**: 可以使用 Facebook 登入後獲得的 token 訪問和測試所有受保護的 API 端點。

### Tests for User Story 2 (REQUIRED - based on spec) ⚠️

- [X] T017 [P] [US2] Contract test for all protected API endpoints in tests/contract/test_protected_endpoints_contract.go
- [X] T018 [P] [US2] Integration test for API endpoints with Facebook JWT token in tests/integration/api_integration_test.go
- [X] T019 [P] [US2] Integration test for all user endpoints in tests/integration/user_api_integration_test.go
- [X] T020 [P] [US2] Integration test for all admin endpoints in tests/integration/admin_api_integration_test.go
- [X] T021 [P] [US2] Integration test for all organizer endpoints in tests/integration/organizer_api_integration_test.go
- [X] T022 [P] [US2] Integration test for all review endpoints in tests/integration/review_api_integration_test.go

### Implementation for User Story 2

- [X] T023 [P] [US2] Create test data generator for API testing in tests/testutils/test_data.go
- [X] T024 [US2] Implement user API endpoints test with Facebook JWT in tests/integration/api_integration_test.go
- [X] T025 [US2] Implement admin API endpoints test with Facebook JWT in tests/integration/api_integration_test.go
- [X] T026 [US2] Implement organizer API endpoints test with Facebook JWT in tests/integration/api_integration_test.go
- [X] T027 [US2] Implement review API endpoints test with Facebook JWT in tests/integration/api_integration_test.go
- [X] T028 [US2] Add test for permission checking of different user roles in tests/integration/api_integration_test.go
- [X] T029 [US2] Add test for expired JWT token handling in tests/integration/api_integration_test.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - 本地環境測試設置 (Priority: P3)

**Goal**: 在本地開發環境中，需要有一套完整的測試設置，讓開發者能夠輕鬆執行 Facebook 登入和 API 測試，確保測試環境與生產環境的一致性。

**Independent Test**: 可以在乾淨的本地環境中設置和運行測試套件，驗證測試環境的完整性和可用性。

### Tests for User Story 3 (REQUIRED - based on spec) ⚠️

- [X] T030 [P] [US3] Test environment setup validation in tests/e2e/env_setup_test.go
- [X] T031 [P] [US3] Test suite execution validation in tests/e2e/test_suite_validation_test.go
- [X] T032 [P] [US3] Performance test for complete flow in tests/performance/fb_login_performance_test.go

### Implementation for User Story 3

- [X] T033 [P] [US3] Create test setup script in scripts/test_setup.sh
- [X] T034 [US3] Implement complete test suite runner in tests/e2e/complete_flow_test.go
- [X] T035 [US3] Add test result reporting mechanism in tests/testutils/result_reporter.go
- [X] T036 [US3] Add test timeout and cleanup in tests/testutils/test_cleanup.go
- [X] T037 [US3] Create README for local test execution in tests/README.md
- [X] T038 [US3] Add test environment validation in tests/e2e/env_setup_test.go

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T039 [P] Documentation updates for test setup in tests/README.md
- [X] T040 Code cleanup and refactoring across test files
- [X] T041 Performance optimization for test execution
- [X] T042 [P] Additional edge case tests in tests/e2e/edge_case_test.go
- [X] T043 Security validation tests for JWT and OAuth tokens
- [X] T044 Run quickstart.md validation for complete test suite

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 (needs Facebook login to get JWT token)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Contract test for Facebook OAuth endpoints in tests/contract/test_fb_oauth_contract.go"
Task: "Integration test for Facebook OAuth flow in tests/integration/fb_auth_integration_test.go"
Task: "Unit test for JWT token generation in tests/unit/jwt_token_test.go"

# Launch all helpers for User Story 1 together:
Task: "Create Facebook OAuth test helpers in tests/testutils/fb_test_helpers.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2 (after US1 foundation)
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence