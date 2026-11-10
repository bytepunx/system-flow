---
id: T-2845
type: task
nature: remediation
title: Close I-0329 with what fixed it
status: done
parent: S-0652
owner: alex
created: 2026-11-10T07:01:10Z
updated: 2026-11-10T08:10:17Z
transitions:
  - to: ready
    at: 2026-11-10T08:10:07Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:10:07Z
    by: system-flow
  - to: done
    at: 2026-11-10T08:10:17Z
    by: agent-S-0652
stream: S-0652
tags: [issues]
touches: [design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md, design/issues/summary.md]
after: [T-2841, T-2843]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 392
      cache_read: 525133
      cache_write: 3950
      cost: 0.1445
---
# T-2845 Close I-0329 with what fixed it

## Work

Close I-0329 with `flai issue close I-0329 --reason`, saying what fixed each test: the dev server test's deadline now outlasts the fake vite's start, and the owners' MCP test server sets a slow threshold no call reaches. `flai issue close` regenerates `design/issues/summary.md`.

It waits for T-2841 and T-2843, because the reason names what those two fixed.

## Done when

- I-0329 is closed with a reason naming both fixes.
- `design/issues/summary.md` no longer lists I-0329, and `flai check --strict` is clean.

## Notes
