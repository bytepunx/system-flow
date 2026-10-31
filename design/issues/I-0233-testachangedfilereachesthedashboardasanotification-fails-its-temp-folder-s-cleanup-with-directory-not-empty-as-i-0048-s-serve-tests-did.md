---
id: I-0233
title: "TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with \"directory not empty\", as I-0048's serve tests did"
class: defect
status: closed
count: 2
cost: 10m
first_reported: 2026-10-29T16:37:28Z
last_reported: 2026-10-29T17:47:15Z
updated: 2026-10-31T09:50:17Z
---

# I-0233 TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did

## Description
TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did

## Instances

### 2026-10-29T16:37:28Z
Story: S-0537.
2026-10-29: S-0537's close-out stopped at the go-test tier after 8m5s on flai/internal/serve TestAChangedFileReachesTheDashboardAsANotification (serve_test.go:337). The failure: "TempDir RemoveAll cleanup: unlinkat /tmp/TestAChangedFileReachesTheDashboardAsANotification.../001/serve: directory not empty". Something still writes to the serve folder after the test returns, as I-0048 found for the stub agents (closed by S-0122's newAgentLab wait). S-0537 changed flai/internal/serve/orchestrate.go and orchestrate_reasons.go, the prime forms of the orchestrator's reasons, and not the dashboard notification path the test covers. The close-out was run again.

### 2026-10-29T17:47:15Z
Story: S-0418.
S-0418's close-out stopped at the go-test tier after 7m26s on its new loop test, TestServeStartsADefinedRoleOnAnEventAndOnASchedule: "TempDir RemoveAll cleanup: unlinkat /tmp/TestServeStartsADefinedRoleOnAnEventAndOnASchedule...: directory not empty". Cause found: Run's halt waits for a served project's client (r.done) but not for its loop goroutine or its watcher's looks, so a look in flight after Run returns writes inbox cursors, reason queues, and the journal into the folder being removed. S-0418's T-2343 has Run wait for both once the client stops, before the project is done; S-0543 may need only to confirm it.

## Remediation

Story S-0543 remediates this issue, created from it at 2026-10-29T17:25:41Z.
Closed 2026-10-31T09:50:17Z: S-0418's fe589e207 made Run wait for each served project's client, loop of looks, and watcher before it returns, which the test's cleanup waits for before its temp folder is removed; S-0548 T-2362 replaced the test's 80 ms sleep with a wait for the watcher's first look. GOMAXPROCS=2 go test -race -count=10 passed it 10 of 10, and 50 of 50 in a further run.
