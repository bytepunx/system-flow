---
id: I-0050
title: flai serve 1.23.0, which has holds, started an agent for a story the board showed held
class: defect
status: open
count: 2
cost: 8m
first_reported: 2026-09-29T02:06:24Z
last_reported: 2026-09-29T02:13:28Z
updated: 2026-09-29T02:13:28Z
---

# I-0050 flai serve 1.23.0, which has holds, started an agent for a story the board showed held

## Description
flai serve 1.23.0, which has holds, started an agent for a story the board showed held

## Instances

### 2026-09-29T02:06:24Z
2026-09-29: S-0144 entered ready at 01:10:17Z, claiming flai/cmd, while S-0138 (in progress since 01:07:54Z) claimed flai/cmd/prime.go. flai board and inbox showed S-0144 held (overlap), yet the host's flai serve 1.23.0 started agent-S-0144 for it with the instruction to work it to review. Unlike I-0049 the host flai postdates holds, so the agent start path appears not to consult them; not investigated. The agent narrowed the claim to what it changes and worked it; the remaining overlap (docs/users, design/system/flai-cli.md) fails the smoke tier's strict repository check until S-0138 is accepted.

### 2026-09-29T02:13:28Z
2026-09-29: S-0145 entered ready at 01:57:54Z claiming flai/cmd, while S-0138 (in progress) claimed flai/cmd/prime.go. flai board and inbox showed S-0145 held (overlap), yet the host's flai serve 1.23.0 (pid 2077672, running since 01:06:55Z) started agent-S-0145 for it at 02:10:59Z with the instruction to work it to review. Second instance in one evening, same serve. The story is research whose output is an ADR draft and touches no code, so the agent narrowed the claim to design/adrs and design/issues, which cleared the hold, and pulled it.

## Remediation
