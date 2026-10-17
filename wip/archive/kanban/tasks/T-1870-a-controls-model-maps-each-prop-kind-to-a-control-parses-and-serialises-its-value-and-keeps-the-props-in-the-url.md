---
id: T-1870
type: task
nature: feature
title: A controls model maps each prop kind to a control, parses and serialises its value, and keeps the props in the URL
status: done
parent: S-0440
owner: alex
created: 2026-10-16T21:35:26Z
updated: 2026-10-17T02:04:32Z
transitions:
  - to: ready
    at: 2026-10-17T01:57:25Z
    by: system-flow
  - to: in-progress
    at: 2026-10-17T01:57:26Z
    by: system-flow
  - to: done
    at: 2026-10-17T02:04:31Z
    by: system-flow
stream: S-0440
tags: [dashboard]
touches: [flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/explorer/controls.test.ts, flai/internal/components/testdata/harness/src/lib/Stateful.svelte]
usage:
  source: log
  seconds: 425
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 12919
      cache_read: 1158722
      cache_write: 44707
      cost: 0.757
---
# T-1870 A controls model maps each prop kind to a control, parses and serialises its value, and keeps the props in the URL

## Work

`controls.ts`: `controlFor(prop)` by kind, `parse(prop, text)` with its error, `initial(prop)` from the default as written where it is a literal, and `toQuery` and `fromQuery` for a component's props in the URL (JSON, base64url-encoded, under one parameter).

Waits for nothing.

## Done when

Tests cover every kind, a default that is not a literal, invalid JSON, and the URL's round trip. `flai test flaiover/src/lib/explorer` passes.

## Notes

None.
