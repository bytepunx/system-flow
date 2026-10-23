---
id: D-0011
title: A host action that stops its own answering container is killed mid-run by that connection's own cancellation
kind: cause
sources: [I-0035]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd, flaiover/src, design/system]
date: 2026-09-22
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0011 A host action that stops its own answering container is killed mid-run by that connection's own cancellation

I-0035 was closed on 2026-09-22.

The cause:

```text
A host action that stops its own answering container is killed mid-run by that connection's own cancellation
```

The fix:

```text
S-0081, commit 735702d: dashboard.restart, dashboard.upgrade, and dashboard.stop run with their own background context (spec.detachTimeout, internal/hostapi), verified fixed against real Docker
```

Source: I-0035.
