---
title: Code quality
updated: 2026-11-09
audience: agent
order: 60
status: active
topics: [all]
roles: [story, verify, analyze]
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
  - Before moving to review, have the whole suite run through the close-out, which runs `flai verify`: ask flai serve for its job and wait for its `verified` event, or, where no flai serve runs, run it yourself, as `delegation.md` says; do not hand that run to a sub-agent.
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
- Go: standard library `testing` only, no assertion framework; `go test -race ./...`; golangci-lint v2 config in `flai/.golangci.yaml`.
- The I/O abstraction for git and docker is `flai/internal/execx.Runner`. cobra and goccy/go-yaml are used directly as the CLI and serialization layers; they are the boundary, not behind it.
- Fixture projects live under `internal/*/testdata/`; the round-trip test in `internal/workitem` and the render test in `cmd` run against this repository and the template and must stay green.
- Tiers here: `make test` runs `go test -race -short ./...` in `flai/` (behavior; tests that need real git or the monorepo skip under `-short`); `make integration` runs the full `go test -race ./...` including the git-clone and monorepo round-trip tests; `make smoke` runs `scripts/template-test.sh` (render the template, `flai check --strict` on the result). `flai/tests/integration/` holds the tests that need something beside the code, each skipping without it: `mcp_http_test.go` a running dashboard, `golangci_lock_test.go` golangci-lint, `forge_github_test.go` a GitHub repository and token, and `opencode_test.go` node and OpenCode. The integration tier runs them with the rest of `./...`. There is no `tests/smoke/`: the smoke tier is its scripts. Elsewhere the tiers are selected by `-short` and by script. CI runs all three in `system-flow-check.yml`. The paid gateway tests, `TestGatewayOpenRouter` and `TestGatewayLiteLLM`, run in no tier but the `gateway` tier, `scripts/gateway-smoke.sh`, which is manual: only by name, `flai test gateway`, or on the host's schedule, never in a close-out, under `--all`, or in the merge queue's batch, and the smoke tier no longer runs them (S-0548). They are the only check that the guard and `permission_prompt` hold through a gateway (S-0351, S-0356); a run spends about 0.09 USD at OpenRouter, about 0.32 USD a check in Claude Code's count, which `FLAI_TEST_GATEWAY_BUDGET` caps. `flai test` and the close-out's `go-test` tier run the tests of the changed packages and of the packages that import them with `-short`, skipping `^TestGateway` too (I-0221). The integration tier, `scripts/integration.sh`, runs them without `-short`. The close-out of a story that touches or changes `flai/` or `template/` runs `integration` and `smoke` too, and the merge queue's batch runs them again ([ADR-0233](../adrs/0233-while-the-merge-queue-can-batch-the-close-out-runs-a-queue-tier-whose-when.md)). So the tests that skip under `-short` are the integration tier's own: those that need real git or the monorepo, among them the four that check, lint, and prime the whole monorepo, `TestMonorepoIsClean`, `TestRepositoryLintsClean`, `TestPrimeOnThisRepository`, and `TestPrimeComesInPartsOnThisRepository`, which took minutes under `-race` and test the repository rather than the change. This reverses I-0201's choice: a test you gate under `-short` first runs in the integration tier, in the close-out of a story that touches `flai/` or `template/`, under `flai test --all`, or in the batch, so gate only a test that needs real git or the monorepo (S-0548). Under `flai verify`, in a close-out or a batch, `TestMonorepoIsClean` skips, and so does `TestRepositoryLintsClean` unless the stories change `flai/internal/mdlint`, the markdownlint configuration, or `scripts/lint-md.sh`: the smoke tier's `scripts/check.sh` and `scripts/lint-md.sh` run the check and the lint there, once.
