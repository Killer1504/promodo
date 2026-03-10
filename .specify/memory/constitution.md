<!--
  SYNC IMPACT REPORT
  Version change: 0.0.0 → 1.0.0
  Bump rationale: Initial constitution creation (MAJOR - first governance adoption)
  
  Modified principles: N/A (initial creation)
  Added sections:
    - Core Principles (7 principles)
    - Technology Stack & Constraints
    - Development Workflow & Quality Gates
    - Governance
  Removed sections: N/A
  
  Templates requiring updates:
    - .specify/templates/plan-template.md ✅ aligned (Constitution Check section exists)
    - .specify/templates/spec-template.md ✅ aligned (user stories + requirements structure)
    - .specify/templates/tasks-template.md ✅ aligned (phased structure matches principles)
  
  Follow-up TODOs: None
-->

# DemoWails Constitution

## Core Principles

### I. Go-Backend Ownership (NON-NEGOTIABLE)

All business logic, file system access, OS-level operations, and data
processing MUST reside in Go. The frontend is a rendering and interaction
layer only. This separation is the architectural backbone of every Wails
application.

- **Go owns**: data models, services, file I/O, system calls, validation,
  and state that persists beyond the UI lifecycle.
- **Frontend owns**: rendering, animations, user input capture, and
  transient UI state (form drafts, modals, scroll position).
- No business rule MUST ever be enforced solely in JavaScript/TypeScript.
  The Go layer MUST be the single source of truth.
- **Rationale**: Go's type safety, concurrency model, and compiled
  performance make it the correct place for critical logic. Keeping the
  frontend thin ensures the app can be tested and verified without a
  browser context.

### II. Typed Contract Bridge

Every function exposed from Go to the frontend via Wails bindings MUST
have explicit input/output types. The auto-generated TypeScript bindings
MUST be treated as the contract between layers.

- Bound Go methods MUST use well-defined structs for request/response,
  never raw `interface{}` or untyped maps.
- Frontend code MUST import and use the generated TypeScript models from
  `wailsjs/go/` — never manually re-declare types.
- When a Go struct changes, the frontend MUST be updated to match before
  the change is considered complete.
- **Rationale**: Type drift between Go and JS is the #1 source of
  runtime errors in hybrid desktop apps. Strict typing at the bridge
  eliminates an entire class of bugs.

### III. Security-First Desktop Mindset

Desktop applications run with the user's full OS permissions. Every
feature MUST be designed with this threat model in mind.

- All user input received from the frontend MUST be validated and
  sanitized in Go before processing.
- File paths MUST be resolved and checked against allowed directories to
  prevent path-traversal attacks.
- Sensitive data (tokens, credentials, API keys) MUST never be stored in
  plain text; use OS-level secure storage or encrypted config.
- External process execution (`exec.Command`) MUST use explicit argument
  lists, never shell string interpolation.
- **Rationale**: Unlike web apps behind a server, desktop apps are a
  direct attack surface on the user's machine. Security is
  non-negotiable.

### IV. Frontend Quality & Accessibility

The web frontend MUST meet modern UI/UX standards. Users expect desktop
applications to feel native, responsive, and polished.

- All interactive elements MUST be keyboard-navigable and provide visible
  focus indicators.
- UI MUST be responsive within the application window (handle resize
  gracefully down to a defined minimum size).
- Loading and error states MUST be handled explicitly — no blank screens
  or silent failures.
- Animations MUST use GPU-accelerated properties (`transform`, `opacity`)
  and respect `prefers-reduced-motion`.
- Color palettes MUST maintain WCAG AA contrast ratios at minimum.
- **Rationale**: A desktop app that feels like a broken web page erodes
  user trust. Quality UI is a functional requirement, not a nice-to-have.

### V. Test-First Verification

Code MUST be verifiable without manual intervention. Automated tests are
the primary proof that the system works as designed.

- Go backend logic MUST have unit tests covering core services and edge
  cases. Use `go test` with table-driven tests.
- Integration tests MUST cover the Wails binding layer — verify that
  bound methods return expected results for known inputs.
- Frontend components with complex interaction logic SHOULD have tests
  (Vitest / Testing Library or equivalent).
- The Red-Green-Refactor cycle is the preferred workflow: write a failing
  test, make it pass, then clean up.
- **Rationale**: In a dual-language system, the seams between Go and JS
  are fragile. Tests at these boundaries catch regressions early.

### VI. Simplicity & YAGNI

Start with the simplest solution that meets the requirement. Complexity
MUST be justified and documented.

- Prefer standard library over third-party packages when the standard
  library solution is adequate.
- Do not introduce abstractions (interfaces, plugins, middleware layers)
  until a concrete second use case exists.
- Every added dependency MUST be justified — evaluate maintenance burden,
  binary size impact, and license compatibility.
- **Rationale**: Desktop apps ship as a single binary. Every dependency
  increases build complexity, attack surface, and maintenance cost.
  Simplicity is a feature.

### VII. Observability & Error Handling

All errors MUST be surfaced, logged, and handled gracefully. Silent
failures are forbidden.

- Go services MUST return structured errors with context
  (`fmt.Errorf("operation %s: %w", name, err)`).
- The frontend MUST display user-friendly error messages — never raw
  stack traces or Go error strings.
- Structured logging MUST be used in Go (e.g., `slog`) with consistent
  log levels (DEBUG, INFO, WARN, ERROR).
- Critical errors MUST be logged to a file for post-mortem analysis.
  Log rotation or size limits MUST be configured.
- **Rationale**: Desktop apps run on diverse user machines. Without
  structured logs and clear error surfaces, debugging production issues
  becomes impossible.

## Technology Stack & Constraints

- **Runtime**: Go 1.21+ with Wails v2/v3 framework
- **Frontend**: HTML + CSS + JavaScript/TypeScript (framework choice
  determined per-feature — vanilla, React, Svelte, or Vue are acceptable)
- **Build Target**: Windows (primary), macOS and Linux (secondary)
- **Packaging**: Single binary via `wails build`; NSIS installer for
  Windows distribution
- **State Management**: Go-side persistence (JSON files, SQLite, or
  bolt/bbolt); frontend uses only transient UI state
- **Styling**: Vanilla CSS preferred for maximum control; Tailwind
  acceptable if explicitly chosen per-project. No UI component libraries
  without explicit justification.
- **Minimum Window Size**: MUST be defined per app (recommended: 800×600)

## Development Workflow & Quality Gates

### Pre-Commit Checklist

1. `go vet ./...` — static analysis passes
2. `go test ./...` — all tests pass
3. Frontend linting passes (ESLint or equivalent)
4. No hardcoded secrets, file paths, or environment-specific values
5. Generated Wails bindings are up-to-date (`wails generate module`)

### Code Review Standards

- Every change touching Go ↔ JS boundary MUST include verification that
  both sides compile and types align.
- New bound methods MUST include at least one Go unit test.
- UI changes MUST include a screenshot or screen recording demonstrating
  the change.
- Performance-sensitive changes MUST include benchmark results.

### Build & Release

- Development: `wails dev` for hot-reload during development.
- Production: `wails build` produces the release binary.
- Version tagging follows MAJOR.MINOR.PATCH semantic versioning.
- Release notes MUST document user-facing changes and known issues.

## Governance

This constitution is the supreme authority for development decisions in
this project. All code, reviews, and architectural choices MUST comply
with the principles above.

- **Amendments**: Any principle change requires documented rationale, a
  version bump, and a migration plan for existing code that violates the
  new rule.
- **Versioning**: Constitution follows semantic versioning —
  MAJOR for principle removals/redefinitions, MINOR for additions,
  PATCH for clarifications.
- **Compliance**: Every pull request and code review MUST verify
  adherence. Violations MUST be flagged and resolved before merge.
- **Complexity Justification**: Any deviation from Principle VI
  (Simplicity) MUST be documented in the relevant plan or task file
  with explicit rationale.
- **Guidance**: Use this document as the primary reference during
  development. Runtime guidance and agent-specific instructions are
  secondary to these principles.

**Version**: 1.0.0 | **Ratified**: 2026-03-10 | **Last Amended**: 2026-03-10
