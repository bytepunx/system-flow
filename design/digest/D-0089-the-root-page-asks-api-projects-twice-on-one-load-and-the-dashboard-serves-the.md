---
id: D-0089
title: The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially
kind: cause
sources: [I-0033]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/routes/+page.svelte, flaiover/src/routes/+layout.svelte, flaiover/src/lib/project.svelte.ts, flaiover/src/routes/api/projects, flaiover/src/lib/project.svelte.test.ts, flaiover/src/routes/overview.svelte.test.ts, design/system/flaiover-dashboard.md, design/issues/I-0033-the-root-page-asks-api-projects-twice-on-one-load-and-the-dashboard-serves-the-two-calls-serially.md, design/issues/summary.md]
date: 2026-10-01
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0089 The root page asks /api/projects twice on one load, and the dashboard serves the two calls serially
