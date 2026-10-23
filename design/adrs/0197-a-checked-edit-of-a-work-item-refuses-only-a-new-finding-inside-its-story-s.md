---
id: ADR-0197
title: "A checked edit of a work item refuses only a new finding inside its story's scope, and leaves out what a check scoped to that story leaves out"
status: proposed
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0089]
topics: [cli]
---

# ADR-0197 A checked edit of a work item refuses only a new finding inside its story's scope, and leaves out what a check scoped to that story leaves out

## Context

Every checked edit of a work item, `flai edit`, `flai criteria tick`, the MCP tools `item_edit` and `criteria_tick`, and the dashboard's edits, saves through `itemedit.Apply`. It runs `flai check` on the whole repository before and after the write, and refuses the edit, putting every file back, on any finding the second run has that the first did not.

That diff is not scoped. A finding another writer brings between the two runs refuses an edit that did not cause it. I-0182 has two instances:

1. **Another story just started.** `criteria_tick` on S-0474 was refused on two `narrative.state` findings on `wip/agents/S-0471.md`, the narrative of a story whose agent had just started it. [ADR-0125](0125-a-close-out-records-no-narrative-state-finding-on-another-open-story-s.md) leaves such a finding out of a check scoped to a story, and S-0474's own `flai task done` passed at the same moment.
2. **Another agent was writing a file.** `criteria_tick` on S-0466 was refused on a `markdown.MD012` on T-2042, a task of S-0467 another agent was writing, untracked, in the main checkout. A check scoped to S-0466 marks it a note outside the story ([ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)).

`flai task done`, `flai verify`, and `flai check --story` scope their check to the story; a checked edit did not.

## Decision

A checked edit of a work item refuses only a new finding inside its story's scope, and leaves out what a check scoped to that story leaves out.

1. **The second run is scoped to the item's story.** `itemedit.Apply` scopes the check it runs with the change in place to the edited item's story, the story itself or a task's parent, with `check.ScopeToStories`, giving as the changed paths every file the edit wrote: the item, the parents whose lists it changed, a retitled narrative, the files relinked, and `wip/agents/index.md`.
2. **What the scope leaves out refuses nothing.** A new finding the scoped check drops, such as a `narrative.state` or markdown finding on another open story's narrative, is neither a refusal nor a warning.
3. **What the scope marks outside is a warning.** A new finding on a path outside the story, such as a markdown finding on another story's task, is returned among the edit's warnings, and the edit stands.
4. **What is inside still refuses.** A new finding inside the story's scope, or an error on the item's own file, refuses the edit as before; a new `wip.overlap`, `touches.folder`, or `wip.after-hub` stays a warning.
5. **An epic has no story.** An edit of an epic refuses a new finding only on a path the edit wrote, and returns every other new finding as a warning.
6. **The other checked writers keep their diff.** `flai story new` and `flai task new`, a document save, and `flai adr new` still diff two unscoped runs; no instance asks for more.

This refines ADR-0085, whose rule that a finding outside the story is a note now reaches a checked edit as a warning, and [ADR-0089](0089-acceptance-criteria-are-ticked-through-flai-by-number-by-the-story-s-agent-once.md), whose tick through flai is the edit I-0182 saw refused.

## Consequences

- A story's agent ticks its criteria, and edits its story and tasks, while other agents start stories and write files in the main checkout; I-0182's cause no longer occurs.
- An edit can no longer be refused for a finding it did not bring, but neither can one that brings a finding on another story's file be refused for it: it gets a warning, which the caller shows, and the main checkout's `flai check` still reports the finding.
- Each checked edit reads the story's tasks, threads, and conversations once more to build the scope.

## Alternatives considered

- **Scope only to the paths the edit wrote, for every item.** It would not leave out what a check scoped to the story leaves out by rule, such as `item.archive` or `board.wip-limit` outside the story, and would differ from `flai task done`'s scope for no reason.
- **Retry the edit when a new finding appears.** It hides the race rather than removing it, and a slow writer outlasts any retry.
- **Scope the other checked writers too.** No instance asks for it; each can follow when one does.
