---
id: T-0935
type: task
nature: improvement
title: An ADR extends ADR-0083 to issues and to the story the issue step makes from one
status: backlog
parent: S-0227
owner: alex
created: 2026-10-05T05:45:19Z
updated: 2026-10-05T05:45:19Z
transitions: []
stream: S-0227
tags: [flai]
touches: [design/adrs, design/system/continuous-improvement.md]
---
# T-0935 An ADR extends ADR-0083 to issues and to the story the issue step makes from one

## Work

Write the ADR with `flai adr new`, refining ADR-0083 and ADR-0079, and the design it decides. It waits for nothing: the code tasks follow the shape it decides.

- An issue's front matter may carry `usage` with `strategic` only: one entry per kind, in the shape ADR-0083 gives items (`kind`, `seconds`, `estimated: true`, `models`).
- An analyzer activity's apportioned usage goes under the `analyzer` entry of the issues it named, split evenly. An activity that names no issue goes to the project strategic total that S-0226 adds.
- The story the issue step makes carries the issue's strategic usage as its own `strategic` entries, summed up to its epic. The issue keeps its own figures, and `flai stats` counts them once: the issue's until a story is made from it, the story's after.
- Add `usage` to the issue file example in `design/system/continuous-improvement.md`, and say there how the charge and the carry-over work.

## Done when

- The ADR is accepted, links ADR-0083, and names the alternatives it weighed, such as moving the usage off the issue rather than keeping it.
- `design/system/continuous-improvement.md` shows `usage` in the issue file and describes the charge and the carry-over.
- `flai check --strict` and the markdown lint pass.

## Notes
