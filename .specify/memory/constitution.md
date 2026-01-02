# azure-tts Constitution

This constitution extends the [TeleAgent Ecosystem Constitution](../../../.specify/memory/constitution.md).

## Project Overview

`azure-tts` is a Go library package in the TeleAgent gopkg ecosystem, providing reusable functionality for TeleAgent services.

## Core Principles (Inherited + Extended)

### I. Go Best Practices

**Follow Go idioms and conventions:**

- Use `effective go` style guidelines
- Prefer composition over inheritance
- Use interfaces for abstraction
- Error handling: return errors, don't panic
- Use context.Context for cancellation and timeouts

### II. Library Design

**Self-contained, reusable packages:**

- Single responsibility principle
- Minimal external dependencies
- Clear, documented public API
- Internal packages for implementation details
- Semantic versioning for releases

### III. API Stability

**Backward compatibility:**

- Breaking changes require major version bump
- Deprecation warnings before removal
- Migration guides for breaking changes
- Changelog maintained for each release

### IV. Testing Requirements

**Comprehensive test coverage:**

- Unit tests: Table-driven tests with testify
- Example tests for documentation
- Benchmark tests for performance-critical code
- Minimum 70% coverage for new code

### V. Documentation

**Self-documenting code:**

- GoDoc comments for all exported types/functions
- README with usage examples
- Examples in `example/` directory where applicable

## Governance

This constitution extends the ecosystem constitution.
Library-specific rules take precedence for API design decisions.

**Version**: 1.0.0 | **Ratified**: 2026-01-02 | **Last Amended**: 2026-01-02
