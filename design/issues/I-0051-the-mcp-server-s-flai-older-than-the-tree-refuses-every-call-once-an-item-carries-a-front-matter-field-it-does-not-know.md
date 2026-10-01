---
id: I-0051
title: The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know
class: defect
status: closed
count: 1
cost: 3m
first_reported: 2026-09-29T06:50:57Z
last_reported: 2026-09-29T06:50:57Z
updated: 2026-10-01T10:55:14Z
---

# I-0051 The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know

## Description
The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know

## Instances

### 2026-09-29T06:50:57Z
S-0152, 2026-09-29: inbox failed with unknown field "usage" on E-0012 right after S-0143 (which added usage) was accepted; the MCP server runs the installed flai, not the tree's. Worked around with scripts/flai.sh board.

## Remediation
Closed 2026-10-01T10:55:14Z: S-0181: work items, threads, and issues are decoded leniently on the listing paths; an unknown front-matter field is kept and written back unchanged, logged once per file, and reported by flai check as item/threads/issues.unknown-field. flai.minimum in system-flow.yaml, raised by a flai release that changes the fields, makes an older flai say so instead.
