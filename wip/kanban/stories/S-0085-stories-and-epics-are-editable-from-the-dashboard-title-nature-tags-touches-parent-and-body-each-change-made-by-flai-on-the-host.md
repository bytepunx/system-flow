---
id: S-0085
type: story
nature: feature
title: "Stories and epics are editable from the dashboard: title, nature, tags, touches, parent, and body, each change made by flai on the host"
status: review
parent: E-0003
owner: alex
created: 2026-09-20T13:58:34Z
updated: 2026-09-20T15:46:38Z
transitions:
  - to: ready
    at: 2026-09-20T13:58:53Z
    by: alex
  - to: in-progress
    at: 2026-09-20T15:15:52Z
    by: alex
  - to: review
    at: 2026-09-20T15:46:38Z
    by: system-flow
tags: [dashboard, cli]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, docs/users]
---
# S-0085 Stories and epics are editable from the dashboard: title, nature, tags, touches, parent, and body, each change made by flai on the host

## Goal
The designer edits a story or an epic where they read it: its title, nature, tags, what it touches, its parent, and its body (goal, acceptance criteria, notes), from the item's page, without opening the raw file and without a shell. flai makes every change, on the host, so its rules stay the only rules.

## Acceptance criteria
- [x] flai gains what is missing to change an item's own fields: title (front matter, heading, and the file's name kept in step, links to it kept working), nature, tags, and parent, beside `flai touches`; each validates as `flai story new` does, refuses what the workflow forbids (an archived item, a parent of the wrong type), and is a command first, so an agent and a shell can do the same
- [x] The channel offers them as named, typed methods (ADR-0029), with the hash of what was read so a change made meanwhile by an agent is a conflict and not an overwrite, as for documents (ADR-0023)
- [x] An item's page has an edit mode for stories and epics: the fields above as inputs, the body as Markdown with a preview, acceptance criteria tickable in place; saving checks the repository with the change in place and refuses with the findings when the check does, leaves the text in the editor, and commits on its own as the designer with the dashboard's trailer
- [x] What stays flai's is shown and not editable: ID, type, status, transitions, blocked intervals, owner, created and updated; status still changes by moving the card
- [x] An agent working on the story is told: the change shows in `inbox` and wakes `wait_for_events`, saying what changed (criteria, title, touches), and the narrative's links still resolve after a retitle
- [x] Tasks are out of scope here and said to be: they are the agent's to write; whether the designer edits them too is asked, not assumed
- [x] The users' documentation is corrected where it says items can be retitled from the board today, which they cannot, and says how editing works

## Tasks
- T-0295 flai edit changes an item's title, nature, tags, touches, parent, and body in one checked, committed step
- T-0296 An agent is told when someone else edits its story, and what changed
- T-0297 The channel offers item.show and item.edit, and the dashboard has routes for them
- T-0298 An item's page has an edit mode for stories and epics, and criteria can be ticked in place
- T-0299 The design and the users' documentation say how items are edited, and stop saying they can be retitled from the board
- T-0300 Tried end to end in a scratch project and a browser

## Notes
Asked for by the operator on 2026-09-20. What exists today: the document editor (S-0040, ADR-0023) opens an item's file and lets the body be edited as raw Markdown, with flai-owned front matter shown read-only; `flai touches` sets `touches`; nothing changes a title, nature, tags, or parent after creation except a hand edit, and a title lives in three places (front matter, heading, file name). `docs/users/flaiover.md` says "Move, block, and retitle items from the board or with flai"; there is no retitle. Creating items from the dashboard is S-0059; this is its counterpart for changing them.
