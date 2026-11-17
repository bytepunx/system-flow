---
id: D-0521
title: The flaiover-flai tier is selected only by flai/internal/hostapi and flai/cmd/hostapi.go, not by every Go change, since it builds flai and runs the 12 flaiover test files that drive it
kind: decision
sources: [S-0592, I-0304, S-0756]
topics: [flai, flaiover, testing]
paths: [system-flow.yaml, flai/internal/verify/select_test.go, flaiover/src/lib/server/testing.ts, flaiover/src/lib/server/workflow.test.ts]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0521 The flaiover-flai tier is selected only by flai/internal/hostapi and flai/cmd/hostapi.go, not by every Go change, since it builds flai and runs the 12 flaiover test files that drive it

S-0592 added the `flaiover-flai` tier to close I-0304. It runs `vitest related src/lib/server/testing.ts`, which selects the 12 flaiover test files that drive the built flai, since each imports `flaiover/src/lib/server/testing.ts`. `scripts/flaiover-unit.sh` builds `bin/flai` first when any Go file is newer than it.

The tier's paths are `flai/internal/hostapi/**` and `flai/cmd/hostapi.go` alone, as the planner assumed. S-0592's narrative, `## Decisions`:

```text
A tier selected by every Go change would build flai and run 12 test files on every Go change.
```

So a change elsewhere in flai that alters what the host API answers selects none of those tests through this tier. `flai/internal/verify/select_test.go` pins the tier's selection against the repository's manifest. The tier is described in `design/system/devex.md` § Used by this monorepo, which does not give the reason for its narrow paths.

A story that adds or changes a flaiover test of the built flai still gets it run: the `vitest` tier, selected by the changed flaiover file, runs it against the built `bin/flai`. S-0756 added `flaiover/src/lib/server/workflow.test.ts` with Go changes outside hostapi, and its narrative, `## Decisions`, records:

```text
The flaiover test is not selected by the flaiover-flai tier, whose paths are only hostapi's Go files. The vitest tier runs it against the built bin/flai, and a deliberately broken assertion confirmed it runs rather than skips.
```

Sources: S-0592 (wip/archive/agents/S-0592.md, Decisions); I-0304; S-0756 (wip/archive/agents/S-0756.md, Decisions).
