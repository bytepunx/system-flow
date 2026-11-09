---
id: S-0632
type: story
nature: remediation
title: A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it
status: in-progress
owner: alex
created: 2026-11-07T23:40:59Z
updated: 2026-11-09T22:11:55Z
transitions:
  - to: ready
    at: 2026-11-09T22:01:06Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T22:01:12Z
    by: agent-S-0632
tags: [flai, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/handover_test.go, flai/internal/serve/proc_linux.go, flai/internal/serve/proc_unix.go, flai/internal/serve/proc_other.go, flai/internal/serve/proc_linux_test.go, design/system/flai-cli.md, design/issues/I-0299-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md, design/issues/summary.md, design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 730
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-09
      ceremony: 3
      hand_edits: 1
      work: 30
  models:
    - model: claude-opus-5-5
      input: 78
      output: 407
      cache_read: 3899970
      cache_write: 143469
      cost: 1.4784
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-07T23:40:59Z
  value: 50
  by: planner-S-0632
  at: 2026-11-09T21:58:29Z
forecast:
  duration: 45m
  delivery: 2026-11-09T22:53:00Z
  basis: "flai forecast's 22m (median 120 s per unit over 46 done medium-band remediation stories, size 11) doubled for a test that must reproduce an unreaped zombie, two fixes in two areas, and three tasks; delivery is flai's start at 6th in the pull order plus 45m."
  by: planner-S-0632
  at: 2026-11-09T21:58:29Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:01:01Z
---
# S-0632 A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it

## Goal

This story remediates [I-0299](../../../design/issues/I-0299-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md), "A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0299 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0299 is closed with `flai issue close I-0299 --reason` saying what fixed it

## Tasks
- T-2769 A launcher that hands its agent over reaps it in the background instead of releasing it
- T-2770 Owns and Alive count a zombie process as gone on Linux
- T-2771 Describe the reaping and the zombie rule in flai-cli.md and close I-0299

## Notes

Cost of delay inputs set by flai from I-0299. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:21:58Z, 0.6 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0299's one instance and the code it names. Both halves, since either alone leaves a gap:

1. Reap a handed-over child (T-2769). `launcher.await` (`flai/internal/serve/agents.go`), when `handOver` is set, closes the output and calls `cmd.Process.Release()` without `cmd.Wait()`. A short-lived parent (`flai analyze`, `flai plan`, `flai serve agent start`) exits and the child is reparented and reaped. A long-lived parent, `flai mcp` over stdio, which starts a run through `role_run`, `analyze`, `plan`, or `curate`, never reaps it. The fix waits in a background goroutine instead of releasing, and records nothing: the serving flai still settles the run in `settleOrphans`.
2. Count a zombie as gone (T-2770). `Owns` and `Alive` (`flai/internal/serve/proc_unix.go`) read the session and start from `/proc/<pid>/stat` (`proc_linux.go`) but not the state. Signal 0 succeeds on a zombie, so `running()` holds the run open. On Linux, state `Z` (or `X`) makes the process not running. This also covers a zombie whose parent is any other long-lived process.

T-2771 then describes both in `design/system/flai-cli.md` and closes I-0299.

D-0500 (TH-0859) records the operator waiting on this run by hand in S-0584 and naming I-0299 as needing its own story.

Touches, each from where it came:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/serve/agents.go` | design (I-0299 names it) | `await`'s hand-over branch |
| `flai/internal/serve/handover_test.go` | layout | new test file for the reaping, kept apart from `agents_test.go` (co-change 51%) |
| `flai/internal/serve/proc_linux.go` | design (I-0299 names it) | read the state field |
| `flai/internal/serve/proc_unix.go` | design (I-0299 names it) | `Owns` and `Alive` honour the state |
| `flai/internal/serve/proc_other.go` | layout | same `process` signature off Linux |
| `flai/internal/serve/proc_linux_test.go` | layout | new test of the stat parsing on a fixture |
| `design/system/flai-cli.md` | co-change (10%) and design | its `flai serve agent` row describes `handOver` and `Owns` |
| `design/issues/I-0299-…md`, `design/issues/summary.md` | criteria | `flai issue close I-0299` writes both |

No folder touch. No `after`: the story reads no other open story's output.

Figures:

- Forecast: 45m. `flai forecast` gives 22m, the median 120 s per unit of size over 46 done medium-band remediation stories, times size 11. Doubled because the test must reproduce an unreaped zombie in a live process, and the fix spans two areas over three tasks. Delivery is flai's start at 6th in the pull order plus 45m.
- Cost of delay: 50 USD a week, as `flai cod` gives it from the operator's 20m per 168h cycle at 150 USD an hour. It stands. Every hand-over launcher started from `flai mcp` shares the cause, so occurrences may rise, but the inputs are the operator's.
