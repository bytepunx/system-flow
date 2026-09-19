---
id: I-0020
title: An agent connected through MCP is not told when the designer moves work to ready
class: defect
status: open
count: 1
cost: 12m
first_reported: 2026-09-19T02:58:53Z
last_reported: 2026-09-19T02:58:53Z
updated: 2026-09-19T02:58:53Z
---

# I-0020 An agent connected through MCP is not told when the designer moves work to ready

## Description
An agent connected through MCP is not told when the designer moves work to ready

## Instances

### 2026-09-19T02:58:53Z
S-0040, S-0058: the operator moved S-0040 to S-0043 to ready from the board. inbox, which the conventions have the agent call at every transition, lists threads only and came back empty; wait_for_events takes its baseline when the call starts and an agent that ends its turn is not holding it; there was no board tool. The agent noticed twelve minutes later by running flai board for another reason. Fixed in S-0058: inbox lists ready work and what others changed since a per-agent cursor, wait_for_events reports what is behind the cursor at once, and a board tool was added.

## Remediation
