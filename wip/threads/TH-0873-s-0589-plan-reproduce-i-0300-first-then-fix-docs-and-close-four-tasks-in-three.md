---
id: TH-0873
title: "S-0589 plan: reproduce I-0300 first, then fix, docs, and close; four tasks in three layers, 1h30m, 25 USD a week"
anchor:
  path: wip/kanban/stories/S-0589-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md
  item: S-0589
status: open
participants: [planner-S-0589]
created: 2026-11-07T10:07:55Z
updated: 2026-11-07T10:07:55Z
---

# TH-0873 S-0589 plan: reproduce I-0300 first, then fix, docs, and close; four tasks in three layers, 1h30m, 25 USD a week

On wip/kanban/stories/S-0589-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md.

## Entries

### 2026-11-07T10:07:55Z planner-S-0589
**Plan: four tasks in three layers.** No question is open. S-0589 was finalized and pulled while I planned, so agent-S-0589 reviews these tasks.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2573 Reproduce I-0300 in tests and record its cause | `agents_test.go`, `serve_test.go` |
| 2 | T-2574 Fix: a story already in ready starts on the first look after `agent` is turned on (after T-2573) | `agents.go`, `serve.go`, the two tests, `cmd/serve_actions.go` and its test, `hostapi/writes.go` and its test |
| 3 | T-2575 Operators guide and `flai-cli.md` (after T-2574) | `docs/operators/index.md`, `design/system/flai-cli.md` |
| 3 | T-2576 `flai issue close I-0300` (after T-2574) | the I-0300 file, `design/issues/summary.md` |

**Proposed solution.** The issue's guess, that serve starts an agent only when a story enters ready, does not match the code. `(*launcher).look` in `flai/internal/serve/agents.go` runs on every change under `kanban/` and `threads/` and every minute, and reads the `agent` action each time. T-2573 tests three candidate causes:

- Turning the action on writes no path the watcher sees.
- The project loop's `time.After(time.Minute)` in `serve.go` is made again on every pass, so frequent orchestrator or role wakes can starve the minute look.
- A project the board's Import added while serve ran may not have had its looks run.

T-2574 then makes turning `agent` on wake the launcher at once, and fixes whichever cause the tests confirm.

**Figures.**

- Forecast: 1h30m, delivery 2026-11-07T12:00:00Z. flai said 24m. I raised it because the cause is unconfirmed and serve's loop tests are slow on this host (S-0560).
- Cost of delay: 25 USD a week, from your 10m-per-cycle input, kept as flai cod gave it.

**Assumptions.**

- No ADR is needed: ADR-0041 and ADR-0043 already decide that a ready story gets its agent whenever it entered ready.
- No `after` on another story: S-0587's default-harness gate is a different gate, since S-0001 started on your Start agent (TH-0860).
- T-2573's tests stay red until T-2574. Its agent runs `flai test` only on those two files between the two tasks.

**Nothing to split, merge, or drop.**
