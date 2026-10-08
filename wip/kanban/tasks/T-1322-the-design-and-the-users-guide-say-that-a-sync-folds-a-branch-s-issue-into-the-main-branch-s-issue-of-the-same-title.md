---
id: T-1322
type: task
nature: remediation
title: The design and the users' guide say that a sync folds a branch's issue into the main branch's issue of the same title
status: done
parent: S-0326
owner: alex
created: 2026-10-08T00:34:01Z
updated: 2026-10-08T06:35:33Z
transitions:
  - to: ready
    at: 2026-10-08T06:27:47Z
    by: agent-S-0326
  - to: in-progress
    at: 2026-10-08T06:27:48Z
    by: agent-S-0326
  - to: done
    at: 2026-10-08T06:35:33Z
    by: agent-S-0326
stream: S-0326
tags: [docs, design]
touches: [design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1319]
usage:
  source: log
  seconds: 465
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 201
      cache_read: 1634278
      cache_write: 77784
      cost: 0.7639
---
# T-1322 The design and the users' guide say that a sync folds a branch's issue into the main branch's issue of the same title

## Work

Say what T-1319's ADR decides where readers look for it, linking the ADR rather than restating it:

- `design/system/continuous-improvement.md`, beside what `## Summary` says about two stories that each record an issue
- `design/system/flai-cli.md`, in the rows of `flai stream sync` and `flai accept`
- `docs/users/flai.md`, beside its paragraph on the issue summary under the sync

It waits for T-1319 because the ADR decides the behaviour it describes, and runs beside T-1320 and T-1321, which touch code only.

## Done when

- Each of the three documents says that a sync, a task's close, and an acceptance fold an open issue the branch added into the main branch's open issue of the same title, and links the ADR
- `updated` is bumped in each design document's front matter
- `flai test` on the three documents passes the markdown lint

## Notes
