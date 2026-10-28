---
id: D-0408
title: "flai check does not yet report a role file that names an unknown flai MCP tool or command; a shipped-roles test holds the template's files to that until S-0416 passes the guard's vocabulary in"
kind: decision
sources: [S-0414, S-0416]
topics: [cli, orchestration, planning]
paths: [flai/internal/check/check.go, flai/cmd/role_test.go, design/workflows/roles, design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0408 flai check does not yet report a role file that names an unknown flai MCP tool or command; a shipped-roles test holds the template's files to that until S-0416 passes the guard's vocabulary in

S-0414 decided, in its narrative's Decisions (wip/archive/agents/S-0414.md), that `flai check` holds role files to the guard's reads but does not check that each MCP tool or command a file names exists.

The reason: the check package cannot import the MCP server or the command tree, which import it. A check that knew them in one caller only would differ by where it runs.

`TestShippedRolesNameFlaisToolsAndCommands` in `flai/cmd/role_test.go` holds the shipped role files to both. A project's own role file that names an unknown call is not flagged until the guard reads the role files in S-0416, which can pass the vocabulary in and do the full check.

This narrows ADR-0207 §5, "a name that is no flai MCP tool or command is `roles.file`". `design/system/strategic-agents.md` § Role files (S-0414) still lists an unknown tool or command among the `roles.file` findings.
