---
id: TH-0180
title: "S-0295 scope: which of I-0087's directions this story builds"
anchor:
  path: wip/kanban/stories/S-0295-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md
  item: S-0295
status: resolved
participants: [planner-S-0295, alex]
created: 2026-10-06T11:49:42Z
updated: 2026-10-06T12:21:44Z
---

# TH-0180 S-0295 scope: which of I-0087's directions this story builds

On wip/kanban/stories/S-0295-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md.

## Entries

### 2026-10-06T11:49:42Z planner-S-0295
I-0087 leaves the remedy undesigned, so I need your choice of direction before I write S-0295's tasks.

**Recommended: build all three in S-0295, as one new ADR that refines ADR-0046.**

1. **A story in review no longer holds.** Only in-progress stories hold ready ones. A branch in review is finished and synced. When it is accepted, the acceptance notice and `flai stream sync`'s trial merge tell an overlapping story what changed. This change is in `flai/internal/workitem/hold.go`. It alone would have freed 2h45m of S-0220's four hours.
2. **Shared paths do not hold.** A manifest list, defaulting to `docs/users/flai.md`, `docs/users/flai-reference.md` (generated), `design/system/flai-cli.md`, `design/adrs`, `design/issues`, and `template/CHANGELOG.md`, names paths that many stories change in separate sections or new files. An overlap only on those paths does not hold a story. The trial merge at sync still reports a real textual conflict. This is cheaper than splitting `docs/users/flai.md` (1408 lines) by command, and it also covers `design/adrs`, which holds S-0278 and S-0284 today.
3. **Tasks narrow a folder claim.** When a story's tasks name files inside one of its folder touches, the claim holds those files (done tasks included, because the branch changed them), not the folder. The planner and the story's agent are told to declare file-level touches where they can. This helps ready stories with drafted tasks as well as open ones.

**Alternatives:**

- (b) S-0295 builds 2 and 3, and I draft a separate story for 1. 1 is a policy change, and S-0296 (review waits for the operator) is close to it.
- (c) Split `docs/users/flai.md` into one file per command area instead of 2.
- (d) Any subset you name.

Also: is the default shared list in 2 right, or do you want it empty and set per project?

Until you answer, I am setting S-0295's cost of delay value (780 USD a week, from flai's inputs). I'll write the touches, forecast, and tasks once the scope is settled.

### 2026-10-06T12:11:08Z alex
I would allow for operators on other projects to edit the list. let’s add tasks for extending flais cli, http, mcp APIs for checking and managing that lists entries, ui in project settings for changing it, and glob interpretation to make it easy for the operator

### 2026-10-06T12:18:29Z planner-S-0295
Thank you. S-0295 is planned on that answer. This thread is the plan's summary.

**Tasks, in seven layers:**

1. T-1023: an ADR refines ADR-0046 with the three rules, and the design says so. It proposes the key `claims.shared`, the command `flai shared list|add|remove|check`, the host method `settings.shared`, and the MCP tools `shared_paths` and `shared_paths_edit`.
2. T-1024, T-1025, and T-1026 run together.
   - T-1024: a story in review holds nothing.
   - T-1025: the manifest key, glob validation, a matcher, and one edit function the CLI, HTTP, and MCP share.
   - T-1026: the planner and the story's agent declare file-level touches (conventions with their template copies, and the planner's prompt).
3. T-1027, T-1028, and T-1029 run together.
   - T-1027: tasks narrow a folder claim. It waits for T-1024, which also changes `hold.go`.
   - T-1028: the CLI.
   - T-1029: the MCP tools, with `flai guard` limiting the edit tool to your session.
4. T-1030 and T-1031 run together.
   - T-1030: an overlap inside a shared path does not hold, is not a `wip.overlap`, and does not count as a grown claim.
   - T-1031: the host API reports, checks, and edits the list through `flai shared`.
5. T-1032: the project settings panel lists, adds, removes, and checks the shared paths, with a one-line glob hint.
6. T-1033: the users' and operators' guides and the dashboard design.
7. T-1034: a test rebuilds I-0087's board, fails if any one rule is reverted, and I-0087 is closed.

**Assumptions. Correct any of these:**

- A plain path in the list covers everything below it, `*` stays within one folder, and `**` crosses folders. `design/adrs` and `design/adrs/**` therefore mean the same.
- This project's list is the one I proposed. New projects from the template start with the same list, which you can edit there.
- Only your session may change the list through MCP. A story's agent, its sub-agents, the planner, and the analyzer may only read and check it, because an agent that edits the list could free its own story. The orchestrator gets the edit tool only if the ADR grants it.
- The trial merge at sync and the overlap notice at acceptance still report on shared paths, so a real conflict there is still caught.
- A done task's files stay in the narrowed claim, and a cancelled task's do not.

**Figures:**

- Forecast: 2h. flai said 1h24m; I raised it for the seven mostly sequential layers.
- Delivery: 18:34Z.
- Cost of delay: 780 USD a week, as `flai cod` computed it.
- Touches: 54, files wherever a task can name them. The only folders are `design/adrs`, until `flai adr new` gives the ADR its number, and `flai/internal/guard`.

**Related:** S-0296 (review waits for the operator) is now ready. T-1024 removes most of what its idle board costs, because a story in review no longer blocks the rest. S-0296 still matters for the review queue itself. I changed nothing on it.

### 2026-10-06T12:21:42Z alex
Resolved.

Resolved.

Resolved.

Resolved.

Resolved.

Resolved.

### 2026-10-06T12:21:43Z alex
Resolved.

Resolved.

Resolved.

Resolved.

Resolved.

Resolved.

### 2026-10-06T12:21:44Z alex
Resolved.
