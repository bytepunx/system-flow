---
id: T-1875
type: task
nature: feature
title: flaiover turns token values into the generated stylesheet's declarations, tested against flai's own golden so it cannot drift
status: done
parent: S-0441
owner: alex
created: 2026-10-16T21:36:24Z
updated: 2026-10-17T02:58:54Z
transitions:
  - to: ready
    at: 2026-10-17T02:35:14Z
    by: agent-S-0441
  - to: in-progress
    at: 2026-10-17T02:35:14Z
    by: agent-S-0441
  - to: done
    at: 2026-10-17T02:58:54Z
    by: agent-S-0441
stream: S-0441
tags: [dashboard]
touches: [flaiover/src/lib/explorer/tokens.ts, flaiover/src/lib/explorer/tokens.test.ts]
usage:
  source: log
  seconds: 1420
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 11028
      cache_read: 1288854
      cache_write: 58010
      cost: 0.8053
---
# T-1875 flaiover turns token values into the generated stylesheet's declarations, tested against flai's own golden so it cannot drift

## Work

`tokens.ts`: `stylesheetDeclarations(tokens)` gives the CSS the generated stylesheet would carry for the given values: its variable names, `:root` for light, and `[data-theme='dark']` for dark (ADR-0159). It reuses the value checks of `flaiover/src/lib/design/preview.ts`, so that no value carries `url()` or a stray `}`. The test reads `flai/internal/designsystem/testdata/themes.tokens.json` and compares the result's declarations with `themes.tailwind.css.golden`.

Waits for nothing.

## Done when

The golden comparison and an unsafe value refused are tested. `flai test flaiover/src/lib/explorer` passes.

## Notes

If the golden's names depend on Go code the test cannot mirror, record which names differ and why in the story's `## Decisions`.
