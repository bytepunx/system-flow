---
id: TH-1254
title: "Plan for S-0807: five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0807-a-story-s-move-needs-the-suites-its-transition-s-verify-lists-to-have-passed-at.md
  item: S-0807
status: resolved
participants: [planner-S-0807, orchestrator]
created: 2026-11-16T23:50:22Z
updated: 2026-11-16T23:56:20Z
---

# TH-1254 Plan for S-0807: five tasks in three layers

On wip/kanban/stories/S-0807-a-story-s-move-needs-the-suites-its-transition-s-verify-lists-to-have-passed-at.md.

## Entries

### 2026-11-16T23:50:22Z planner-S-0807
I planned S-0807 as five tasks in three layers. The story now has 34 file touches, with no folder touches, a forecast of 1h30m, and a cost of delay of 140.92 USD a week.

## Tasks and layers

| Layer | Task | What it does | After |
|-------|------|--------------|-------|
| 1 | T-3454 | The verify report records each suite's tiers, head, and base. It merges suites per commit and reuses them across a clean rebase (ADR-0156). It adds the leaf package `flai/internal/verifyrecord`. | none |
| 2 | T-3457 | Adds the rule `verified` to the vocabulary and applies it in `Repo.Move`. | T-3454 |
| 2 | T-3460 | `flai verify` and the close-out job run the suites that the transitions out of the story's state need, and `--all` runs every suite. | T-3454 |
| 3 | T-3462 | A baseline regression test: what the close-out runs and the move to review's `warn` stay as they were before this story. | T-3457, T-3460 |
| 3 | T-3465 | The conventions and their template copies, `workflow.md`, `flai-cli.md`, the user guide, the change-the-board runbook, the reference, and the template changelog. | T-3457, T-3460 |

T-3457 and T-3460 share no path, so they can run together. So can T-3462 and T-3465.

## Assumptions

- The report's types and reader move to a new leaf package, `flai/internal/verifyrecord`. `verify` imports `workitem`, so `Repo.Move` cannot import `verify`. A leaf package also means none of the ten callers of `Move` has to be wired to a callback.
- A run of fewer suites at one commit merges into that commit's record. Today a later run replaces the whole record, which would erase a suite that already passed and make the move refuse.
- A repository with no git (`repo.Git` nil) skips `verified`, as it skips `clean-worktree`. A story with git but no record fails the rule.
- A workflow with no `verify:` on the transitions out of a state runs every suite. Together with S-0806's baseline (`verify: all` at `warn`), that keeps today's close-out. A tier's `queue` and `when` keep their meaning until S-0813.
- `after` stays S-0806 alone. S-0805's `suite` reaches this story through S-0806.

## Notes for you

- `flai check` warns `wip.after-hub` on S-0806: five open stories name it in `after`. S-0807 does read S-0806's parsed `verify:` key, so its `after` stays.
- S-0804, in progress, touches the verify package, the close-out, and the serve job. I set no `after` for that, since a trial merge judges the overlap (ADR-0230).
- I would not split, merge, or drop any story or task.

### 2026-11-16T23:55:04Z orchestrator
Approved. The plan fits S-0807.

- The five tasks cover the five criteria. T-3457 and T-3460 share no path, and neither do T-3462 and T-3465.
- 34 file touches, no folder. The forecast is 1h30m, raised from flai's 1h12m on four sibling stories of E-0034. The value of 140.92 USD a week is flai's figure.
- The leaf package `flai/internal/verifyrecord` is a design choice inside the story's scope.
- A workflow with no `verify:` running every suite, together with the baseline's `verify: all` at `warn`, keeps today's close-out. That is criterion 4 and ADR-0312 §8.
- The overlap with S-0804 is left to the trial merge, as ADR-0230 says.

I finalize S-0807 now. It stays in the backlog, held by `after` on S-0806.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-16T23:56:20Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0807.
