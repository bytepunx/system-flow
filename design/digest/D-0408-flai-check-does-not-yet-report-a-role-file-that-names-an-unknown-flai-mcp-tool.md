---
id: D-0408
title: "flai check does not yet report a role file that names an unknown flai MCP tool or command; a shipped-roles test holds the template's files to that until S-0416 passes the guard's vocabulary in"
kind: decision
sources: [S-0414, S-0416]
topics: [cli, orchestration, planning]
paths: [flai/internal/check/check.go, flai/cmd/role_test.go, flai/internal/roles/roles.go, flai/internal/guard/roles.go, design/workflows/roles, design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0408 flai check does not yet report a role file that names an unknown flai MCP tool or command; a shipped-roles test holds the template's files to that until S-0416 passes the guard's vocabulary in

S-0414 decided, in its narrative's Decisions (wip/archive/agents/S-0414.md), that `flai check` holds role files to the guard's reads but does not check that each MCP tool or command a file names exists.

The reason: the check package cannot import the MCP server or the command tree, which import it. A check that knew them in one caller only would differ by where it runs.

`TestShippedRolesNameFlaisToolsAndCommands` in `flai/cmd/role_test.go` holds the shipped role files to both. A project's own role file that names an unknown call is not flagged.

S-0416 was expected to pass that vocabulary in once the guard read the role files. It was accepted on 2026-10-29 and did not: its Decisions (wip/archive/agents/S-0416.md) do not mention it. In the code after its acceptance, `check.RolesVocabulary` in `flai/internal/check/check.go` still gives `roles.Vocabulary` its `Reads` alone. The guard validates a role's file with an empty `roles.Vocabulary{}` (`flai/internal/guard/roles.go`). A nil `Tools` or `Commands` skips the name check (`flai/internal/roles/roles.go`). So the gap stays open for a later story to close.

This narrows ADR-0207 §5, "a name that is no flai MCP tool or command is `roles.file`". `design/system/strategic-agents.md` § Role files (S-0414) still lists an unknown tool or command among the `roles.file` findings.
