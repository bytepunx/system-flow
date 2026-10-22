---
id: T-2043
type: task
nature: feature
title: flai guard lets a story's sub-agent run flai components list and serve, a read-only one list alone, and the strategic agents list alone
status: done
parent: S-0477
owner: alex
created: 2026-10-22T20:59:53Z
updated: 2026-10-22T21:12:44Z
transitions:
  - to: ready
    at: 2026-10-22T21:00:26Z
    by: agent-S-0477
  - to: in-progress
    at: 2026-10-22T21:00:27Z
    by: agent-S-0477
  - to: done
    at: 2026-10-22T21:12:44Z
    by: agent-S-0477
stream: S-0477
tags: []
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/component-library.md, design/adrs/README.md, flai/internal/guard/call.go, design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md]
usage:
  source: log
  seconds: 737
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 20662
      cache_read: 4870105
      cache_write: 141132
      cost: 2.3803
---

# T-2043 flai guard lets a story's sub-agent run flai components list and serve, a read-only one list alone, and the strategic agents list alone

## Work

`components` is in no list of `flai/internal/guard/guard.go`, so every sub-agent is refused `flai components list` and `flai components serve` (I-0165, I-0156). List only reads the manifest and the components' props; serve runs the sub-project's own Vite dev server on 127.0.0.1, writing only under `.flai-cache` and, with `--story`, the worktree's git-ignored dependencies. Neither changes a work item, a thread, or history (ADR-0060).

- Add `components list` to `cliReads`, so every agent may list.
- Let a story's sub-agent run `components serve` in its `subAgent` rules, but not a read-only sub-agent (its definition grants neither Edit nor Write, ADR-0141: serve writes files outside the temporary folder), and not the planner, the orchestrator, or the analyzer themselves. `components config`, which rewrites `system-flow.yaml`, stays refused.
- Tests in `guard_test.go` that fail without the change: a task sub-agent's list and serve pass, through `flai` and `scripts/flai.sh` alike; a read-only sub-agent's list passes and serve is refused; the strategic agents list but do not serve; config is refused to every sub-agent.
- Record the decision with `flai adr new` (no `--commit`), and update the guard's help in `flai/cmd/guard.go`, the generated `docs/users/flai-reference.md`, the guard row of `design/system/flai-cli.md`, `docs/users/flai.md`'s sub-agents section, and `design/system/component-library.md` where it says who runs serve.

## Done when

- `flai test` passes on the changed paths.
- The tests above fail on main's guard and pass on this one.
- The docs and the ADR say which agents may list and serve.

## Notes
