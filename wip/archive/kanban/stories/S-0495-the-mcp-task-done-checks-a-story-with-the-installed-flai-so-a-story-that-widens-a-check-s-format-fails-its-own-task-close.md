---
id: S-0495
type: story
nature: improvement
title: The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close
status: done
owner: alex
created: 2026-10-24T20:28:57Z
updated: 2026-10-25T21:19:07Z
transitions:
  - to: ready
    at: 2026-10-25T20:35:10Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-25T20:35:15Z
    by: agent-S-0495
  - to: review
    at: 2026-10-25T20:58:51Z
    by: agent-S-0495
  - to: done
    at: 2026-10-25T21:19:07Z
    by: orchestrator
tags: [cli, mcp]
touches: [design/adrs/README.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/agent-narrative.md, design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md, design/issues/summary.md, docs/operators/settings.md, docs/users/flai.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/buildinfo/task-done-rules.txt, system-flow.yaml, design/adrs/0201-when-the-manifest-names-flai-command-task-done-s-check-step-runs-that-command-s.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1426
  priced_by: harness
  turns:
    - day: 2026-10-25
      hand_edits: 1
      work: 37
  models:
    - model: claude-opus-5-5
      input: 152
      output: 55247
      cache_read: 8061445
      cache_write: 432093
      cost: 5.4281
  strategic:
    - kind: orchestrator
      seconds: 2143
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 109
          output: 1655
          cache_read: 9461464
          cache_write: 102048
          cost: 2.0969
        - model: claude-sonnet-5-5
          input: 8
          output: 65
          cache_read: 77041
          cache_write: 39091
          cost: 0.0794
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-24T20:28:57Z
  value: 50
  by: planner-S-0495
  at: 2026-10-25T20:30:34Z
forecast:
  duration: 1h30m
  delivery: 2026-10-26T11:30:00Z
  basis: "flai forecast gave 24m (89 s per unit over 95 done improvement stories on claude-opus-5-5 in the large band, size 16) and delivery 2026-10-26T10:12Z, 21st in the pull order; raised to 1h30m because nine of the fourteen touches are one-paragraph document edits while the change is an ADR, a manifest key, a subprocess path in task done's check step with a fake-flai test, and a possible rules bump; delivery keeps the forecast's start and cycle factor 1.19."
  by: planner-S-0495
  at: 2026-10-25T20:34:31Z
finalized:
  by: orchestrator
  at: 2026-10-25T20:35:04Z
---
# S-0495 The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close

## Goal

This story remediates [I-0194](../../../design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md), "The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0194 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0194 is closed with `flai issue close I-0194 --reason` saying what fixed it

## Tasks
- T-2108 An ADR settles, from I-0194's instances, how task done checks a story whose tree builds the flai that knows its format
- T-2109 The manifest names, beside flai.minimum, the command that runs the flai the tree builds, with its settings row, and this repository sets it to scripts/flai.sh
- T-2110 task done's check step runs the flai the tree builds when the manifest names it, in the story's worktree, says which flai checked, and stops when that flai cannot run
- T-2111 The design and the guides say that task done's check runs the flai the tree builds when the manifest names it
- T-2112 I-0194 is closed with the reason naming the fix

## Notes

Cost of delay inputs set by flai from I-0194. time_lost_per_cycle 20m: 4m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-10-22T22:06:51Z, 1.9 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0495 on 2026-10-25. The remedy the tasks build, which T-2108's ADR settles first: a manifest key beside `flai.minimum`, `flai.command` as recommended, names the command that runs the flai the tree builds (`scripts/flai.sh` here), and task done's check step runs it in the story's worktree in place of the in-process check, so the rules that check a branch are the branch's own.

Touches, file by file. Declared: none; the story was made from I-0194 with no touches, so `flai touches suggest` refused to run and co-change gave nothing. Design (I-0194's instances, ADR-0160, ADR-0180, `flai-cli.md` § Versions: the host's flai and the tree) and layout (the code the explorer located):

- `flai/internal/taskdone/taskdone.go`, `taskdone_test.go`: the check step, `run.check`, calls `check.RunScoped` in-process; the tests fake a tree flai as `mcpserver/installed_test.go` fakes the installed one.
- `flai/internal/manifest/manifest.go`, `manifest_test.go`: the key, beside `Minimum`.
- `flai/internal/buildinfo/task-done-rules.txt`: ADR-0180 raises it when what task done checks changes; T-2108's ADR decides.
- `system-flow.yaml`: this repository sets the key.
- `docs/operators/settings.md`: every manifest key has a row (`TestSettingsIndexIsComplete`); `design/system/project-manifest.md` describes it.
- `design/system/flai-cli.md`, `design/system/agent-narrative.md`, `docs/users/flai.md`: where task done and the MCP `task_done` are described.
- `design/adrs/README.md`: the index `flai adr new` updates; the ADR's own file cannot be named before its slug is known, and `design/adrs` is in `claims.shared`, so no folder touch is kept and `flai task done` adds the file at T-2108's close (ADR-0144).
- `design/issues/I-0194-….md`, `design/issues/summary.md`: the close.

Forecast: `flai forecast` gave 24m, 89 s per unit over 95 done improvement stories on claude-opus-5-5 in the large band times size 16 (2 criteria, 14 touches), delivery 2026-10-26T10:12Z, 21st in the pull order behind twenty backlog stories with an in-progress limit of 5. Raised to 1h30m: the size counts touches, and nine of the fourteen are one-paragraph document edits, while the change is an ADR, a new manifest key with its settings row, a subprocess path in task done's check step with a fake-flai test, and a possible rules bump, across five tasks in four layers. Delivery 2026-10-26T11:30Z keeps the forecast's start and its cycle factor of 1.19 over the raised duration.

Cost of delay: `flai cod` gives 50.00 USD a week from the inputs flai set from I-0194, 20m lost per 168h cycle at 150 USD an hour, one cycle a week; written as is. The inputs are not changed.

### Accepted by the orchestrator

- Verified: 96ea776de2c3028c4f986b306378b053e5ea5805
- Batch: 20261025-210729
- Enqueued verified: a406b12064cb6fadb3a3416ed5bb235dbb4be52c
- At: 2026-10-25T21:19:07Z

Verdict: Accept S-0495: both criteria are met by in-touches changes at a406b12, with a reproducing test, ADR-0201 and the docs, and a passing verify.
- 1: flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/buildinfo/task-done-rules.txt, system-flow.yaml, design/adrs/0201-when-the-manifest-names-flai-command-task-done-s-check-step-runs-that-command-s.md, design/adrs/README.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/settings.md
- 2: design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md, design/issues/summary.md
