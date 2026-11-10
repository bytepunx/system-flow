---
id: TH-1032
title: "Plan for S-0678: keep agent.status's trial merges across calls, time Activity's parts, and measure 13 at once, three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md
  item: S-0678
status: open
participants: [planner-S-0678]
created: 2026-11-10T10:32:34Z
updated: 2026-11-10T10:32:34Z
---

# TH-1032 Plan for S-0678: keep agent.status's trial merges across calls, time Activity's parts, and measure 13 at once, three tasks in two layers

On wip/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md.

## Entries

### 2026-11-10T10:32:34Z planner-S-0678
S-0678 is planned: three tasks in two layers, with touches, a forecast of 1h30m, and a cost of delay of 0.04 USD a week.

## Proposed solution

I-0313's one instance left this repository's `agent.status` at about 1.1 s with 13 at once. Alone, `serve.Activity` took 540 to 600 ms. S-0602 has since kept the threads and conversations parsed, which took 250 to 390 ms of that. What is left is the holds with their trial merges, and `elsewhere`. `Activity` makes a new trial judge on every call, so each call runs `git --version`, a `git rev-parse` per branch, and a `git merge-tree` per overlap, though no branch moved.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3004 agent.status times serve.Activity's parts as phases: `activity.runs`, `activity.holds`, `activity.trial`, `activity.elsewhere` | `flai/internal/serve/agents_activity.go`, `activity_test.go`, `flai/cmd/serve_actions.go`, `serve_actions_test.go` |
| 1 | T-3008 The trial judge keeps its verdicts per root and commit pair across calls, checks git's version once a process, and reads branch heads with one `for-each-ref` | `flai/internal/storygit/trial.go`, `trial_test.go` |
| 2 | T-3013 Measure 13 calls at once before and after as S-0603 did, record it in `design/system/server-performance.md`, and close I-0313 | `design/system/server-performance.md`, the I-0313 file, `design/issues/summary.md` |

T-3004 and T-3008 share no path, so they run together. T-3013 waits for both: it measures T-3008's fix with T-3004's phases.

## Assumptions

- The 200 ms target is judged as S-0603 judged it: an in-process program, 13 calls released together, on the operator's host. That host has run at a load of about 100 on 24 cores. If what is left is the wait for a core, which the phases will show, T-3013 closes I-0313 and says so. If a part of `Activity` still takes most of a call, it adds a task for that part and leaves the issue open.
- A `git merge-tree` of two given commits always gives the same paths, so a verdict kept by its commits needs no other check. The resolved-file filter is still applied on each call.
- No `after`: S-0602, which this builds on, is done.
- Cost of delay stays at `flai cod`'s 0.04 USD a week, from flai's input of 1s per cycle. I set no input.

## Proposals

None to split, merge, or drop.
