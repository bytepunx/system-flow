---
id: I-0297
title: CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles
class: impression
status: closed
count: 1
first_reported: 2026-11-07T08:17:37Z
last_reported: 2026-11-07T08:17:37Z
updated: 2026-11-10T05:57:08Z
usage:
  strategic:
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759049
          cache_write: 13854
          cost: 0.278
---

# I-0297 CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles

## Description
CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles

## Instances

### 2026-11-07T08:17:37Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: code-quality.md says pin versions; the runtime pins are exact, the build's own inputs are not; signing is one long-lived key with no transparency log (ADR-0070).

## Impact
every uses: in .github/workflows names a tag (actions/checkout@v7, goreleaser/goreleaser-action@v7 with version '~> v2' at release-flai.yml:81, sigstore/cosign-installer@v4.1.2, docker/* v3 to v6); flai.yml, flaiover.yml, system-flow-check.yml have no permissions: block; release-flai.yml:8-11 and release-flaiover.yml:9-11 grant contents/packages/id-token write at the top level; golang:1.26-alpine and alpine:3.24 (bedrock Dockerfile:11-15,31), node:24-alpine, grafana/otel-lgtm:latest (compose.yaml:29), litellm:main-stable (scripts/litellm.sh:87); flaiover/pnpm-lock.yaml and flaiover/package-lock.json both tracked, every package.json range a caret

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0630 remediates this issue, created from it at 2026-11-07T23:40:47Z.
Closed 2026-11-10T05:57:08Z: S-0630 fixed it. Every workflow action is pinned to its commit SHA with its version in a comment, and GoReleaser is v2.18.1 in release-flai.yml (T-2774). Every workflow grants contents: read at the top level, and the release workflows grant their writes only on the jobs that use them. Every base and local-stack image names an exact version: node:24.22.2-alpine3.24, by digest in flaiover/Dockerfile, golang:1.26.0-alpine3.23, alpine:3.24.4, grafana/otel-lgtm:0.40.0, litellm v1.108.4. Dependabot watches the agent image folders (T-2778). Both package-lock.json files are removed, so flaiover/pnpm-lock.yaml is the one lockfile (T-2779). scripts/ci-pins.sh, run by smoke.sh, fails any of these again, and on the story's base it reports all 51 instances (T-2781).
