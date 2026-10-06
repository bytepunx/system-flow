---
id: I-0068
title: The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file
class: efficiency
status: open
count: 7
cost: 3m
first_reported: 2026-10-03T21:21:24Z
last_reported: 2026-10-06T18:12:10Z
updated: 2026-10-06T18:12:10Z
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

### 2026-10-05T06:26:24Z
Story: S-0217.
S-0217's MCP prime was 164,036 characters; Claude Code saved it to a file instead of returning it, so the agent primed from the story and its task bodies

### 2026-10-05T07:10:40Z
Story: S-0218.
S-0218's MCP prime pack was 190,501 characters; read back from the saved file with jq

### 2026-10-06T02:59:10Z
Story: S-0219.
S-0219's MCP prime pack was 176,888 characters; Claude Code saved it to a file instead of returning it

### 2026-10-06T18:12:10Z
Story: S-0226.
S-0226's MCP prime pack was 183,880 characters; read back from the saved file with jq

## Remediation

Story S-0261 remediates this issue, created from it at 2026-10-04T20:34:55Z.
