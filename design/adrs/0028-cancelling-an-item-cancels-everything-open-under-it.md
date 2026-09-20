---
id: ADR-0028
title: Cancelling an item cancels everything open under it
status: accepted
date: 2026-09-20
supersedes: []
superseded_by: []
refines: [ADR-0004]
---

# ADR-0028 Cancelling an item cancels everything open under it

## Context

[ADR-0004](0004-workflow-states-and-transitions.md) fixes the transitions: an item goes to `cancelled` from `backlog`, `ready`, or `in-progress`, and from `review` only to `done` or back to `in-progress`. Each transition is validated for one item; a parent's children are consulted only when the parent goes to `done`.

So cancelling an epic changed the epic alone. Its stories stayed in backlog, ready, or in progress under a parent that would never be delivered, and an agent could still pull one. One level down it was the same: on 2026-09-20 the operator cancelled S-0069 from the board while it was in progress, and its two open tasks stayed open until they were cancelled by hand. The operator asked for the rule: "when an epic is cancelled, all of it's stories and tasks should also be cancelled" (S-0070).

A story in `review` is the case the existing rules do not cover. Its work is finished and waits for the operator on a branch; ADR-0004 gives the operator two answers, accept it or send it back, and no third.

## Decision

Cancelling an item cancels everything open under it, in the same operation: an epic takes its stories that are not done or cancelled and their tasks that are not done or cancelled; a story takes its open tasks. Each child gets its own `cancelled` transition with the parent's actor and timestamp, and a note under `## Notes` naming the item whose cancellation caused it (`E-0007 cancelled: <why>`). Done, cancelled, and archived items are not touched. Every transition is validated before any file is written, so a refusal changes nothing.

A cancellation that follows a parent's may take an item out of `review`. Cancelling an item in `review` directly stays refused, as ADR-0004 has it: the operator accepts it or sends it back. The exception exists because the parent's cancellation is already the operator's answer for everything under it, and the alternative, refusing to cancel an epic until each review under it has been resolved, makes the operator do by hand what they just asked for. What is shown before the cancellation names the items in review, so the answer is given knowingly.

Nothing of git is touched. A cancelled story's branch, worktree, and narrative stay where they are, and the command says so; unmerged work is the operator's to keep or delete.

Every door uses the one rule: `flai move`, MCP `item_move`, and the dashboard, which runs `flai move`. `flai check` reports an open item under a cancelled parent, which is how a tree cancelled by an older flai or edited by hand is caught.

## Consequences

- One cancellation can change many files. `flai move` lists them first and asks on a terminal; `--dry-run` shows the list; the board shows the same list in its confirmation.
- `cancelled` stays terminal. A story wanted after all is written again; there is no undo of a cascade.
- The metrics in `design/system/metrics.md` are unchanged: a cascaded item is cancelled at the time of the cascade, by its own transition, like any other.
- An agent whose story is cancelled under it learns it from `inbox` and `wait_for_events`, which name the cause; the convention tells it to stop and leave the branch alone.
- The reason is recorded once per item, in its notes, so an archived task still says why it ended without its parent at hand.

## Alternatives considered

- Refuse to cancel a parent while anything under it is open: the state before, with an error message. It keeps each cancellation a deliberate act and makes an abandoned epic cost one command per item.
- Cancel children but refuse when any story is in review: safe, and it blocks the operator on exactly the stories they have already decided not to want.
- Allow `review` to `cancelled` for every item, directly: simpler to state, and it changes ADR-0004's promise that finished work gets an explicit answer. The cascade is the narrower change.
- Leave children as they are and have readers treat an open item under a cancelled parent as cancelled: every reader (board, stats, inbox, agents) would have to know the rule, and the files would say something that is not true.
