---
id: I-0230
title: "flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null"
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-29T07:33:59Z
last_reported: 2026-10-29T07:33:59Z
updated: 2026-10-29T07:33:59Z
---

# I-0230 flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null

## Description
flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null

## Instances

### 2026-10-29T07:33:59Z
On 2026-10-29 flai guard refused the orchestrator `flai order --by throughput --apply --json 2>&1` (and the same with cd and a jq pipe) and `flai board --json 2>/dev/null`, saying the orchestrator never does it, while the same `flai order --by throughput --apply` without the redirection passed. The refusal names the command the role file allows, so the redirection alone turns a permitted call into a refused one. The orchestrator's log lists the same call under Refusals since 2026-10-28T12:04Z.

## Remediation
