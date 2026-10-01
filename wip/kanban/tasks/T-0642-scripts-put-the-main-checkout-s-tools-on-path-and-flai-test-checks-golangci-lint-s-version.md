---
id: T-0642
type: task
nature: remediation
title: Scripts put the main checkout's tools on PATH and flai-test checks golangci-lint's version
status: ready
parent: S-0183
owner: arobson
created: 2026-10-01T08:34:18Z
updated: 2026-10-01T08:34:25Z
transitions:
  - to: ready
    at: 2026-10-01T08:34:25Z
    by: agent-S-0183
stream: S-0183
tags: []
touches: [scripts/env.sh, scripts/flai-test.sh, scripts/install-tools.sh]
---
# T-0642 Scripts put the main checkout's tools on PATH and flai-test checks golangci-lint's version

## Work
scripts/env.sh puts the main checkout's bin/ after the worktree's own and drops /usr/local/go/bin. scripts/install-tools.sh installs into the main checkout's bin/. scripts/flai-test.sh prints which golangci-lint it runs and its version, and stops with a clear message on a v1.

## Done when
From a story worktree with no bin/golangci-lint, scripts/flai-test.sh runs the main checkout's v2; with only a v1 on PATH it fails saying so.

## Notes
