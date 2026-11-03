---
id: D-0474
title: The strategic agents' .claude/agents links dangle on every fresh checkout and story worktree until flai role definition --write runs there, so a walk of the repository's markdown leaves a dangling link out, as markdownlint-cli2 does
kind: decision
sources: [S-0572, S-0553]
topics: [testing, ci, markdown, conventions]
paths: [".claude/agents/planner.md", ".claude/agents/orchestrator.md", ".claude/agents/analyzer.md", ".claude/agents/curator.md", flai/internal/mdlint/repo_test.go]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0474 The strategic agents' .claude/agents links dangle on every fresh checkout and story worktree until flai role definition --write runs there, so a walk of the repository's markdown leaves a dangling link out, as markdownlint-cli2 does

Since S-0553 ([ADR-0236](../adrs/0236-a-built-in-role-s-claude-code-definition-is-generated-into-flai-cache-agents-by.md)), `.claude/agents/<name>.md` for the planner, the orchestrator, the analyzer, and the curator are committed relative links into `.flai-cache/agents/`. That folder is gitignored, and only `flai role definition --write` writes it. So the links dangle on any fresh checkout, CI's included, and in a story worktree, until that command runs there.

`TestRepositoryLintsClean` in `flai/internal/mdlint/repo_test.go` read them and failed the integration tier and the `flai` workflow. S-0572 fixed it, from its Decisions (`wip/archive/agents/S-0572.md`):

```text
T-2448: a dangling link is left out of the walk, not failed and not treated as a finding: markdownlint-cli2 in `system-flow-check.yml` passes on the same tree, and the test's claim is that mdlint reports nothing markdownlint does not. A link that resolves is still read, as it is in the main checkout once `flai role definition --write` has run. The links themselves stay as ADR-0236 made them.
```

Code or a test that walks the repository's files must expect these four links to dangle. It leaves them out and keeps reading a link that resolves. Claude Code itself only hides an agent whose link dangles (ADR-0236 § Context). In the main checkout, the operator runs `flai role definition --write` after an acceptance (D-0455).

Sources: S-0572, S-0553.
