---
id: I-0200
title: "TestSyncListsPathsChangedOutsideTheClaim fails on main since S-0468: its fixture gives a task the folder touch design, which ADR-0188 now refuses"
class: defect
status: closed
count: 1
cost: 30m
first_reported: 2026-10-23T18:05:26Z
last_reported: 2026-10-23T18:05:26Z
updated: 2026-10-23T18:42:05Z
---

# I-0200 TestSyncListsPathsChangedOutsideTheClaim fails on main since S-0468: its fixture gives a task the folder touch design, which ADR-0188 now refuses

## Description
TestSyncListsPathsChangedOutsideTheClaim fails on main since S-0468: its fixture gives a task the folder touch design, which ADR-0188 now refuses

## Instances

### 2026-10-23T18:05:26Z
Found running flai/cmd's tests on main after S-0468 was accepted by hand at 197afc747. The test (flai/cmd/stream_sync_test.go, from S-0131) runs flai task new --touches design; the fixture project's design folder holds design/README.md from the template, so S-0468's rule refuses it: "T-0001's touch design is a folder holding 1 file: a task takes no folder touch (ADR-0188)". It fails the same way at S-0468's own pre-rebase tip 02851c72b, so it is S-0468's regression, not the rebase's. It is the only failing test in flai/cmd (go test ./cmd/... on main). Likely fix: give the task the file the test commits, --touches design/note.md, which keeps the test's point that the task's touches cover that path. Until fixed, every close-out whose go-test tier runs flai/cmd, and every merge-queue batch, stops on it.

## Remediation
Closed 2026-10-23T18:42:05Z: S-0429 (T-2089): the fixture's task now touches design/note.md, the file it changes, instead of the folder design that ADR-0188 refuses
