---
id: I-0418
title: The release build's goreleaser hooks run go mod tidy and the whole test suite with the cosign key and its password in the environment
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-27T22:47:08Z
last_reported: 2026-11-27T22:47:08Z
updated: 2026-11-27T22:47:08Z
---

# I-0418 The release build's goreleaser hooks run go mod tidy and the whole test suite with the cosign key and its password in the environment

## Description
The release build's goreleaser hooks run go mod tidy and the whole test suite with the cosign key and its password in the environment

## Instances

### 2026-11-27T22:47:08Z
Report: design/analysis/2026-11-27-all.md.
Read in the workflow and the goreleaser file today; no occurrence observed. A compromised dependency or test could read the signing key during a tag build, and tidy at release time means the signed binary is not necessarily built from the tagged tree.

## Impact
flai/.goreleaser.yaml:16-19 before.hooks: go mod tidy, go test ./...; .github/workflows/release-flai.yml:81-91 runs goreleaser/goreleaser-action with COSIGN_PRIVATE_KEY and COSIGN_PASSWORD in the step's env, so the hooks, every dependency's init, and every test run with the release key readable, on a tagged tree whose tests the close-out and the merge queue's batch already ran; go mod tidy may rewrite go.mod and go.sum after the tag. ADR-0070 confines the key to the release environment; the hooks run inside it. Every uses: is pinned by SHA and every workflow has permissions: contents: read. No instance observed.

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
