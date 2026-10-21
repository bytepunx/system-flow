---
id: I-0165
title: flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-10-19T05:38:19Z
last_reported: 2026-10-19T05:38:19Z
updated: 2026-10-21T11:59:27Z
---

# I-0165 flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent

## Description
flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent

## Instances

### 2026-10-19T05:38:19Z
Story: S-0449.
S-0449's T-1906 sub-agent was refused both `flai components list` (a read) and `flai components serve --story S-0449`, through `bin/flai` and `scripts/flai.sh` alike: `components` is not in the guard's read list (`cliReads` in `flai/internal/guard/guard.go`). The pass it had written then had to be run by the story's agent, in its own context. A dev server started for the story's own worktree on 127.0.0.1 changes nothing in the repository; letting a task sub-agent list and serve components would keep such a pass out of the story agent's context.

## Remediation

Story S-0477 remediates this issue, created from it at 2026-10-21T11:59:27Z.
