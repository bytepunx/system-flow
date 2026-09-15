---
id: I-002
title: Go toolchain is not on PATH in the agent shell
class: blocker
status: open
count: 1
cost: 3m
first_reported: 2026-09-15T12:00:57Z
last_reported: 2026-09-15T12:00:57Z
updated: 2026-09-15T22:40:33Z
---

# I-002 Go toolchain is not on PATH in the agent shell

## Description
`go` lives in `/usr/local/go/bin`, which is not on the PATH the agent shell inherits. Every Go command fails with "command not found" until PATH is prefixed by hand.

## Instances
### 2026-09-15T12:00:57Z
S-007: `go test` failed mid-story; found the binary, prefixed PATH on every subsequent call.

## Remediation
Scripts under `scripts/` prefix PATH with `/usr/local/go/bin` and `~/go/bin` themselves. Close when the operator adds it to the shell profile or the scripts cover every Go entry point.
