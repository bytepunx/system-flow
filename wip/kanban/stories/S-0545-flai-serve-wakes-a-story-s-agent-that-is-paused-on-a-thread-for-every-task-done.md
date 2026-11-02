---
id: S-0545
type: story
nature: improvement
title: flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply
status: backlog
owner: alex
created: 2026-10-29T18:11:17Z
updated: 2026-11-02T21:20:41Z
transitions: []
tags: [flai, template]
topics: [cli]
touches: [design/adrs, design/adrs/README.md, design/system/workflow.md, design/system/agent-coordination.md, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, docs/users/flai.md, template/CHANGELOG.md, design/issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md, design/issues/summary.md]
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
      seconds: 451
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 92
          output: 34163
          cache_read: 6862084
          cache_write: 194118
          cost: 3.609
    - kind: orchestrator
      seconds: 104
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 33
          output: 194
          cache_read: 1053883
          cache_write: 165172
          cost: 0.2882
        - model: claude-sonnet-5-5
          input: 4
          output: 20
          cache_read: 47900
          cache_write: 18402
          cost: 0.0466
cost_of_delay:
  inputs:
    time_lost_per_cycle: 8m
    by: flai
    at: 2026-10-29T18:11:17Z
  value: 20
  by: planner-S-0545
  at: 2026-11-02T20:05:05Z
forecast:
  duration: 50m
  delivery: 2026-11-03T03:28:00Z
  basis: "Its own forecast of 50m; 11th in the pull order with an in-progress limit of 5, behind S-0510, S-0554, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0520, S-0241 and S-0543."
  by: flai
  at: 2026-11-02T21:20:41Z
finalized:
  by: orchestrator
  at: 2026-11-02T20:07:21Z
---
# S-0545 flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply

## Goal

This story remediates [I-0232](../../../design/issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md), "flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0232 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0232 is closed with `flai issue close I-0232 --reason` saying what fixed it

## Tasks
- T-2440 Decide in an ADR refining ADR-0154 that a task-done notice does not start an agent paused on its own question, and say so in the design
- T-2441 Mark the notice flai task done sends as a notice in its conversation entry
- T-2442 Teach the work-management convention, its template copy, and the user guide that a notice waits for a paused agent's answer
- T-2443 flai serve does not start an agent paused on its own question for a conversation that holds only task-done notices, with a test reproducing I-0232
- T-2444 Close I-0232 with the reason naming the ADR and the change that fixed it

## Notes

Cost of delay inputs set by flai from I-0232. time_lost_per_cycle 8m: 4m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-29T15:42:30Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Where the cause is. `resume` in `flai/internal/serve/agents.go` starts an agent that ended asking once `answered` returns something. `answered` falls through to `messaged`, which returns any conversation that awaits the story's reply to a message written after the run ended, "even while its question to the operator is unanswered". ADR-0154 §3 keeps a blocked story's agent asleep for messages, but the interval flai serve opens when an agent ends asking (`blockAsked`, ADR-0205) is expressly not such a block, so every task-done notice from `flai/internal/taskdone/taskdone.go` wakes the paused agent. The other story's agent writes the notice through `messages.Notify`, so nothing in the conversation tells it apart from a message that asks something. The same wake still happens: S-0510, paused on TH-0814, was started again for MS-0323 and MS-0315 on 2026-11-02.

Proposed fix, for the story's agent to settle in an ADR before building it: mark a task-done notice as a notice in its entry, as ADR-0154 marks a settled reply, and have flai serve pass over a conversation whose entries to the story since the run ended are only notices while the run's own question to the operator is unanswered. The answer starts the agent, and its first `inbox` lists the notices, which it then answers. ADR-0154 rejected holding a notice to a story with no commits because its agent may be about to build on the changed shape; a story paused on its question builds nothing until the answer, so the new ADR refines ADR-0154 for that case only. A message that asks something, a share request, and flai's overlap question about a held ready story (ADR-0134) still wake it.

Tags: `flai` and `template`, the two components it delivers to (the serve and messages code, and the baseline convention).

Touches, and where each came from. `flai touches suggest S-0545` listed nothing, since the story declared no touches.

| Touch | Source |
|-------|--------|
| `design/adrs` | design: the new ADR, whose file name its title sets, so no task can name it yet; kept as a folder touch for that reason. It is in `claims.shared`, so it holds no ready story |
| `design/adrs/README.md` | design: the generated ADR index |
| `design/system/workflow.md` | design: § Transitions and who makes them, the "An answered agent" rule this changes |
| `design/system/agent-coordination.md` | design: the messages between stories it describes |
| `flai/internal/messages/messages.go`, `messages_test.go` | layout: the entry heading and its marks, where a notice mark goes |
| `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` | layout: the sender of the notice |
| `flai/internal/serve/agents.go`, `agents_test.go` | layout: `resume`, `answered`, and `messaged`, the cause |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md` | design: the baseline rules on answering a task-done notice and on messages while paused |
| `docs/users/flai.md` | design: § Messages between stories |
| `template/CHANGELOG.md` | layout: the template's convention changes |
| `design/issues/I-0232-…md`, `design/issues/summary.md` | criteria: closing I-0232 |

No `after`: the story reads no other open story's output.

Forecast: 50m, adjusted. `flai forecast` gave 28m: 91 s per unit of size over 117 large improvement stories on claude-opus-5-5, times size 18 (2 criteria, 16 touches). S-0372, which built ADR-0154's settled mark across the same messages, serve, convention, and docs files, took 55m of agent time, and this story has the same shape with one rule less; 50m sits between. The delivery moves by the same 22m, to 2026-11-03T02:05:00Z, behind the 14 stories ahead of it in the pull order.

Cost of delay: 20 USD a week, as `flai cod` gives it, unchanged: 8m lost per 168h cycle at 150 USD an hour. The input is I-0232's 4m per wake over 2 instances; the S-0510 wakes of 2026-11-02 suggest it is higher, but the input is the operator's to change.

Tasks and layers:

| Layer | Task | After | Why |
|-------|------|-------|-----|
| 1 | T-2440 ADR and design | none | decides the mark and the rule |
| 2 | T-2441 notice mark | T-2440 | builds the mark the ADR decides |
| 2 | T-2442 convention and docs | T-2440 | states the ADR's rule; shares no path with T-2441 |
| 3 | T-2443 serve rule and I-0232 test | T-2441 | reads the mark |
| 4 | T-2444 close I-0232 | T-2442, T-2443 | closes once the cause and the rules are changed |
