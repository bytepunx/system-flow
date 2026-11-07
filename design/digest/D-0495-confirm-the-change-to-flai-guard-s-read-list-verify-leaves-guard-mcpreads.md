---
id: D-0495
title: "Confirm the change to flai guard's read list: verify leaves guard.MCPReads: Recommended: yes"
kind: answer
sources: [TH-0849, I-0240, S-0556]
topics: [guard, testing, flai, cli, go, template, conventions, code, mcp, orchestration, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/internal/roles/baseline/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/issues/I-0240-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md, design/issues/summary.md, design/issues/I-0277-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md, flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md, docs/users/flai.md]
date: 2026-11-07
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0495 Confirm the change to flai guard's read list: verify leaves guard.MCPReads: Recommended: yes

The operator, alex, confirmed on TH-0849 at 2026-11-07T06:51:54Z the orchestrator's recommendation to narrow `flai guard`'s handling of the close-out job, which S-0556 then shipped and which closed I-0240.

What was asked, by agent-S-0556:

```text
Do you confirm that `verify` leaves `guard.MCPReads`, and that `flai verify --job` is refused to sub-agents? I recommend yes.
```

The answer, the orchestrator's recommendation the operator confirmed:

```text
**Recommended: yes.** Confirm that `verify` leaves `guard.MCPReads` and that sub-agents are refused `flai verify --job`.
```

What changed (I-0240's fix): the MCP tool `verify` left `guard.MCPReads`. `flai guard` refuses every sub-agent, and every strategic agent but the orchestrator, that tool and `flai verify --job`, since the close-out job commits in the story's worktree (ADR-0229). The story's agent asks for the job; the orchestrator holds `verify` by name in its role file and in `guard.MCPOrchestrates`. A reader asking at a commit the record already holds uses `flai verify --last`, which runs nothing, and plain `flai verify` still passes as a read. The two calls D-0446 asked S-0556 to keep both stay. A project whose `orchestrator.md` predates this gets `roles.definition` from `flai check` until `flai upgrade` replaces the baseline (template/CHANGELOG.md).

Sources: TH-0849, I-0240, S-0556.
