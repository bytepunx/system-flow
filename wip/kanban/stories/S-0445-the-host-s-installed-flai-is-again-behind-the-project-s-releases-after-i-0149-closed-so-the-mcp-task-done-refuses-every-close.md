---
id: S-0445
type: story
nature: improvement
title: The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close
status: review
owner: alex
created: 2026-10-17T16:03:15Z
updated: 2026-10-21T00:51:45Z
transitions:
  - to: ready
    at: 2026-10-20T20:40:53Z
    by: alex
  - to: in-progress
    at: 2026-10-20T20:42:49Z
    by: agent-S-0445
  - to: review
    at: 2026-10-21T00:51:45Z
    by: agent-S-0445
tags: [flai]
touches: [design/adrs/README.md, flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/buildinfo_test.go, flai/cmd/version.go, flai/cmd/version_test.go, flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/adrs/0180-flai-task-done-refuses-an-older-flai-only-when-the-newest-flai-release-raised.md, flai/internal/buildinfo/task-done-rules.txt, docs/users/flai-reference.md, flai/cmd/task_done.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3034
  priced_by: harness
  turns:
    - day: 2026-10-20
      test_runs: 2
      hand_edits: 2
      work: 83
    - day: 2026-10-21
      ceremony: 3
      work: 24
  models:
    - model: claude-opus-5-5
      input: 342
      output: 94594
      cache_read: 26255635
      cache_write: 689388
      cost: 12.1844
  strategic:
    - kind: orchestrator
      seconds: 6547
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 260
          output: 5216
          cache_read: 92416184
          cache_write: 76508
          cost: 23.1873
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-17T16:03:15Z
  value: 37.5
  by: planner-S-0445
  at: 2026-10-18T02:31:13Z
forecast:
  duration: 1h
  delivery: 2026-10-20T22:00:00Z
  basis: "flai forecast's 35m (87 s per unit over 78 large improvement stories, size 24) raised to 1h: seven tasks in five layers run mostly one at a time with an ADR first, and S-0450, six tasks of the same shape, took about an hour from start to review on 2026-10-20."
  by: planner-S-0445
  at: 2026-10-20T20:42:43Z
finalized:
  by: alex
  at: 2026-10-20T20:40:48Z
---
# S-0445 The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

## Goal

This story remediates [I-0157](../../../design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md), "The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0157 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0157 is closed with `flai issue close I-0157 --reason` saying what fixed it

## Tasks
- T-1920 An ADR refines ADR-0142: task done refuses only when the newest flai release's task-done rules number is above the running flai's
- T-1921 flai carries a task-done rules number in buildinfo, and flai version --json prints it
- T-1922 flai release writes the task-done rules number into each flai tag's annotation
- T-1923 flai task done refuses only when the newest flai tag's task-done rules number is above the running flai's, and keeps the minor rule for a tag without one
- T-1924 The MCP task_done's check that the installed flai has caught up compares task-done rules numbers, read from its flai version --json
- T-1925 The CLI design, the narrative design, and the user and operator docs describe the task-done rules number and when task done refuses
- T-1926 I-0157 is closed with a reason naming the task-done rules number that fixed it

## Notes

Cost of delay inputs set by flai from I-0157. time_lost_per_cycle 15m: 3m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-10-17T01:51:51Z, 0.6 days before this story; under one cycle counts as one).

### Planning

The fix planned is option B of TH-0506. `flai task done` refuses only when the newest flai tag's task-done rules number is above the running flai's. Each release records that number in its tag's annotation. That thread is still unanswered. Both recommendations on it, the planner's and the orchestrator's, name B for this story, and the operator finalized the story at 2026-10-20T20:40Z. T-1920's ADR therefore comes first, so the operator still accepts the decision before the code lands. If the operator answers A or C instead, rewrite T-1920 and drop T-1921 to T-1925.

Touches. The story declared none, and `flai touches suggest` had nothing to start from. Every touch below is a file; no folder touch is kept.

| Touch | From |
|-------|------|
| `design/adrs/README.md` | design: ADR-0142 and ADR-0160 are refined; the new ADR file is named by `flai adr new` and widened by `flai task done` |
| `flai/internal/buildinfo/buildinfo.go`, `buildinfo_test.go` | layout: the leaf package that `release`, `taskdone`, and `mcpserver` can all import |
| `flai/cmd/version.go`, `flai/cmd/version_test.go` (new) | layout: `flai version --json`, which ADR-0160's check reads |
| `flai/internal/release/release.go`, `release_test.go` | layout: `Tag` and `TagPending` make the annotated flai tags |
| `flai/internal/release/outdated.go`, `outdated_test.go` | design: ADR-0142 point 1, `FlaiOutdated` |
| `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` | design: ADR-0142 point 3, `refuseOutdated` |
| `flai/internal/mcpserver/task.go`, `task_test.go`, `installed.go`, `installed_test.go` | design: ADR-0160, `closeWithInstalled` |
| `design/system/flai-cli.md`, `design/system/agent-narrative.md` | design: the sections that describe the refusal |
| `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/runbooks/update.md` | design: the user and operator passages on the refusal |
| I-0157's file, `design/issues/summary.md` | goal: criterion 2 closes the issue |

Forecast: 1h, up from `flai forecast`'s 35m. That figure is 87 s per unit over 78 large improvement stories, at size 24. The tasks run in five layers, mostly one at a time, with an ADR first. S-0450, six tasks of the same shape, took about an hour from start to review on 2026-10-20.

Cost of delay: 37.50 USD a week, as `flai cod` gives from the inputs: 15m lost per 168h cycle at 150 USD an hour. The inputs stand as flai set them from five occurrences. I-0157 now counts 14, so they understate the cost, but they are not the planner's to change.
