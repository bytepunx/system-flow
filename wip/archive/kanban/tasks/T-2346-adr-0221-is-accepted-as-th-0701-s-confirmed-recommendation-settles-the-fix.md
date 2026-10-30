---
id: T-2346
type: task
nature: feature
title: ADR-0221 is accepted, as TH-0701's confirmed recommendation settles the fix
status: done
parent: S-0529
owner: alex
created: 2026-10-30T18:41:01Z
updated: 2026-10-30T18:41:18Z
transitions:
  - to: ready
    at: 2026-10-30T18:41:07Z
    by: system-flow
  - to: in-progress
    at: 2026-10-30T18:41:08Z
    by: system-flow
  - to: done
    at: 2026-10-30T18:41:18Z
    by: agent-S-0529
stream: S-0529
tags: []
touches: [design/adrs/0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md, design/adrs/README.md]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1335
      cache_read: 356383
      cache_write: 14711
      cost: 0.1935
---

# T-2346 ADR-0221 is accepted, as TH-0701's confirmed recommendation settles the fix

## Work

On 2026-10-30, flai serve confirmed the orchestrator's recommendation on TH-0701 after `confirm_after`, and the orchestrator resolved the thread. The fix S-0529 built stands. Accept ADR-0221 with `flai adr accept 0221`. It was written `proposed` until this confirmation.

It waits for nothing.

## Done when

- ADR-0221's status is `accepted` in its front matter and in `design/adrs/README.md`.

## Notes
