---
id: T-1539
type: task
nature: feature
title: The site menu has a fourth group, Design, whose first page is /design
status: done
parent: S-0382
owner: alex
created: 2026-10-13T17:47:03Z
updated: 2026-10-15T01:46:17Z
transitions:
  - to: ready
    at: 2026-10-15T01:34:24Z
    by: agent-S-0382
  - to: in-progress
    at: 2026-10-15T01:34:26Z
    by: agent-S-0382
  - to: done
    at: 2026-10-15T01:46:17Z
    by: agent-S-0382
stream: S-0382
tags: [dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte, flaiover/src/lib/components/SiteMenu.svelte.test.ts]
usage:
  source: log
  seconds: 711
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 9770
      cache_read: 1481176
      cache_write: 42809
      cost: 0.7421
---
# T-1539 The site menu has a fourth group, Design, whose first page is /design

## Work

Add the group `design`, labelled Design, after Host in `SITE_MENU`, with the page Tokens at `/design`, and its key to `GroupKey` and `PageKey`. `locate` puts `/design` and every path below it in the group. Update the label and locate tests in `sitemenu.test.ts` and the component test in `SiteMenu.svelte.test.ts`. It waits for nothing.

## Done when

- The menu shows Design with its page, keeps `?project=`, and the menu tests pass.

## Notes

- Layer 1. Later stories of E-0020 add pages to the group.
