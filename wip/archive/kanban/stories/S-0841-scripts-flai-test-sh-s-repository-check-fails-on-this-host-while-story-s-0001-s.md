---
id: S-0841
type: story
nature: remediation
title: scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged
status: done
owner: alex
created: 2026-11-17T08:32:52Z
updated: 2026-11-17T12:57:32Z
transitions:
  - to: ready
    at: 2026-11-17T12:39:16Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T12:53:43Z
    by: agent-S-0841
  - to: review
    at: 2026-11-17T12:56:07Z
    by: agent-S-0841
  - to: done
    at: 2026-11-17T12:57:32Z
    by: orchestrator
tags: [issues, testing]
topics: [testing]
touches: [design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 154
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 3
      work: 20
  models:
    - model: claude-opus-5-5
      input: 52
      output: 8730
      cache_read: 3631784
      cache_write: 191284
      cost: 2.4314
  strategic:
    - kind: planner
      seconds: 60
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 2
          output: 5
          cache_read: 251314
          cache_write: 6008
          cost: 0.0176
        - model: claude-opus-5-5
          input: 45
          output: 5867
          cache_read: 3834449
          cache_write: 52718
          cost: 1.3929
    - kind: orchestrator
      seconds: 44
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 93
          cache_read: 990545
          cache_write: 30423
          cost: 0.3197
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-17T08:32:52Z
  value: 12.5
  by: planner-S-0841
  at: 2026-11-17T12:36:04Z
forecast:
  duration: 11m
  delivery: 2026-11-17T12:56:00Z
  basis: "Its own forecast of 11m; 1st in the pull order with an in-progress limit of 3, behind S-0806, S-0834 and S-0838."
  by: flai
  at: 2026-11-17T12:39:22Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:38:54Z
---
# S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged

## Goal

This story remediates [I-0370](../../../design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md), "scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0370 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0370 is closed with `flai issue close I-0370 --reason` saying what fixed it

## Tasks
- T-3593 Confirm S-0781's fix keeps I-0370's cause from occurring, and record the evidence in the story's Notes
- T-3594 Close I-0370 with a reason naming S-0781's fix and the branches deleted on TH-1190

## Notes

Cost of delay inputs set by flai from I-0370. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:39:18Z, 3.5 days before this story; under one cycle counts as one).

### Planning

Proposed solution: build nothing new. S-0781 already fixed I-0370's cause.

- The four branches were never real story branches. They came from S-0777's discarded test run, which set TMPDIR inside its worktree. Test temp projects with no `git init` of their own let git walk up, so `flai story start` in the tests made `story/S-0001`, `S-0002`, `S-0003`, and `S-0005` in the main repository (D-0617).
- S-0781 (I-0371) set `GIT_CEILING_DIRECTORIES` in the TestMain of `flai/cmd` and `flai/internal/mcpserver`. It added `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch` in `flai/cmd/story_start_test.go`, which reproduces the leak. That test meets criterion 1's "with a test that reproduces it".
- The operator approved deleting the four branches on TH-1190 (D-0625). On 2026-11-17 the main checkout has no `story/S-000*` branch, and its unscoped `flai check --strict` reports no `story.unaccepted`.
- The planner's explorer checked the other test packages that run git: adr, check, conflictmark, hostapi, itemedit, mergequeue, pending, preview, serve, statsread, storygit, taskdone, threads, verify, and workitem. In each, every test that writes refs runs `git init` in its own temp folder first. Only read-only probes run in temp folders that were never initialised (pending, statsread, storygit), and those write no branch.

So the story verifies the fix (T-3593), then closes the issue (T-3594).

The unscoped check still fails on this host for a reason outside I-0370: a `threads.archived` on TH-1232, which is open on the archived S-0748. Only the operator or the thread's opener clears it. This story leaves it alone.

Touches, both from the design: criterion 2's `flai issue close` changes the issue file and regenerates the summary.

- `design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md`: the issue the story closes.
- `design/issues/summary.md`: flai regenerates it on close.

`flai touches suggest` listed co-change paths: `flai/cmd/cmd_test.go`, `flai/cmd/story_start.go`, `flai/cmd/task_done_test.go`, `flai/internal/mcpserver/story_start_test.go`, and `flai/internal/storystart/start_test.go`. All of them come from S-0781's own fix. None is kept, since no change to them is planned. If T-3593 finds a package that still leaks, it adds that package's test files to its own touches. No folder touch is kept. The story names nothing in `after`.

Tasks: T-3593 in layer 1, then T-3594 in layer 2, after T-3593, because the closing reason cites T-3593's verification.

Forecast 11m, kept as `flai forecast` gave it. The figure is the median of 24 done remediation stories, times size 4 (2 criteria, 2 touches). The story writes no code, and those medians already include a close-out. Delivery 2026-11-17T20:32:00Z, as flai gave it from the pull order.

Cost of delay 12.50 USD a week, as `flai cod` worked it out from flai's input: 5m of time lost per 168h cycle at 150 USD an hour. The input is the operator's, so it stands. Since S-0781, the cause no longer recurs, so the delay costs less than the figure says. Changing the input is the operator's call.

### Verification (T-3593)

What the story's agent observed on 2026-11-17, in the main checkout and on `story/S-0841`:

| Check | Result |
|-------|--------|
| `git branch --list 'story/S-000*'` in the main checkout | none |
| `git worktree list`, entries marked `prunable` | none |
| `scripts/flai.sh check --strict`, unscoped, in the main checkout | 0 `story.unaccepted`; exit 1 only on `threads.archived` for TH-1232, outside I-0370 |
| `GIT_CEILING_DIRECTORIES` in the test setup | set in `flai/cmd/cmd_test.go` (TestMain) and `flai/internal/mcpserver/server_test.go` (TestMain), both from S-0781's commit 637bed28f |
| Reproducer test | `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch` in `flai/cmd/story_start_test.go`: runs `flai story start` in a temporary project inside an outer repository and fails if the outer repository gains a `story/*` branch; it needs real git, so it skips under `-short` and runs in the integration tier |

S-0781 was accepted at cafa1f05c, so every close-out since that touches `flai/` runs the reproducer in its integration tier. The cause no longer occurs, and no package needs a change.

### Accepted by the orchestrator

- Verified: 8e81a39f21f0613a9a1f47115925dd466120520a
- At: 2026-11-17T12:57:32Z

Verdict: pass. Both criteria are met at 8e81a39f21f0613a9a1f47115925dd466120520a; S-0781's fix already removed I-0370's cause and the story closes the issue.
- 1: design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md
- 2: design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md, design/issues/summary.md
