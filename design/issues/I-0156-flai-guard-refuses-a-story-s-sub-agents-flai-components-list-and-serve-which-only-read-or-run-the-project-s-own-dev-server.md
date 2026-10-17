---
id: I-0156
title: flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-10-17T00:41:49Z
last_reported: 2026-10-17T00:41:49Z
updated: 2026-10-17T00:41:49Z
---

# I-0156 flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server

## Description
flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server

## Instances

### 2026-10-17T00:41:49Z
Story: S-0438.
T-1862's sub-agent could not run flai components serve for the hand run: cliReads in flai/internal/guard/guard.go has no components entry. It rebuilt what Serve writes by hand, and the story's agent ran the command itself.

## Remediation
