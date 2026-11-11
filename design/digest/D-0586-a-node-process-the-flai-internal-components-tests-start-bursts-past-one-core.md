---
id: D-0586
title: "A node process the flai/internal/components tests start bursts past one core for a few seconds at GOMAXPROCS=1, and NODE_OPTIONS=--v8-pool-size=1 with UV_THREADPOOL_SIZE=1 does not bound it"
kind: cause
sources: [S-0722]
topics: [go, testing, flai]
paths: [flai/internal/components, scripts/integration.sh, scripts/test.sh]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0586 A node process the flai/internal/components tests start bursts past one core for a few seconds at GOMAXPROCS=1, and NODE_OPTIONS=--v8-pool-size=1 with UV_THREADPOOL_SIZE=1 does not bound it

S-0722 put the close-out's Go tests on one core with `GOMAXPROCS=1` (D-0584). The one-core runs still peaked above one core, and its agent traced the peaks and tried to bound them. From S-0722's Decisions (`wip/archive/agents/S-0722.md`):

```text
Peaks after: 9.4, 6.1, and 6.0 cores over 1 s and 4.9, 3.5, and 3.3 over 5 s, all from one node process the flai/internal/components tests start (Vite and the svelte props script), which GOMAXPROCS does not bound. A sampler on run 2 (.flai-cache/measure/peaks.py) found the tree above 1.5 cores in 4 seconds of 507, every one of them that node process. NODE_OPTIONS=--v8-pool-size=1 with UV_THREADPOOL_SIZE=1 did not remove it (components alone: peak 5.4 against 5.9), so neither is set: one close-out's Go tests hold about one core but for a node start-up of a few seconds.
```

`design/system/workflow.md` § "What a close-out's Go tests run (S-0722)" records the burst; this entry records that the two node settings were tried and left out, so a later story that wants a hard bound on a close-out's cores looks elsewhere, such as a CPU set for the process tree, rather than trying them again.

Source: S-0722.
