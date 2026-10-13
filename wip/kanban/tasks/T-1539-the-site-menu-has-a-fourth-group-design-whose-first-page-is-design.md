---
id: T-1539
type: task
nature: feature
title: The site menu has a fourth group, Design, whose first page is /design
status: backlog
parent: S-0382
owner: alex
created: 2026-10-13T17:47:03Z
updated: 2026-10-13T17:47:03Z
transitions: []
stream: S-0382
tags: [dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte.test.ts]
---
# T-1539 The site menu has a fourth group, Design, whose first page is /design

## Work

Add the group `design`, labelled Design, after Host in `SITE_MENU`, with the page Tokens at `/design`, and its key to `GroupKey` and `PageKey`. `locate` puts `/design` and every path below it in the group. Update the label and locate tests in `sitemenu.test.ts` and the component test in `SiteMenu.svelte.test.ts`. It waits for nothing.

## Done when

- The menu shows Design with its page, keeps `?project=`, and the menu tests pass.

## Notes

- Layer 1. Later stories of E-0020 add pages to the group.
