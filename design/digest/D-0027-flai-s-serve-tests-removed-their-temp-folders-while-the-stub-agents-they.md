---
id: D-0027
title: flai's serve tests removed their temp folders while the stub agents they started still wrote to them, so a release run failed
kind: cause
sources: [I-0048]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover, flai/internal/hostapi]
date: 2026-09-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0027 flai's serve tests removed their temp folders while the stub agents they started still wrote to them, so a release run failed

I-0048 was closed on 2026-09-26.

The cause:

```text
flai's serve tests removed their temp folders while the stub agents they started still wrote to them, so a release run failed
```

The fix:

```text
S-0122 T-0450 (cbd211a): newAgentLab waits for its stubs and the launcher before its temp folders go; 20 race runs pass
```

Source: I-0048.
