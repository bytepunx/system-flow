---
id: T-1545
type: task
nature: feature
title: A token editor changes each kind of token in place and marks a bad value and a failing contrast beside it
status: done
parent: S-0383
owner: alex
created: 2026-10-13T17:47:36Z
updated: 2026-10-15T02:38:18Z
transitions:
  - to: ready
    at: 2026-10-15T02:18:54Z
    by: agent-S-0383
  - to: in-progress
    at: 2026-10-15T02:18:54Z
    by: agent-S-0383
  - to: done
    at: 2026-10-15T02:38:18Z
    by: agent-S-0383
stream: S-0383
tags: [dashboard]
touches: [flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts]
after: [T-1543, T-1544]
usage:
  source: log
  seconds: 1164
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 143
      output: 98970
      cache_read: 9014475
      cache_write: 353067
      cost: 5.7281
---
# T-1545 A token editor changes each kind of token in place and marks a bad value and a failing contrast beside it

## Work

`TokenEditor.svelte` edits one namespace of tokens bound to the edit model: a colour with a picker and a text field per theme, a length or scale step with a number and a unit, a font family with a text field. It marks a value the model rejects and each failing contrast pair the preview reports, with its ratio and theme, in words as well as colour. Component tests cover each kind of field and both marks. It waits for T-1543 and T-1544, whose model and contrast check it shows.

## Done when

- `TokenEditor.svelte.test.ts` passes for each field kind and both marks.

## Notes

- Layer 2.
