---
id: T-2602
type: task
nature: feature
title: LogFeed shows a step's status and follows its log live through /api/kluster/log
status: done
parent: S-0595
owner: alex
created: 2026-11-07T22:07:23Z
updated: 2026-11-07T22:36:57Z
transitions:
  - to: ready
    at: 2026-11-07T22:10:07Z
    by: agent-S-0595
  - to: in-progress
    at: 2026-11-07T22:32:39Z
    by: agent-S-0595
  - to: done
    at: 2026-11-07T22:36:57Z
    by: agent-S-0595
stream: S-0595
tags: [dashboard]
touches: [flaiover/src/lib/components/LogFeed.svelte, flaiover/src/lib/components/LogFeed.svelte.test.ts, design/components/flaiover/src/lib/components/LogFeed.examples.json]
after: [T-2601]
usage:
  source: log
  seconds: 258
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 14320
      cache_read: 6999653
      cache_write: 55229
      cost: 2.0465
---
# T-2602 LogFeed shows a step's status and follows its log live through /api/kluster/log

## Work

ADR-0263 §4 and §5. A component, `flaiover/src/lib/components/LogFeed.svelte`, for one step of a namespace: props `namespace`, `step`, `label`, `outcome` (`pending`, `running`, `passed`, `failed`, `skipped`, or empty), `startedAt`, `endedAt`, `exitCode`, `open` (bindable). It shows the step's name and a status badge from the design tokens (success, danger, info, warning, text-muted roles), its duration, and an expandable panel (`<details>` or a button with `aria-expanded`) that, while open, reads `/api/kluster/log` from offset 0 and then from `next`, about every second while the outcome is `running` (and once more after it ends), appending the text, clearing on `reset`, keeping at most the last 256 KB, and staying scrolled to the end unless the reader has scrolled up (as `AgentStream.svelte` does). It stops polling when closed or unmounted. A read error is shown word for word. Write its examples file (`components.examples`) with pending, running, passed, and failed examples. Waits for T-2601, which adds the route and `KlusterLog` type.

## Done when

- `LogFeed.svelte.test.ts`, with fake timers and a mocked `$lib/api`, covers: no read while closed; a read from 0 on open; follow-on reads from `next` while running; the stop after it ends; `reset` clearing the text; the badge per outcome. It fails without the component.
- `flai test` on the changed paths passes, `components.examples` included.

## Notes
