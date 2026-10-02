---
title: Code quality
updated: 2026-10-02
audience: agent
order: 60
status: active
topics: [all]
roles: [verify]
---

# Code quality

How to write maintainable code, and how it is tested.

## Rules

- Observe SOLID principles:
  - SRP at every level: functions, modules, or services, in ascending order.
  - No direct dependency on I/O libraries, depend on an idealized abstraction.
- Model behavior in modules that can be tested without concrete I/O dependencies.
- Transport modules deliver functionality by integrating behavioral modules and their dependency adapters.
- Tests come in three increasingly expensive tiers, each with its own command and folder:
  - Each tier runs only once the prior passes.
  - Behavior tests:
    - cover a module's public API, with no I/O
    - take seconds to run
    - use the language's native test layout
    - run with `make test`
  - Integration tests:
    - one service or component against real adapters
    - located in `<project>/tests/integration/`
    - run with `make integration`
  - Smoke tests:
    - end-to-end paths through the full system from the consumer's perspective
    - located in `tests/smoke/` at the repository root
    - run with `make smoke`
    - run once before review and in CI
  - A change is not done until affected tiers pass.
  - A story is not done until all three suites pass.
  - New or changed behavior gets a test that fails without it; a fixed bug gets a test that reproduces it.
  - While working, run only the tests for what you changed, for quick feedback.
  - Before moving to review, a verifier runs the whole suite, as `delegation.md` says; do not run it yourself as well. Where the harness has no verifier, run it yourself.
- Lint clean. Run the project's linter as configured; fix findings rather than suppressing them. A suppression needs a comment saying why.
- Report test and lint results as they are. "Tests pass" means you ran them and saw them pass in this environment.
- Keep changes small and reviewable:
  - one story at a time
  - one logical concern per commit where possible
- Fix root causes:
  - Do not hide errors with retries or over-broad catches.
  - Special cases should only be introduced when specified in requirements.
- Abstractions:
  - Abstract I/O and external dependencies with adapters.
  - Delay speculative abstractions, generalize after the second use arrives.
- No dead code, no commented-out code, no TODO comments: work left to do is a task or a story.
- Dependencies:
  - Adding one needs a reason in `design/tech`.
  - Include version, use, and alternatives considered.
  - Pin versions.
  - Prefer the standard library where possible.
- Consistency & Idioms:
  - Use dominant language idioms.
  - Use consistent naming, structure, error handling, formatting.
- Public interfaces get a one-line doc comment saying what, not how.
- Errors carry actionable context:
  - what was attempted
  - on what
  - what to do
  - no silent swallowing
- Do not optimize without a measurement.
- Do not leave a known performance cliff unmentioned.
- Fixtures and test data belong in the source tree, never gitignored.
- Read the existing code before changing it. Do not guess an API; open the file.

## When in doubt

- If you would not want to review it, do not submit it.
- If a test is hard to write, revisit the design; note it in the narrative.

<!-- system-flow:end-of-baseline -->

## Project additions
