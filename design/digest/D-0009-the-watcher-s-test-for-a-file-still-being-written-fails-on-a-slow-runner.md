---
id: D-0009
title: The watcher's test for a file still being written fails on a slow runner
kind: cause
sources: [I-0030]
topics: [flai, cli, go, code, all]
paths: [flai/internal]
date: 2026-09-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0009 The watcher's test for a file still being written fails on a slow runner

I-0030 was closed on 2026-09-21.

The cause:

```text
The watcher's test for a file still being written fails on a slow runner
```

The fix:

```text
S-0093: the debounce is a unit fed snapshots and the test helper writes in one step; no watcher test asserts a count a stalled goroutine can change
```

Source: I-0030.
