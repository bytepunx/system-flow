---
id: I-0343
title: "TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with \"directory not empty\", as I-0048's and I-0233's serve tests did"
class: defect
status: closed
count: 1
cost: 12m
first_reported: 2026-11-10T10:41:09Z
last_reported: 2026-11-10T10:41:09Z
updated: 2026-11-11T18:23:47Z
---

# I-0343 TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with "directory not empty", as I-0048's and I-0233's serve tests did

## Description
TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with "directory not empty", as I-0048's and I-0233's serve tests did

## Instances

### 2026-11-10T10:41:09Z
Story: S-0668.
S-0668's close-out (job S-0668-20261110T103052Z, at 88639b6) stopped at the integration tier after 8m21s on one failure: `flai/internal/serve TestARestartStartsWhatIsReadyAndNothingTwice: TempDir RemoveAll cleanup: unlinkat /tmp/TestARestartStartsWhatIsReadyAndNothingTwice4134983737/001: directory not empty`. S-0668 changes nothing under flai/internal/serve. The test (agents_test.go, newAgentLab) starts stub agents, and something still writes into the lab's temp folder when t.TempDir's cleanup runs, the cause I-0048 and I-0233 recorded for other serve tests. The close-out was asked for again; the cost is the lost run.

## Remediation

Story S-0715 remediates this issue, created from it at 2026-11-10T10:55:22Z.
Closed 2026-11-11T18:23:47Z: S-0715 (T-3112): the agent lab in flai/internal/serve/agents_test.go now makes every extra launcher through lab.launcher(), and its cleanup, settle, waits until no launcher of the lab waits for a PID, logging what it still waits for past its deadline. TestARestartStartsWhatIsReadyAndNothingTwice's second launcher, which records the missed story's run under the project root after the stub has exited, is waited for before t.TempDir removes the folder. TestTheLabSettlesEveryLauncherItMade fails when settle ignores a second launcher.
