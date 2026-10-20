---
id: S-0398
type: story
nature: feature
title: flai queue lists the stories in review that would batch, in the order they entered review, and why each other story would not, and the board, the MCP board, and the dashboard show the queue
status: done
parent: E-0022
owner: alex
created: 2026-10-16T04:57:36Z
updated: 2026-10-19T18:58:59Z
transitions:
  - to: ready
    at: 2026-10-19T04:52:30Z
    by: alex
  - to: in-progress
    at: 2026-10-19T14:34:41Z
    by: agent-S-0398
  - to: review
    at: 2026-10-19T18:57:39Z
    by: agent-S-0398
  - to: done
    at: 2026-10-19T18:58:59Z
    by: orchestrator
tags: [cli, dashboard]
topics: [orchestration]
touches: [design/adrs, design/system/workflow.md, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/root.go, flai/cmd/board.go, flai/cmd/board_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go, flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, flai/internal/mcpserver/folder.go, design/components/flaiover/src/lib/components/BoardCard.examples.json, design/components/flaiover/src/lib/components/MergeQueue.examples.json, flaiover/src/lib/server/board.ts, design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4217
  priced_by: harness
  turns:
    - day: 2026-10-19
      ceremony: 5
      hand_edits: 3
      work: 60
  models:
    - model: claude-opus-5-5
      input: 532
      output: 176893
      cache_read: 26965234
      cache_write: 1558850
      cost: 18.7101
  strategic:
    - kind: orchestrator
      seconds: 10803
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 224
          output: 3834
          cache_read: 29309384
          cache_write: 70247
          cost: 7.3658
        - model: claude-sonnet-5-5
          input: 26
          output: 130
          cache_read: 428837
          cache_write: 137234
          cost: 0.4633
forecast:
  duration: 47m
  delivery: 2026-10-19T05:54:00Z
  basis: "Its own forecast of 47m; 5th in the pull order with an in-progress limit of 5, behind S-0449, S-0404, S-0405, S-0403 and S-0402."
  by: flai
  at: 2026-10-19T04:55:09Z
finalized:
  by: alex
  at: 2026-10-19T04:52:24Z
---
# S-0398 flai queue lists the stories in review that would batch, in the order they entered review, and why each other story would not, and the board, the MCP board, and the dashboard show the queue

## Goal

The operator sees, before anything is batched, which stories in review a merge queue would take and in what order, and why each other story in review would not be taken. Nothing is verified, merged, or moved: this story is the queue's state and its view, read-only, and the ADR that the rest of E-0022 builds on.

## Acceptance criteria

- [x] `flai queue`, and `flai queue --json`, list the stories in review that are eligible, in the order they last entered review: the last `flai verify` passed at the head of the story's branch, every acceptance criterion is ticked, no thread on the story or its tasks is open, and no path is recorded as undoing main's change (ADR-0155). Each other story in review is listed after them with every reason it is not eligible, each by name.
- [x] Under `acceptance.mode: pull-request` (ADR-0165), `flai queue` says the queue does not run in that mode, because the forge's required checks gate each story (ADR-0169), and lists no story as eligible. `flai board`, the MCP `board`, and the dashboard say the same in place of the queue.
- [x] The queue's state is one file in the main checkout's `.flai-cache`, read through one package: the stories enqueued, with who enqueued each and when, the batch under verification and its step, and each story's last outcome. This story writes nothing to it; with no file, the queue shows no batch and no outcomes.
- [x] `flai board`, the MCP `board`, and the host read `board.get` carry the queue: the eligible stories in order, the reasons of the others, and the batch and its step. The dashboard's board marks each card in review as eligible or not, with its reasons, and shows the batch and its step above the review lane.
- [x] An ADR records the queue as TH-0465 settles it: what runs it, the host action and the permission it runs under, that it runs under local acceptance only, eligibility, the batch's upper size, and that a story undoing main's change is never batched. `design/system/workflow.md` § What waits on an acceptance, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, and `docs/users/flai.md` describe `flai queue` and the view.

## Tasks
- T-1638 An ADR records the merge queue: what runs it, its host action and permission, eligibility, the batch's upper size, and that a story undoing main's change is never batched
- T-1639 The mergequeue package works out which stories in review are eligible, in the order they entered review, with each other story's reasons, and reads the queue's state
- T-1655 flai queue lists the eligible stories in order, the others with their reasons, and the batch and its step, as text and as JSON
- T-1656 flai board, the MCP board, and the host read board.get carry the queue
- T-1667 The dashboard's board marks each card in review as eligible or not with its reasons, and shows the queue's batch and step above the review lane
- T-1675 The workflow, the CLI and dashboard designs, and the user guides describe flai queue and the queue on the board

## Notes

- Drafted by the planner for E-0022, first of its stories. The rest of the epic writes the state this story reads: S-0399 enqueues and batches, and S-0400 lands and bisects.
- The design answers it rests on are asked on TH-0465 (questions 2 to 4). The ADR records the answers given there.
- The queue runs under local acceptance only, as proposed on TH-0466 after S-0395 was accepted. ADR-0169 blocks every acceptance under pull-request acceptance on `forge_stale` when the kept forge state was read at a commit other than the branch's head. A batch rebases each story, so a queue landing would always stop there. In that mode the forge's CI already verifies each story, and S-0397 merges through the remote.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Design: `design/system/workflow.md`, `flai-cli.md`, `flaiover-dashboard.md`, and `docs/users/flai.md`, which the epic's last criterion names.
  - Layout: the new package `flai/internal/mergequeue` and `flai/cmd/queue.go`. The board's three surfaces, as the explorer found them: `flai/cmd/board.go`, `flai/internal/mcpserver/server.go`, and `board.get` in `flai/internal/hostapi/hostapi.go`. The dashboard board and its cards.
  - Co-change: `flai touches suggest` added `docs/users/flai-reference.md`, regenerated when a command is added, and `docs/users/flaiover.md`, the dashboard's guide.
  - Folder touch kept: `design/adrs`, because the ADR's file name is known only once `flai adr new` writes it. It is a shared path, so it holds no story.
  - The local-only rule reads `manifest.Acceptance.PullRequest()`, which S-0395's line of E-0021 added. It needs no new touch.
- **Forecast.** flai's 47m is kept: the view reads what flai already reads, through one new package. The local-only rule is one check and a message.
- **Cost of delay.** Waits on TH-0465 question 1: the epic has no inputs, so `flai cod` refuses. It is written once the inputs are set.

### Accepted by the orchestrator

- Verified: 0414c5e0b2e53cfdda7a7e9427776bd181651cfe
- At: 2026-10-19T18:58:59Z

Verdict: every criterion 1 to 5 is met by changed files and ticked; ADR-0173 is accepted as TH-0516 settled it; flai verify passed every step at head 0414c5e.
- 1: flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/root.go, docs/users/flai-reference.md
- 2: flai/internal/mergequeue/eligible.go, flai/cmd/queue.go, flai/cmd/board.go, flai/cmd/board_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts
- 3: flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go
- 4: flai/cmd/board.go, flai/cmd/board_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go, flaiover/src/lib/server/board.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/BoardCard.examples.json, design/components/flaiover/src/lib/components/MergeQueue.examples.json
- 5: design/adrs/0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md, design/adrs/README.md, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md
