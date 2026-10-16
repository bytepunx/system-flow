---
id: S-0401
type: story
nature: feature
title: A batch settles a rebase that stops inside the manifest's shared paths by the regenerate-and-merge steps and never sends a story back for it, and tells the story's agent of a conflict outside them by message
status: backlog
parent: E-0022
owner: alex
created: 2026-10-16T04:58:20Z
updated: 2026-10-16T06:07:38Z
transitions: []
tags: [cli]
topics: [orchestration]
touches: [design/adrs, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, design/system/agent-coordination.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md]
after: [S-0399]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 50m
  delivery: 2026-10-16T12:28:00Z
  basis: "Its own forecast of 50m; 11th in the pull order with an in-progress limit of 5, behind S-0426, S-0427, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399 and S-0400."
  by: flai
  at: 2026-10-16T06:07:38Z
---
# S-0401 A batch settles a rebase that stops inside the manifest's shared paths by the regenerate-and-merge steps and never sends a story back for it, and tells the story's agent of a conflict outside them by message

## Goal

A batch never stops, and never sends a story back, because two stories both changed a file every story changes, such as `design/issues/summary.md`, `design/adrs/README.md`, or `design/system/flai-cli.md`. The queue settles such a stop itself, as `flai stream sync` and `flai accept` already settle some of them. A conflict outside those paths is real work: the queue tells the story's agent, as a trial merge does, and leaves the story out of the batch.

## Acceptance criteria

- [ ] The batch's rebase of each story runs the steps `flai stream sync` and `flai accept` run on a stop: `summary.md` written again (ADR-0098), issue files merged by their instances and folded (ADR-0126), and the ADR index merged row by row (ADR-0157). A stop they settle is continued, and the story stays in the batch.
- [ ] A stop on another path inside `claims.shared`, such as `design/system/flai-cli.md` or `docs/users/flai.md`, is settled as the ADR decides, and never sends the story back: the story stays enqueued and is shown waiting with the paths until it is settled. The ADR, refining ADR-0098, ADR-0126, and ADR-0157, records how, with the cases it cannot settle.
- [ ] A stop on a path outside `claims.shared` leaves the story out of the batch, and is told to its agent by a message in the conversation of the two stories whose changes conflict, by `flai`, `about` the paths, as ADR-0121's trial-merge conflict is. A conflict with main itself, the story ahead having landed, sends the story back with the paths as the reason.
- [ ] `design/system/workflow.md` § Branches and collisions, `design/system/agent-coordination.md`, and `docs/users/flai.md` describe what the queue settles and what it reports.

## Tasks
- T-1645 An ADR records how the queue settles a rebase that stops inside the shared paths, and what it does with a stop outside them
- T-1646 The batch's rebase of each story runs the summary, issue, and ADR index merge steps that flai stream sync runs on a stop
- T-1661 storygit merges a shared markdown file both sides changed as the ADR decides, and leaves a stop it cannot settle
- T-1670 The batch settles a shared-path stop with the markdown merge, keeps a story it cannot settle enqueued, and tells a conflict outside the shared paths by message
- T-1678 The workflow, the coordination design, and the CLI docs describe what the queue settles in the shared paths and what it reports outside them

## Notes

- Drafted by the planner for E-0022, after S-0399, whose batch rebase this story extends. It can run beside S-0400 only once their claims no longer overlap: both change `flai/internal/serve/mergequeue.go` and `design/system/workflow.md`.
- I-0150's third story stopped mid-rebase on `design/system/flai-cli.md`, a shared path that none of the three existing merge steps covers. That case is why the second criterion asks for more than the steps of today.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - the rebase and its merge steps: `flai/internal/storygit/sync.go` (`Rebase`, `continueOverGenerated`)
    - the markdown merge: a new `flai/internal/storygit/mdmerge.go`
    - the batch: `flai/internal/mergequeue/batch.go`
    - the driver: `flai/internal/serve/mergequeue.go`
    - the messages: `messages.Notify`, called and not changed.
  - Design: the epic's criterion 5, ADR-0098, ADR-0121, ADR-0126, and ADR-0157; `design/system/agent-coordination.md`, which the epic names.
  - Co-change: `flai touches suggest` ranked `design/system/flai-cli.md` first, at 59%. It is added for the sync's described behaviour.
  - Folder touch kept: `design/adrs`, for the new ADR.
- **Forecast.** Raised from flai's 28m to 50m. A three-way markdown merge by section is new logic with many edge cases, beyond reusing the existing steps.
- **Cost of delay.** Waits on TH-0465 question 1.
