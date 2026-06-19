# Tasks: Battleship Game Engine

**Input**: Design documents from `/specs/001-battleship-game-engine/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are OPTIONAL - only included where explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project structure per plan.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan: `src/models/`, `src/services/`, `src/adapters/api/`, `src/adapters/db/`, `src/lib/`, `tests/unit/`, `tests/integration/`, `tests/contract/`
- [ ] T002 Initialize Go project with `go mod init battleship-game-engine` and add dependencies (gin, gorm, postgres, zap, otel, prometheus, jwt, limiter, validator)
- [ ] T003 [P] Configure linting (golangci-lint) and formatting (gofmt) tools in `.golangci.yml` and `Makefile`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T004 Setup database schema and migrations framework: Create `src/adapters/db/models.go` with GORM models (Game, Board, Ship, Player) and `src/adapters/db/migrate.go` for migration management
- [ ] T005 [P] Implement authentication/authorization framework: Create `src/adapters/gin/middleware/jwt.go` with JWT middleware and `src/adapters/gin/middleware/authorization.go` with player access control
- [ ] T006 [P] Setup API routing and middleware structure: Create `src/adapters/gin/router.go` with middleware chain (security headers, rate limiting, error handler, tracing, correlation ID, logging, recovery)
- [ ] T007 Create base models/entities that all stories depend on: Create `src/models/game.go`, `src/models/board.go`, `src/models/ship.go`, `src/models/player.go` with core business rules
- [ ] T008 Configure error handling and logging infrastructure: Create `src/lib/logger/logger.go` with zap configuration and `src/lib/validation/validation.go` for input validation
- [ ] T009 Setup environment configuration management: Create `src/config/config.go` for loading environment variables (database URL, JWT secret, rate limits, etc.)
- [ ] T010 [P] Create health check endpoints: Create `src/adapters/gin/handlers/health.go` with `/health/live` and `/health/ready` endpoints

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Single Player Battleship Game (Priority: P1) 🎯 MVP

**Goal**: Enable a single player to start a new Battleship game, take turns shooting at the computer's ships, and see visual feedback until all enemy ships are sunk

**Independent Test**: Call startGame() with default 8x8 board, repeatedly call shoot() with valid coordinates until all ships are sunk, verifying all required data is returned correctly

### Implementation for User Story 1

- [ ] T011 [P] [US1] Create Ship entity model in `src/models/ship.go` with placement logic and hit tracking
- [ ] T012 [P] [US1] Create Board entity model in `src/models/board.go` with cell state management (space/O/X)
- [ ] T013 [P] [US1] Create Game entity model in `src/models/game.go` with turn management and win detection
- [ ] T014 [US1] Implement GameService in `src/services/game_service.go` with startGame() and shoot() methods
- [ ] T015 [US1] Implement Game API handlers in `src/adapters/gin/handlers/game.go` for POST /games and POST /games/:game_id/shoot
- [ ] T016 [US1] Add input validation for board size (5-100) and coordinates (0 to N-1) in `src/lib/validation/game.go`
- [ ] T017 [US1] Add logging for game operations in `src/services/game_service.go` (INFO for operations, ERROR for failures)
- [ ] T018 [US1] Implement ship placement retry logic (max 100 attempts) in `src/models/ship.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Display Game Board State (Priority: P1)

**Goal**: Enable users to see the current state of the game board showing which cells have been targeted and whether those shots were hits or misses

**Independent Test**: Verify the hits/misses array is properly formatted and displays correctly as a 2D grid showing O for misses and X for hits

### Implementation for User Story 2

- [ ] T019 [P] [US2] Create Board display helper in `src/lib/display/board_display.go` for converting board state to 2D character grid
- [ ] T020 [US2] Implement board state retrieval endpoint in `src/adapters/gin/handlers/game.go` for GET /games/:game_id
- [ ] T021 [US2] Add board visualization to game state response in `src/adapters/gin/handlers/game.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Two Player Alternate Turns (Priority: P2)

**Goal**: Enable two players to take turns playing Battleship on the same terminal, with each player seeing their own board and the opponent's targeted cells

**Independent Test**: Start a 2-player game, verify turn alternation works correctly, and ensure each player sees only their own board state

### Implementation for User Story 3

- [ ] T022 [P] [US3] Update Game entity to support multiple players in `src/models/game.go` with player-specific board tracking
- [ ] T023 [P] [US3] Update Ship entity for multi-player support in `src/models/ship.go` with player association
- [ ] T024 [US3] Implement turn alternation logic in `src/services/game_service.go` for shoot() method
- [ ] T025 [US3] Add turn enforcement middleware in `src/adapters/gin/middleware/turn.go` to reject out-of-turn shots
- [ ] T026 [US3] Update game state endpoint to return player-specific board in `src/adapters/gin/handlers/game.go`

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently

---

## Phase 6: User Story 4 - Game Statistics Display (Priority: P2)

**Goal**: Enable users to view game statistics at any time during gameplay without interrupting the game flow

**Independent Test**: Call gameStats() at various points in the game and verify all metrics are accurate

### Implementation for User Story 4

- [ ] T027 [P] [US4] Create GameStats entity in `src/models/stats.go` with metrics calculation
- [ ] T028 [US4] Implement StatsService in `src/services/stats_service.go` with gameStats() method
- [ ] T029 [US4] Implement game statistics endpoint in `src/adapters/gin/handlers/game.go` for GET /games/:game_id/stats
- [ ] T030 [US4] Add statistics logging in `src/services/stats_service.go` for observability

**Checkpoint**: At this point, User Stories 1-4 should all work independently

---

## Phase 7: User Story 5 - Configurable Game Parameters (Priority: P2)

**Goal**: Enable users to specify board size and number of players when starting a game

**Independent Test**: Call startGame() with various parameters and verify the board is created with the specified dimensions

### Implementation for User Story 5

- [ ] T031 [P] [US5] Update Game entity for configurable board size in `src/models/game.go` with board_rows and board_columns fields
- [ ] T032 [US5] Implement board size validation in `src/lib/validation/game.go` (min 5x5, max 100x100)
- [ ] T033 [US5] Update GameService to accept board size parameters in `src/services/game_service.go`
- [ ] T034 [US5] Update API request body for configurable parameters in `src/adapters/gin/handlers/game.go`

**Checkpoint**: At this point, User Stories 1-5 should all work independently

---

## Phase 8: User Story 6 - Play Again Flow (Priority: P2)

**Goal**: Enable users to start a new game without restarting the application after completing a game

**Independent Test**: Complete a game, verify the play-again prompt is shown, and verify new games start cleanly

### Implementation for User Story 6

- [ ] T035 [US6] Implement game reset logic in `src/services/game_service.go` for clearing previous game state
- [ ] T036 [US6] Add victory detection and play-again prompt in `src/services/game_service.go`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T037 [P] Documentation updates in `docs/api.md` for all endpoints
- [ ] T038 [P] Code cleanup and refactoring across all layers
- [ ] T039 [P] Performance optimization for ship placement algorithm in `src/models/ship.go`
- [ ] T040 [P] Additional unit tests in `tests/unit/` for edge cases (EC-001 through EC-013)
- [ ] T041 [P] Security hardening: input sanitization and access control review
- [ ] T042 Run quickstart.md validation scenarios 1-5
- [ ] T043 [P] Configure log rotation in `src/lib/logger/logger.go`
- [ ] T044 [P] Setup Prometheus metrics collection in `src/lib/metrics/metrics.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 9)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 board state
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 game state
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 game state
- **User Story 5 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 game state
- **User Story 6 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 victory detection

### Within Each User Story

- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all models for User Story 1 together:
Task: "Create Ship entity model in src/models/ship.go"
Task: "Create Board entity model in src/models/board.go"
Task: "Create Game entity model in src/models/game.go"
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
6. Add User Story 5 → Test independently → Deploy/Demo
7. Add User Story 6 → Test independently → Deploy/Demo
8. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (models + service + endpoint)
   - Developer B: User Story 2 (board display + endpoint)
   - Developer C: User Story 3 (multiplayer support)
   - Developer D: User Story 4 (statistics)
   - Developer E: User Story 5 (configurable parameters)
   - Developer F: User Story 6 (play again flow)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (if tests requested)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
