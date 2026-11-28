---
id: T-3635
type: task
nature: improvement
title: "Project additions: documentation keeps no running list of sibling stories, and the planner gives each sibling of one document its own section"
status: backlog
parent: S-0856
owner: alex
created: 2026-11-27T23:58:36Z
updated: 2026-11-27T23:58:36Z
transitions: []
stream: S-0856
tags: [conventions]
touches: [design/conventions/documentation.md, design/conventions/strategic-agents.md]
after: [T-3632]
---
# T-3635 Project additions: documentation keeps no running list of sibling stories, and the planner gives each sibling of one document its own section

## Work

- In `design/conventions/documentation.md`, below the marker, add a project addition: describe each step or behaviour in its own sentence, paragraph, or section, in the present tense; keep no sentence that lists which stories of an epic or ADR have landed and which are still to come, since every sibling then rewrites that one line and each acceptance stops the next on a rebase conflict (I-0387, I-0414). Name the advisory `docs.unstarted-story` that `flai check` gives such a line.
- In `design/conventions/strategic-agents.md`, below the marker, add a project addition for the planner: when two or more stories of one epic or ADR change the same design or docs document, give each its own section or paragraph to write, named in its criteria or Notes, and never seed a "still to come" list for them to edit in turn.
- Waits for the check task, because both additions name the rule it adds and should match its name and wording.

## Done when

- Both files carry the addition below `<!-- system-flow:end-of-baseline -->`, and nothing above the marker changed.
- Each names `docs.unstarted-story` as the check task wrote it.
- `flai test design/conventions/documentation.md design/conventions/strategic-agents.md` passes.

## Notes
