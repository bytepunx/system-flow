---
id: S-0823
type: story
nature: remediation
title: A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back
status: backlog
owner: alex
created: 2026-11-17T07:13:19Z
updated: 2026-11-17T08:32:48Z
transitions: []
tags: [flai]
touches: [flai/internal/sendback/sendback.go, flai/internal/sendback/sendback_test.go, flai/cmd/move.go, flai/cmd/move_sendback_test.go, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0397-a-story-the-operator-moves-from-review-to-in-progress-by-hand-gets-no-agent-and.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 35
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 8
          output: 33
          cache_read: 767864
          cache_write: 7717
          cost: 0.0529
        - model: claude-opus-5-5
          input: 41
          output: 4673
          cache_read: 3285718
          cache_write: 47221
          cost: 1.1361
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-11-17T07:13:19Z
  value: 150
  by: planner-S-0823
  at: 2026-11-17T08:30:28Z
forecast:
  duration: 36m
  delivery: 2026-11-17T15:00:00Z
  basis: "Its own forecast of 36m; 14th in the pull order with an in-progress limit of 3, behind S-0752, S-0806, S-0827, S-0754, S-0811, S-0809, S-0756, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815 and S-0822."
  by: flai
  at: 2026-11-17T08:32:48Z
---
# S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back

## Goal

This story remediates [I-0397](../../../design/issues/I-0397-a-story-the-operator-moves-from-review-to-in-progress-by-hand-gets-no-agent-and.md), "A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back". The issue recommends this solution:

A move of a story from a state of kind `waiting` back to a working state by the operator's hand, from the dashboard or `flai move`, with a reason, is a send back: flai opens the send back's thread on the story with the reason as its text, as `sendback.Thread` does for `flai accept` and the rebase and resolve job, so that flai serve starts the story's agent on it. Where the operator gives no reason, the dashboard asks for one, or the board shows the story as in progress with no agent and no thread awaiting it, with what to do. Say which in the story, and bring the workflow design and the user guide in step.

## Acceptance criteria
- [ ] The cause I-0397 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0397 is closed with `flai issue close I-0397 --reason` saying what fixed it

## Tasks
- T-3541 flai move along the send-back transition by hand opens the send back's thread with the move's reason
- T-3543 flai serve starts the agent of a story sent back by hand and tells it who sent it back and why
- T-3545 The workflow design and the user guides name the send back by hand, and I-0397 is closed

## Notes

Cost of delay inputs set by flai from I-0397. time_lost_per_cycle 1h carried over from I-0397's Impact section. Its Impact time_lost_per_cycle was taken rather than the 2h30m derived from its cost and count.

### Planning

Which of the issue's two ways: the thread. The baseline's send-back transition, review to in-progress, already lists the `reason-given` rule, so `flai move` refuses it without `--reason`, and the dashboard's move dialog already asks for one (`AsksReason` in `flai/internal/hostapi/writes_describe.go`). Every send back by hand has a reason, so no board marker is needed. A project's workflow that drops `reason-given` gets the thread with a text saying no reason was given.

What counts as a send back by hand: a story's move along the transition with the `send-back` effect, made through `flai move` (which the dashboard's `item.move` runs with `--by` the operator) by anyone but `flai`, whose own send backs open their thread first, and the story's own agent, which is working it. The MCP `item_move` is left alone: the operator does not move through it.

Touches, each from:

- `flai/internal/sendback/sendback.go`, `sendback_test.go`: design (the issue names `sendback.Thread`); the shared helper for a send back by hand goes here, since `workitem` cannot import `sendback`.
- `flai/cmd/move.go`, `flai/cmd/move_sendback_test.go` (new): layout; `flai move` runs `repo.TransitionAll`, and the dashboard's `item.move` runs `flai move` (`flai/internal/hostapi/writes_items.go`). The new test reproduces I-0397.
- `docs/users/flai-reference.md`: layout; generated from `flai move`'s help, which says what the move now does.
- `flai/internal/harness/harness.go`, `harness_test.go`: layout; `harness.SentBack` words the agent's prompt, and its `Acceptance` text speaks of a rebase the acceptance could not settle, which a send back by hand is not.
- `flai/internal/serve/agents_restart.go`, `agents_start.go`, `agents_test.go`: layout; `startAcceptSentBack` finds the `Acceptance sent back:` thread by `flai` and starts the agent, `resolveAcceptSentBack` resolves it once the story is back in review, and `agents_start.go` writes the journal's words.
- `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flaiover.md`: design and documentation; the review to in-progress row of Transitions and who makes them, the pull order's send-back paragraph, `flai move`, and the dashboard's move.
- The I-0397 issue file and `design/issues/summary.md`: declared by criterion 2, written by `flai issue close`.

`flai touches suggest` listed only hub files at 11% or less (`design/system/flaiover-dashboard.md`, `docs/operators/index.md`, `template/CHANGELOG.md`); none is changed by this fix, so none was added. No folder touch.

Forecast: 36m from `flai forecast` over 16 touches, kept: the change is one helper, two callers' words, and docs, in line with the 79 remediation stories it is measured on.

Cost of delay: 150 USD a week from `flai cod`, kept: 1h lost per 168h cycle, the operator's input from I-0397, at 150 USD an hour.

No `after`: the story reads no output of another open story.
