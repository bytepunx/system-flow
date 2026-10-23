---
id: I-0201
title: flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-23T18:27:27Z
last_reported: 2026-10-23T18:27:27Z
updated: 2026-10-23T18:27:27Z
---

# I-0201 flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch

## Description
flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch

## Instances

### 2026-10-23T18:27:27Z
Story: S-0487.
S-0487's T-2078 sub-agent added two real-git tests in flai/internal/verify (realRepo, skipped under -short). flai test on the changed paths passed, but go-test runs with -short, so neither ran; the integration tier that would run them is all_only and queue, left to the merge queue's batch. One of them failed when run without -short (its fixture's worktree path S-004 is not the canonical S-0004, so issues.Generated wrote the summary in the main checkout). The story's agent found it only by running go test -run without -short by hand, which the delegation rules forbid. A way to run the non-short tests of the packages a change touches, such as flai test --integration with paths, would catch this before the batch.

## Remediation
