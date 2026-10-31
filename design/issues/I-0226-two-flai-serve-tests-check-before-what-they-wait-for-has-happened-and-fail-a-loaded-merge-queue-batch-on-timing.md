---
id: I-0226
title: Two flai serve tests check before what they wait for has happened, and fail a loaded merge queue batch on timing
class: defect
status: closed
count: 2
cost: 18m
first_reported: 2026-10-28T23:45:11Z
last_reported: 2026-10-30T14:35:46Z
updated: 2026-10-31T09:50:12Z
---

# I-0226 Two flai serve tests check before what they wait for has happened, and fail a loaded merge queue batch on timing

## Description
Two flai serve tests check before what they wait for has happened, and fail a loaded merge queue batch on timing

## Instances

### 2026-10-28T23:45:11Z
Story: S-0507.
Merge queue batch 20261028-232929 sent S-0507 back on TH-0724, though S-0507 changes nothing in flai/internal/serve. TestAStorySentBackIsStartedAfreshOrBeforeARestart wanted exactly one journal entry once the stub agent ended, and got the thread flai serve opens for an agent that ended with its story in progress as a second. TestRequestMergeIsAnsweredByFlaiServesMergePass read the merge's log line after waiting for its journal entry alone, and found the log empty. Both pass alone and failed in no earlier batch; that batch's mcpserver suite took 679 s beside them. S-0507's T-2317 makes the first check for one start and the second wait for the log line.

### 2026-10-30T14:35:46Z
Story: S-0524.
S-0524's close-out stopped at go-test on TestAStorySentBackIsStartedAfreshOrBeforeARestart/a_failed_run_from_before_the_send_back (flai/internal/serve/agents_test.go:2778): "never happened: it is started for the send back". waitFor gives it 5 s, and the host's load average was 23.7. S-0524 changes only flai/internal/serve/plan.go in serve, the planner's trigger, and the same package passed in its flai test run 40 minutes earlier. The close-out was run again; 10m lost.

## Remediation
Closed 2026-10-31T09:50:12Z: S-0548 T-2362: the send-back subtest raced its own stub, which could end before the same look restarted what ended (2 failures in 100 runs reproduced, 0 in 200 after); it now holds the stub. waitFor and the merge test's waits allow 30 s, and the merge test waits for its log line. GOMAXPROCS=2 go test -race -count=10 on TestAStorySentBackIsStartedAfreshOrBeforeARestart, TestRequestMergeIsAnsweredByFlaiServesMergePass, and TestAChangedFileReachesTheDashboardAsANotification passed 30 of 30 at load about 40.
