---
id: D-0649
title: "flai serve replaces a run held on its image job in one write: keep in agents_state.go leaves out a replaced imageHeld() run, so the PID-0 held run never joins a role's past runs (I-0369)"
kind: cause
sources: [I-0369, S-0795, T-3392, T-3393]
topics: [go, code, testing, flai]
paths: [flai/internal/serve/agents_state.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go, flai/internal/serve/container_test.go]
date: 2026-11-16
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0649 flai serve replaces a run held on its image job in one write: keep in agents_state.go leaves out a replaced imageHeld() run, so the PID-0 held run never joins a role's past runs (I-0369)

I-0369 was closed on 2026-11-16 by S-0795. `TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone` failed under load in the race run at `container_test.go:1618`, "the held run is kept as a past run": a planner's run held on its image job (`ImageHeld` true, PID 0) was found among `PastPlans`. It failed in S-0777's and S-0790's race runs, and on main, never with the test run alone.

The cause was in flai's code, not the test. Replacing a held run with its real start was two writes: `put` recorded the replaced run among the role's past runs, and a second write, `dropPast`, took it out again. A reader between the two saw the held run as a past run. From I-0369's close reason:

```text
T-3392 (587b9d4) made a held run's replacement one write: keep in flai/internal/serve/agents_state.go leaves out a replaced run that is imageHeld(), so put never records the PID-0 held run among a role's past runs, and the second write, dropPast, is gone.
```

S-0795's Decisions (`wip/archive/agents/S-0795.md`) say why `dropPast` could go and what still holds: "with keep leaving out an imageHeld() run, no write puts a held run among the past runs ... The stop and failure paths clear ImageHeld before they put, and a run recorded again with the same Started and PID is `same`, so they still keep exactly what they kept." `container_test.go` did not change; its check was right.

`TestAHeldRunItsStartReplacesIsNoPastRun` reproduces it without load and fails on the old code. Under load (one CPU, `GOMAXPROCS=1`, 40 busy loops, `-race`) main's code failed 22 of 100 runs and the fix passed 100 of 100 (T-3393's Notes; D-0492 says how such a load is made).

So a later change to how flai serve moves a run to the past runs keeps the move one write, and a stopped or failed held run clears `ImageHeld` before it is put.

Sources: I-0369 (Remediation); S-0795 (`wip/archive/agents/S-0795.md`, Decisions); T-3392; T-3393.
