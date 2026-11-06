---
id: T-2465
type: task
nature: feature
title: A save in one section of the host settings leaves unsaved edits in the other sections as typed
status: done
parent: S-0574
owner: alex
created: 2026-11-06T18:02:17Z
updated: 2026-11-06T18:10:05Z
transitions:
  - to: ready
    at: 2026-11-06T18:03:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-06T18:05:15Z
    by: system-flow
  - to: done
    at: 2026-11-06T18:10:05Z
    by: agent-S-0574
stream: S-0574
tags: []
touches: [flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts]
after: [T-2463]
usage:
  source: log
  seconds: 290
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 8930
      cache_read: 3905257
      cache_write: 61397
      cost: 1.4238
---
# T-2465 A save in one section of the host settings leaves unsaved edits in the other sections as typed

## Work

`SettingsPanel.svelte` reads the settings again after every save and `fill()` rewrites every section's fields, so unsaved edits in other sections revert. Keep the seed `fill` last wrote, rebase the fields on each read with `rebase`, and let the saved section's own fields take the read.

Waits for T-2463, which writes `rebase`.

## Done when

- [ ] A component test edits one section, saves another, and the first keeps what was typed.
- [ ] The saved section shows what flai answered, a reset included.

## Notes
