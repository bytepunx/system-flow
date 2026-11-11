---
id: D-0595
title: flai serve's agent-lab tests clean up only after every launcher the lab made has stopped waiting for a PID, since a restart's second launcher records its run under the temp folder after the stub exits
kind: cause
sources: [I-0343, S-0715, T-3112]
topics: [testing, go, flai, code]
paths: [flai/internal/serve/agents_test.go, flai/internal/serve/agents_start.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0595 flai serve's agent-lab tests clean up only after every launcher the lab made has stopped waiting for a PID, since a restart's second launcher records its run under the temp folder after the stub exits

I-0343 recorded `TestARestartStartsWhatIsReadyAndNothingTwice` failing its temp folder's cleanup with "directory not empty", as I-0048's (D-0027) and I-0233's (D-0427) serve tests had.

S-0715 confirmed the cause from the code: "launcher.await (flai/internal/serve/agents_start.go) deletes a PID from waiting only after end() has recorded the run, while the stub writes ended-<pid> before it exits; settle checked only lab.l.waiting, so the restart test's second launcher, fresh, could still be writing under the lab's root when TempDir's cleanup ran" (wip/archive/agents/S-0715.md, Decisions).

The fix (T-3112): the agent lab in `flai/internal/serve/agents_test.go` now makes every extra launcher through `lab.launcher()`, and its cleanup, `settle`, waits until no launcher of the lab waits for a PID, logging what it still waits for past its deadline. `TestTheLabSettlesEveryLauncherItMade` fails when `settle` ignores a second launcher.

So a serve test that starts agents through a launcher other than the lab's own makes it with `lab.launcher()`, so that cleanup waits for it.

Sources: I-0343, S-0715 (T-3112).
