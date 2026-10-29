---
id: D-0427
title: flai serve's Run returned before a served project's loop goroutine and its watcher's looks ended, so a look in flight wrote into a test's temp folder after the test and its cleanup failed with directory not empty
kind: cause
sources: [S-0418, T-2343, I-0233]
topics: [go, flai, code]
paths: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0427 flai serve's Run returned before a served project's loop goroutine and its watcher's looks ended, so a look in flight wrote into a test's temp folder after the test and its cleanup failed with directory not empty

S-0418's close-out stopped at the go-test tier on its new loop test, TestServeStartsADefinedRoleOnAnEventAndOnASchedule, with "TempDir RemoveAll cleanup: ... directory not empty". S-0537's close-out had stopped the same way on TestAChangedFileReachesTheDashboardAsANotification (I-0233).

The cause, as S-0418 recorded it on I-0233: "Run's halt waits for a served project's client (r.done) but not for its loop goroutine or its watcher's looks, so a look in flight after Run returns writes inbox cursors, reason queues, and the journal into the folder being removed." S-0418's T-2343 made `Run` wait for both once the client stops, before the project is done (wip/archive/agents/S-0418.md, Decisions).

A later change that adds a goroutine or a look to `flai serve`'s per-project loop has `Run` wait for it too, or the serve tests' temp folders fail their cleanup again. I-0048 (D-0027) was the same failure from stub agents still writing. S-0543, I-0233's remediation story, was left to confirm the fix.

Sources: S-0418 (T-2343), I-0233.
