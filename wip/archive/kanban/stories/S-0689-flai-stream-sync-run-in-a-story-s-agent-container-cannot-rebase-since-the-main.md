---
id: S-0689
type: story
nature: remediation
title: flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached
status: done
owner: alex
created: 2026-11-10T08:49:00Z
updated: 2026-11-11T21:55:28Z
transitions:
  - to: ready
    at: 2026-11-11T08:53:38Z
    by: alex
  - to: in-progress
    at: 2026-11-11T14:50:40Z
    by: agent-S-0689
  - to: review
    at: 2026-11-11T16:13:51Z
    by: agent-S-0689
  - to: in-progress
    at: 2026-11-11T18:17:00Z
    by: flai
  - to: review
    at: 2026-11-11T18:28:58Z
    by: agent-S-0689
  - to: done
    at: 2026-11-11T21:55:28Z
    by: alex
tags: [flai, git, container]
topics: [git]
touches: [flai/cmd/verify.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/cmd/branch.go, flai/cmd/branch_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/container_test.go, flai/internal/workitem/repo.go, flai/internal/taskdone/taskdone_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/agent-adapters.md, docs/operators/runbooks/agent-container.md, design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md, design/issues/summary.md, flai/internal/taskdone/taskdone.go, design/issues/I-0351-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5727
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 4
      work: 81
  models:
    - model: claude-opus-5-5
      input: 458
      output: 171076
      cache_read: 40850477
      cache_write: 1900225
      cost: 23.5046
  strategic:
    - kind: planner
      seconds: 22
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 1358
          cache_read: 323499
          cache_write: 8355
          cost: 0.1553
    - kind: orchestrator
      seconds: 52
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 125
          cache_read: 653765
          cache_write: 119453
          cost: 0.2173
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
  delivery: 2026-11-11T17:12:00Z
  basis: "Its own forecast of 1h30m; 3rd in the pull order with an in-progress limit of 3, behind S-0694, S-0710, S-0719, S-0724 and S-0716."
  by: flai
  at: 2026-11-11T14:23:13Z
finalized:
  by: alex
  at: 2026-11-11T08:53:19Z
---
# S-0689 flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached

## Goal

This story remediates [I-0272](../../../design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md), "flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0272 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0272 is closed with `flai issue close I-0272 --reason` saying what fixed it

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

No MCP tool syncs alone; adding one is left out of this story. No `stream_sync` MCP tool exists, yet `RequireMain`'s error and three documents name one. The replan of 2026-11-11 gives dropping it from the error to T-3026 and from the documents to T-3033.

Touches, all files, no folder:

- From the code layout, found by the explorer: `flai/cmd/verify.go` (the container helper beside `orchestratorInRoleContainer`), `flai/cmd/stream.go`, `flai/cmd/stream_sync.go`, `flai/cmd/stream_sync_test.go`, `flai/cmd/task_done.go`, `flai/cmd/task_done_test.go`, `flai/cmd/branch.go`, `flai/cmd/branch_test.go`, `flai/internal/storygit/sync.go`, `flai/internal/storygit/sync_test.go`, `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `flai/internal/harness/container.go`, and `flai/internal/harness/container_test.go`.
- From the code layout, added at the replan: `flai/internal/workitem/repo.go`, whose `RequireMain` error names the `stream_sync` tool, and `flai/internal/taskdone/taskdone_test.go`, which pins that text.
- From the design and docs that describe the container harness: `docs/users/flai.md`, `docs/users/flai-reference.md` (regenerated from the help text), `design/system/flai-cli.md`, and `docs/operators/runbooks/agent-container.md`. Added at the replan: `design/system/agent-adapters.md`, whose MCP server row and open question name `stream_sync`.
- From the criteria: the I-0272 file and `design/issues/summary.md`, which `flai issue close` writes.
- The story declared none. `flai touches suggest` listed co-changed paths at 10 percent or less, none tied to this change, so none was added.

Since the first plan, S-0688, S-0679, S-0686, S-0674, and S-0665 changed `flai/internal/harness/container.go` and `flai/cmd/task_done.go`. None changes where the sync refuses or what the container prompt says, so the tasks stand.

Forecast: `flai forecast` gave 40m. It is raised to 1h30m for four tasks across four packages, a reproduction on a read-only git common dir, and the integration and smoke tiers a `flai/` change selects at close-out, which took 8m14s and 20m44s in the close-outs of S-0672 and S-0669. The two added files are small text changes and add nothing. The delivery moves by the same 50m from flai's 11:44Z.

Cost of delay: 30 USD a week, as `flai cod` worked it from the operator's input of 12m lost per 168h cycle at 150 USD an hour. It stands: the issue's two instances give no reason to change it.

No `after`: the story reads no other open story's output.
- 2026-11-11T18:17:00Z: moved to in-progress: the rebase of story/S-0689 onto main stopped on conflicts in design/system/flai-cli.md, which the acceptance of S-0689 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0689
