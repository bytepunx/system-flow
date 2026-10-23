---
id: D-0002
title: Charts page keeps showing the previous chart when switching chart type until a manual refresh
kind: cause
sources: [I-0014]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/lib/components/Chart.svelte, flaiover/src/routes/charts]
date: 2026-09-18
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0002 Charts page keeps showing the previous chart when switching chart type until a manual refresh

I-0014 was closed on 2026-09-18.

The cause:

```text
Charts page keeps showing the previous chart when switching chart type until a manual refresh
```

The fix:

```text
fixed in S-0045: the chart redraw effect now tracks its option; empty report lists no longer throw
```

Source: I-0014.
