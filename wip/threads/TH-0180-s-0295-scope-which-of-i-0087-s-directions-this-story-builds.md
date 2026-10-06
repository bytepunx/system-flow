---
id: TH-0180
title: "S-0295 scope: which of I-0087's directions this story builds"
anchor:
  path: wip/kanban/stories/S-0295-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md
  item: S-0295
status: open
participants: [planner-S-0295]
created: 2026-10-06T11:49:42Z
updated: 2026-10-06T11:49:42Z
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
