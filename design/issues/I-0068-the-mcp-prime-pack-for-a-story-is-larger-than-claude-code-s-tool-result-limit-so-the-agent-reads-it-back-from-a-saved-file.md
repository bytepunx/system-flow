---
id: I-0068
title: The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file
class: efficiency
status: open
count: 3
cost: 3m
first_reported: 2026-10-03T21:21:24Z
last_reported: 2026-10-04T20:35:56Z
updated: 2026-10-04T20:35:56Z
---

# I-0068 The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file

## Description
The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file

## Instances

### 2026-10-03T21:21:24Z
Story: S-0205.
MCP prime for S-0205 returned 142,861 characters against an 80KB budget (exceeded: briefs; the JSON encoding adds to the budgeted text). Claude Code refused to show the result inline and saved it to a file, so the agent had to extract the conventions with jq into a temporary file and read that.

### 2026-10-04T20:07:01Z
Story: S-0210.
S-0210's prime pack was 152 KB, saved to a file and read back with jq.

### 2026-10-04T20:35:56Z
Story: S-0211.
S-0211's pack was 126029 bytes against an 81920 budget (exceeded: briefs); MCP prime returned it as a saved file, read back with jq

## Remediation

Story S-0261 remediates this issue, created from it at 2026-10-04T20:34:55Z.
