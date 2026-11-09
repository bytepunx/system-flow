---
id: I-0230
title: "flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null"
class: efficiency
status: closed
count: 2
cost: 5m
first_reported: 2026-10-29T07:33:59Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-09T09:06:09Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
---

# I-0230 flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null

## Description
flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null

## Instances

### 2026-10-29T07:33:59Z
On 2026-10-29 flai guard refused the orchestrator `flai order --by throughput --apply --json 2>&1` (and the same with cd and a jq pipe) and `flai board --json 2>/dev/null`, saying the orchestrator never does it, while the same `flai order --by throughput --apply` without the redirection passed. The refusal names the command the role file allows, so the redirection alone turns a permitted call into a refused one. The orchestrator's log lists the same call under Refusals since 2026-10-28T12:04Z.

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: the order_ready permission's one call was refused 36 times for its stderr redirection; guard.go:1722 itself says a descriptor duplication writes no file.

## Impact
wip/agents/orchestrator.md Refusals: 105 entries 2026-10-25 to 2026-11-07, 42 this week; 36 of flai order --by throughput --apply with 2>&1 or --json 2>&1, 11 of flai serve status with a redirection, 5 of flai board/queue --json 2>/dev/null

- time_lost_per_cycle: 1h

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0613 remediates this issue, created from it at 2026-11-07T23:38:58Z.
Closed 2026-11-09T09:06:09Z: S-0613 (T-2706, 89cb155): every flai guard check that judges a shell command, the orchestrator's, a role file's, the planner's, the analyzer's, the curator's, the host-configuration and shared-paths checks, and a sub-agent's, now judges each simple command's arguments, its redirections and their targets left out, so 2>&1 or 2>/dev/null neither allows nor refuses a command; tests in orchestrate_test.go, roles_test.go, and hostconfig_test.go reproduce the instances. The 11 refusals of flai serve status had another cause, the guard's reads not listing it, recorded as I-0318.
