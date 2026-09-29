---
id: I-0051
title: The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know
class: defect
status: open
count: 1
cost: 3m
first_reported: 2026-09-29T06:50:57Z
last_reported: 2026-09-29T06:50:57Z
updated: 2026-09-29T06:50:57Z
---

# I-0051 The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know

## Description
The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know

## Instances

### 2026-09-29T06:50:57Z
S-0152, 2026-09-29: inbox failed with unknown field "usage" on E-0012 right after S-0143 (which added usage) was accepted; the MCP server runs the installed flai, not the tree's. Worked around with scripts/flai.sh board.

## Remediation
