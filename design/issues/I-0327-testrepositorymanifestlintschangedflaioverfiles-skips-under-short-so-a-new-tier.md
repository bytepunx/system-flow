---
id: I-0327
title: TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-11-09T23:09:47Z
last_reported: 2026-11-09T23:09:47Z
updated: 2026-11-09T23:09:47Z
---

# I-0327 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Description
TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Instances

### 2026-11-09T23:09:47Z
Story: S-0644.
S-0644 added the design-routes tier, which a change to flaiover/src/routes selects. flai test on select_test.go and system-flow.yaml passed: the go-test tier runs with -short, and TestRepositoryManifestLintsChangedFlaioverFiles in flai/internal/verify/select_test.go skips under -short ("integration: reads the monorepo"). That test pins every tier a change to flaiover/src/routes/+page.svelte and flaiover/package.json selects, so it failed only in the close-out's integration tier, after 7m30s. The other repository-manifest tests beside it are not gated under -short, so that make test catches a manifest change; this one reads only the manifest and the file tree, needs no git, and could drop its gate as they did.

## Remediation
