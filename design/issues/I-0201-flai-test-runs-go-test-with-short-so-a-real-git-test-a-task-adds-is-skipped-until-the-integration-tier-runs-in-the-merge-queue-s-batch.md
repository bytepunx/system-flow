---
id: I-0201
title: flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch
class: defect
status: closed
count: 2
cost: 20m
first_reported: 2026-10-23T18:27:27Z
last_reported: 2026-10-23T19:51:37Z
updated: 2026-10-24T00:04:37Z
---

# I-0201 flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch

## Description
flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch

## Instances

### 2026-10-23T18:27:27Z
Story: S-0487.
S-0487's T-2078 sub-agent added two real-git tests in flai/internal/verify (realRepo, skipped under -short). flai test on the changed paths passed, but go-test runs with -short, so neither ran; the integration tier that would run them is all_only and queue, left to the merge queue's batch. One of them failed when run without -short (its fixture's worktree path S-004 is not the canonical S-0004, so issues.Generated wrote the summary in the main checkout). The story's agent found it only by running go test -run without -short by hand, which the delegation rules forbid. A way to run the non-short tests of the packages a change touches, such as flai test --integration with paths, would catch this before the batch.

### 2026-10-23T19:51:37Z
Story: S-0489.
S-0489's T-2084 changed `config.Keys` and the long help of `flai dashboard restart`, `status`, and `upgrade`, which make up `docs/users/flai-reference.md`. `flai test` on `flai/cmd` and the close-out both passed: `TestReferenceIsCurrent` skips under `-short`, and the close-out left the integration tier to the merge queue. The queue's batch 20261023-194240 then failed S-0489 at integration with "docs/users/flai-reference.md is stale: run make flai-reference" (TH-0625). That cost a send-back, a regeneration, a second close-out, and a second batch.

## Remediation

Story S-0490 remediates this issue, created from it at 2026-10-23T20:20:26Z.
Closed 2026-10-24T00:04:37Z: S-0490: this repository's `go-test` tier runs `go test -json -race -count=1 -skip ^TestGateway {packages}`, without `-short`, so `flai test` and the close-out run the changed packages' real-git and monorepo tests rather than leaving them to the merge queue's batch; `TestRepositoryManifestRunsChangedGoPackagesWithoutShort` in `flai/internal/verify/select_test.go` pins the tier's arguments.
