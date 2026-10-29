---
id: D-0411
title: "Until S-0419 generates the strategic agents' definitions from their role files, a change to a strategic agent's prompt is made in both .claude/agents/<name>.md and design/workflows/roles/<name>.md, and in their template copies"
kind: decision
sources: [S-0502]
topics: [orchestration, planning, analysis, template, conventions]
paths: [".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: [T-2330]
---

# D-0411 Until S-0419 generates the strategic agents' definitions from their role files, a change to a strategic agent's prompt is made in both .claude/agents/<name>.md and design/workflows/roles/<name>.md, and in their template copies

S-0414's acceptance (ADR-0207) put each strategic agent's definition prompt into its role file under `design/workflows/roles/` as well, while ADR-0207 § 4 keeps `.claude/agents/*.md` as the definitions until S-0419 generates them. S-0502 changed the orchestrator's prompt and recorded: "S-0414's acceptance (ADR-0207) put the orchestrator's definition prompt in `design/workflows/roles/orchestrator.md` too; S-0419 will generate `.claude/agents/orchestrator.md` from it. Both role files gained the same item 11, and T-2145's touches widened to them."

So until S-0419 lands, a story that changes a planner's, orchestrator's, analyzer's, or curator's prompt changes four files: the `.claude/agents/<name>.md` pair, which is protected and asked of the operator (ADR-0106), and the `design/workflows/roles/<name>.md` pair. It also noted that a permission which gates no call of the role, such as `confirm_after`, is not named in the role file's permissions list.

Source: S-0502 (wip/archive/agents/S-0502.md, Decisions).
