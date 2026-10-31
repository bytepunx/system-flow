---
id: D-0411
title: "Until S-0419 generates the strategic agents' definitions from their role files, a change to a strategic agent's prompt is made in both .claude/agents/<name>.md and design/workflows/roles/<name>.md, and in their template copies"
kind: decision
sources: [S-0502, S-0536, S-0537, TH-0760, S-0419]
topics: [orchestration, planning, analysis, template, conventions]
paths: [".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md]
date: 2026-10-28
status: absorbed
absorbed_by: "design/system/strategic-agents.md#role-files-s-0414"
asked_again: []
cited: [S-0419, S-0537, S-0553, T-2313, T-2330, T-2338, T-2350, T-2352, T-2353, TH-0755, TH-0760, TH-0768]
---

# D-0411 Until S-0419 generates the strategic agents' definitions from their role files, a change to a strategic agent's prompt is made in both .claude/agents/<name>.md and design/workflows/roles/<name>.md, and in their template copies

S-0414's acceptance (ADR-0207) put each strategic agent's definition prompt into its role file under `design/workflows/roles/` as well, while ADR-0207 § 4 keeps `.claude/agents/*.md` as the definitions until they are generated from the role files. S-0502 changed the orchestrator's prompt and recorded: "S-0414's acceptance (ADR-0207) put the orchestrator's definition prompt in `design/workflows/roles/orchestrator.md` too; S-0419 will generate `.claude/agents/orchestrator.md` from it. Both role files gained the same item 11, and T-2145's touches widened to them."

S-0419 does not generate `.claude/agents/*.md`. On TH-0760 the operator confirmed option A (D-0429): S-0419 passes each built-in role's file body, with its additions, as the definition in `--agents`, and writes nothing under `.claude/`. Generating and checking the `.claude/agents` definitions is the follow-up story D-0351 named. The orchestrator's confirmed recommendation said: "Until the generation story lands, D-0411 still holds: a change to a baseline prompt is made in both `.claude/agents/<name>.md` and the role file."

So until that generation story lands, a story that changes a planner's, orchestrator's, analyzer's, or curator's prompt changes four files: the `.claude/agents/<name>.md` pair, which is protected and asked of the operator (ADR-0106), and the `design/workflows/roles/<name>.md` pair. It also noted that a permission which gates no call of the role, such as `confirm_after`, is not named in the role file's permissions list.

S-0536 changed the orchestrator's prompt the same way and put the protected pair in a task of its own, so the rest could close while the write waited on the operator: "The `.claude/agents/orchestrator.md` pair moved from T-2303 to its own task, T-2330. T-2303's code and role files could then close while the protected write waits on TH-0737." The operator answered TH-0737, and T-2330 wrote both copies.

S-0537 did the same for the orchestrator's priming step: "The protected `.claude/agents/orchestrator.md` pair is T-2338, a task of its own after T-2313 (D-0411), so that the rest can close while the write waits on the operator."

S-0419 was accepted on 2026-10-31 as option A: each built-in role's definition prompt is now its role file's body, and the description, tools, and model still come from `.claude/agents/<name>.md`. The generation story is S-0553, which E-0024's replan wrote after the operator asked for it on TH-0760; until it lands, this entry's rule holds, and `design/system/strategic-agents.md` § Role files (S-0414) states it. flai's embedded baseline copy under `flai/internal/roles/baseline/` must match the template's role file as well (D-0433).

Source: S-0502 (wip/archive/agents/S-0502.md, Decisions); S-0536 (wip/archive/agents/S-0536.md, Decisions); S-0537 (wip/archive/agents/S-0537.md, Decisions); TH-0760; S-0419 (wip/archive/agents/S-0419.md, Decisions).
