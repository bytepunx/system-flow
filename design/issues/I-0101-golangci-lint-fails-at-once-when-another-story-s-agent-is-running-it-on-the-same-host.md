---
id: I-0101
title: golangci-lint fails at once when another story's agent is running it on the same host
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-07T00:34:48Z
last_reported: 2026-10-07T00:34:48Z
updated: 2026-10-07T01:07:14Z
---

# I-0101 golangci-lint fails at once when another story's agent is running it on the same host

## Description
golangci-lint fails at once when another story's agent is running it on the same host

## Instances

### 2026-10-07T00:34:48Z
Story: S-0273.
S-0273's T-1069 sub-agent ran golangci-lint while another agent on the host was linting: it exited 3 with 'parallel golangci-lint is running'. The flai test tier passes --allow-parallel-runners; scripts/flai-test.sh, and so the close-out, does not, so parallel stories can still fail their close-out on it.

## Remediation

Story S-0308 remediates this issue, created from it at 2026-10-07T01:07:14Z.
