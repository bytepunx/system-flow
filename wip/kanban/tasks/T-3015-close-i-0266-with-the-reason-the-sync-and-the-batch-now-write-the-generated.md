---
id: T-3015
type: task
nature: improvement
title: Close I-0266 with the reason the sync and the batch now write the generated files again at the head
status: done
parent: S-0702
owner: alex
created: 2026-11-10T10:31:44Z
updated: 2026-11-11T09:18:22Z
transitions:
  - to: ready
    at: 2026-11-11T09:18:10Z
    by: agent-S-0702
  - to: in-progress
    at: 2026-11-11T09:18:10Z
    by: agent-S-0702
  - to: done
    at: 2026-11-11T09:18:22Z
    by: agent-S-0702
stream: S-0702
tags: [issues]
touches: [design/issues/I-0266-flai-check-finds-issues-summary-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-3005, T-3009, T-3011]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 530344
      cache_write: 867
      cost: 0.1966
---
# T-3015 Close I-0266 with the reason the sync and the batch now write the generated files again at the head

## Work

Run `flai issue close I-0266 --reason "<reason>"`. The reason names the new ADR and what fixed the cause: a rebase that wrote a generated file again at a stop now writes it again at the branch's head and commits it. That holds in `flai stream sync` and in the merge queue's batch, so a later commit that adds an issue file no longer leaves main's summary stale. `flai issue close` writes `design/issues/summary.md` again.

It waits for T-3005, T-3009, and T-3011, the fix and its docs, which the reason names.

## Done when

- I-0266's status is closed, with the reason.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` reports no `issues.summary`.

## Notes
