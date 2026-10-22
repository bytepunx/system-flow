---
id: I-0190
title: PageSidebar's link entries have no page using them since the charts page dropped its sidebar
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-10-22T20:35:40Z
last_reported: 2026-10-22T20:35:40Z
updated: 2026-10-22T20:35:40Z
---

# I-0190 PageSidebar's link entries have no page using them since the charts page dropped its sidebar

## Description
PageSidebar's link entries have no page using them since the charts page dropped its sidebar

## Instances

### 2026-10-22T20:35:40Z
Story: S-0484.
S-0484 moved the chart list from the charts page's PageSidebar to the nav rail. The link entries S-0461 added to PageSidebar (an entry's href, the link kind in $lib/pagesidebar.svelte.ts, and their tests) now serve no page: the item, review, and Docs pages use only pane and scroll entries. The story left PageSidebar untouched, as planned; removing the link kind, or keeping it for a coming page, is a choice for a story of its own.

## Remediation
