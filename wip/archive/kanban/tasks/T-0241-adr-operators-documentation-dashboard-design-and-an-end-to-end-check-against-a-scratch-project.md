---
id: T-0241
type: task
nature: feature
title: ADR, operators documentation, dashboard design, and an end to end check against a scratch project
status: done
parent: S-0064
owner: alex
created: 2026-09-19T10:18:11Z
updated: 2026-09-19T10:28:45Z
transitions:
  - to: ready
    at: 2026-09-19T10:24:52Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:24:52Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:28:45Z
    by: system-flow
stream: S-0064
tags: []
---

# T-0241 ADR, operators documentation, dashboard design, and an end to end check against a scratch project

## Work
An ADR recording the mount layout and what is and is not protected, refining ADR-0018 and ADR-0022. `docs/operators/index.md` and `flaiover-dashboard.md` say what the container can and cannot write, including what stays open and why. Build the branch's image and run it with the branch's `flai dashboard` against the scratch project: a hook and a config setting written from inside the container fail, an acceptance with a story branch and worktree, a document save, a board move, and a push with a throwaway key all still work. Tick the criteria for what was seen; close I-0022 if it is closed.

## Done when
- The end to end check is recorded in the narrative
- The ADR and the documents say what is closed and what is not
- The criteria are ticked and `flai check --strict` is clean

## Notes
