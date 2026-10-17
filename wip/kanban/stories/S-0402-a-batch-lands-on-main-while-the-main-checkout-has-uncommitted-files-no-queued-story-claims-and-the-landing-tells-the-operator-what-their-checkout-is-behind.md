---
id: S-0402
type: story
nature: feature
title: A batch lands on main while the main checkout has uncommitted files no queued story claims, and the landing tells the operator what their checkout is behind
status: backlog
parent: E-0022
owner: alex
created: 2026-10-16T04:58:32Z
updated: 2026-10-17T07:43:16Z
transitions: []
tags: [cli, dashboard]
topics: [orchestration]
touches: [design/adrs, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/users/flai.md]
after: [S-0400]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 45m
  delivery: 2026-10-17T13:01:00Z
  basis: "Its own forecast of 45m; 11th in the pull order with an in-progress limit of 5, behind S-0395, S-0396, S-0397, S-0393, S-0394, S-0241, S-0398, S-0399, S-0400 and S-0401."
  by: flai
  at: 2026-10-17T07:43:16Z
---
# S-0402 A batch lands on main while the main checkout has uncommitted files no queued story claims, and the landing tells the operator what their checkout is behind

## Goal

Files the operator or the analyzer left uncommitted in the main checkout no longer stop an acceptance for a day. A batch lands on main while the checkout is dirty, so long as no queued story's change would overwrite those files, and the operator is told what their checkout is behind and how to catch it up, instead of being asked to commit or stash first.

## Acceptance criteria

- [ ] An ADR decides how a batch lands on main while the main checkout has uncommitted files, and what the checkout is left as, with the alternatives weighed: among them landing through the batch's worktree and moving main's branch with git's own check, which refuses only a file the change would overwrite.
- [ ] Uncommitted files in the main checkout that no story of the batch changes stop no landing, tested with a dirty analysis report and issue files, as on TH-0462.
- [ ] An uncommitted file that a story of the batch changes stops that story's landing only: the queue shows the file and the story, and lands the rest.
- [ ] After a landing that left the checkout's files behind main, `flai queue`, `flai board`, and the dashboard name the commits and paths the checkout is behind and the command that catches it up, until it has.
- [ ] `flai accept` by hand keeps its refusal and `--yes` as today. `design/system/workflow.md`, `docs/operators/index.md`, and `docs/users/flai.md` describe the landing on a dirty checkout.

## Tasks
- T-1647 An ADR decides how a batch lands on main while the main checkout has uncommitted files, and what the checkout is left as
- T-1648 storygit tells which uncommitted files in the main checkout a change between two commits would overwrite
- T-1662 A batch lands while the main checkout is dirty, stops only the story whose change would overwrite a dirty file, and records what the checkout is behind
- T-1671 flai queue, flai serve's log, and the dashboard name what the main checkout is behind and the command that catches it up
- T-1679 The workflow, the CLI and dashboard designs, and the operator and user docs describe landing on a dirty checkout and catching it up

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
