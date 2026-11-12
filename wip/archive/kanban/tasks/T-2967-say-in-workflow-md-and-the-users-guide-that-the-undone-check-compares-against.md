---
id: T-2967
type: task
nature: remediation
title: Say in workflow.md and the users guide that the undone check compares against the commit the rebase took
status: done
parent: S-0712
owner: alex
created: 2026-11-10T10:28:18Z
updated: 2026-11-11T22:21:18Z
transitions:
  - to: ready
    at: 2026-11-11T22:12:59Z
    by: agent-S-0712
  - to: in-progress
    at: 2026-11-11T22:12:59Z
    by: agent-S-0712
  - to: done
    at: 2026-11-11T22:21:18Z
    by: agent-S-0712
stream: S-0712
tags: [docs, git]
touches: [design/system/workflow.md, docs/users/flai.md]
after: [T-2963]
usage:
  source: log
  seconds: 499
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 4849
      cache_read: 1217599
      cache_write: 52030
      cost: 0.6821
---
# T-2967 Say in workflow.md and the users guide that the undone check compares against the commit the rebase took

## Work

Two documents describe how the undone check works: the S-0371 paragraph in `design/system/workflow.md` under "Branches and collisions (ADR-0019)", and the "Main's change undone." bullet under "Story branches" in `docs/users/flai.md`. Add one sentence to each, as T-2963 leaves the code: the sync takes main's head once, as a commit, before it rebases. It reads the old base against that commit, rebases onto it, compares against it, and settles the branch's changelog lines under its `## Unreleased`. So a commit main gains while the sync runs is never counted as a change the branch undoes (I-0338, I-0335). Link I-0338 and I-0335 from the design paragraph as it links I-0124. ADR-0155 is accepted and stays as it is, since its text still fits.

It waits for T-2963: it describes the code that task leaves.

## Done when

- Both documents say the check compares against the commit the rebase took, and why.
- `workflow.md` has its `updated` date bumped.
- `flai test design/system/workflow.md docs/users/flai.md` passes the markdown lint.

## Notes
