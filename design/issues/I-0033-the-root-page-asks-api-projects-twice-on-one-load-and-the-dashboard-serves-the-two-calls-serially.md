---
id: I-0033
title: The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-09-20T22:21:35Z
last_reported: 2026-09-20T22:21:35Z
updated: 2026-10-01T09:23:22Z
---

# I-0033 The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially

## Description
The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially

## Instances

### 2026-09-20T22:21:35Z
Found verifying S-0080 T-0309's glance feature: +layout.svelte's projectState.refresh() and +page.svelte's loadGlances() each fetch /api/projects independently. Measured serial handling server-side (or possibly connection reuse client-side, not isolated): N concurrent browser fetches to the route took N * ~3s (the glance timeout), not ~3s together, so a stalled project's cost is paid once per caller rather than once per page load. A real page load settles in ~6s instead of ~3s. Not a T-0309 defect (still well under the ~15s bar the task's done-when needs), worth a follow-up: dedupe the two fetches, or find why concurrent calls to the same route serialize.

## Remediation
Closed 2026-10-01T09:23:22Z: S-0186: projectState.refresh() shares the /api/projects request in flight, so one root page load makes one request and the page reads its glances from projectState; and it asks with fetch's cache: 'no-store', because the serial calls were Chromium's HTTP cache lock holding a second GET of /api/projects until the first answered (measured: 3 concurrent GETs of one URL 1/2/3 s, 1 s together past the cache). The dashboard and flai serve concurrent calls side by side.
