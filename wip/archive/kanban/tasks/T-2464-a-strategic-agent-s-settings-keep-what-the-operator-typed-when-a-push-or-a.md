---
id: T-2464
type: task
nature: feature
title: A strategic agent's settings keep what the operator typed when a push or a reconnect reads them again
status: done
parent: S-0574
owner: alex
created: 2026-11-06T18:02:17Z
updated: 2026-11-06T18:09:52Z
transitions:
  - to: ready
    at: 2026-11-06T18:03:08Z
    by: system-flow
  - to: in-progress
    at: 2026-11-06T18:05:14Z
    by: system-flow
  - to: done
    at: 2026-11-06T18:09:52Z
    by: agent-S-0574
stream: S-0574
tags: []
touches: [flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts, flaiover/src/lib/form.test.ts, flaiover/src/lib/form.ts]
after: [T-2463]
usage:
  source: log
  seconds: 278
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 44
      output: 10610
      cache_read: 4639893
      cache_write: 72947
      cost: 1.6916
---
# T-2464 A strategic agent's settings keep what the operator typed when a push or a reconnect reads them again

## Work

`StrategicSettings.svelte` copies every new read of the settings into its form (the `$effect` on `rows`), so a change to `system-flow.yaml` or a reconnect of the live events (`askAgain`) refills it while the operator types. Rebase the form on the new read with `rebase`, keeping the refusals while edits stay. A save still takes what flai wrote, since `before` is the saved form.

Waits for T-2463, which writes `rebase`.

## Done when

- [ ] A component test types into a field, gives the component a new read, and the typed value stays while untouched fields show the read.
- [ ] After a save and the read that follows, the form shows the read and is not dirty.

## Notes
