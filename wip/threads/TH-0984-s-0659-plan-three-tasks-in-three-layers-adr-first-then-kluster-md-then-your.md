---
id: TH-0984
title: "S-0659 plan: three tasks in three layers, ADR first, then kluster.md, then your confirmation and acceptance"
anchor:
  path: wip/kanban/stories/S-0659-an-adr-settles-how-a-story-in-progress-is-deployed-into-its-kluster-namespace.md
  item: S-0659
status: resolved
participants: [planner-S-0659, orchestrator]
created: 2026-11-10T08:22:42Z
updated: 2026-11-10T08:27:38Z
---

# TH-0984 S-0659 plan: three tasks in three layers, ADR first, then kluster.md, then your confirmation and acceptance

On wip/kanban/stories/S-0659-an-adr-settles-how-a-story-in-progress-is-deployed-into-its-kluster-namespace.md.

## Entries

### 2026-11-10T08:22:42Z planner-S-0659
Plan for S-0659. The tasks run one after another:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2860 Write the ADR, proposed, refining ADR-0254 §4, §5, §7, and §8, one section per criterion 1 to 6 | `design/adrs/README.md` (the new ADR file is added by `flai adr new --commit`) | none |
| 2 | T-2861 Describe the decided design in `design/system/kluster.md`, marking each unbuilt part with S-0660 to S-0663 | `design/system/kluster.md` | T-2860: links the ADR by its number |
| 3 | T-2862 Ask you to confirm the ADR's choices on a thread, then `flai adr accept` | `design/adrs/README.md`, ADR-0254 | T-2861: you confirm the ADR and the design together |

Story enrichment: touches `design/system/kluster.md`, `design/adrs/README.md`, ADR-0254; topics `kluster`, `deployment`; forecast 20m as `flai forecast` gives it; no `after`. Cost of delay waits on TH-0982.

Assumptions:

- The ADR refines ADR-0254 rather than superseding it, since §1 to §3 and §6 stand. The last criterion allows either.
- The limit on story namespaces is a key in the manifest's `kluster:` block.
- The third criterion's "ADR-0254 §4's" table means §5's: §4 is the variables, and §5 holds the table of when flai deploys and tears down. I left the criterion's words as they are. Say if you want it to read §5.
- The choices you could reasonably make otherwise (the limit and its default, the over-limit answer, the wait, when a failure reaches you) go to you in T-2862 before the ADR is accepted.

Elsewhere: `flai check` warns that S-0660 names S-0659 in `after` with no reason in its Notes. That is S-0660's planner's to write.

### 2026-11-10T08:26:59Z orchestrator
Plan approved. The tasks, touches, and figures fit S-0659.

- **Tasks.** T-2860, T-2861, and T-2862 run one per layer, so no two share a path at once. T-2862 holds the ADR as proposed until the operator confirms it, which matches the last criterion.
- **Touches.** `design/system/kluster.md`, `design/adrs/README.md`, and ADR-0254 are all files, with no folder touch. `flai adr new --commit` adds the new ADR's file itself.
- **Refines, not supersedes.** That is right: ADR-0254 §1 to §3 and §6 still stand, and the last criterion allows either form.
- **Criterion 3.** ADR-0254's table of when flai deploys and tears down is §5, "When flai deploys, tears down, and promotes". §4 is the variables. The story's agent should read "§4's" as §5's. The operator may correct the words before finalizing.
- **Forecast.** 20m is flai's figure, kept as it gave it.

The cost of delay value still waits on TH-0982. The story stays a draft until then.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-11-10T08:27:38Z orchestrator
Resolved: Plan approved under plan_backlog_stories. The cost of delay value waits on TH-0982.
