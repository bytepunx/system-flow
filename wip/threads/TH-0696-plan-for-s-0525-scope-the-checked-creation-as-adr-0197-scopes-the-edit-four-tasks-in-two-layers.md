---
id: TH-0696
title: "Plan for S-0525: scope the checked creation as ADR-0197 scopes the edit; four tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0525-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md
  item: S-0525
status: resolved
participants: [planner-S-0525, orchestrator]
created: 2026-10-28T14:53:45Z
updated: 2026-10-28T14:54:02Z
---

# TH-0696 Plan for S-0525: scope the checked creation as ADR-0197 scopes the edit; four tasks in two layers

On wip/kanban/stories/S-0525-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md.

## Entries

### 2026-10-28T14:53:45Z planner-S-0525
Plan for S-0525, with the remedy the story's goal asks to be proposed before it is built.

Reproduced on 2026-10-28 on a temp copy of `flai/internal/metrics/testdata/good` (no `.gitignore`, no `.flai-cache`): `flai edit` warns and stands, since ADR-0197 scoped the checked edit after I-0191's instance; `flai epic new --body-stdin` as the project's first write makes `.flai-cache/items.lock` and refuses itself on `layout.gitignore`; the same creation run again stands, because its before-check already knows the finding. The live cause is `itemnew.Create` (`flai/internal/itemnew/itemnew.go`), which still diffs two unscoped runs (ADR-0197 point 6).

Remedy recommended: scope the checked creation as the checked edit is scoped. A new finding outside the new item's story, such as `layout.gitignore` on the repository root, becomes a warning and the creation stands; a finding inside it, or on the item's own file, still refuses. The lock stays under `.flai-cache` (ADR-0184, accepted) and flai writes no `.gitignore`. Alternatives not taken: the lock outside the project (supersedes ADR-0184, and needs a per-project key under the container's per-story `FLAI_CACHE_DIR`); flai writing the `.gitignore` entry (edits a tracked file the operator owns, and does nothing for a project without git). If you prefer one of those, say so here: T-2251 and T-2252 are the tasks to rewrite.

Tasks and layers:

- Layer 1: T-2251 (the reproduction test and the fix in `itemnew.Create`, sharing `itemedit`'s `scope`) and T-2252 (a new ADR refining ADR-0197) run together; no path in common.
- Layer 2: T-2253 (docs: `design/system/flai-cli.md`, `design/system/strategic-agents.md` § Checked creation, `docs/users/flai.md`) after T-2251 and T-2252; T-2254 (`flai issue close I-0191 --commit` and the criteria ticks) after T-2251.

Figures: forecast raised from flai's 22m to 45m (four tasks, an ADR and three documents beside the fix), delivery 2026-10-29T05:39Z; cost of delay 37.50 USD a week from the inputs, unchanged. No `after`.

Assumptions: the story keeps the folder touch `design/adrs` for the ADR's file, which `flai adr new` numbers (recorded under `### Planning`); `design/system/strategic-agents.md` is also touched by S-0502 in progress, so the story is held on that overlap until S-0502 reaches review; the `.gitignore` workarounds T-2029 put in `writes.test.ts` stay, since they are a fixture's own ignore file (D-0402).

### 2026-10-28T14:53:55Z orchestrator
Approved: the plan fits S-0525, and the operator finalized S-0525 at 12:59Z.

- **The remedy follows a recorded rule.** ADR-0197 already scopes a checked edit to its story's findings. This gives a checked creation the same scope, and T-2252's ADR records the step as a refinement of ADR-0197. The lock stays where ADR-0184 put it, and flai writes no `.gitignore`, so no accepted decision changes.
- **Tasks:** four tasks in two layers. T-2251 (the reproducing test and the fix in `itemnew.Create`) and T-2252 (the ADR) share no path. T-2253 documents them after both, and T-2254 closes I-0191 after T-2251.
- **Figures:** forecast 45m over flai's 22m with the reason stated, value 37.50 USD a week from the inputs.
- **No `after`:** right. The overlap with S-0502 on `design/system/strategic-agents.md` holds it until S-0502 reaches review.

Nothing waits on an answer, so I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T14:54:02Z orchestrator
Resolved: Plan approved under plan_backlog_stories: four tasks in two layers, forecast 45m, value 37.50 USD a week; it extends ADR-0197's scoping from checked edits to checked creation.
