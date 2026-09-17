---
id: T-0008
type: task
nature: feature
title: Go module and cobra root command
status: done
parent: S-0004
owner: agent
created: 2026-09-15T16:24:42Z
updated: 2026-09-15T16:28:37Z
transitions:
  - to: ready
    at: 2026-09-15T16:24:42Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:27:19Z
    by: agent
  - to: done
    at: 2026-09-15T16:28:37Z
    by: agent
stream: S-0004
tags: [cli]
---

# T-0008 Go module and cobra root command

## Work
Create flai/ with go.mod (module github.com/bytepunx/system-flow/flai, go 1.26), main.go, cmd/root.go with global --config flag and --json convention, and a README.

## Done when
go build ./... succeeds and flai --help lists commands.

## Notes
