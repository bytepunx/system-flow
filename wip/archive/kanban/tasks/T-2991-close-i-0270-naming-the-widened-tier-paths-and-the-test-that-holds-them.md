---
id: T-2991
type: task
nature: remediation
title: Close I-0270 naming the widened tier paths and the test that holds them
status: done
parent: S-0696
owner: alex
created: 2026-11-10T10:30:03Z
updated: 2026-11-11T08:11:44Z
transitions:
  - to: ready
    at: 2026-11-11T08:11:20Z
    by: agent-S-0696
  - to: in-progress
    at: 2026-11-11T08:11:20Z
    by: agent-S-0696
  - to: done
    at: 2026-11-11T08:11:44Z
    by: agent-S-0696
stream: S-0696
tags: [issues]
touches: [design/issues/I-0270-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md, design/issues/summary.md]
after: [T-2987, T-2990]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 689
      cache_read: 546574
      cache_write: 16004
      cost: 0.2511
---
# T-2991 Close I-0270 naming the widened tier paths and the test that holds them

## Work

- Run `flai issue close I-0270 --reason "<reason>"`. The reason should say that the `integration` and `smoke` tiers' `paths` in `system-flow.yaml` now hold `template/**`, so a template-only change selects both in the close-out and in the merge queue's batch. It should name the test T-2987 added in `flai/internal/verify/select_test.go`.
- flai writes the issue file and regenerates `design/issues/summary.md`.
- Waits for T-2987 and T-2990, because the reason names the fix they land and the documents that describe it.

## Done when

- `flai show I-0270` reads closed, with that reason.
- `design/issues/summary.md` no longer lists I-0270 as open.
- `flai check --strict` is clean on the story.

## Notes
