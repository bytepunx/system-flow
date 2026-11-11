---
id: D-0427
title: flai serve's Run returned before a served project's loop goroutine and its watcher's looks ended, so a look in flight wrote into a test's temp folder after the test and its cleanup failed with directory not empty
kind: cause
sources: [S-0418, T-2343, I-0233, S-0548, S-0543, T-2435]
topics: [go, flai, code]
paths: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: [S-0543, S-0690, T-2984, TH-0813]
---

# D-0427 flai serve's Run returned before a served project's loop goroutine and its watcher's looks ended, so a look in flight wrote into a test's temp folder after the test and its cleanup failed with directory not empty

S-0418's close-out stopped at the go-test tier on its new loop test, TestServeStartsADefinedRoleOnAnEventAndOnASchedule, with "TempDir RemoveAll cleanup: ... directory not empty". S-0537's close-out had stopped the same way on TestAChangedFileReachesTheDashboardAsANotification (I-0233).

The cause, as S-0418 recorded it on I-0233: "Run's halt waits for a served project's client (r.done) but not for its loop goroutine or its watcher's looks, so a look in flight after Run returns writes inbox cursors, reason queues, and the journal into the folder being removed." S-0418's T-2343 made `Run` wait for both once the client stops, before the project is done (wip/archive/agents/S-0418.md, Decisions).

S-0548 closed I-0233 on 2026-10-31 and confirmed the fix. Its T-2362 found "fe589e207 (S-0418) already covers I-0233's cleanup". It also replaced the notification test's 80 ms sleep with a wait for the watcher's first look, using a probe file. "GOMAXPROCS=2 go test -race -count=10 passed it 10 of 10, and 50 of 50 in a further run."

S-0543, the story made from I-0233, confirmed the fix again under heavy load and changed no code. Its T-2435 found that "All 40 runs of the two named tests and all 3 whole-package runs passed under `-race` at load averages of 54 to 71 on 24 CPUs. So the 'fix the writer' branch of T-2435 has nothing to fix." (wip/archive/agents/S-0543.md, Decisions)

A later change that adds a goroutine or a look to `flai serve`'s per-project loop must have `Run` wait for it too. Otherwise the serve tests' temp folders fail their cleanup again. A test that needs the watcher running waits for its first look, and does not sleep. I-0048 (D-0027) was the same failure, from stub agents that were still writing.

Sources: S-0418 (T-2343), I-0233, S-0548 (T-2362), S-0543 (T-2435).
