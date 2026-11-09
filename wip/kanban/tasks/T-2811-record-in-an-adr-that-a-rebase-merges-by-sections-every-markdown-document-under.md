---
id: T-2811
type: task
nature: improvement
title: Record in an ADR that a rebase merges by sections every markdown document under the design and docs folders
status: backlog
parent: S-0637
owner: alex
created: 2026-11-09T22:09:09Z
updated: 2026-11-09T22:09:09Z
transitions: []
stream: S-0637
tags: [adr, design]
touches: [design/adrs/0273-a-rebase-merges-by-sections-every-markdown-document-under-the-design-and-docs.md, design/adrs/README.md]
---
# T-2811 Record in an ADR that a rebase merges by sections every markdown document under the design and docs folders

## Work

Write an ADR with `flai adr new` that refines ADR-0175. It decides that a rebase that stops on a markdown file merges it by sections when the file is inside `claims.shared` or under the manifest's design or docs folder, the issue files and the generated indexes keeping their own steps. Files under the wip folder, and any other path, stay out: flai writes the work items and narratives there.

Give its context from I-0315's instances: S-0599, S-0598, and S-0603 were each sent back for `design/system/server-performance.md`, where E-0033's stories each added a cause section after the S-0596 Stories table. Every resolution kept both sections. Name the alternatives: adding each such file to `claims.shared`, which is the operator's list and is known only after the conflict, and a planning rule that sends each story's section to its own file. Say that the merge queue's outcome table is unchanged: a file the merge settles stops nothing, and one it does not settle outside `claims.shared` is still `left_out` or `sent_back`. Say that the trial merge still reports these files, since a section merge can fail.

This task waits for nothing. If `flai adr new` gives another number or slug than this task's touch, update the touch of the story and of this task with `flai touches`.

## Done when

- The ADR exists, refines ADR-0175, and `flai adr index` has written it into `design/adrs/README.md`.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
