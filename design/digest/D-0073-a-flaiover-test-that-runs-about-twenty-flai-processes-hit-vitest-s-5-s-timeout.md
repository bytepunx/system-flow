---
id: D-0073
title: A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host
kind: cause
sources: [I-0042]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flai/cmd, flai/internal/serve, flai/internal/hostapi, flai/internal/harness, flai/internal/config, flai/internal/manifest, flai/internal/itemedit, design/adrs, design/system, design/issues, docs/operators, docs/users]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0073 A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host

I-0042 was closed on 2026-10-01.

The cause:

```text
A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host
```

The fix:

```text
Fixed by S-0116 (7510387): the test has a 30 s timeout.
```

Source: I-0042.
