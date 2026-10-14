---
id: S-0372
type: story
nature: improvement
title: A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
status: in-progress
owner: alex
created: 2026-10-12T16:52:49Z
updated: 2026-10-14T20:58:12Z
transitions:
  - to: ready
    at: 2026-10-14T20:56:28Z
    by: alex
  - to: in-progress
    at: 2026-10-14T20:56:35Z
    by: agent-S-0372
tags: [flai, messages, serve]
topics: [cli, conventions]
touches: [design/adrs/README.md, design/system/agent-coordination.md, design/system/workflow.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/folder.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md, design/adrs/0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1015
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-14
      hand_edits: 1
      work: 26
  models:
    - model: claude-opus-5-5
      input: 148
      output: 887
      cache_read: 5438232
      cache_write: 462084
      cost: 2.7165
  strategic:
    - kind: orchestrator
      seconds: 290
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 63
          output: 1095
          cache_read: 18045511
          cache_write: 26564
          cost: 4.5546
cost_of_delay:
  inputs:
    time_lost_per_cycle: 36m
    by: flai
    at: 2026-10-12T16:52:49Z
  value: 90
  by: planner-S-0372
  at: 2026-10-12T17:24:12Z
forecast:
  duration: 1h
  delivery: 2026-10-14T22:05:00Z
  basis: "Its own forecast of 1h; 5th in the pull order with an in-progress limit of 5, behind S-0304, S-0371, S-0305, S-0306 and S-0241."
  by: flai
  at: 2026-10-14T20:55:58Z
finalized:
  by: alex
  at: 2026-10-14T20:56:25Z
---
# S-0372 A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

## Goal

This story remediates [I-0125](../../../design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md), "A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0125 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0125 is closed with `flai issue close I-0125 --reason` saying what fixed it

## Tasks
- T-1482 ADR: a reply may settle a conversation, which then awaits neither story, and flai serve wakes no agent for a settled reply or for a message while its story is blocked
- T-1483 flai task done leaves the paths sync regenerates out of its changed-paths notice
- T-1484 A settled entry makes a conversation await neither story until the next message
- T-1485 flai message reply --settled and message_reply settled, with inbox and wait_for_events saying no reply is owed
- T-1486 flai serve starts no agent again for a settled entry, nor for a message while its story is blocked
- T-1487 Design, convention, template, and user docs describe the settled reply, the blocked hold on wakes, and the notice without generated paths
- T-1488 Close I-0125 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0125. time_lost_per_cycle 36m: 9m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T10:21:38Z, 4.3 days before this story; under one cycle counts as one).

### Planning

The proposed remedy, one change per kind of instance in I-0125:

| Instance | Cause | Remedy | Task |
|----------|-------|--------|------|
| 1, 2 | Every reply awaits the other story, so "agreed" wakes the other agent | A settled reply awaits neither story until the next message, and wakes no agent but one waiting on that reply | T-1484, T-1485, T-1486 |
| 3 | The task-done notice names `design/issues/summary.md`, which sync regenerates | The notice leaves out the paths sync regenerates | T-1483 |
| 4 | A blocked story's agent is started for every notice and closing entry | flai serve starts a blocked story's agent only for the answer to its thread | T-1486 |

Holding a notice to a story with no commits on its branch was considered and not adopted. Such a story's agent may still be about to build on the shape the commit changed.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flai/internal/messages/messages.go`, `messages_test.go` | layout: `Conversation.Awaiting` and `Reply` (explorer) |
| `flai/cmd/message.go`, `message_test.go` | layout: `flai message reply` |
| `flai/internal/mcpserver/messages.go`, `messages_test.go`, `folder.go` | layout: `message_reply` and the server's instructions |
| `flai/internal/inbox/inbox.go`, `inbox_test.go` | layout: `messages` summary and the `message` event text |
| `flai/internal/serve/agents.go`, `agents_test.go` | layout: `answered`, `messaged`, `resume`, `look` |
| `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` | layout: `(*run).tell`, the notice |
| `design/adrs/README.md` | design: the new ADR's index. `flai task done` widens the touches with the ADR's file, whose number is not known yet |
| `design/system/agent-coordination.md`, `workflow.md`, `agent-narrative.md` | design: the sections on messages, the notice, and wakes |
| `design/system/flai-cli.md` | co-change (26% of commits on these paths) and design: its `flai message` row |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design: the baseline rules on answering a message and a notice |
| `docs/users/flai.md`, `docs/users/flai-reference.md` | design: Messages between stories, and the generated reference |
| `design/issues/I-0125-…md`, `design/issues/summary.md` | declared by criterion 2 |

No declared touch existed. `flai touches suggest` needed paths to start from, so it was run after the layout touches were set. Its other co-change paths were left out because no task changes them. These include `flaiover-dashboard.md`, `docs/users/flaiover.md`, and `mcpserver/server.go`. `Messages.svelte` already shows nothing for an empty `awaiting`, and `endWhy` follows `Awaiting`.

Overlaps: S-0367, in progress, claims `flai/internal/taskdone/taskdone.go`, `flai/cmd` as a folder, `design/adrs/README.md`, `template/CHANGELOG.md`, and `docs/users/flai.md`. It holds this story in ready until it is accepted. E-0019 names `flai/internal/serve` as a folder.

Forecast: 1h, raised from flai's 34m. flai's figure is 78 s per unit of size over 61 large improvement stories, at size 26. This story has seven tasks in five layers. It spans four packages and a new ADR, and it changes a baseline convention in two copies. Each of the single-feature message stories S-0330 to S-0334 and S-0338 took 32 to 45 minutes of agent time. The delivery moves by the same 26 minutes, to 2026-10-13T02:58:00Z.

Cost of delay: 90 USD a week, as `flai cod` gives it from flai's input: 36m lost per 168h cycle at 150 USD an hour. It stands. The input counts four occurrences, while the fourth instance alone reports 12 restarts, so it is low. The inputs are not the planner's to change.
