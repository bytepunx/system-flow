---
id: ADR-0076
title: "An epic follows its stories: to ready and in-progress with the first, to review with the last open one, to done when that one is accepted"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0004]
---

# ADR-0076 An epic follows its stories: to ready and in-progress with the first, to review with the last open one, to done when that one is accepted

## Context

ADR-0004 gives every item type one state machine, and until S-0200 an epic's status was moved by hand. It drifted from its stories: E-0015 was still in ready with two of its stories done. The planner and orchestrator of E-0016 need epics whose state means something. The designer decided on 2026-10-02 that an epic follows its children.

## Decision

An epic's stories, cancelled ones left out and archived ones counted, put it in a state: done when all are done, review when all are in review or done, in-progress when any has started (in-progress, review, or done), ready when any is ready, else backlog. When a story moves, its epic moves in the same write toward that state, but only in the direction the story moved:

- forward when the story moved forward or was cancelled, so the epic goes to ready with its first ready story, to in-progress with its first started one, and to review with its last open one;
- back when the story moved back or came back from cancelled, and only as far as no other story holds it.

The epic walks the allowed transitions one at a time, each with the story's actor and time and a note naming the story, as a cancellation's cascade does (ADR-0028). Every step is validated before any file is written. It goes to done only in acceptance: `flai accept`, `flai move <story> done`, and the dashboard's acceptance walk it to done with its last open story and archive it, with its cancelled stories, in the one acceptance commit. An epic that is done, cancelled, or archived does not follow. A story given another epic with `flai edit --parent` moves the epic it joined as though it had entered from backlog, and the one it left as though it had been cancelled out of it. `flai move`, MCP's `item_move`, and `flai accept` report the epic's move. The inbox names the story an epic's move followed, read from the transitions' shared time and actor; the front matter records nothing new. `flai check` warns (`epic.lags-stories`, advisory) about an epic that its stories put further on than its status.

## Consequences

- An epic's state is a summary of its stories' and is no longer the operator's to set by hand, though `flai move` still moves it: a forward story move never pulls an epic back, so an epic moved ahead stays there until a story moves back.
- An epic accepted with its last story counts toward the next publish like an epic accepted by hand (git.md's release rules).
- An epic in review, with all its stories in review or done, cannot be cancelled directly, as a story in review cannot: it is moved back to in-progress first, or its stories are accepted.
- Epics from before this change, or moved by hand, lag until a story moves; `flai check` finds them, and `--strict` passes over the warning because only the operator moves an epic.
- The inbox's attribution of an epic's move to a story is inferred, like a cascade's cause: a hand move of the epic in the same second by the same actor reads as following.

## Alternatives considered

- Recompute the epic's state from its stories on every move, in either direction: it would pull an epic moved ahead by hand back on the next forward move of any story, against the story's own direction.
- Derive the epic's status at read time and store nothing: metrics read transitions (ADR-0004), so an epic would have no history of its own.
- Let the epic go to done when its last story is cancelled with the rest done: done is acceptance, the operator's act, so the epic stops at review.
