---
id: I-0190
title: PageSidebar's link entries have no page using them since the charts page dropped its sidebar
class: impression
status: closed
count: 1
cost: 5m
first_reported: 2026-10-22T20:35:40Z
last_reported: 2026-10-22T20:35:40Z
updated: 2026-11-07T05:17:03Z
---

# I-0190 PageSidebar's link entries have no page using them since the charts page dropped its sidebar

## Description
PageSidebar's link entries have no page using them since the charts page dropped its sidebar

## Instances

### 2026-10-22T20:35:40Z
Story: S-0484.
S-0484 moved the chart list from the charts page's PageSidebar to the nav rail. The link entries S-0461 added to PageSidebar (an entry's href, the link kind in $lib/pagesidebar.svelte.ts, and their tests) now serve no page: the item, review, and Docs pages use only pane and scroll entries. The story left PageSidebar untouched, as planned; removing the link kind, or keeping it for a coming page, is a choice for a story of its own.

## Remediation

Story S-0559 remediates this issue, created from it at 2026-11-02T00:04:42Z.
Closed 2026-11-07T05:17:03Z: S-0559 removed PageSidebar's link entries: SidebarEntry has no href, and the link kind, hashless, and the component's link branch are gone, with their tests and example. A @ts-expect-error test in pagesidebar.svelte.test.ts fails svelte-check if an entry takes an href again.
