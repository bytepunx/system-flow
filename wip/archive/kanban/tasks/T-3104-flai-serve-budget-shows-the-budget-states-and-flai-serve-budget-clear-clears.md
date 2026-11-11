---
id: T-3104
type: task
nature: remediation
title: flai serve budget shows the budget states and flai serve budget clear clears one by hand
status: done
parent: S-0720
owner: alex
created: 2026-11-11T08:55:06Z
updated: 2026-11-11T22:55:19Z
transitions:
  - to: ready
    at: 2026-11-11T22:36:30Z
    by: agent-S-0720
  - to: in-progress
    at: 2026-11-11T22:36:30Z
    by: agent-S-0720
  - to: done
    at: 2026-11-11T22:55:19Z
    by: agent-S-0720
stream: S-0720
tags: [flai, cli, serve, budget]
touches: [flai/cmd/serve_budget.go, flai/cmd/serve_budget_test.go, flai/cmd/serve.go, docs/users/flai-reference.md, flai/internal/guard/hostconfig.go, flai/internal/guard/guard_tables.go, flai/internal/guard/hostconfig_test.go]
after: [T-3098]
usage:
  source: log
  seconds: 1129
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 73
      output: 20932
      cache_read: 5733880
      cache_write: 96948
      cost: 2.1537
---
# T-3104 flai serve budget shows the budget states and flai serve budget clear clears one by hand

## Work

- Add `flai serve budget`, in a new `flai/cmd/serve_budget.go`, registered beside the other `serve` subcommands in `flai/cmd/serve.go`: it lists each state in the host's `budget.json`, its key, limit, since, until, and whether a probe runs, with `--json`.
- Add `flai serve budget clear [<provider>]`: with a key, it clears that state; with none, every state. It calls the clear T-3098 adds, with why `by hand`, and says what it cleared.
- Cover both in `flai/cmd/serve_budget_test.go`: the list, a clear by key, a clear of all, and a key that names no state.
- Regenerate `docs/users/flai-reference.md` and the flag table of `docs/operators/settings.md` with `make flai-reference`, since the new commands change both.
- Waits for T-3098: it calls the clear T-3098 adds.

## Done when

- `flai serve budget` shows a held state and `flai serve budget clear <key>` clears it, in tests.
- `TestReferenceIsCurrent` and the settings index tests pass.
- `flai test` on the changed paths passes.

## Notes
