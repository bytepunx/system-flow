---
id: T-0843
type: task
nature: improvement
title: The design, scripts README, operator docs, Makefile help, and close-out header describe the tiers as flai-test.sh now runs them
status: cancelled
parent: S-0267
owner: alex
created: 2026-10-05T00:18:19Z
updated: 2026-10-05T00:20:01Z
transitions:
  - to: cancelled
    at: 2026-10-05T00:20:01Z
    by: agent-S-0267
stream: S-0267
tags: [docs, scripts]
touches: [design/system/devex.md, scripts/README.md, docs/operators/index.md, Makefile, scripts/close-out.sh]
after: [T-0842]
---
# T-0843 The design, scripts README, operator docs, Makefile help, and close-out header describe the tiers as flai-test.sh now runs them

## Work

Say what T-0842 changed everywhere the tiers are described:

- `design/system/devex.md`, the "Test tiers" row: `make test` and `make integration` are unchanged for one tier. `scripts/flai-test.sh`, which the close-out runs when a branch changes `flai/`, runs the full Go run once, then flaiover's unit tests, then smoke, and skips the `-short` run.
- `scripts/README.md`: correct the `flai-test.sh` row, which now says "then all three tiers in order". Add a row for the flaiover unit script if T-0842 added one, and say in the `test.sh` row that it calls that script.
- `docs/operators/index.md`: the `flai serve checks set --name flai -- scripts/flai-test.sh` example near "checks set". Say in a sentence what that check runs, so an operator knows it runs the Go tests once. Change nothing else there unless it describes the tiers.
- `Makefile`: the `flai-test` help text, "Lint plus all three test tiers in order".
- `scripts/close-out.sh`: the header comment, "flai/ runs flai-test.sh (lint, every tier, ...)". Do not change the script's behaviour.

It waits for T-0842, because the words describe what T-0842 settles, such as whether a new script exists and what it is called.

## Done when

- [ ] each of the five files describes `flai-test.sh` as running the Go tests once, the full run, and `make test` and `make integration` as unchanged
- [ ] `scripts/lint-md.sh` passes and `flai check --strict` reports nothing for the files changed
- [ ] `scripts/close-out.sh` behaves as before: only its comments changed, as `git diff` shows

## Notes
- 2026-10-05T00:20:01Z: moved to cancelled: duplicate of T-0848; its close-out.sh header touch is folded into T-0848
