---
id: T-2811
type: task
nature: improvement
title: Record in an ADR that a rebase merges by sections every markdown document under the design and docs folders
status: done
parent: S-0637
owner: alex
created: 2026-11-09T22:09:09Z
updated: 2026-11-10T07:58:59Z
transitions:
  - to: ready
    at: 2026-11-10T07:58:06Z
    by: agent-S-0637
  - to: in-progress
    at: 2026-11-10T07:58:06Z
    by: agent-S-0637
  - to: done
    at: 2026-11-10T07:58:59Z
    by: agent-S-0637
stream: S-0637
tags: [adr, design]
touches: [design/adrs/0276-a-rebase-that-stops-on-a-markdown-file-under-the-manifest-s-design-or-docs.md, design/adrs/README.md]
usage:
  source: log
  seconds: 53
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 56
      cache_read: 1381622
      cache_write: 9232
      cost: 0.5118
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
