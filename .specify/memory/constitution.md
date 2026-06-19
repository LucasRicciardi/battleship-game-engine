<!--
  Sync Impact Report:
  - Version change: 2.1.1 → 2.2.0
  - Modified principles:
    - I. Code Quality (expanded clean code principles from generic reference to specific principles)
  - Added sections:
    - Clean Code Principles section with 7 specific principles
    - Architecture & Structure section (reorganized existing content)
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

### I. Code Quality

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

- **Clean Architecture**: Engine modules MUST be organized into clear layers (core, services, rendering, input) with dependencies flowing inward
- **Module Cohesion**: Each module MUST have a single, well-defined responsibility; modules that grow beyond 500 lines MUST be reviewed for refactoring
- **Documentation**: All public APIs MUST include documentation describing parameters, return values, and side effects
- **Error Handling**: All critical operations MUST include explicit error handling with descriptive error messages

**Rationale**: Clean code is non-negotiable for this project. A battleship game engine requires predictable, maintainable code to support complex game states, AI opponents, and multiplayer coordination. Clean code principles ensure the codebase remains readable, testable, and extendable as complexity grows. Following these principles reduces cognitive load, minimizes bugs, and enables rapid iteration without fear of breaking existing functionality.

### II. Testing Standards

Testing is MANDATORY for all game engine functionality. Three layers of testing are required:

- **Unit Tests**: All core logic (game state, ship placement, turn management) MUST have unit tests
- **Integration Tests**: API endpoints, data persistence, and inter-service communication MUST have integration tests verifying end-to-end flows
- **E2E Tests**: All public API endpoints MUST have E2E tests verifying complete request/response cycles including error scenarios

**Test-First Development**: New features MUST have tests written before implementation (TDD); tests MUST fail before implementation begins.

**CI Integration**: All tests (unit, integration, E2E) MUST run automatically on pull requests; no merges allowed with failing tests.

**Rationale**: Game engines are foundational infrastructure - bugs in core logic can cascade through all UI consumers. Comprehensive testing at all three layers catches issues early and enables confident refactoring. E2E tests specifically verify that all API endpoints behave correctly in realistic scenarios.

### III. Test Coverage

Test coverage is MANDATORY and MUST meet the following standards:

- **Coverage Targets**: All modules MUST achieve ≥90% code coverage
- **Coverage Types**: Coverage MUST include unit tests, integration tests, and E2E tests
- **Coverage Verification**: Coverage reports MUST be generated and reviewed on all pull requests
- **Coverage Maintenance**: New code MUST maintain or improve overall coverage; significant drops require explicit justification
- **Coverage Reporting**: Coverage metrics MUST be publicly visible in the repository

**Rationale**: High test coverage is not just about numbers - it ensures that the engine's logic is well-tested and that critical paths are protected from regressions. Coverage targets provide measurable quality gates while maintaining flexibility for genuine edge cases.

### IV. User Experience Consistency

While this engine does not implement UI directly, all public APIs MUST provide a consistent experience for UI consumers:

- **Deterministic Behavior**: Game mechanics MUST be deterministic; identical inputs MUST produce identical outputs (critical for multiplayer sync and replay)
- **API Response Time**: All API endpoints MUST respond within 200ms for typical operations (board queries, state updates)
- **Error Messages**: All API errors MUST include structured error codes, human-readable messages, and suggested recovery actions
- **Consistent Data Shapes**: All API responses MUST follow consistent data models and naming conventions
- **Cross-Platform**: Engine MUST work consistently for all UI consumers (web, mobile, desktop) without platform-specific quirks

**Rationale**: This engine is foundational infrastructure for UI applications. Consistent, predictable APIs enable UI developers to build reliable experiences without workarounds.

### V. Performance Requirements

Performance is critical for responsive API responses:

- **API Response Time**: All endpoints MUST respond within 200ms for typical operations (board queries, state updates)
- **Memory Limits**: Game state MUST fit within 50MB RAM for typical matches (2 players, standard board)
- **Load Time**: Initial engine module load MUST complete within 1 second
- **Scalability**: Engine MUST support up to 4 players per match without performance degradation
- **Profiling**: Performance-critical paths MUST include timing instrumentation for ongoing monitoring

**Rationale**: This engine powers UI applications - slow API responses directly impact end-user experience. Responsive APIs enable snappy, interactive game interfaces.

## Governance

This constitution supersedes all other development practices for the Battleship Game Engine project.

**Amendment Process**:
1. Propose amendment with rationale and impact analysis
2. Review effects on all dependent templates and documentation
3. Update version number per semantic versioning (MAJOR.MINOR.PATCH)
4. Document all changes in this file's history

**Compliance**:
- All pull requests MUST include a Constitution Compliance section
- Violations MUST be justified with explicit trade-off documentation
- Quarterly reviews MUST assess adherence and propose improvements

**Versioning Policy**:
- MAJOR: Principle removals, fundamental paradigm shifts
- MINOR: New principles, material expansion of existing principles
- PATCH: Clarifications, wording improvements, typo fixes

**Version**: 2.2.0 | **Ratified**: 2026-06-18 | **Last Amended**: 2026-06-19
