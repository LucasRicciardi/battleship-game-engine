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

## Phase 1: Setup (Project Initialization)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project structure per implementation plan in `src/`, `tests/`, `docker-compose.yml`, `Dockerfile`
- [X] T002 Initialize Go project with `go mod init battleship-game-engine` and add dependencies (gin, gorm, postgres, zap, otel, prometheus, jwt, limiter, validator)
- [X] T003 [P] Configure linting (golangci-lint) and formatting (gofmt) tools in `.golangci.yml` and `Makefile`
- [X] T004 [P] Create `Makefile` with common commands (build, test, lint, run, migrate)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T005 Setup database schema and migrations framework in `src/adapters/db/migrate/`
- [X] T006 [P] Implement logging infrastructure in `lib/logger/logger.go` using zap with layer-based log levels
- [X] T007 [P] Implement tracing infrastructure in `lib/tracing/tracing.go` using OpenTelemetry SDK
- [X] T008 [P] Implement metrics infrastructure in `lib/metrics/metrics.go` using Prometheus client library
- [X] T009 [P] Create correlation ID middleware in `src/adapters/gin/middleware/correlationid.go`
- [X] T010 [P] Create health check endpoints in `src/adapters/gin/handlers/health.go` (`/health/live`, `/health/ready`)
- [X] T011 [P] Create security headers middleware in `src/adapters/gin/middleware/securityheaders.go`
- [X] T012 [P] Create rate limiting middleware in `src/adapters/gin/middleware/ratelimit.go`
- [X] T013 [P] Create error handler middleware in `src/adapters/gin/middleware/errorhandler.go`
- [X] T014 [P] Create JWT authentication middleware in `src/adapters/gin/middleware/jwt.go`
- [X] T015 [P] Create player authorization middleware in `src/adapters/gin/middleware/authorization.go`
- [X] T016 Create base models/entities in `src/models/` (Game, Board, Ship, Player, GameStats - database-independent)
- [X] T017 Create database models in `src/adapters/db/models.go` (GameDB, BoardDB, ShipDB, PlayerDB with GORM tags)
- [X] T018 Create mapper functions in `src/adapters/db/mappers.go` (ToGame, ToGameDB, etc.)
- [X] T019 Create repository interface in `src/models/game_repository.go` and implementation in `src/adapters/db/game_repository.go`
- [X] T020 Configure environment configuration management in `config/config.go` and `config/.env.example`
- [X] T021 Setup API routing structure in `src/adapters/gin/router.go` with middleware chain and route groups
- [X] T022 Create validation package in `src/validation/api.go` with request structs and struct tags

---

## Phase 3: User Story 1 - Single Player Battleship Game (Priority: P1) 🎯 MVP

**Goal**: Implement core game functionality for single-player Battleship with startGame() and shoot() functions

**Independent Test**: Call startGame() with default 8x8 board, repeatedly call shoot() with valid coordinates until all ships are sunk, verifying all required data is returned correctly

### Implementation for User Story 1

- [X] T023 [P] [US1] Create Ship entity in `src/models/ship.go` with id, type, length, positions, hits, sunk fields
- [X] T024 [P] [US1] Create Board entity in `src/models/board.go` with cells tracking (space=untargeted, O=miss, X=hit)
- [X] T025 [P] [US1] Create Game entity in `src/models/game.go` with board state, ships, turn tracking, status
- [X] T026 [US1] Implement Ship placement logic in `src/services/ship_placer.go` using rejection sampling (100 attempts max)
- [X] T027 [US1] Implement Game service in `src/services/game_service.go` with startGame() and shoot() functions
- [X] T028 [US1] Implement Game controller in `src/adapters/gin/controllers/game_controller.go` with HTTP handlers
- [X] T029 [US1] Add API endpoint for startGame in `src/adapters/gin/routes.go` (POST /api/v1/games)
- [X] T030 [US1] Add API endpoint for shoot in `src/adapters/gin/routes.go` (POST /api/v1/games/:game_id/shoot)
- [X] T031 [US1] Add input validation for board size (5x5 to 100x100) in `src/validation/api.go`
- [X] T032 [US1] Add input validation for coordinates (0 to rows-1, 0 to columns-1) in `src/validation/api.go`
- [X] T033 [US1] Add duplicate shot detection in `src/services/game_service.go`
- [X] T034 [US1] Add boundary validation for shots in `src/services/game_service.go`
- [X] T035 [US1] Add logging for game operations in `src/services/game_service.go` (INFO/WARN levels)
- [X] T036 [US1] Add metrics collection for game starts and shots in `lib/metrics/metrics.go`

---

## Phase 4: User Story 2 - Display Game Board State (Priority: P1)

**Goal**: Implement board state display functionality showing hits and misses

**Independent Test**: Verify hits/misses array is properly formatted and displays correctly as 2D grid showing O for misses and X for hits

### Implementation for User Story 2

- [X] T037 [US2] Implement board state retrieval in `src/services/game_service.go` (getGameState function)
- [X] T038 [US2] Add API endpoint for getGameState in `src/adapters/gin/routes.go` (GET /api/v1/games/:game_id)
- [X] T039 [US2] Add board visualization helper in `src/adapters/gin/helpers/board_visualizer.go`
- [X] T040 [US2] Add logging for board state retrieval in `src/services/game_service.go`

---

## Phase 5: User Story 3 - Two Player Alternate Turns (Priority: P2)

**Goal**: Implement multiplayer support with turn alternation

**Independent Test**: Start 2-player game, verify turn alternation works correctly, ensure each player sees only their own board state

### Implementation for User Story 3

- [X] T041 [US3] Update Game entity to support multiple players in `src/models/game.go` (numPlayers, currentPlayer fields)
- [X] T042 [US3] Update Ship entity to support per-player ships in `src/models/ship.go`
- [X] T043 [US3] Update Board entity to support per-player boards in `src/models/board.go`
- [X] T044 [US3] Implement turn alternation logic in `src/services/game_service.go`
- [X] T045 [US3] Add out-of-turn shot rejection in `src/services/game_service.go`
- [X] T046 [US3] Add player authorization check in `src/adapters/gin/middleware/authorization.go`
- [X] T047 [US3] Add API endpoint for shoot with player_id in `src/adapters/gin/routes.go`
- [X] T048 [US3] Add logging for multiplayer operations in `src/services/game_service.go`

---

## Phase 6: User Story 4 - Game Statistics Display (Priority: P2)

**Goal**: Implement gameStats() function for displaying game metrics

**Independent Test**: Call gameStats() at various points in the game and verify all metrics are accurate

### Implementation for User Story 4

- [X] T049 [US4] Create GameStats entity in `src/models/game_stats.go` with turns, hits, misses, shipsRemaining
- [X] T050 [US4] Implement gameStats() function in `src/services/game_service.go`
- [X] T051 [US4] Add API endpoint for gameStats in `src/adapters/gin/routes.go` (GET /api/v1/games/:game_id/stats)
- [X] T052 [US4] Add logging for stats retrieval in `src/services/game_service.go`

---

## Phase 7: User Story 5 - Configurable Game Parameters (Priority: P2)

**Goal**: Implement configurable board size and player count

**Independent Test**: Call startGame() with various parameters and verify board is created with specified dimensions

### Implementation for User Story 5

- [X] T053 [US5] Update startGame() to accept board_rows and board_columns parameters in `src/services/game_service.go`
- [X] T054 [US5] Update startGame() to accept num_players parameter in `src/services/game_service.go`
- [X] T055 [US5] Add board size validation (min 5x5, max 100x100) in `src/validation/api.go`
- [X] T056 [US5] Add player count validation (1 or 2) in `src/validation/api.go`
- [X] T057 [US5] Add logging for configurable game parameters in `src/services/game_service.go`

---

## Phase 8: User Story 6 - Play Again Flow (Priority: P2)

**Goal**: Implement play-again functionality after game completion

**Independent Test**: Complete a game, verify play-again prompt is shown, verify new games start cleanly

### Implementation for User Story 6

- [X] T058 [US6] Add game completion detection in `src/services/game_service.go`
- [X] T059 [US6] Add victory message generation in `src/services/game_service.go`
- [X] T060 [US6] Add play-again prompt support in `src/adapters/gin/controllers/game_controller.go`
- [X] T061 [US6] Add new game state clearing in `src/services/game_service.go`

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T062 [P] Create comprehensive documentation in `docs/` (API documentation, architecture overview)
- [X] T063 Code cleanup and refactoring across all layers
- [X] T064 Performance optimization across all services (ensure <100ms p95 for typical operations)
- [X] T065 [P] Additional unit tests in `tests/unit/` (target 95%+ for core logic)
- [X] T066 [P] Additional integration tests in `tests/integration/` (target 80%+ for adapters)
- [X] T067 [P] Additional contract tests in `tests/contract/` for all API endpoints
- [X] T068 Security hardening (input validation, access control review)
- [X] T069 Run quickstart.md validation scenarios
- [X] T070 Performance testing (1000 sequential games/minute, <100ms response time)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed) or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 9)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 board state
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 core game logic
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 game state tracking
- **User Story 5 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 startGame implementation
- **User Story 6 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 game completion detection

### Within Each User Story

- Models before services
- Services before controllers/endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: Foundational Phase

```bash
# Launch all middleware tasks in parallel:
Task: "Create correlation ID middleware in src/adapters/gin/middleware/correlationid.go"
Task: "Create health check endpoints in src/adapters/gin/handlers/health.go"
Task: "Create security headers middleware in src/adapters/gin/middleware/securityheaders.go"
Task: "Create rate limiting middleware in src/adapters/gin/middleware/ratelimit.go"
Task: "Create error handler middleware in src/adapters/gin/middleware/errorhandler.go"
Task: "Create JWT authentication middleware in src/adapters/gin/middleware/jwt.go"
Task: "Create player authorization middleware in src/adapters/gin/middleware/authorization.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Single Player Battleship)
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
   - Developer A: User Story 1 (Single Player)
   - Developer B: User Story 2 (Board Display)
   - Developer C: User Story 3 (Two Player)
   - Developer D: User Story 4 (Game Stats)
   - Developer E: User Story 5 (Configurable Parameters)
   - Developer F: User Story 6 (Play Again Flow)
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
