---
id: TH-0801
title: Add the formula role step's charge to S-0425's criterion 2, as confirmed on TH-0800
anchor:
  path: wip/kanban/stories/S-0425-metrics-md-generalizes-strategic-usage-per-kind-to-usage-per-role-by-an-adr-so-a-project-s-own-role-is-charged-on-its-items-reported-by-flai-stats-and-charted-as-the-built-in-roles-are.md
  item: S-0425
status: resolved
participants: [orchestrator, agent-S-0425]
created: 2026-11-02T03:43:11Z
updated: 2026-11-02T18:56:10Z
---

# TH-0801 Add the formula role step's charge to S-0425's criterion 2, as confirmed on TH-0800

On wip/kanban/stories/S-0425-metrics-md-generalizes-strategic-usage-per-kind-to-usage-per-role-by-an-adr-so-a-project-s-own-role-is-charged-on-its-items-reported-by-flai-stats-and-charted-as-the-built-in-roles-are.md.

## Entries

### 2026-11-02T03:43:11Z orchestrator
**Recommended:** append this to the end of S-0425's criterion 2: "including a formula's role step, charged to the instance its run was started for".

**Why:**

- You confirmed this on TH-0800 (D-0457) at 2026-11-02T03:42:04Z. The usage-charge clause came out of S-0422's criterion 2 and moves here.
- S-0422's agent removed the clause from S-0422 and recorded the move in its Notes. S-0425 is finalized, so editing its criteria is yours, and S-0425's criterion 2 does not include the clause yet.
- Nothing needs to change in S-0425's tasks. T-1791 already charges a defined role's run to the item it was started for, and S-0422's T-1775 makes the instance that item.

**Alternative:** S-0425's agent adds the clause when it reviews its tasks on pull. This thread stays open so that its inbox shows it.

S-0425 is held behind S-0423 and S-0424, so nothing is urgent.

### 2026-11-02T18:53:52Z agent-S-0425
Done: S-0425's criterion 2 now ends "including a formula's role step, charged to the instance its run was started for", as the operator confirmed on TH-0800 (D-0457). T-1791 covers it with no change to the plan: a defined role's run is charged under the `item` rule to the item it was started for, which T-1775 makes the instance.

Source: design/digest/D-0457-criterion-2-s-usage-charge-of-a-formula-s-role-run-needs-s-0425-s-per-role.md

### 2026-11-02T18:56:10Z orchestrator
Resolved: Done: agent-S-0425 added the formula role step's charge to S-0425's criterion 2, as confirmed on TH-0800 (D-0457).
