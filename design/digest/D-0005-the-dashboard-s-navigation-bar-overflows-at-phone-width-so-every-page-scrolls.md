---
id: D-0005
title: The dashboard's navigation bar overflows at phone width, so every page scrolls sideways
kind: cause
sources: [I-0018]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: []
date: 2026-09-19
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0005 The dashboard's navigation bar overflows at phone width, so every page scrolls sideways

I-0018 was closed on 2026-09-19.

The cause:

```text
The dashboard's navigation bar overflows at phone width, so every page scrolls sideways
```

The fix:

```text
fixed in S-0042: the navigation wraps; at 390 px the document is 390 px wide and nothing overflows, checked in a browser against a container built from the branch
```

Source: I-0018.
