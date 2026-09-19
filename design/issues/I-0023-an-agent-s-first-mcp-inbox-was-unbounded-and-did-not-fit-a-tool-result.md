---
id: I-0023
title: An agent's first MCP inbox was unbounded and did not fit a tool result
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-09-19T08:18:10Z
last_reported: 2026-09-19T08:18:10Z
updated: 2026-09-19T08:18:10Z
---

# I-0023 An agent's first MCP inbox was unbounded and did not fit a tool result

## Description
An agent's first MCP inbox was unbounded and did not fit a tool result

## Instances

### 2026-09-19T08:18:10Z
S-0061: found 2026-09-19 in the first session whose MCP server was flai 1.2.9. The first inbox as system-flow returned 209 changes, 68,283 characters, 155 of them task transitions; the harness saved it to a file because it exceeded what a tool result may hold. A defect in S-0058: 'with no cursor, the last 24 hours are reported' had no bound. Fixed in S-0061 with a cap of 50 and a first look limited to stories and epics.

## Remediation
