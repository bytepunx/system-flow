---
id: TH-0599
title: "Plan for S-0468: the remedy proposed from I-0152's instances, six tasks in five layers"
anchor:
  path: wip/kanban/stories/S-0468-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md
  item: S-0468
status: open
participants: [planner-S-0468]
created: 2026-10-22T21:07:12Z
updated: 2026-10-22T21:07:12Z
---

# TH-0599 Plan for S-0468: the remedy proposed from I-0152's instances, six tasks in five layers

On wip/kanban/stories/S-0468-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md.

## Entries

### 2026-10-22T21:07:12Z planner-S-0468
**The remedy proposed.** I-0152's instance is E-0016's planned stories: folder touches (`flai/cmd`, `flai/internal/*`, `flaiover/src/routes`, `docs/users`) and `after` hubs (S-0218, S-0223). What stands already asks for files but enforces nothing: the convention and `planPrompt` say name files (S-0295), ADR-0096 narrows a story's folder to the files its tasks name, ADR-0134 shares. ADR-0096's narrowing works only when tasks name files, and nothing makes them. So the tasks build, through flai's write paths rather than words alone:

1. A task takes no folder touch (a folder of the tree with files in it, read as the claim reads it; a path absent from the tree is a file to be made; a folder inside `claims.shared` is left alone). `flai task new`, `flai edit --touches`, `flai touches`, `item_new`, and `item_edit` refuse it, writing nothing, with the folder, its file count, and what to do.
2. The planner's folder touch on a story is refused unless the story's `### Planning` names it; the operator's and the story's agent's story touches are not refused.
3. `flai check` reports `touches.folder`, a warning, on a story in progress whose claim still keeps a folder, with the ready stories it holds; and `wip.after-hub` when three or more open stories name one story in `after`. A close-out's scoped check drops another story's, recording no issue.
4. The planner's convention, prompt, and `planner.md` say `after` on a story names only a story whose accepted change the story needs on main; the pull order sets preference.

**Tasks and layers** (`after` in brackets):

- Layer 1: T-2055 the ADR and the three design sections.
- Layer 2: T-2056 [T-2055] `workitem.FolderTouch` and the task refusal through every write path.
- Layer 3, together: T-2057 [T-2056] the planner's story refusal with the Planning note; T-2058 [T-2056] the two check rules and the scope.
- Layer 4: T-2059 [T-2057, T-2058] conventions and template copies, `planPrompt`, both `planner.md`, user docs, template version and changelog.
- Layer 5: T-2060 [T-2059] close I-0152, tick the criteria.

**Figures.** Touches: 35, files but `design/adrs` (the ADR's name is `flai adr new`'s; inside `claims.shared`, so it holds nothing). Forecast: flai's 55m raised to 2h, reason in the story's `### Planning`; delivery 2026-10-22T23:30:00Z. Cost of delay: 6600 USD a week from the 44h input, unchanged.

**Assumptions, for the operator or the orchestrator to confirm or redirect.**

- A refusal, not a warning, for a task's folder touch. The convention already asked and was not followed; a refusal is at the write and testable. If too strict for the story's agent, the fallback is a `flai check` warning on tasks alone, and T-2056 shrinks to the helper.
- The planner alone is refused a story's folder touch. The operator's words stay theirs; a story's agent's story touches are narrowed by its tasks under ADR-0096 and reported by `touches.folder` while in progress.
- The hub threshold is three, fixed in the ADR, not a manifest setting.
- An `after` hold is not released when the hub reaches review: the hub's change is not on main, so the dependent branch cannot build on it. Left as an alternative in the ADR.
- Most existing tests create tasks in bare test repositories whose folders hold no files, so the refusal reaches few of them; T-2056 fixes those it does.
- `.claude/agents/planner.md` is in scope because `design/system/strategic-agents.md` says it, the convention, and `planPrompt` say the same, and it lacks the folder rule today; its write goes through the permission prompt, and the story is then the operator's to accept (ADR-0106).

Nothing in the plan waits on an answer; the story's agent reviews the tasks when it pulls the story.
