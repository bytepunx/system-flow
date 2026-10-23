---
id: D-0006
title: A dashboard container can leave a git hook that runs on the host as the operator
kind: cause
sources: [I-0022]
topics: [flai, cli, go, code, all]
paths: [flai/cmd, docs/operators]
date: 2026-09-19
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0006 A dashboard container can leave a git hook that runs on the host as the operator

I-0022 was closed on 2026-09-19.

The cause:

```text
A dashboard container can leave a git hook that runs on the host as the operator
```

The fix:

```text
S-0064, ADR-0027: flai dashboard mounts .git/hooks, .git/config and .git/info read-only over the clone; a hook and settings written from inside the container fail, tried end to end. Takes effect for a dashboard once it is restarted by a flai that has the fix.
```

Source: I-0022.
