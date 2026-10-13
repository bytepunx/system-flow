---
id: I-0146
title: "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged"
class: defect
status: open
count: 1
cost: 8m
first_reported: 2026-10-13T19:01:05Z
last_reported: 2026-10-13T19:01:05Z
updated: 2026-10-13T19:01:05Z
---

# I-0146 TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

## Description
TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

## Instances

### 2026-10-13T19:01:05Z
Story: S-0378.
S-0378's close-out stopped at the integration tier on `flai/internal/channel/channel_test.go:884` TestTheGateIsAskedOnEveryAttempt: "logged 1 refusals, want the second said again after a connection". S-0378 does not touch flai/internal/channel, and `flai test flai/internal/channel` passed at once on the same commit. The test waits for the state to read refused, then reads the log at once; the client appears to set the state before it logs the refusal, so the log can lag. Remedy: wait for the second log line (poll gateRefusals with a deadline, as waitState does), or log before publishing the state.

## Remediation
