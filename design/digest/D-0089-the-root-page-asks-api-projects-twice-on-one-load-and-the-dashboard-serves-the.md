---
id: D-0089
title: The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially
kind: cause
sources: [I-0033]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/routes/+page.svelte, flaiover/src/routes/+layout.svelte, flaiover/src/lib/project.svelte.ts, flaiover/src/routes/api/projects, flaiover/src/lib/project.svelte.test.ts, flaiover/src/routes/overview.svelte.test.ts, design/system/flaiover-dashboard.md, design/issues/I-0033-the-root-page-asks-api-projects-twice-on-one-load-and-the-dashboard-serves-the-two-calls-serially.md, design/issues/summary.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0089 The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially

I-0033 was closed on 2026-10-01.

The cause:

```text
The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially
```

The fix:

```text
S-0186: projectState.refresh() shares the /api/projects request in flight, so one root page load makes one request and the page reads its glances from projectState; and it asks with fetch's cache: 'no-store', because the serial calls were Chromium's HTTP cache lock holding a second GET of /api/projects until the first answered (measured: 3 concurrent GETs of one URL 1/2/3 s, 1 s together past the cache). The dashboard and flai serve concurrent calls side by side.
```

Source: I-0033.
