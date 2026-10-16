---
id: T-1867
type: task
nature: feature
title: The Design group in the site menu gains the Components page at /design/components
status: backlog
parent: S-0439
owner: alex
created: 2026-10-16T21:34:51Z
updated: 2026-10-16T21:34:51Z
transitions: []
stream: S-0439
tags: [dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte]
---
# T-1867 The Design group in the site menu gains the Components page at /design/components

## Work

Add the page key `components`, labelled Components, at `/design/components`, after Design system in the Design group of `SITE_MENU`. Add it to `PageKey`, to the `HREF` map in `SiteMenu.svelte`, and to `sitemenu.test.ts`, whose `locate` assertions gain the new path.

Waits for nothing.

## Done when

`flai test flaiover/src/lib/sitemenu.ts flaiover/src/lib/components/SiteMenu.svelte` passes, with svelte-check clean on the `HREF` map.

## Notes

None.
