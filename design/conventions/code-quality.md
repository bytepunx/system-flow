---
title: Code quality
updated: 2026-09-15
audience: agent
order: 60
status: active
---

# Code quality

What every code change carries with it, and what it never carries.

## Rules

- Single responsibility applies at every level: separate concerns belong in separate functions, modules, or services, in that ascending order.
- Modules never depend directly on third-party libraries for I/O. They depend on an abstraction that expresses the ideal operation, so the concrete dependency can change without rewriting behavior.
- Each component presents a public API that implements its behavior and can be tested without concrete external dependencies. That is what behavior testing means here.
- Integration tests exercise each service in isolation against real adapters.
- Smoke tests exercise end-to-end paths through the system, are automated, and run in CI so breaking changes are caught before merge.
- Tests accompany the change. New behavior gets a test that fails without it; a fixed bug gets a test that reproduces it. Run the whole suite, not just the new test, before reporting.
- Lint clean. Run the project's linter as configured; fix findings rather than suppressing them. A suppression needs a comment saying why.
- Report test and lint results as they are. "Tests pass" means you ran them and saw them pass in this environment.
- Keep changes small and reviewable: one story's worth, one concern per commit where possible. If a diff needs a tour to review, split it.
- Fix the root cause. Do not paper over an error with a retry, a broad catch, or a special case unless the story says that is the fix.
- No speculative abstractions beyond the I/O boundaries above. Build what the story needs; generalize when the second use arrives.
- No dead code, no commented-out code, no TODO without an item ID.
- Fewest dependencies. Adding one needs a reason in `design/tech` with version, use, and alternatives considered. Pin versions. Prefer the standard library.
- Match the surrounding code: naming, structure, error handling, formatting. Consistency beats preference.
- Public interfaces get a one-line doc comment saying what, not how.
- Errors carry context the reader can act on: what was attempted, on what, and what to do. No silent swallowing.
- Do not optimize without a measurement. Do not leave a known performance cliff unmentioned.
- Fixtures and test data live under the source tree in a folder that is not git-ignored. Verify by cloning if unsure.
- Read the existing code before changing it. Do not guess an API; open the file.

## When in doubt

- If you would not want to review it, do not submit it.
- If a test is hard to write, the design is telling you something; note it in the narrative.

<!-- system-flow:end-of-baseline -->

## Project additions
- Go: standard library `testing` only, no assertion framework; `go test -race ./...`; golangci-lint v2 config in `flai/.golangci.yaml`.
- The I/O abstraction for git and docker is `flai/internal/execx.Runner`. cobra and goccy/go-yaml are used directly as the CLI and serialization layers; they are the boundary, not behind it.
- Fixture projects live under `internal/*/testdata/`; the round-trip test in `internal/workitem` and the render test in `cmd` run against this repository and the template and must stay green.
- Smoke test for the standard is `scripts/template-test.sh`: render the template, run `flai check --strict` on the result. CI runs it in `system-flow-check.yml`.
