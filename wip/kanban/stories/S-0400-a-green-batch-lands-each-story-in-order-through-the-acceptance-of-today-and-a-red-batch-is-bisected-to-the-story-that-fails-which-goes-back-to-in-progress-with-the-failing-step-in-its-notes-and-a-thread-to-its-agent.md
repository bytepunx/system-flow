---
id: S-0400
type: story
nature: feature
title: A green batch lands each story in order through the acceptance of today, and a red batch is bisected to the story that fails, which goes back to in-progress with the failing step in its notes and a thread to its agent
status: backlog
parent: E-0022
owner: alex
created: 2026-10-16T04:58:07Z
updated: 2026-10-17T16:03:13Z
transitions: []
tags: [cli, dashboard]
topics: [orchestration]
touches: [design/adrs, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/operators/index.md]
after: [S-0399]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h15m
  delivery: 2026-10-17T20:21:00Z
  basis: "Its own forecast of 1h15m; 8th in the pull order with an in-progress limit of 5, behind S-0395, S-0396, S-0397, S-0394, S-0241, S-0398 and S-0399."
  by: flai
  at: 2026-10-17T16:03:13Z
---
# S-0400 A green batch lands each story in order through the acceptance of today, and a red batch is bisected to the story that fails, which goes back to in-progress with the failing step in its notes and a thread to its agent

## Goal

A verified batch ends in acceptances or in a send-back, with no person in between. When the batch passes, each of its stories is accepted in order through the acceptance of today, so that its transitions, acceptance commit, archive, and evidence read as they do now. When it fails, the queue finds the story that breaks it, lands the stories before it, and sends that story back to its agent with the failing step.

## Acceptance criteria

- [ ] A green batch lands each story in queue order through the same flow as `flai accept`: merge, done, archive, commit, by whoever enqueued it, with the orchestrator's evidence under `### Accepted by the orchestrator` when it enqueued. Each story's done transition, acceptance commit, archive, and `flai stats` figures are as an acceptance by hand gives them, tested against a story accepted by hand.
- [ ] The acceptance takes the batch's verification in place of the story head's: the verified commit is the story's rebased head inside a batch that passed. An ADR, refining ADR-0093 and ADR-0032, records this and that the queue may send a story back.
- [ ] A red batch is bisected: the queue verifies main with the first half of the batch, then halves again, until it finds the first story with which the batch fails. The stories before it land; the stories after it stay enqueued for the next batch. Each bisection step is kept in the queue's state and shown.
- [ ] The failing story moves from review to in-progress, by `flai serve`, with a reason naming the failing tier and its first findings, appended to the story's notes as a send-back's is, and a thread on the story to its agent saying what failed and at which commit.
- [ ] When main moves while a batch is under verification, by an acceptance by hand, the batch is dropped and built again from the new main.
- [ ] `design/system/workflow.md` (review to done, review to in-progress), `design/system/strategic-agents.md` § Accepting a story, `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` describe landing and bisection.

## Tasks
- T-1642 An ADR records that the queue's acceptance takes a passing batch's verification for the story's rebased head, and that the queue may send a story back
- T-1643 The acceptance preview and flai accept take a passing batch's verification of the story's rebased head in place of the story's own verify
- T-1644 A red batch is bisected to the first story with which it fails, with each step kept in the queue's state
- T-1659 A green batch lands each story in queue order through the acceptance flow, as whoever enqueued it
- T-1660 The dashboard's queue panel shows each bisection step and each story's outcome: landed, sent back, or still enqueued
- T-1669 flai serve lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves
- T-1677 The workflow, the strategic agents' design, and the CLI, dashboard, and operator docs describe landing and bisection

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
