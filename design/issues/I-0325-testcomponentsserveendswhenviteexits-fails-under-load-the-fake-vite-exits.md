---
id: I-0325
title: "TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers"
class: defect
status: closed
count: 1
cost: 8m
first_reported: 2026-11-09T22:51:04Z
last_reported: 2026-11-09T22:51:04Z
updated: 2026-11-10T07:25:06Z
---

# I-0325 TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers

## Description
TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers

## Instances

### 2026-11-09T22:51:04Z
Story: S-0622.
S-0622's third close-out, job S-0622-20261109T224304Z, stopped at integration after 7m11s at a7269472e. flai/cmd/components_test.go:748 TestComponentsServeEndsWhenViteExits failed: "the Vite dev server in .../ui exited before it answered at <http://127.0.0.1:40602/__flai/components/>: exit status 2", and the fake vite printed "VITE ready". The test wants the fake to answer first and then crash. Under load it crashed before flai's readiness probe got an answer. The test file has not changed since 2026-10-31, and S-0622 changes no components code: only two orchestrator test files. Like I-0311 (TestStartDevServerStopsViteWhenItsContextEnds), it is a fake-vite timing race. The host was running five stories in progress, their close-outs, planners, and the orchestrator.

## Remediation

Story S-0650 remediates this issue, created from it at 2026-11-10T05:12:13Z.
Closed 2026-11-10T07:25:06Z: S-0650 (T-2836): the crash-mode fake vite in flai/cmd/components_test.go exits 2 only after it has answered the harness's path once, instead of on a 300 ms timer that raced StartDevServer's 150 ms poll. FLAI_FAKE_COMPONENTS_VITE_DELAY delays the fake's listen, and TestComponentsServeEndsWhenASlowViteExits uses it (600 ms) to reproduce the early exit without load: it failed with "exited before it answered" under the old timer and passes with the fix.
