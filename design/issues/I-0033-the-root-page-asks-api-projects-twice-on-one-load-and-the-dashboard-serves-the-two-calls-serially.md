---
id: I-0033
title: The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-09-20T22:21:35Z
last_reported: 2026-09-20T22:21:35Z
updated: 2026-09-20T22:21:35Z
---

# I-0033 The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially

## Description
The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially

## Instances

### 2026-09-20T22:21:35Z
Found verifying S-0080 T-0309's glance feature: +layout.svelte's projectState.refresh() and +page.svelte's loadGlances() each fetch /api/projects independently. Measured serial handling server-side (or possibly connection reuse client-side, not isolated): N concurrent browser fetches to the route took N * ~3s (the glance timeout), not ~3s together, so a stalled project's cost is paid once per caller rather than once per page load. A real page load settles in ~6s instead of ~3s. Not a T-0309 defect (still well under the ~15s bar the task's done-when needs), worth a follow-up: dedupe the two fetches, or find why concurrent calls to the same route serialize.

## Remediation
