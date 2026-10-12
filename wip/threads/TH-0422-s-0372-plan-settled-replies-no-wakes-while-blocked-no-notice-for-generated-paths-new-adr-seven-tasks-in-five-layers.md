---
id: TH-0422
title: "S-0372 plan: settled replies, no wakes while blocked, no notice for generated paths (new ADR), seven tasks in five layers"
anchor:
  path: wip/kanban/stories/S-0372-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md
  item: S-0372
status: open
participants: [planner-S-0372, orchestrator]
created: 2026-10-12T17:24:30Z
updated: 2026-10-12T17:24:50Z
---

# TH-0422 S-0372 plan: settled replies, no wakes while blocked, no notice for generated paths (new ADR), seven tasks in five layers

On wip/kanban/stories/S-0372-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md.

## Entries

### 2026-10-12T17:24:30Z planner-S-0372
**Recommendation:** approve the remedy below and the change to the baseline of `work-management.md` that it brings. S-0372 stays a draft until it is finalized.

## The remedy, from I-0125's instances

- **Instances 1 and 2.** `flai message reply --settled`, or `message_reply` with `settled: true`, writes an entry that awaits neither story. A settled entry wakes no agent, except one whose run ended waiting on that reply. The conversation stays open, and the next message awaits the other story again.
- **Instance 3.** `flai task done` leaves the paths sync regenerates, such as `design/issues/summary.md`, out of its notice.
- **Instance 4.** While a story is blocked, flai serve starts its agent only for the answer to the thread it asked. A message or a closed conversation does not start it.
- **Not adopted.** Holding notices to a story with no commits. Its agent may still build on what changed.

## Needs your agreement

- ADR-0120 defines `awaiting` as the story that did not write the last entry, and the settled mark changes that. T-1482 records the change in a new ADR that refines ADR-0120 and ADR-0108.
- The baseline rule on a task-done notice says "Reply on the conversation either way, since it awaits your story until you do". T-1487 changes it to "reply settled when the change does not break your work", in `design/conventions/work-management.md` and its template copy.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1482 ADR | `design/adrs/README.md`, widened with the ADR file |
| 1 | T-1483 notice leaves out generated paths | `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` |
| 2 | T-1484 settled mark and `Awaiting` (after T-1482) | `flai/internal/messages/messages.go`, `messages_test.go` |
| 3 | T-1485 CLI `--settled`, MCP `settled`, inbox text (after T-1484) | `flai/cmd/message.go`, `message_test.go`, `flai/internal/mcpserver/messages.go`, `messages_test.go`, `folder.go`, `flai/internal/inbox/inbox.go`, `inbox_test.go` |
| 3 | T-1486 serve wakes (after T-1484) | `flai/internal/serve/agents.go`, `agents_test.go` |
| 4 | T-1487 design, convention, template, user docs (after T-1483, T-1485, T-1486) | `agent-coordination.md`, `workflow.md`, `agent-narrative.md`, `flai-cli.md`, two copies of `work-management.md`, `template/CHANGELOG.md`, `docs/users/flai.md`, `flai-reference.md` |
| 5 | T-1488 close I-0125 (after T-1487) | I-0125's file, `design/issues/summary.md` |

All touches name files.

## Figures

- **Forecast: 1h, up from flai's 34m.** The story has seven tasks across four packages, an ADR, and a convention change. Comparable message stories S-0330 to S-0334 took 32 to 45 minutes each.
- **Delivery:** 2026-10-13T02:58:00Z.
- **Cost of delay: 90 USD a week**, from flai's 36m input. The input counts four occurrences, but instance 4 alone reports 12 restarts, so the figure is low.

## Assumptions

- The dashboard needs no change. `Messages.svelte` already shows no "awaits" label for an empty `awaiting`.
- `flai stats`' conversation metrics (ADR-0138) do not count settled entries. That would need its own ADR.
- S-0367 claims `taskdone.go` and `flai/cmd`, so it holds S-0372 in ready until it is accepted.

### 2026-10-12T17:24:50Z orchestrator (recommendation)
**Recommendation:** approve the plan as drafted, including the new ADR refining ADR-0120 and ADR-0108, and the baseline change to `work-management.md` in both copies.

## Why

- **The remedy matches the instances.** I-0125 records settled conversations that keep waking each agent in turn. A settled mark that awaits neither story stops that and keeps the conversation open for the next real message. It still wakes an agent whose run ended waiting on that very reply, so nothing that needs an answer is lost.
- **The baseline rule causes the loop.** § Rules says "Reply on the conversation either way, since it awaits your story until you do". That reply is what wakes the other side. "Reply settled when the change does not break your work" keeps the duty to answer and ends the ping-pong.
- **No wakes while blocked** matches the blocked rule in the same section: a blocked story waits on the operator's answer, not on its peers.
- **Figures.** Every touch names a file. Raising flai's 34m forecast to 1h has a stated reason: seven tasks, an ADR, and a convention change. The 90 USD a week is flai cod's, from flai's 36m input.

## Why this is yours

Both points change recorded decisions: ADR-0120's definition of `awaiting`, and a baseline convention rule, which the template releases to every project. I leave S-0372 a draft until you answer. It is held behind S-0367 either way.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/conventions/work-management.md § Rules
