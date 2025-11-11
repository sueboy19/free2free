---
description: "Task list for Test Fix implementation"
---

# Tasks: Test Fix

**Input**: Design documents from `/specs/003-test-fix/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `handlers/`, `middleware/`, `routes/`, `tests/` at repository root
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project structure per implementation plan for test fixes
- [X] T002 Set up test environment configuration files with proper SESSION_KEY for local testing
- [X] T003 [P] Create test utilities for session handling validation in tests/testutils/session_validation.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Create base session initialization middleware in middleware/session.go
- [X] T005 [P] Implement proper session handling for test environments in tests/testutils/test_session.go
- [X] T006 Update JWT claims structure to match actual implementation in utils/auth.go
- [X] T007 Create test server with all documented routes in tests/testutils/test_server.go
- [X] T008 Create environment setup for OAuth flow testing in tests/testutils/env_setup.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Fix Session Handling in Authentication Flow (Priority: P1) 🎯 MVP

**Goal**: As a user, I need to have stable authentication sessions, so that I can securely access the application without experiencing unexpected disconnections or errors during my session.

**Independent Test**: Can be fully tested by executing authentication flows and verifies that session state is properly maintained across requests, delivering secure and stable user sessions.

### Tests for User Story 1 (REQUIRED - based on spec) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T009 [P] [US1] Unit test for session initialization in tests/unit/session_init_test.go
- [X] T010 [P] [US1] Integration test for authentication endpoint session handling in tests/integration/session_handling_test.go

### Implementation for User Story 1

- [X] T011 [P] [US1] Fix session initialization in handlers/auth_handlers.go
- [X] T012 [US1] Update logout handler to properly handle session in handlers/auth_handlers.go
- [X] T013 [US1] Implement nil session check in middleware/auth.go
- [X] T014 [US1] Add session validation before access in utils/auth.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Fix Missing API Routes (Priority: P2)

**Goal**: As a user, I need to access all documented API functionality, so that I can use all the features that are supposed to be available in the application.

**Independent Test**: Can be fully tested by calling all documented API endpoints and verifies that they return appropriate responses instead of error messages, delivering complete API functionality.

### Tests for User Story 2 (REQUIRED - based on spec) ⚠️

- [X] T015 [P] [US2] Contract test for profile endpoint in tests/contract/profile_endpoint_contract.go
- [X] T016 [US2] Integration test for profile endpoint access in tests/integration/profile_endpoint_test.go

### Implementation for User Story 2

- [X] T017 [P] [US2] Register profile endpoint in routes/user.go
- [X] T018 [US2] Implement profile handler in handlers/user_handlers.go
- [X] T019 [US2] Add profile endpoint to test server router in tests/testutils/test_server.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Fix Test Execution Issues (Priority: P1)

**Goal**: As a quality assurance stakeholder, I need all tests to run successfully, so that we can validate the system functionality and ensure it meets quality standards.

**Independent Test**: Can be fully tested by running the test suite and verifies that all tests execute successfully, delivering confidence in system functionality.

### Tests for User Story 3 (REQUIRED - based on spec) ⚠️

- [X] T020 [P] [US3] Unit test for JWT claims structure in tests/unit/jwt_claims_test.go
- [X] T021 [P] [US3] Unit test for User model fields in tests/unit/user_model_test.go
- [X] T022 [US3] Integration test compilation validation in tests/integration/build_validation_test.go
- [X] T040 [P] [US3] Unit test for authentication error feedback messages in tests/unit/auth_feedback_test.go

### Implementation for User Story 3

- [X] T023 [P] [US3] Fix JWT claims structure in utils/jwt.go
- [X] T024 [US3] Update User model fields to match implementation in models/user.go
- [X] T025 [US3] Fix Activity model fields to match implementation in models/activity.go
- [X] T026 [US3] Update test data to match actual model structures in tests/testutils/test_data.go
- [X] T041 [US3] Implement standardized error response messages in handlers/auth_handlers.go
- [X] T042 [US3] Add comprehensive error handling middleware in middleware/error_handler.go

**Checkpoint**: At this point, User Stories 1, 2 AND 3 should all work independently

---

## Phase 6: User Story 4 - Fix OAuth Authentication Flow (Priority: P1)

**Goal**: As a user, I need the OAuth authentication flow to work correctly, so that I can securely log in using my preferred social media account and access protected resources.

**Independent Test**: Can be fully tested by executing OAuth login flows and verifies that sessions are properly established and maintained, delivering secure authentication functionality.

### Tests for User Story 4 (REQUIRED - based on spec) ⚠️

- [X] T027 [P] [US4] Integration test for OAuth flow completion in tests/integration/oauth_flow_test.go
- [X] T028 [P] [US4] Contract test for OAuth endpoints in tests/contract/oauth_endpoints_contract.go
- [X] T029 [US4] Unit test for OAuth session establishment in tests/unit/oauth_session_test.go

### Implementation for User Story 4

- [X] T030 [P] [US4] Fix OAuth session handling in handlers/auth_handlers.go
- [X] T031 [US4] Update OAuth callback to properly establish sessions in handlers/auth_handlers.go
- [X] T032 [US4] Add OAuth environment validation in utils/auth.go
- [X] T033 [US4] Update OAuth validation for test environments in tests/testutils/oauth_helpers.go

**Checkpoint**: All user stories should now be integrated and working together

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T034 [P] Documentation updates for session handling in docs/session_handling.md
- [X] T035 Code cleanup and refactoring of authentication logic across handlers
- [X] T036 Performance validation for OAuth flow completion time
- [X] T037 [P] Additional edge case tests for session handling in tests/e2e/session_edge_cases_test.go
- [X] T038 Security validation tests for session management
- [X] T039 Run quickstart.md validation for complete test suite

## Phase 8: Implementation Verification and Test Fixes (2025-11-11)

**Purpose**: Fix compilation and integration test issues discovered during implementation verification

- [X] T040 Fix model field mismatches in integration tests (Provider → SocialProvider, etc.)
- [X] T041 Fix type mismatches between uint and int64 in test utilities
- [X] T042 Fix unused variable declarations in auth integration tests
- [X] T043 Update migration test to use correct model fields and types
- [X] T044 Verify basic project build and core test functionality
- [X] T045 Correct JWT token validation in Facebook auth integration tests

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → P4)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 4 (P1)**: Can start after Foundational (Phase 2) - May depend on US1 (session handling fixes)

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
Task: "Unit test for session initialization in tests/unit/session_init_test.go"
Task: "Integration test for authentication endpoint session handling in tests/integration/session_handling_test.go"

# Launch all implementation for User Story 1 together:
Task: "Fix session initialization in handlers/auth_handlers.go"
Task: "Update logout handler to properly handle session in handlers/auth_handlers.go"
Task: "Implement nil session check in middleware/auth.go"
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
5. Add User Story 4 → Test independently → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Session handling - Priority P1)
   - Developer B: User Story 2 (Missing routes - Priority P2)
   - Developer C: User Story 3 (Test execution - Priority P1)
   - Developer D: User Story 4 (OAuth flow - Priority P1)
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