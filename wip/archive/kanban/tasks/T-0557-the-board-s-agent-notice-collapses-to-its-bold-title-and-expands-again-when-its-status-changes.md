---
id: T-0557
type: task
nature: feature
title: The board's agent notice collapses to its bold title and expands again when its status changes
status: done
parent: S-0150
owner: alex
created: 2026-09-29T19:05:18Z
updated: 2026-09-29T19:08:40Z
transitions:
  - to: ready
    at: 2026-09-29T19:05:33Z
    by: agent-S-0150
  - to: in-progress
    at: 2026-09-29T19:05:33Z
    by: agent-S-0150
  - to: done
    at: 2026-09-29T19:08:40Z
    by: agent-S-0150
stream: S-0150
tags: []
usage:
  source: log
  seconds: 187
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 29
      output: 7430
      cache_read: 1141068
      cache_write: 37821
      cost: 0.6795
---

# T-0557 The board's agent notice collapses to its bold title and expands again when its status changes

## Work

- `flaiover/src/lib/components/HostAgentNotice.svelte`: a caret button beside the bold title toggles the rest of the notice; collapsed, only the title shows. The ended case gets a bold title too (`Agent ended`).
- The collapse is remembered per browser against what the notice says (running, last, waiting), so a reload keeps it collapsed and any change to that status opens it again. The 15 s re-ask, which returns an equal status, does not.
- Tests in `HostAgentNotice.svelte.test.ts`.

## Done when

- The toggle collapses and expands the notice, with `aria-expanded` true when open.
- A changed status opens a collapsed notice; an unchanged one leaves it collapsed.
- `npm run check`, lint, and the flaiover tests pass.

## Notes
