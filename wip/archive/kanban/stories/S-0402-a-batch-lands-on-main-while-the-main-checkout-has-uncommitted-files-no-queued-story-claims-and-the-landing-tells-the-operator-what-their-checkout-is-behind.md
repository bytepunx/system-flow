---
id: S-0402
type: story
nature: feature
title: A batch lands on main while the main checkout has uncommitted files no queued story claims, and the landing tells the operator what their checkout is behind
status: done
parent: E-0022
owner: alex
created: 2026-10-16T04:58:32Z
updated: 2026-10-20T05:05:10Z
transitions:
  - to: ready
    at: 2026-10-19T04:53:22Z
    by: alex
  - to: in-progress
    at: 2026-10-19T23:52:18Z
    by: agent-S-0402
  - to: review
    at: 2026-10-20T05:03:39Z
    by: agent-S-0402
  - to: done
    at: 2026-10-20T05:05:10Z
    by: orchestrator
tags: [cli, dashboard]
topics: [orchestration]
touches: [design/adrs, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/users/flai.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/issues/I-0168-a-story-asked-to-share-a-path-agrees-in-a-settled-reply-but-records-no-share-so-the-asking-story-s-task-stays-refused-on-the-overlap.md, flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept_markdown.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/progress.go, design/components/flaiover/src/lib/components/MergeQueue.examples.json, flai/cmd/board.go, flai/cmd/queue_held_test.go]
after: [S-0400]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3295
  priced_by: harness
  turns:
    - day: 2026-10-19
      test_runs: 1
      hand_edits: 1
      work: 44
    - day: 2026-10-20
      ceremony: 5
      test_runs: 5
      hand_edits: 3
      work: 83
  models:
    - model: claude-opus-5-5
      input: 398
      output: 133608
      cache_read: 36623031
      cache_write: 745154
      cost: 15.5044
  strategic:
    - kind: orchestrator
      seconds: 27137
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 160
          output: 3486
          cache_read: 35283062
          cache_write: 42092
          cost: 8.8561
        - model: claude-sonnet-5-5
          input: 18
          output: 83
          cache_read: 355259
          cache_write: 65519
          cost: 0.3444
forecast:
  duration: 45m
  delivery: 2026-10-20T00:30:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 5, behind S-0404 and S-0403."
  by: flai
  at: 2026-10-19T23:41:20Z
finalized:
  by: alex
  at: 2026-10-19T04:53:13Z
---
# S-0402 A batch lands on main while the main checkout has uncommitted files no queued story claims, and the landing tells the operator what their checkout is behind

## Goal

Files the operator or the analyzer left uncommitted in the main checkout no longer stop an acceptance for a day. A batch lands on main while the checkout is dirty, so long as no queued story's change would overwrite those files, and the operator is told what their checkout is behind and how to catch it up, instead of being asked to commit or stash first.

## Acceptance criteria

- [x] An ADR decides how a batch lands on main while the main checkout has uncommitted files, and what the checkout is left as, with the alternatives weighed: among them landing through the batch's worktree and moving main's branch with git's own check, which refuses only a file the change would overwrite.
- [x] Uncommitted files in the main checkout that no story of the batch changes stop no landing, tested with a dirty analysis report and issue files, as on TH-0462.
- [x] An uncommitted file that a story of the batch changes stops that story's landing only: the queue shows the file and the story, and lands the rest.
- [x] `flai queue`, `flai board`, and the dashboard name each story a dirty file in the main checkout holds back, the files, and the command that frees them (commit them, or `git stash push -- <files>`), until they are clean. The held story stays enqueued and lands with the next batch.
- [x] `flai accept` by hand keeps its refusal and `--yes` as today. `design/system/workflow.md`, `docs/operators/index.md`, and `docs/users/flai.md` describe the landing on a dirty checkout.

## Tasks
- T-1647 An ADR decides how a batch lands on main while the main checkout has uncommitted files, and what the checkout is left as
- T-1648 storygit tells which uncommitted files in the main checkout a change between two commits would overwrite
- T-1662 A batch lands while the main checkout is dirty, stops only the story whose change would overwrite a dirty file, and records what the checkout is behind
- T-1671 flai queue, flai serve's log, and the dashboard name what the main checkout is behind and the command that catches it up
- T-1679 The workflow, the CLI and dashboard designs, and the operator and user docs describe landing on a dirty checkout and catching it up
- T-1911 flai accept --batch leaves the main checkout's uncommitted files outside wip out of its refusal, its lint fix, and its commit

## Notes

- Drafted by the planner for E-0022, after S-0400, whose landing this story changes.
- TH-0462 is the case: S-0391 waited because the analyzer's report, its issues, and a design system's files were uncommitted on main, none of them S-0391's.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - the dirty-checkout check: `storygit.DirtyOutsideWip` in `flai/internal/storygit/storygit.go`, which the preview reads, and `flai/cmd/accept.go`, which refuses unless `--yes`
    - the landing: `flai/internal/mergequeue/land.go`, from S-0400
    - the view: `flai/cmd/queue.go` and `MergeQueue.svelte`.
  - Design: the epic's criterion 6.
  - Co-change: `flai touches suggest` added `design/system/flai-cli.md` and `design/system/flaiover-dashboard.md`.
  - `flai/internal/storygit/storygit_test.go` is a new file; the package's tests are split by topic today.
  - Folder touch kept: `design/adrs`, for the new ADR.
- **Forecast.** Raised from flai's 34m to 45m. The ADR weighs three ways to land on a dirty checkout, and the tests cover untracked and overwritten files.
- **Cost of delay.** Waits on TH-0465 question 1.

### Criterion 4

- Reworded on TH-0521, which the operator confirmed on 2026-10-20: a landing stays in the main checkout with git's own check, so it never leaves the checkout's files behind main. Criterion 4 was "After a landing that left the checkout's files behind main, `flai queue`, `flai board`, and the dashboard name the commits and paths the checkout is behind and the command that catches it up, until it has."

### Accepted by the orchestrator

- Verified: edad54dadc5b068042b519052252121d3a93f977
- At: 2026-10-20T05:05:10Z

Verdict: every criterion 1 to 5, as reworded on TH-0521, is met by changed files and ticked. ADR-0178 is accepted and records the landing in the main checkout with git's own check, so main never moves without the checkout. A story whose change would overwrite a dirty file stays enqueued and is named with the files and the freeing command. flai verify passed every tier at head edad54d.
- 1: design/adrs/0178-the-merge-queue-lands-each-story-in-the-main-checkout-with-git-s-own-check.md, design/adrs/README.md
- 2: flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_markdown.go, flai/cmd/accept_batch_dirty_test.go, flai/internal/mergequeue/land_test.go
- 3: flai/internal/mergequeue/land.go, flai/internal/mergequeue/progress.go, flai/internal/storygit/storygit.go, flai/internal/mergequeue/land_test.go, flai/internal/storygit/storygit_test.go
- 4: flai/internal/mergequeue/eligible.go, flai/cmd/queue.go, flai/cmd/board.go, flai/cmd/queue_held_test.go, flai/cmd/queue_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json
- 5: flai/cmd/accept.go, flai/cmd/accept_batch_dirty_test.go, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/users/flai.md
