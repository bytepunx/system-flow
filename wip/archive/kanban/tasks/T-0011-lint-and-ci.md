---
id: T-0011
type: task
nature: feature
title: Lint config and CI workflow
status: done
parent: S-0004
owner: agent
created: 2026-09-15T16:24:42Z
updated: 2026-09-15T16:32:33Z
transitions:
  - to: ready
    at: 2026-09-15T16:24:42Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:31:14Z
    by: agent
  - to: done
    at: 2026-09-15T16:32:33Z
    by: agent
stream: S-0004
tags: [cli]
---

# T-0011 Lint config and CI workflow

## Work
flai/.golangci.yaml, root .github/workflows/flai.yml running golangci-lint, go test -race, go build on changes under flai/.

## Done when
Workflow file valid; golangci-lint run and go test -race pass locally.

## Notes
