# Feature Specification: Battleship Game Engine

**Feature Branch**: `001-battleship-game-engine`

**Created**: 2026-06-19

**Status**: Draft

**Input**: User description: "Battleship Game Engine Tier: 3-Advanced Battleship Game Engine (BGE) implements the classic turn-based board game as a package thats separated from any presentation layer. This is a type of architectural pattern that useful in many application since it allows any number of apps to utilize the same service.

The BGE itself is invoked through a series of function calls rather than through directly coupled end user actions. In this respect using the BGE is is similar to using an API or a series of routes exposed by a web server.

This challenge requires that you develop the BGE and a very thin, text-based presentation layer for testing thats separate from the engine itself. Due to this the User Stories below are divided two sets - one for the BGE and one for the text-based presentation layer. BGE is responsible for maintaining game state. BGE follows a success/failure result pattern for all operations - functions return objects with `{ success: boolean, data?: ..., error?: ... }` where `error` contains descriptive error information when `success` is `false`. User Stories BGE Caller can invoke a startGame() function to begin a 1-player game. This function will generate an 8x8 game board consisting of 3 ships having a width of one square and a length of: Destroyer: 2 squares Cruiser: 3 squares Battleship: 4 squares startGame() will randomly place these ships on the board in any direction and will return an array representing ship placement. Caller can invoke a shoot() function passing the target row and column coordinates of the targeted cell on the game board. shoot() will return indicators representing if the shot resulted in a hit or miss, the number of ships left (i.e. not yet sunk), the ship placement array, and an updated hits and misses array. Cells in the hits and misses array will contain a space if they have yet to be targeted, O if they were targeted but no part of a ship was at that location, or X if the cell was occupied by part of a ship. Text-based Presentation Layer User can see the hits and misses array displayed as a 2 dimensional character representation of the game board returned by the startGame() function. User can be prompted to enter the coordinates of a target square on the game board. User can see an updated hits and misses array display after taking a shot. User can see a message after each shot indicating whether the shot resulted in a hit or miss. User can see an congratulations message after the shot that sinks the last remaining ship. User can be prompted to play again at the end of each game. Declining to play again stops the game. Bonus features BGE Caller can specify the number of rows and columns in the game board as a parameter to the startGame() function. Caller can invoke a gameStats() function that returns a Javascript object containing metrics for the current game. For example, number of turns played, current number of hits and misses, etc. Caller can specify the number of players (1 or 2) when calling startGame() which will generate one board for each player randomly populated with ships. shoot() will accept the player number the shot is being made for along with the coordinates of the shot. The data it returns will be for that player. Text-based Presentation Layer User can see the current game statics at any point by entering the phrase stats in place of target coordinates. (Note that this requires the gameStats() function in the BGE) User can specify a two player game is to be played, with each player alternating turns in the same terminal session (Note that this requires corresponding Bonus Features in the BGE) User can see the player number in prompts associated with the inputs in each turn. User can see both players boards at the end of each turn. Useful links and resources Battleship Game (Wikipedia) Battleship Game Rules (Hashas)"

## Clarifications

### Session 2026-06-19

- Q: For the `startGame()` and `shoot()` functions, how should errors and success results be returned? → A: Return success/failure objects with `result: { success: boolean, data?: ..., error?: ... }` pattern
- Q: For two-player games, how should turn alternation and board visibility work? → A: Engine tracks turn order internally; each `shoot()` call specifies which player is shooting; engine returns only that player's board state
- Q: What should the exact return structure of `startGame()` and `shoot()` functions look like? → A: `startGame()` returns `{ success: boolean, data: { ships: Ship[], board: string[][] }, error?: string }`; `shoot()` returns `{ success: boolean, data: { hit: boolean, shipSunk?: string, shipsRemaining: number, board: string[][] }, error?: string }`
- Q: For board boundaries and coordinate validation, what should be the exact range? → A: Zero-based indexing (0 to rows-1, 0 to columns-1) - consistent with array indexing
- Q: For duplicate shots at already-targeted cells, how should the system respond? → A: Return failure with `success: false` and descriptive error - reject duplicate shots
- Q: For the `startGame()` function's return value, what exactly should the `ships: Ship[]` array contain? → A: Each Ship object contains: `{ id: string, type: string, length: number, positions: [{row, col}[]], hits: number, sunk: boolean }` - complete ship state with unique ID for tracking
- Q: How should the engine handle out-of-turn shots in multiplayer games? → A: Reject out-of-turn shots with a clear error message - the caller is responsible for turn management

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Single Player Battleship Game (Priority: P1)

A single player wants to start a new Battleship game, take turns shooting at the computer's ships, and see visual feedback about their shots until all enemy ships are sunk.

**Why this priority**: This is the core minimum viable product - without this basic flow, there is no game. This represents the fundamental value proposition of the Battleship Game Engine.

**Independent Test**: Can be fully tested by calling startGame() with default parameters, then repeatedly calling shoot() with valid coordinates until all ships are sunk, verifying all required data is returned correctly.

**Acceptance Scenarios**:

1. **Given** no active game exists, **When** startGame() is called with default 8x8 board, **Then** engine creates a new game with 3 ships (Destroyer: 2 squares, Cruiser: 3 squares, Battleship: 4 squares) randomly placed on the board and returns ship placement data

2. **Given** an active game exists, **When** shoot() is called with valid row and column coordinates, **Then** engine returns shot result (hit/miss), number of ships remaining, updated ship placement, and updated hits/misses array

3. **Given** a shot results in all parts of a ship being hit, **When** the ship is sunk, **Then** engine indicates the ship was sunk and decrements the ships remaining count

4. **Given** all ships are sunk, **When** the final shot is processed, **Then** engine indicates victory and no ships remain

5. **Given** a shot is fired at a cell that was previously targeted, **When** shoot() is called with duplicate coordinates, **Then** engine returns an error indicating the cell was already targeted

---

### User Story 2 - Display Game Board State (Priority: P1)

Users need to see the current state of the game board showing which cells have been targeted and whether those shots were hits or misses.

**Why this priority**: Without visual feedback, players cannot make informed decisions about where to shoot next. This is essential for usability.

**Independent Test**: Can be tested independently by verifying the hits/misses array is properly formatted and displays correctly as a 2D grid showing O for misses and X for hits.

**Acceptance Scenarios**:

1. **Given** a game exists with no shots fired, **When** the board is displayed, **Then** all cells show as untargeted (space character)

2. **Given** shots have been fired at various cells, **When** the board is displayed, **Then** targeted cells show O for misses and X for hits

3. **Given** multiple shots have been fired, **When** the board is displayed, **Then** all shot results are visible in their correct positions on the grid

---

### User Story 3 - Two Player Alternate Turns (Priority: P2)

Two players want to take turns playing Battleship on the same terminal, with each player seeing their own board and the opponent's targeted cells.

**Why this priority**: This is a popular variation that enables peer-to-peer gameplay without requiring separate devices. Important for extending the engine's utility.

**Independent Test**: Can be tested independently by starting a 2-player game, verifying turn alternation works correctly, and ensuring each player sees only their own board state.

**Acceptance Scenarios**:

1. **Given** startGame() is called with 2 players, **When** the game initializes, **Then** two independent boards are created, each with randomly placed ships

2. **Given** it is Player 1's turn, **When** shoot() is called with player number 1 and target coordinates, **Then** engine processes the shot for Player 1's board and returns results for that player only

3. **Given** Player 1 has fired, **When** turn alternates to Player 2, **Then** engine expects next shoot() call to include player number 2

4. **Given** both players are active, **When** turn is alternated, **Then** each player can see only their own board state

5. **Given** a player has sunk all opponent ships, **When** victory is detected, **Then** engine declares that player the winner

---

### User Story 4 - Game Statistics Display (Priority: P2)

Users want to view game statistics at any time during gameplay without interrupting the game flow.

**Why this priority**: This is a convenience feature that enhances user experience by providing visibility into game progress without requiring players to remember counts.

**Independent Test**: Can be tested independently by calling gameStats() at various points in the game and verifying all metrics are accurate.

**Acceptance Scenarios**:

1. **Given** a game is in progress, **When** gameStats() is called, **Then** engine returns number of turns played, total hits, total misses, and ships remaining

2. **Given** multiple players are active, **When** gameStats() is called for a specific player, **Then** engine returns stats for that player's board only

3. **Given** stats are requested at game start, **When** no shots have been fired, **Then** engine returns zero counts for all shot-related metrics

4. **Given** stats are requested after multiple turns, **When** shots have been fired, **Then** engine returns accurate counts matching actual game state

---

### User Story 5 - Configurable Game Parameters (Priority: P2)

Users want flexibility in game setup, including board size and number of players.

**Why this priority**: This feature allows the engine to support different difficulty levels and gameplay variations while maintaining a clean, consistent API.

**Independent Test**: Can be tested independently by calling startGame() with various parameters and verifying the board is created with the specified dimensions.

**Acceptance Scenarios**:

1. **Given** startGame() is called with custom row and column counts, **When** the game initializes, **Then** engine creates a board with the specified dimensions

2. **Given** startGame() is called with 1 player, **When** the game initializes, **Then** engine creates one board with randomly placed ships

3. **Given** startGame() is called with 2 players, **When** the game initializes, **Then** engine creates two independent boards, each with randomly placed ships

4. **Given** startGame() is called with default parameters, **When** no parameters are specified, **Then** engine creates an 8x8 board with default settings

---

### User Story 6 - Play Again Flow (Priority: P2)

After completing a game, users want the option to start a new game without restarting the application.

**Why this priority**: This is a standard expectation for turn-based games and improves user experience by providing seamless continuation.

**Independent Test**: Can be tested independently by completing a game, verifying the play-again prompt is shown, and verifying new games start cleanly.

**Acceptance Scenarios**:

1. **Given** all ships are sunk, **When** the victory message is displayed, **Then** engine prompts user to play again with accepted responses: "yes"/"no", "y"/"n" (case-insensitive)

2. **Given** play again prompt is shown, **When** user confirms ("yes" or "y"), **Then** engine starts a new game with fresh board state

3. **Given** play again prompt is shown, **When** user declines ("no" or "n"), **Then** engine ends the session and returns control to caller

4. **Given** a new game is started, **When** startGame() is called, **Then** all previous game state is cleared and fresh ships are placed

---

### Edge Cases

- **EC-001**: When a player tries to shoot at coordinates outside the board boundaries, the engine MUST reject the shot with a descriptive error

- **EC-002**: When invalid coordinate formats are provided (non-numeric input, negative numbers), the engine MUST validate inputs and return a descriptive error

- **EC-003**: When shoot() is called before startGame(), the engine MUST reject the shot with a descriptive error indicating no active game exists

- **EC-004**: When duplicate shots are fired at the same cell, the engine MUST reject with a descriptive error that includes whether the previous shot was a hit or miss

- **EC-005**: When all possible shots are fired but ships remain unsunk, the engine MUST continue accepting shots until all ships are sunk (no board-full condition)

- **EC-006**: When concurrent calls to startGame() occur, the engine MUST handle each game independently with unique game IDs

- **EC-007**: When gameStats() is called when no game is active, the engine MUST return a failure response with a descriptive error

- **EC-008**: When ship placement conflicts occur (ships overlapping or extending beyond board boundaries), the engine MUST retry up to 100 attempts then fail with a descriptive error indicating the specific conflict type

- **EC-009**: When a player shoots out of turn in multiplayer games, the engine MUST reject with a descriptive error indicating whose turn it is

- **EC-010**: When a board size smaller than 5×5 is specified, the engine MUST fail startGame() with a descriptive error

- **EC-011**: When a board size larger than 100×100 is specified without performance justification, the engine MUST fail startGame() with a descriptive error

- **EC-012**: When a ship length of 0 is specified, the engine MUST fail startGame() with a descriptive error

- **EC-013**: When a ship length exceeds board dimensions, the engine MUST fail startGame() with a descriptive error

## Requirements *(mandatory)*

### Functional Requirements

#### Core Game Engine

- **FR-001**: System MUST provide a startGame() function that initializes a new Battleship game with default 8x8 board when no parameters are provided

- **FR-002**: System MUST place exactly 3 ships on the board at game start: Destroyer (2 squares), Cruiser (3 squares), Battleship (4 squares)

- **FR-003**: System MUST randomly place ships in either horizontal or vertical orientation on the board

- **FR-004**: System MUST return ship placement data from startGame() containing all ship positions with complete state (id, type, length, positions, hits, sunk) in `data.ships` field of the response

- **FR-005**: System MUST provide a shoot() function that accepts row and column coordinates as parameters

- **FR-006**: System MUST return shot result (hit or miss) from shoot() for each target coordinate in `data.hit` field of the response

- **FR-007**: System MUST return number of ships remaining (not yet sunk) from shoot()

- **FR-008**: System MUST return updated ship placement data from shoot() in `data.ships` field of the response

- **FR-009**: System MUST return updated hits/misses array from shoot() in `data.board` field of the response

- **FR-010**: System MUST represent untargeted cells in hits/misses array as space characters

- **FR-011**: System MUST represent miss cells (no ship at location) in hits/misses array as "O"

- **FR-012**: System MUST represent hit cells (ship occupied) in hits/misses array as "X"

- **FR-013**: System MUST detect when all parts of a ship are hit and indicate the ship is sunk

- **FR-014**: System MUST decrement ships remaining count when a ship is sunk

- **FR-015**: System MUST declare victory when all ships are sunk

- **FR-016**: System MUST reject duplicate shots at previously targeted cells with appropriate error

- **FR-017**: System MUST reject shots at coordinates outside board boundaries with appropriate error

- **FR-018**: System MUST validate board minimum size is 5×5 for standard 3-ship configuration

- **FR-019**: System MUST validate board maximum size is 100×100 without performance justification

- **FR-020**: System MUST validate ship minimum length is 1 square

- **FR-021**: System MUST validate ship maximum length does not exceed board dimensions

- **FR-022**: System MUST use rejection sampling algorithm for ship placement with a maximum retry limit of 100 attempts

- **FR-023**: System MUST fail startGame() if ships cannot be placed within the retry limit

#### Game Statistics

- **FR-024**: System MUST provide a gameStats() function that returns game metrics

- **FR-025**: System MUST return number of turns played from gameStats()

- **FR-026**: System MUST return total hit count from gameStats()

- **FR-027**: System MUST return total miss count from gameStats()

- **FR-028**: System MUST return ships remaining count from gameStats()

#### Multiplayer Support

- **FR-029**: System MUST support startGame() with a players parameter (1 or 2) that validates input is an integer and rejects non-integer values

- **FR-030**: System MUST create independent board state for each player when players=2

- **FR-031**: System MUST accept player number as parameter in shoot() for multiplayer games

- **FR-032**: System MUST return data for specified player only when shoot() is called with player number

- **FR-033**: System MUST track turn order internally and increment turn counter on each valid shoot() call

- **FR-034**: System MUST declare winner when one player sinks all opponent ships

- **FR-035**: System MUST maintain independent turn counts per player in gameStats() for multiplayer games

- **FR-036**: System MUST reject shoot() calls that do not match the current player's turn with a descriptive error

#### Configurable Game Parameters

- **FR-037**: System MUST accept rows parameter in startGame() to configure board height

- **FR-038**: System MUST accept columns parameter in startGame() to configure board width

- **FR-039**: System MUST validate rows and columns are positive integers

- **FR-040**: System MUST validate ship placement fits within configured board dimensions

#### Game Flow Control

- **FR-041**: System MUST prompt user to play again after victory is declared

- **FR-042**: System MUST start new game when user confirms play again

- **FR-043**: System MUST end session when user declines play again

- **FR-044**: System MUST clear all previous game state when startGame() is called

#### Input Validation

- **FR-045**: System MUST validate row and column are valid integers before processing

- **FR-046**: System MUST validate row is within board boundaries (0 to rows-1)

- **FR-047**: System MUST validate column is within board boundaries (0 to columns-1)

- **FR-048**: System MUST reject calls to shoot() before startGame() is called

#### Error Handling

- **FR-049**: System MUST return all operations using success/failure result pattern: `{ success: boolean, data?: ..., error?: ... }` where `error` contains descriptive error information when `success` is `false`. Error responses MUST include an `error_code` field with standardized error codes (e.g., `INVALID_BOARD_SIZE`, `INVALID_COORDINATES`, `GAME_NOT_FOUND`). Error messages MUST include operation name, invalid parameter name, and suggested recovery action. Error messages MUST NOT include stack traces, file paths, or internal state information.

- **FR-050**: System MUST maintain valid game state after failed operations (no partial updates)

#### Error Response Format

- **Success Response**: `{ success: true, data: { ... }, error: null, error_code: null }`
- **Failure Response**: `{ success: false, data: null, error: string, error_code: string }`

### Key Entities

- **Game**: Represents an active Battleship game session containing board state, ship positions, hit/miss tracking, turn information, and player data

- **Board**: An N×M grid representing the game area where ships are placed and shots are fired; each cell tracks target status (space for untargeted, "O" for miss, "X" for hit)

- **Ship**: A vessel with properties `{ id: string, type: string, length: number, positions: {row,col}[], hits: number, sunk: boolean }` placed horizontally or vertically on the board; uniquely identified by ID for tracking across game operations

- **HitMissCell**: A cell in the hit/miss tracking grid containing space (untargeted), "O" (miss - no ship), or "X" (hit - ship occupied)

- **Player**: A game participant with independent board state and ship positions; tracks turn order and victory status

- **GameStats**: A data structure containing metrics about current game state including turns played, hits, misses, and ships remaining

- **TurnState**: Tracks which player's turn it is; engine enforces turn order and rejects shots from the wrong player

### API Response Structures

- **startGame() success**: `{ success: true, data: { ships: Ship[], board: string[][] }, error: null, error_code: null }`

- **startGame() failure**: `{ success: false, data: null, error: string, error_code: string }`

- **shoot() success**: `{ success: true, data: { hit: boolean, shipSunk?: string, shipsRemaining: number, board: string[][] }, error: null, error_code: null }`

- **shoot() failure**: `{ success: false, data: null, error: string, error_code: string }`

- **gameStats() success**: `{ success: true, data: GameStats, error: null, error_code: null }`

- **gameStats() failure**: `{ success: false, data: null, error: string, error_code: string }`

### Ship Placement Rules

- **SPR-001**: Engine MUST use rejection sampling algorithm for ship placement with a maximum retry limit of 100 attempts

- **SPR-002**: Engine MUST fail `startGame()` if ships cannot be placed within the retry limit

- **SPR-003**: Engine MUST ensure ships do not overlap or extend beyond board boundaries

- **SPR-004**: Engine MUST validate that board minimum size is 5×5 for standard 3-ship configuration

- **SPR-005**: Engine MUST validate that board maximum size is 100×100 without performance justification

- **SPR-006**: Engine MUST validate that ship minimum length is 1 square

- **SPR-007**: Engine MUST validate that ship maximum length does not exceed board dimensions

- **SPR-008**: Engine MUST assign unique IDs to each Ship object for tracking across game operations

- **SPR-009**: Engine MUST track hit count and sunk status for each Ship object

- **SPR-010**: Engine MUST return complete Ship state (including positions, hits, and sunk status) in `startGame()` response

### Turn Management

- **TM-001**: Engine MUST track which player's turn it is for multiplayer games

- **TM-002**: Engine MUST reject `shoot()` calls that do not match the current player's turn with a descriptive error

- **TM-003**: Engine MUST automatically advance to the next player after a valid `shoot()` call completes

- **TM-004**: Engine MUST maintain independent turn counts per player for statistics

- **TM-005**: Engine MUST declare a winner when one player sinks all opponent ships

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A single player can complete a full game (start to victory) in under 5 minutes average time

- **SC-002**: System responds to all API calls (startGame, shoot, gameStats) within 100ms for standard board sizes

- **SC-003**: System supports at least 1000 sequential games per minute through API calls

- **SC-004**: All valid inputs are processed correctly with 100% accuracy for hit/miss detection and ship sinking

- **SC-005**: Game board displays correctly for all supported board sizes (minimum 5×5, maximum 50×50)

- **SC-006**: Two-player games alternate turns reliably with no state corruption between players

- **SC-007**: Game can restart at least 50 times per session without errors or state leakage

## Assumptions

- Default board size is 8×8 unless specified otherwise

- Default number of players is 1 unless specified otherwise

- Ships are placed randomly at game start with no manual placement option

- Game uses zero-based indexing for rows and columns (0 to N-1)

- Default ship types are Destroyer (2 squares), Cruiser (3 squares), Battleship (4 squares)

- Each ship occupies exactly one row or column (no diagonal placement)

- Ships cannot overlap or extend beyond board boundaries

- Game session is sequential (no concurrent turns)

- Text-based presentation layer handles all user interaction (input/output)

- No persistence of game state between sessions (each startGame() is independent)

- No time limits on turns or overall game duration

- No AI opponent difficulty settings (AI uses default strategy)

- No save/load functionality for ongoing games

- No network multiplayer support (local only)

- No undo functionality for shots

- Ship placement algorithm uses rejection sampling with a maximum retry limit of 100 attempts

- Board minimum size is 5×5; boards smaller than this will fail to place all ships

- Board maximum size is 100×100; larger boards require performance justification

- Ship minimum length is 1 square; ships cannot be 0 length

- Ship maximum length is min(board_rows, board_columns) - 1; ships cannot exceed board dimensions

- No hint or suggestion system for players

## Out of Scope

- AI opponent with advanced strategy or difficulty settings

- Network multiplayer support (local terminal only)

- Save/load game state to disk

- Undo functionality for shots

- Hint or suggestion system for players

- Time limits on turns

- Multiple ship types beyond the default 3 (no carriers, submarines, etc.)

- Custom ship placement (all placement is random)

- Sound effects or visual enhancements (text-only)

- Game history or replay functionality

- Statistics tracking across multiple games (per-session only)

- Tutorial or help system

- Multi-word coordinate inputs (e.g., "A1" - only numeric row/column)

- Game customization (board size, ship types are runtime parameters only)
