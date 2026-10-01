---
id: I-0002
title: Go toolchain is not on PATH in the agent shell
class: blocker
status: closed
count: 1
cost: 3m
first_reported: 2026-09-15T18:00:57Z
last_reported: 2026-09-15T18:00:57Z
updated: 2026-10-01T08:00:48Z
---

# I-0002 Go toolchain is not on PATH in the agent shell

## Description
`go` lives in `/usr/local/go/bin`, which is not on the PATH the agent shell inherits. Every Go command fails with "command not found" until PATH is prefixed by hand.

## Instances

### 2026-09-15T18:00:57Z
S-0007: `go test` failed mid-story; found the binary, prefixed PATH on every subsequent call.

## Remediation
Scripts under `scripts/` prefix PATH with `/usr/local/go/bin` and `~/go/bin` themselves. Close when the operator adds it to the shell profile or the scripts cover every Go entry point.
Closed 2026-10-01T08:00:48Z: No longer applies: /usr/local/go is gone and go 1.26.4 is /opt/homebrew/bin/go on the normal PATH (triage 2026-10-01). The dead PATH entry in scripts/env.sh goes with S-0183.
