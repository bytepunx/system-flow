---
id: D-0004
title: Accepting a story from the dashboard fails because the story worktree's gitdir is a host path the container does not have
kind: cause
sources: [I-0017]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd, flai/internal/check, flaiover, design/adrs]
date: 2026-09-18
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0004 Accepting a story from the dashboard fails because the story worktree's gitdir is a host path the container does not have

I-0017 was closed on 2026-09-18.

The cause:

```text
Accepting a story from the dashboard fails because the story worktree's gitdir is a host path the container does not have
```

The fix:

```text
fixed by S-0050 (flai 1.2.3): flai dashboard mounts the repository at its host path. Confirmed on 2026-09-18 when the operator accepted S-0048 from the board after the dashboard was restarted by the new flai
```

Source: I-0017.
