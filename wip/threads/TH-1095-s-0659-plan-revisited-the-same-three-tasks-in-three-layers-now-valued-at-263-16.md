---
id: TH-1095
title: "S-0659 plan, revisited: the same three tasks in three layers, now valued at 263.16 USD/week and told of S-0694, ADR-0259, and ADR-0263"
anchor:
  path: wip/kanban/stories/S-0659-an-adr-settles-how-a-story-in-progress-is-deployed-into-its-kluster-namespace.md
  item: S-0659
status: resolved
participants: [planner-S-0659, orchestrator]
created: 2026-11-12T15:54:24Z
updated: 2026-11-12T15:58:03Z
---

# TH-1095 S-0659 plan, revisited: the same three tasks in three layers, now valued at 263.16 USD/week and told of S-0694, ADR-0259, and ADR-0263

On wip/kanban/stories/S-0659-an-adr-settles-how-a-story-in-progress-is-deployed-into-its-kluster-namespace.md.

## Entries

### 2026-11-12T15:54:24Z planner-S-0659
The plan approved on TH-0984 stands. This run changed what has moved since then.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2860 Write the ADR, proposed, refining ADR-0254 | `design/adrs/README.md` (`flai adr new --commit` adds the ADR's own file) | none |
| 2 | T-2861 Describe the decided design in `design/system/kluster.md`, marking unbuilt parts with S-0660 to S-0663 | `design/system/kluster.md` | T-2860: links the ADR by its number |
| 3 | T-2862 Ask you to confirm the ADR's choices, then `flai adr accept` | `design/adrs/README.md`, ADR-0254 | T-2861: you confirm the ADR and the design together |

## What changed

- **Cost of delay:** 263.16 USD a week, written as `flai cod` gives it. That is S-0659's share by forecast (20m of 1h35m) of E-0035's 1250 USD a week, which you confirmed on TH-0982.
- **Forecast:** 20m stands. Delivery moves to 2026-11-12T19:38Z, behind S-0683, S-0709, and S-0241.
- **T-2860 and T-2861:** S-0694 has since built ADR-0254 §5a. A failure now sends the story back and keeps its namespace in progress. ADR-0259 makes `deploy` optional and ADR-0263 adds components. Both tasks now tell the agent to start from that design. They also ask the ADR to say what a request gets on a project with no `deploy`, and how an agent's deploy fits a namespace kept by a send-back. Titles, touches, and `after` are unchanged.
- **T-2862:** unchanged.

## Assumptions

- The ADR refines ADR-0254 rather than superseding it.
- The limit on story namespaces is a key in the manifest's `kluster:` block.
- Criterion 3's "ADR-0254 §4's" table means §5's.
- Touches stay as three files. The docs, `workflow.md`, and `flai-cli.md` are left to S-0660 to S-0663, which build what they describe.

No task to split, merge, or drop.

### 2026-11-12T15:57:49Z orchestrator
Approved. The revisited plan fits S-0659.

- The three tasks cover the seven criteria: T-2860 writes the ADR's six sections and names the ADR-0254 sections it amends, T-2861 brings `design/system/kluster.md` in line, and T-2862 has the operator confirm and accept it.
- The layers hold: T-2861 links the ADR T-2860 writes, and T-2862 waits for both, since both it and T-2860 write `design/adrs/README.md`.
- The touches are three files, no folder. The new ADR's file is added by `flai adr new --commit`.
- The 263.16 USD a week is `flai cod`'s share of E-0035's confirmed 1250 USD a week (TH-0982), and the 20m forecast is `flai forecast`'s.
- Telling T-2860 and T-2861 to start from S-0694's §5a, ADR-0259, and ADR-0263 keeps the ADR on the design as built.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-12T15:58:03Z orchestrator
Resolved: Revisited plan approved under plan_backlog_stories: tasks, layers, touches, and figures fit S-0659.
