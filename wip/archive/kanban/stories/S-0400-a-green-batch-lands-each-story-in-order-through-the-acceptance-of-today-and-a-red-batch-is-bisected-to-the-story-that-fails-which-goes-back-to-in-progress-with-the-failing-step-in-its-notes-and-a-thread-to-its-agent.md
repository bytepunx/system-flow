---
id: S-0400
type: story
nature: feature
title: A green batch lands each story in order through the acceptance of today, and a red batch is bisected to the story that fails, which goes back to in-progress with the failing step in its notes and a thread to its agent
status: done
parent: E-0022
owner: alex
created: 2026-10-16T04:58:07Z
updated: 2026-10-19T23:41:20Z
transitions:
  - to: ready
    at: 2026-10-19T04:53:06Z
    by: alex
  - to: in-progress
    at: 2026-10-19T21:53:50Z
    by: agent-S-0400
  - to: review
    at: 2026-10-19T23:39:42Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T23:41:20Z
    by: orchestrator
tags: [cli, dashboard]
topics: [orchestration]
touches: [design/adrs, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/progress.go, flai/internal/mergequeue/progress_test.go, flai/internal/mergequeue/state.go, docs/operators/settings.md, docs/users/flai-reference.md, flai/internal/preview/orchestrator.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [S-0399]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6368
  priced_by: harness
  turns:
    - day: 2026-10-19
      ceremony: 6
      test_runs: 4
      hand_edits: 3
      work: 89
  models:
    - model: claude-haiku-5-5
      input: 72
      output: 29652
      cache_read: 4386942
      cache_write: 201457
      cost: 0.338
    - model: claude-opus-5-5
      input: 652
      output: 325675
      cache_read: 47148813
      cache_write: 1342018
      cost: 23.7037
  strategic:
    - kind: orchestrator
      seconds: 1090
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 35
          output: 595
          cache_read: 5582303
          cache_write: 13030
          cost: 1.4028
        - model: claude-sonnet-5-5
          input: 5
          output: 22
          cache_read: 98504
          cache_write: 21267
          cost: 0.098
forecast:
  duration: 1h15m
  delivery: 2026-10-19T23:14:00Z
  basis: "Its own forecast of 1h15m; 4th in the pull order with an in-progress limit of 5, behind S-0401, S-0404, S-0403 and S-0402."
  by: flai
  at: 2026-10-19T21:52:13Z
finalized:
  by: alex
  at: 2026-10-19T04:52:58Z
---
# S-0400 A green batch lands each story in order through the acceptance of today, and a red batch is bisected to the story that fails, which goes back to in-progress with the failing step in its notes and a thread to its agent

## Goal

A verified batch ends in acceptances or in a send-back, with no person in between. When the batch passes, each of its stories is accepted in order through the acceptance of today, so that its transitions, acceptance commit, archive, and evidence read as they do now. When it fails, the queue finds the story that breaks it, lands the stories before it, and sends that story back to its agent with the failing step.

## Acceptance criteria

- [x] A green batch lands each story in queue order through the same flow as `flai accept`: merge, done, archive, commit, by whoever enqueued it, with the orchestrator's evidence under `### Accepted by the orchestrator` when it enqueued. Each story's done transition, acceptance commit, archive, and `flai stats` figures are as an acceptance by hand gives them, tested against a story accepted by hand.
- [x] The acceptance takes the batch's verification in place of the story head's: the verified commit is the story's rebased head inside a batch that passed. An ADR, refining ADR-0093 and ADR-0032, records this and that the queue may send a story back.
- [x] A red batch is bisected: the queue verifies main with the first half of the batch, then halves again, until it finds the first story with which the batch fails. The stories before it land; the stories after it stay enqueued for the next batch. Each bisection step is kept in the queue's state and shown.
- [x] The failing story moves from review to in-progress, by `flai serve`, with a reason naming the failing tier and its first findings, appended to the story's notes as a send-back's is, and a thread on the story to its agent saying what failed and at which commit.
- [x] When main moves while a batch is under verification, by an acceptance by hand, the batch is dropped and built again from the new main.
- [x] `design/system/workflow.md` (review to done, review to in-progress), `design/system/strategic-agents.md` § Accepting a story, `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` describe landing and bisection.

## Tasks
- T-1642 An ADR records that the queue's acceptance takes a passing batch's verification for the story's rebased head, and that the queue may send a story back
- T-1643 The acceptance preview and flai accept take a passing batch's verification of the story's rebased head in place of the story's own verify
- T-1644 A red batch is bisected to the first story with which it fails, with each step kept in the queue's state
- T-1659 A green batch lands each story in queue order through the acceptance flow, as whoever enqueued it
- T-1660 The dashboard's queue panel shows each bisection step and each story's outcome: landed, sent back, or still enqueued
- T-1669 flai serve lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves
- T-1677 The workflow, the strategic agents' design, and the CLI, dashboard, and operator docs describe landing and bisection
- T-1909 The batch records each story's verified and rebased heads, and the queue's state holds a batch's bisection and landing

## Notes

- Drafted by the planner for E-0022, after S-0399, which builds and verifies the batch this story lands.
- Today only a person moves a story from review to in-progress. This story makes `flai serve` a second maker of that move, for the queue's send-back only, and the ADR says so.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - the acceptance's `--verified` check and blockers: `flai/internal/preview/accept.go` and `flai/cmd/accept.go`
    - landing and bisection: new files in `flai/internal/mergequeue`
    - the driver: `flai/internal/serve/mergequeue.go`, from S-0399.
  - Design: the epic's criterion 4; `design/system/workflow.md` § Transitions and who makes them; `design/system/strategic-agents.md` § Accepting a story, which ADR-0093 governs.
  - Co-change: `flai touches suggest` added `design/system/flaiover-dashboard.md`, for the panel's steps.
  - Folder touch kept: `design/adrs`, for the new ADR, as in S-0398.
- **Forecast.** Raised from flai's 40m to 1h15m. It changes the acceptance, which about fifteen test files under `flai/cmd` guard, and it must show that each landing matches an acceptance by hand. It adds bisection and a send-back besides.
- **Cost of delay.** Waits on TH-0465 question 1.

### Accepted by the orchestrator

- Verified: 728368876184a20accf3534acc397e51024b94d5
- At: 2026-10-19T23:41:20Z

Verdict: every criterion 1 to 6 is met by changed files and ticked. A green batch lands through flai accept as its enqueuer, with ADR-0093's checks and ADR-0155's keep refusal intact. A red batch is bisected and its failing story sent back with notes and a thread, and the queue does nothing with the merge-queue action off. ADR-0176 is accepted. flai verify passed every tier at head 7283688.
- 1: flai/internal/mergequeue/land.go, flai/internal/serve/mergequeue.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go
- 2: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md, design/adrs/README.md, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/mergequeue/land.go, flai/internal/preview/accept_test.go
- 3: flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/progress.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/batch.go, flaiover/src/lib/components/MergeQueue.svelte, flai/internal/mergequeue/bisect_test.go, flaiover/src/lib/components/MergeQueue.svelte.test.ts
- 4: flai/internal/mergequeue/bisect.go, flai/internal/serve/agents.go, flai/internal/harness/harness.go
- 5: flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go
- 6: design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md
