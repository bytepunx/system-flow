---
id: TH-0223
title: "S-0300 plan revisited: T-1123 added for the planner texts the plan missed"
anchor:
  path: wip/kanban/stories/S-0300-the-planner-agent-should-be-available-on-epic-work-items.md
  item: S-0300
status: resolved
participants: [planner-S-0300, agent-S-0300]
created: 2026-10-06T23:08:06Z
updated: 2026-10-06T23:28:59Z
---

# TH-0223 S-0300 plan revisited: T-1123 added for the planner texts the plan missed

On wip/kanban/stories/S-0300-the-planner-agent-should-be-available-on-epic-work-items.md.

## Entries

### 2026-10-06T23:08:06Z planner-S-0300
I revisited S-0300's plan while its agent works it. Layer 1 is done.

**One task added: T-1123.** Three texts still say that a planner run on an epic drafts only the stories. No task covered them:

- `flai/cmd/plan.go`, the help for `flai plan`;
- `flai/internal/mcpserver/plan.go`, the description of the MCP tool `plan`;
- `design/system/conventions.md`, the `plan` row of the roles table.

T-1123 updates all three and regenerates `docs/users/flai-reference.md` from the help.

**Tasks and layers:**

- **Layer 1, done:** T-1047 (design and convention), T-1048 (epic prompt), T-1049 (item page test).
- **Layer 2:** T-1050 (both `planner.md` files, waits on your allow), T-1051 (user and operator docs), T-1123 (help, MCP description, conventions design, reference). Each waits for T-1047, and no two share a path.

**Figures:**

- **Touches:** I added those four files to the story and kept every touch it already had. `flaiover/src` stays a folder touch because it was declared.
- **Forecast:** 45m, up from 40m for T-1123. flai gives 24m. Delivery is 2026-10-06T23:50Z.
- **Cost of delay:** 600 USD a week stands, from your 4h a cycle.

**Assumptions:**

1. The help and the MCP description should match the new behaviour now, in this story. Otherwise `flai plan --help` would contradict the convention once S-0300 lands.
2. S-0261, which is in progress, claims `flai/internal/mcpserver` and `docs/users/flai-reference.md`. T-1123 changes one constant there and regenerates the reference, so if they conflict, regenerating settles it.
3. ADR-0075 and ADR-0082 also say that an epic's planner drafts stories. They stay as they are, because accepted ADRs are not edited and this extends them without reversing them.

I would not split, merge, or drop any task. No answer is needed unless you disagree with an assumption.

### 2026-10-06T23:10:51Z planner-S-0300
I checked the plan again at 23:12Z, after the story moved to in-progress and T-1051 was done. Nothing changes:

- **Tasks:** none added, split, merged, or dropped. Layer 1 (T-1047, T-1048, T-1049) and T-1051 are done. T-1050 and T-1123 remain in layer 2.
- **Touches:** unchanged. A `grep` over `design/system`, `docs`, and the planner's code paths found no other text saying that an epic's planner drafts only stories. `flai-cli.md` and `flaiover-dashboard.md` lead the co-change list, but they only say that the planner starts for an epic or a story.
- **Forecast:** still 45m, delivery 23:50Z. flai's own figure rose from 24m to 29m because the touches grew to 18, so I updated the basis.
- **Cost of delay:** 600 USD a week stands.

No answer is needed.

### 2026-10-06T23:28:59Z agent-S-0300
Resolved: T-1123 kept and done by agent-S-0300: flai plan's help, the MCP plan tool's description, and conventions.md's plan row now say an epic's planner drafts its stories' tasks, and flai-reference.md is regenerated.
