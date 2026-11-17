---
id: I-0327
title: TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier
class: efficiency
status: closed
count: 1
cost: 10m
first_reported: 2026-11-09T23:09:47Z
last_reported: 2026-11-09T23:09:47Z
updated: 2026-11-17T08:45:09Z
---

# I-0327 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Description
TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Instances

### 2026-11-09T23:09:47Z
Story: S-0644.
S-0644 added the design-routes tier, which a change to flaiover/src/routes selects. flai test on select_test.go and system-flow.yaml passed: the go-test tier runs with -short, and TestRepositoryManifestLintsChangedFlaioverFiles in flai/internal/verify/select_test.go skips under -short ("integration: reads the monorepo"). That test pins every tier a change to flaiover/src/routes/+page.svelte and flaiover/package.json selects, so it failed only in the close-out's integration tier, after 7m30s. The other repository-manifest tests beside it are not gated under -short, so that make test catches a manifest change; this one reads only the manifest and the file tree, needs no git, and could drop its gate as they did.

## Remediation

Story S-0693 remediates this issue, created from it at 2026-11-10T08:49:04Z.

Story S-0827 remediates this issue, created from it at 2026-11-17T08:17:17Z.
Closed 2026-11-17T08:45:09Z: S-0827 (T-3523) removed the testing.Short() gate from TestRepositoryManifestLintsChangedFlaioverFiles in flai/internal/verify/select_test.go. The test reads only system-flow.yaml and the file tree, keeps its monorepo-not-present skip, and now runs in make test and flai test's go-test tier, so a tier it must list fails there rather than first in the close-out's integration tier.
