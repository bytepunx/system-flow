---
id: I-0229
title: "Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words"
class: defect
status: closed
count: 4
cost: 17m
first_reported: 2026-10-29T03:27:16Z
last_reported: 2026-11-10T11:28:58Z
updated: 2026-11-12T16:34:45Z
---

# I-0229 Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words

## Description
Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words

## Instances

### 2026-10-29T03:27:16Z
Story: S-0519.
Merge queue batch 20261029-031456 failed S-0519 at integration (TH-0730) on flai/internal/guard/roles_test.go:139 TestADefinedRoleIsHeldToItsFile, a package S-0519 does not touch. The test fails on main alone (flai test flai/internal/guard in the main checkout at c24575c2e). S-0535 (ba74a6f94) changed the to-backlog condition's Allows in flai/internal/roles/conditions.go to "on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog". S-0416, accepted after it (afbc75b15), added a test that expects the old "(a move to backlog and nowhere else)". Neither story's close-out saw the other's change, and the queue sent back the next story it batched. S-0519 T-2323 updates the test's expected text. The queue then failed a story that did not cause the break, since every batch carries main's failure.

### 2026-11-10T11:08:33Z
Story: S-0665.
Again on 2026-11-10: S-0668 added TestStreamStateWritesContextDecisionsAndQuestions in flai/internal/mcpserver with t.Setenv, and S-0684 added TestNoTestSetsTheProcess, which refuses a test that sets the environment. Each passed its own close-out; accepted apart, they left main's mcpserver tests red, and S-0665's go-test tier failed on it after a sync. S-0665's T-3079 made the test call t.Parallel, since streamFixture already gives the server its session through Options.

### 2026-11-10T11:10:19Z
Story: S-0708.
S-0708's close-out (job S-0708-20261110T110353Z, at 0cbda19) stopped at integration on main's own failure. S-0668 added TestStreamStateWritesContextDecisionsAndQuestions in flai/internal/mcpserver/stream_sections_test.go, with t.Setenv and no t.Parallel(). S-0684 added TestNoTestSetsTheProcess, which fails on such a test. The two were accepted separately. S-0708 changes only flai/internal/serve/gateway_test.go and two docs, but the integration tier runs the whole module. S-0665's T-3079 fixes the test on S-0665's branch, so S-0708 cannot pass its close-out until that fix is on main.

### 2026-11-10T11:28:58Z
Story: S-0666.
S-0666's close-out (job S-0666-20261110T112451Z, at 7a284e2) stopped at integration on main's own failure: TestNoTestSetsTheProcess finds that S-0668's TestStreamStateWritesContextDecisionsAndQuestions calls t.Setenv and not t.Parallel. S-0666 does not touch flai/internal/mcpserver. S-0665's T-3079 (651768278, in review) fixes it. As TH-1049 settled for S-0708, S-0666 waits for S-0665's acceptance.

## Remediation

Story S-0683 remediates this issue, created from it at 2026-11-10T08:37:08Z.
Closed 2026-11-12T16:34:45Z: S-0683 (ADR-0291): every acceptance but the merge queue's landing is now refused, with the blocker verify_behind_main, while the paths main brought in since the story's last passing flai verify send back a tier that verify ran (verifyBehindMainGate in flai/internal/preview/accept.go, with verify.MainPathsSince and verify.StaleTiers). TestAnAcceptanceIsRefusedWhileMainOutranTheStorysVerify in flai/cmd reproduces this issue with real git. TestADefinedRoleIsHeldToItsFile now reads the conditions' Allows from internal/roles instead of copying them.
