---
id: T-1876
type: task
nature: feature
title: The harness lays a tokens message over the project's stylesheet in one style element, and an empty message removes it
status: done
parent: S-0441
owner: alex
created: 2026-10-16T21:36:27Z
updated: 2026-10-17T02:58:57Z
transitions:
  - to: ready
    at: 2026-10-17T02:35:15Z
    by: agent-S-0441
  - to: in-progress
    at: 2026-10-17T02:35:16Z
    by: agent-S-0441
  - to: done
    at: 2026-10-17T02:58:57Z
    by: agent-S-0441
stream: S-0441
tags: [cli]
touches: [flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go]
usage:
  source: log
  seconds: 1421
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 55
      output: 21266
      cache_read: 2485444
      cache_write: 111868
      cost: 1.5529
---
# T-1876 The harness lays a tokens message over the project's stylesheet in one style element, and an empty message removes it

## Work

- `harness.mjs`: on `tokens`, write its CSS into one `<style data-flai-tokens>` appended last in `<head>`, so it wins over the project's stylesheet. Keep it last after a hot update adds styles, and remove it on an empty message.
- `integration_test.go`: the harness page serves the element's code. The cascade itself is checked by hand in T-1878.

Waits for nothing in this story.

## Done when

`flai test --all flai/internal/components` passes.

## Notes

None.
