---
id: T-011
type: task
nature: feature
title: Lint config and CI workflow
status: done
parent: S-004
owner: agent
created: 2026-09-15T17:05:00Z
updated: 2026-09-15T17:52:00Z
transitions:
  - to: ready
    at: 2026-09-15T17:05:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:38:00Z
    by: agent
  - to: done
    at: 2026-09-15T17:52:00Z
    by: agent
stream: S-004
tags: [cli]
---

# T-011 Lint config and CI workflow

## Work
flai/.golangci.yaml, root .github/workflows/flai.yml running golangci-lint, go test -race, go build on changes under flai/.

## Done when
Workflow file valid; golangci-lint run and go test -race pass locally.

## Notes
