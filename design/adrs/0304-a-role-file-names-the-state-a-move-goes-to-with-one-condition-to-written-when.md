---
id: ADR-0304
title: "A role file names the state a move goes to with one condition, to, written when: {to: <state>}, which takes any state the workflow has, and to-review and to-backlog stay its aliases"
status: accepted
date: 2026-11-15
supersedes: []
superseded_by: []
refines: [ADR-0207, ADR-0211, ADR-0301]
---

# ADR-0304 A role file names the state a move goes to with one condition, to, written when: {to: <state>}, which takes any state the workflow has, and to-review and to-backlog stay its aliases

## Context

A role file says which state a move goes to with two conditions, `to-review` and `to-backlog` ([ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) §5, [ADR-0211](0211-the-orchestrator-s-trigger-vocabulary-gains-the-events-moved-thread-planned.md)). `to-review` qualifies a `moved` trigger entry. `to-backlog` qualifies a `moved` trigger entry and the `item_move` and `flai move` calls, where it allows a move to backlog and nowhere else. Both name a state of the baseline board. [ADR-0301](0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md) makes the board a workflow file whose states a project names, and §9 leaves to S-0748 the one condition that names any of them. A role cannot react today to a move into a state of a project's own, such as one added between `in-progress` and `review`.

A role file's condition is a name in `when`, and one condition stands on an entry (ADR-0207 §5). No condition so far takes an argument.

## Decision

**A role file names the state a move goes to with one condition, `to`, which takes the state's name as its argument: `when: {to: <state>}`. On a `moved` trigger entry it starts the role only for a move into that state, any state the workflow has. On an `item_move` or `flai move` call it allows only a move into that state, and a role file may name there only the workflow's initial state. `to-review` and `to-backlog` stay valid as aliases of `{to: review}` and `{to: backlog}`.**

1. **The form.** `when` takes a condition's name, as before, or a mapping of one condition's name to its argument. `to` is the only condition that takes one, and it needs one: `when: to` alone, a mapping of more than one key, and an argument on any other condition are `roles.file`. The entries read `{event: moved, when: {to: review}}` and `{call: item_move, when: {to: backlog}}`.
2. **Any state the workflow has.** The argument is checked against the states of the workflow the project runs, the baseline flai carries until a project's own file is read (S-0755). A state the workflow lacks is `roles.file`, naming the states it has. `flai check` and `flai guard` read the file through the same check, so a file that names a state the workflow lacks is refused as a whole, as any other `roles.file` is.
3. **On a trigger.** `{event: moved, when: {to: <state>}}` holds when the moved item's new state is `<state>`. The state's kind does not matter: a role may start on a move into a queue, an active, a waiting, an accepted, or a cancelled state.
4. **On a call.** `{call: item_move, when: {to: <state>}}` and `{call: flai move, when: {to: <state>}}` allow a move into `<state>` and nowhere else. Only the workflow's initial state, `backlog` in the baseline, may be named there, so that the condition grants no more than `to-backlog` did. Any other state on a call is `roles.grant`: a move to ready is `to-ready`'s under `promote_to_ready`, and no setting allows a move into an active, waiting, or accepted state (ADR-0207 §5). The never-granted `item_move` to done and the story's agent's moves stay clear of `to` as they were of `to-backlog`.
5. **The old names.** `to-review` and `to-backlog` keep working as before: flai reads `when: to-review` as `{to: review}` and `when: to-backlog` as `{to: backlog}` when it parses the file, so the guard and `flai serve` read `to` alone. A file that names them gets no finding. They are not listed among the conditions; the documents name them as aliases and say what to write instead.
6. **The baseline role files** in the template, in flai's embedded copies, and in this project write `{to: review}` and `{to: backlog}`. Their triggers and grants are what they were.

## Consequences

- A project's role starts on a move into any state its workflow has, with no new condition and no flai release.
- A project's role file that names `to-review` or `to-backlog` needs no edit, and `flai upgrade` brings the new form into the baseline role files with the template's front matter (ADR-0207 §8).
- A role file that names a state its workflow later drops is refused until it is edited, as a file naming an unknown event is.
- `to-ready` and `accept-as-self` stay their own conditions, gated by their permissions. A move past the initial state on a call is never `to`'s.
- `when` may now be a mapping; a reader of role files that took it as a string reads the condition's name and its argument apart.

## Alternatives considered

- **A condition per state, `to-<state>`.** The old names would be two of them, and nothing would change in the baseline files. A condition's name would then be open-ended, the table of conditions could not list them, and `to-ready`, gated by its permission, would read as one of them while meaning something else.
- **A `to` key on the entry beside `when`**, such as `{event: moved, to: review}`. Shorter, but it puts a condition outside `when`, where ADR-0207 keeps every condition, and lets an entry carry two.
- **A `roles.file` finding on the old names**, naming `{to: review}` or `{to: backlog}`. A project's own role file would be refused by the guard until edited, for a change of spelling alone.
- **`to` on a call for any queue state.** It would let a role move a story to ready without `promote_to_ready`. A move past the initial state stays a gated or a forbidden call.
