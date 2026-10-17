---
id: T-1867
type: task
nature: feature
title: The Design group in the site menu gains the Components page at /design/components
status: done
parent: S-0439
owner: alex
created: 2026-10-16T21:34:51Z
updated: 2026-10-17T01:18:26Z
transitions:
  - to: ready
    at: 2026-10-17T01:11:17Z
    by: agent-S-0439
  - to: in-progress
    at: 2026-10-17T01:11:17Z
    by: agent-S-0439
  - to: done
    at: 2026-10-17T01:18:26Z
    by: agent-S-0439
stream: S-0439
tags: [dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte]
usage:
  source: log
  seconds: 429
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2584
      cache_read: 364449
      cache_write: 9750
      cost: 0.1815
---
# T-1867 The Design group in the site menu gains the Components page at /design/components

## Work

Add the page key `components`, labelled Components, at `/design/components`, after Design system in the Design group of `SITE_MENU`. Add it to `PageKey`, to the `HREF` map in `SiteMenu.svelte`, and to `sitemenu.test.ts`, whose `locate` assertions gain the new path.

Waits for nothing.

## Done when

`flai test flaiover/src/lib/sitemenu.ts flaiover/src/lib/components/SiteMenu.svelte` passes, with svelte-check clean on the `HREF` map.

## Notes

None.
