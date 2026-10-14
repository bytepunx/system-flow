---
id: I-0146
title: "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged"
class: defect
status: open
count: 2
cost: 8m
first_reported: 2026-10-13T19:01:05Z
last_reported: 2026-10-14T22:49:23Z
updated: 2026-10-14T23:09:14Z
---

# I-0146 TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

## Description
TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

## Instances

### 2026-10-13T19:01:05Z
Story: S-0378.
S-0378's close-out stopped at the integration tier on `flai/internal/channel/channel_test.go:884` TestTheGateIsAskedOnEveryAttempt: "logged 1 refusals, want the second said again after a connection". S-0378 does not touch flai/internal/channel, and `flai test flai/internal/channel` passed at once on the same commit. The test waits for the state to read refused, then reads the log at once; the client appears to set the state before it logs the refusal, so the log can lag. Remedy: wait for the second log line (poll gateRefusals with a deadline, as waitState does), or log before publishing the state.

### 2026-10-14T22:49:23Z
Story: S-0371.
S-0371's close-out stopped at the integration tier on TestTheGateIsAskedOnEveryAttempt (channel_test.go:884: logged 1 refusals, want the second said again after a connection). S-0371 does not touch flai/internal/channel; the test passed five times alone. The close-out ran again.

## Remediation

Story S-0389 remediates this issue, created from it at 2026-10-14T23:09:14Z.
