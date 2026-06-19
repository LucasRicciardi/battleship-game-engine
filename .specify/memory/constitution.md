<!--
  Sync Impact Report:
  - Version change: 2.2.0 → 2.3.0
  - Modified principles:
    - I. Clean Architecture (new - supersedes Code Quality as primary principle)
    - II. Code Quality (renumbered - retained clean code principles)
  - Added sections:
    - Architecture Layers (Entities, Use Cases, Interface Adapters, Frameworks & Drivers)
    - The Dependency Rule
    - Crossing Boundaries
    - Data Across Boundaries
  - Removed sections: N/A
  - Templates requiring updates:
    - ⚠ pending: .specify/templates/plan-template.md (Constitution Check section)
    - ⚠ pending: .specify/templates/spec-template.md (requirement alignment)
    - ⚠ pending: .specify/templates/tasks-template.md (task categorization)
  - Follow-up TODOs:
    - TODO(PROJECT_NAME): Set to "Battleship Game Engine" when project naming is finalized
    - TODO(RATIFICATION_DATE): Set to 2026-06-18 (project initialization date)
-->

# Battleship Game Engine Constitution

## Core Principles

### I. Clean Architecture

This project follows Uncle Bob's Clean Architecture principles. The architecture is divided into concentric circles with dependencies flowing inward only (The Dependency Rule).

#### Architecture Layers

##### Entities (Innermost - Enterprise Business Rules)

Entities encapsulate enterprise-wide business rules that are independent of any particular application. These are the most stable and least likely to change.

- Entities MUST encapsulate fundamental business rules (e.g., ship placement logic, turn management, win conditions)
- Entities MUST NOT depend on frameworks, UI, or database concerns
- Entities SHOULD be reusable across multiple applications in the enterprise

##### Use Cases (Application Business Rules)

Use cases contain application-specific business rules that orchestrate data flow between Entities and external layers.

- Use cases MUST orchestrate data flow to/from Entities
- Use cases MUST implement system use cases (e.g., "place ship", "fire shot", "get game state")
- Use cases MUST NOT depend on UI, database, or framework details
- Changes to use cases SHOULD NOT affect Entities

##### Interface Adapters (Gateway Layer)

Interface Adapters convert data between formats convenient for Use Cases/Entities and formats required by external agencies (databases, web frameworks, APIs).

- All framework-specific code (web frameworks, database libraries) MUST reside in this layer
- All database queries, SQL, and persistence logic MUST be restricted to this layer
- All UI components (controllers, views, presenters) MUST reside in this layer
- Adapters MUST convert data to/from formats convenient for inner layers
- No inner layer code SHOULD know about database or UI implementation details

##### Frameworks & Drivers (Outermost - External Agencies)

The outermost layer contains frameworks, tools, and implementation details.

- This layer is where all concrete details reside (web framework, database, APIs)
- This layer SHOULD contain only glue code connecting to inner layers
- Outer layer SHOULD NOT impact inner layers

#### The Dependency Rule

Source code dependencies can ONLY point INWARD. Nothing in an inner circle can know anything about something in an outer circle.

- Names declared in outer circles (functions, classes, variables) MUST NOT be mentioned by inner circle code
- Data formats from outer circles MUST NOT be used by inner circles
- This rule applies to all boundaries between layers

#### Crossing Boundaries

Boundary crossings use the Dependency Inversion Principle with dynamic polymorphism:

- Use cases call interfaces (Output Ports) defined in inner circles
- Outer layers implement these interfaces
- This creates opposing source code dependencies while maintaining control flow

#### Data Across Boundaries

Simple data structures are passed across boundaries:

- Use basic structs, data transfer objects, or function arguments
- Never pass Entities or database rows across boundaries
- Never pass framework-specific data structures inward
- Data is converted to formats most convenient for the receiving layer

---

### II. Code Quality

All game engine code MUST adhere to clean code principles and the following quality standards:

#### Clean Code Principles (MUST follow all):

- **Meaningful Naming**: All identifiers (variables, functions, classes, modules) MUST use clear, intention-revealing names that follow language conventions; abbreviations are prohibited unless universally understood (e.g., "ID", "URL")
- **Small Functions**: Functions MUST be kept small (typically <20 lines); functions exceeding 50 lines MUST be refactored into smaller, well-named functions
- **Single Responsibility**: Each module, class, and function MUST have exactly one reason to change; modules that serve multiple purposes MUST be split
- **DRY Principle**: Duplicate code is prohibited; all repeated logic MUST be extracted into reusable functions/components with clear boundaries
- **Boy Scout Rule**: Developers MUST leave the codebase cleaner than they found it; no new technical debt allowed without explicit justification
- **Testability**: All code MUST be designed for testing; dependencies MUST be injectable, and modules MUST be unit-testable without external services
- **Clear Intent**: Code MUST communicate its purpose through structure and naming; complex logic MUST include inline comments explaining "why" not "what"

#### Architecture & Structure:

- **Module Cohesion**: Each module MUST have a single, well-defined responsibility; modules that grow beyond 500 lines MUST be reviewed for refactoring
- **Documentation**: All public APIs MUST include documentation describing parameters, return values, and side effects
- **Error Handling**: All critical operations MUST include explicit error handling with descriptive error messages

**Rationale**: Clean code is non-negotiable for this project. A battleship game engine requires predictable, maintainable code to support complex game states, AI opponents, and multiplayer coordination. Clean code principles ensure the codebase remains readable, testable, and extendable as complexity grows. Following these principles reduces cognitive load, minimizes bugs, and enables rapid iteration without fear of breaking existing functionality.

### III. Testing Standards

Testing is MANDATORY for all game engine functionality. Three layers of testing are required:

- **Unit Tests**: All core logic (game state, ship placement, turn management) MUST have unit tests
- **Integration Tests**: API endpoints, data persistence, and inter-service communication MUST have integration tests verifying end-to-end flows
- **E2E Tests**: All public API endpoints MUST have E2E tests verifying complete request/response cycles including error scenarios

**Test-First Development**: New features MUST have tests written before implementation (TDD); tests MUST fail before implementation begins.

**CI Integration**: All tests (unit, integration, E2E) MUST run automatically on pull requests; no merges allowed with failing tests.

**Rationale**: Game engines are foundational infrastructure - bugs in core logic can cascade through all UI consumers. Comprehensive testing at all three layers catches issues early and enables confident refactoring. E2E tests specifically verify that all API endpoints behave correctly in realistic scenarios.

### IV. Test Coverage

Test coverage is MANDATORY and MUST meet the following standards:

- **Coverage Targets**: All modules MUST achieve ≥90% code coverage
- **Coverage Types**: Coverage MUST include unit tests, integration tests, and E2E tests
- **Coverage Verification**: Coverage reports MUST be generated and reviewed on all pull requests
- **Coverage Maintenance**: New code MUST maintain or improve overall coverage; significant drops require explicit justification
- **Coverage Reporting**: Coverage metrics MUST be publicly visible in the repository

**Rationale**: High test coverage is not just about numbers - it ensures that the engine's logic is well-tested and that critical paths are protected from regressions. Coverage targets provide measurable quality gates while maintaining flexibility for genuine edge cases.

### V. User Experience Consistency

While this engine does not implement UI directly, all public APIs MUST provide a consistent experience for UI consumers:

- **Deterministic Behavior**: Game mechanics MUST be deterministic; identical inputs MUST produce identical outputs (critical for multiplayer sync and replay)
- **API Response Time**: All API endpoints MUST respond within 200ms for typical operations (board queries, state updates)
- **Error Messages**: All API errors MUST include structured error codes, human-readable messages, and suggested recovery actions
- **Consistent Data Shapes**: All API responses MUST follow consistent data models and naming conventions
- **Cross-Platform**: Engine MUST work consistently for all UI consumers (web, mobile, desktop) without platform-specific quirks

**Rationale**: This engine is foundational infrastructure for UI applications. Consistent, predictable APIs enable UI developers to build reliable experiences without workarounds.

### VI. Performance Requirements

Performance is critical for responsive API responses:

- **API Response Time**: All endpoints MUST respond within 200ms for typical operations (board queries, state updates)
- **Memory Limits**: Game state MUST fit within 50MB RAM for typical matches (2 players, standard board)
- **Load Time**: Initial engine module load MUST complete within 1 second
- **Scalability**: Engine MUST support up to 4 players per match without performance degradation
- **Profiling**: Performance-critical paths MUST include timing instrumentation for ongoing monitoring

**Rationale**: This engine powers UI applications - slow API responses directly impact end-user experience. Responsive APIs enable snappy, interactive game interfaces.

## Governance

This constitution supersedes all other development practices for the Battleship Game Engine project.

**Compliance**:
- All pull requests MUST include a Constitution Compliance section
- Violations MUST be justified with explicit trade-off documentation
- Quarterly reviews MUST assess adherence and propose improvements

**Versioning Policy**:
- MAJOR: Principle removals, fundamental paradigm shifts
- MINOR: New principles, material expansion of existing principles
- PATCH: Clarifications, wording improvements, typo fixes

**Version**: 2.3.0 | **Ratified**: 2026-06-18 | **Last Amended**: 2026-06-19
