---
id: T-008
type: task
nature: feature
title: Go module and cobra root command
status: done
parent: S-004
owner: agent
created: 2026-09-15T17:05:00Z
updated: 2026-09-15T17:20:00Z
transitions:
  - to: ready
    at: 2026-09-15T17:05:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:08:00Z
    by: agent
  - to: done
    at: 2026-09-15T17:20:00Z
    by: agent
stream: S-004
tags: [cli]
---

# T-008 Go module and cobra root command

## Work
Create flai/ with go.mod (module github.com/bytepunx/system-flow/flai, go 1.26), main.go, cmd/root.go with global --config flag and --json convention, and a README.

## Done when
go build ./... succeeds and flai --help lists commands.

## Notes
