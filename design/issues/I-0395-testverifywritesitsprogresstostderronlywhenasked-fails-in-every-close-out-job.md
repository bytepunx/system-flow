---
id: I-0395
title: "TestVerifyWritesItsProgressToStderrOnlyWhenAsked fails in every close-out job, whose environment sets FLAI_VERIFY_PROGRESS=1"
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-11-17T02:14:43Z
last_reported: 2026-11-17T02:14:43Z
updated: 2026-11-17T02:15:21Z
---

# I-0395 TestVerifyWritesItsProgressToStderrOnlyWhenAsked fails in every close-out job, whose environment sets FLAI_VERIFY_PROGRESS=1

## Description
TestVerifyWritesItsProgressToStderrOnlyWhenAsked fails in every close-out job, whose environment sets FLAI_VERIFY_PROGRESS=1

## Instances

### 2026-11-17T02:14:43Z
Story: S-0797.
S-0797's close-out job S-0797-20261117T020356Z stopped at the integration tier on flai/cmd/verify_test.go:1502 TestVerifyWritesItsProgressToStderrOnlyWhenAsked, "without the variable: exit 0", with stderr full of progress lines. S-0804 has flai serve's close-out job run flai verify with FLAI_VERIFY_PROGRESS=1 (flai/internal/serve/verifyjob.go), and the go test processes the integration tier starts inherit it. So the test's first run, meant to have no variable, has it. The test fails in every close-out that runs integration, at a head where S-0797 changed nothing in flai/cmd. Fix: the test clears the variable with t.Setenv(closeout.ProgressEnv, "") before its first run. S-0797 makes that fix, since its own close-out cannot pass without it.

## Remediation
Closed 2026-11-17T02:15:21Z: S-0797, T-3491: TestVerifyWritesItsProgressToStderrOnlyWhenAsked clears FLAI_VERIFY_PROGRESS with t.Setenv(closeout.ProgressEnv, "") before its run without the variable, so it no longer inherits the close-out job's FLAI_VERIFY_PROGRESS=1.
