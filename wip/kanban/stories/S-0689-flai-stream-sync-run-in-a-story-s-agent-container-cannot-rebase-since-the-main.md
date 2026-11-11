---
id: S-0689
type: story
nature: remediation
title: flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached
status: backlog
owner: alex
created: 2026-11-10T08:49:00Z
updated: 2026-11-11T08:46:41Z
transitions: []
tags: [flai, git, container]
topics: [git]
touches: [flai/cmd/verify.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/cmd/branch.go, flai/cmd/branch_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, docs/operators/runbooks/agent-container.md, design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-10T08:49:00Z
  value: 30
  by: planner-S-0689
  at: 2026-11-10T10:34:16Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T19:35:00Z
  basis: "Its own forecast of 1h30m; 13th in the pull order with an in-progress limit of 3, behind S-0693, S-0698, S-0702, S-0694, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663 and S-0683."
  by: flai
  at: 2026-11-11T08:46:41Z
---
# S-0689 flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached

## Goal

This story remediates [I-0272](../../../design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md), "flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0272 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0272 is closed with `flai issue close I-0272 --reason` saying what fixed it

## Tasks
- T-3026 flai stream sync and the shell's flai task done refuse at once in a story's agent container, naming the MCP tools that sync on the host
- T-3027 A story's agent started on the container harness is told to close tasks and sync through the MCP tools task_done and verify, not flai stream sync in its shell
- T-3030 A sync whose rebase stops with no conflict says git's error and offers no abort that cannot restore the branch, and a trial merge that cannot run is reported as not run
- T-3033 The container harness's docs say that flai stream sync refuses there and the MCP tools task_done and verify sync on the host, and I-0272 is closed

## Notes

Cost of delay inputs set by flai from I-0272. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T04:41:52Z, 3.2 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0272's two instances. ADR-0246 says `flai stream sync` refuses in a story's container, but the refusal keys on `MainAbsent`, which ADR-0253's mount of `system-flow.yaml` leaves empty, so the sync runs against the read-only `.git`. The fix has three parts:

- The shell's sync refuses there at once, naming the MCP tools `task_done` and `verify`, which sync on the host (T-3026).
- The container harness's start prompt names those tools in place of the shell's sync (T-3027).
- Wherever a sync runs, a rebase that stops with no conflict says git's error with no false abort promise, and a trial merge that cannot run is reported as not run, never as clean (T-3030).

No MCP tool syncs alone; adding one is left out of this story.

Touches, all files, no folder:

- From the code layout, found by the explorer: `flai/cmd/verify.go` (the container helper beside `orchestratorInRoleContainer`), `flai/cmd/stream.go`, `flai/cmd/stream_sync.go`, `flai/cmd/stream_sync_test.go`, `flai/cmd/task_done.go`, `flai/cmd/task_done_test.go`, `flai/cmd/branch.go`, `flai/cmd/branch_test.go`, `flai/internal/storygit/sync.go`, `flai/internal/storygit/sync_test.go`, `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `flai/internal/harness/container.go`, and `flai/internal/harness/container_test.go`.
- From the design and docs that describe the container harness: `docs/users/flai.md`, `docs/users/flai-reference.md` (regenerated from the help text), `design/system/flai-cli.md`, and `docs/operators/runbooks/agent-container.md`.
- From the criteria: the I-0272 file and `design/issues/summary.md`, which `flai issue close` writes.
- The story declared none. `flai touches suggest` listed co-changed paths at 10 percent or less, none tied to this change, so none was added.

Forecast: `flai forecast` gave 39m. It is raised to 1h30m for four tasks across three packages, a reproduction on a read-only git common dir, and the integration and smoke tiers a `flai/` change selects at close-out, which took 8m14s and 20m44s in today's close-outs of S-0672 and S-0669. The delivery moves by the same 51m from flai's 17:22Z.

Cost of delay: 30 USD a week, as `flai cod` worked it from the operator's input of 12m lost per 168h cycle at 150 USD an hour. It stands: the issue's two instances give no reason to change it.

No `after`: the story reads no other open story's output.
