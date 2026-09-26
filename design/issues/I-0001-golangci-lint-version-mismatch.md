---
id: I-0001
title: golangci-lint on the host is v1 but the config is v2
class: efficiency
status: open
count: 4
cost: 5m
first_reported: 2026-09-15T16:31:32Z
last_reported: 2026-09-26T03:11:53Z
updated: 2026-09-26T03:11:53Z
---

# I-0001 golangci-lint on the host is v1 but the config is v2

## Description
The repository lints with golangci-lint v2 (`flai/.golangci.yaml`, `version: "2"`). The host has v1.64 in `~/go/bin`, which cannot read the config. Every session that lints has to install v2 into a scratch location first, and the scratch copy disappears between sessions.

## Instances

### 2026-09-15T16:31:32Z
S-0004: discovered when the first lint run failed on the config; installed v2.5.0 into the session scratchpad.

### 2026-09-15T18:00:57Z
S-0007: scratch binary gone; reinstalled.

### 2026-09-15T18:11:49Z
S-0008: reinstalled again.

### 2026-09-26T03:11:53Z
S-0118, 2026-09-26: the host's golangci-lint is still v1. I installed v2.5.0 into the worktree's bin/ with go install to run make flai-test.

## Remediation
`scripts/install-tools.sh` installs v2 into `bin/` (git-ignored) so it persists across sessions; the operator can run it once. Close when the host has v2 or the script is in routine use.
