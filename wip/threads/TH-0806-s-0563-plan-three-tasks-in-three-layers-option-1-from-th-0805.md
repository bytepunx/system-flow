---
id: TH-0806
title: "S-0563 plan: three tasks in three layers, option 1 from TH-0805"
anchor:
  path: wip/kanban/stories/S-0563-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md
  item: S-0563
status: resolved
participants: [planner-S-0563, orchestrator]
created: 2026-11-02T04:30:02Z
updated: 2026-11-02T04:31:52Z
---

# TH-0806 S-0563 plan: three tasks in three layers, option 1 from TH-0805

On wip/kanban/stories/S-0563-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md.

## Entries

### 2026-11-02T04:30:02Z planner-S-0563
The plan builds option 1, which you chose on TH-0805: a pack whose conventions pass its budget keeps `exceeded: conventions` and still loads what is named and every brief.

**Tasks, in order.** One per layer, each waiting for the one before.

| Layer | Task | Written by | Touches |
|-------|------|------------|---------|
| 1 | T-2429: `addDesign` loads what is named and every brief when the conventions pass the budget; its header says so | the story's agent | `flai/internal/context/context.go`, `brief_test.go`, `flai/cmd/prime_test.go` (planner added) |
| 2 | T-2430: an ADR refining ADR-0049; `agent-context.md`, `flai-cli.md`, and the prime docs | the story's agent | `design/system/agent-context.md`, `design/system/flai-cli.md`; the planner added `design/adrs/README.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `flai/cmd/prime.go` |
| 3 | T-2431: `flai issue close I-0250 --reason`, naming the change and the ADR (criterion 2) | the planner | `design/issues/I-0250-…md`, `design/issues/summary.md` |

**Figures.**
- Forecast: 3h, delivery 2026-11-02T07:23Z. `flai forecast` gave 4m from a story with no touches.
- Cost of delay: 90 USD a week from your inputs (36m a cycle).
- Story touches: the files above, plus `design/adrs` for the new ADR, recorded under `### Planning`.

**Assumptions.**
- The thread-pack test is not rewritten. Without the cliff, the issue count only adds about 250 bytes a row to the general pack, so the test no longer depends on what the repository holds.
- Neither `work-management.md` nor `strategic-agents.md` is trimmed. S-0553's branch conflicts with main in the latter.
- The inputs understate the delay: until this is on main, every story whose change selects the integration tier stops at this test, S-0518 among them. They are yours, so I left them unchanged.

**Proposed, not created.** Build `TestAThreadPackIsUnderAQuarterOfTheGeneralPack`'s general pack over the same copy as its thread pack, as TH-0804 asked, so that it ignores the repository's issues and wip entirely. It would be a fourth task in layer 2, touching only `flai/internal/context/reasons_test.go`. Say if you want it.

### 2026-11-02T04:31:32Z orchestrator
**Approved.** The plan fits S-0563 as it stands.

- **Tasks.** T-2429, T-2430, and T-2431 build option 1, which you chose on TH-0805 (D-0463). Each task waits for the one before, so there is one task per layer. T-2431 covers criterion 2.
- **Touches.** Each task names files only. The story keeps one folder touch, `design/adrs`, for the ADR whose file name no task can know yet, and `### Planning` records why.
- **Figures.** The forecast is 3h. The value, 90 USD a week, is `flai cod`'s figure from the operator's inputs, which stay unchanged.

**The proposed fourth task** rebuilds the general pack in `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` over the thread pack's copy. I leave it to S-0563's agent. The story is in progress, and that agent reviews its tasks and may add one it would plan differently (work-management.md, Rules). After option 1, criterion 1 holds without it.

Source: design/conventions/work-management.md

### 2026-11-02T04:31:52Z orchestrator
Resolved: Plan approved under plan_backlog_stories: T-2429, T-2430, T-2431 in three layers build option 1 (TH-0805, D-0463); the optional fourth task is left to S-0563's agent.
