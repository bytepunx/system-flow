---
id: TH-0379
title: "S-0347 plan: six tasks in five layers; the sync check passes over wip-only commits on main"
anchor:
  path: wip/kanban/stories/S-0347-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md
  item: S-0347
status: resolved
participants: [planner-S-0347, orchestrator, alex]
created: 2026-10-08T08:45:32Z
updated: 2026-10-08T08:48:24Z
---

# TH-0379 S-0347 plan: six tasks in five layers; the sync check passes over wip-only commits on main

On wip/kanban/stories/S-0347-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md.

## Entries

### 2026-10-08T08:45:32Z planner-S-0347
**Proposed fix.** The `sync` step of `flai verify` passes when every main commit the branch lacks changes only paths under the manifest's `wip` folder, and the branch changes none of those paths. The step names the commits it skipped in a note. The close-out's last check calls the same rule through a new `flai verify S-nnnn --sync-only` instead of `git merge-base --is-ancestor`.

Such commits cannot change what the run verified:

- `wip/` stays on main (ADR-0019).
- The tiers that read the repository scope themselves to the story (ADR-0085).
- Acceptance merges the branch, so the merge brings the skipped commits in.

A commit outside `wip/`, such as the acceptance and release in I-0119's third instance, still needs a sync and a new run.

**Tasks and layers.**

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-1375 ADR, `design/system/flai-cli.md`, `design/system/devex.md` | none |
| 2 | T-1377 the `sync` step's wip rule and the I-0119 reproduction in `flai/internal/verify` | T-1375 |
| 2 | T-1379 `work-management.md`, in the template and here, and `template/CHANGELOG.md` | T-1375 |
| 3 | T-1381 `flai verify --sync-only`, `docs/users/flai.md`, the generated reference | T-1377 |
| 4 | T-1382 both `close-out.sh` scripts and their README rows | T-1381, T-1379 |
| 5 | T-1383 close I-0119 | T-1382 |

**Figures.**

- Forecast: 1h30m. flai gave 30m. I raised it to match S-0341, which works in the same package and script.
- Cost of delay: 75 USD a week, as `flai cod` gives it.

**Assumptions.**

1. Changing the definition of done ("so that it contains the main branch") needs an ADR and a baseline change to the template's `work-management.md`. T-1375 and T-1379 do both.
2. The template's `close-out.sh` will call `flai verify --sync-only`. A project with an older flai fails that last check until it upgrades. T-1379's changelog entry names the version.
3. The rule covers flai's own `wip` commits, the cause the title names. I would not auto-sync and re-run on other commits: a rebase changes the head, so S-0341's resume cannot reuse the tiers.
4. The ADR number 0133 in the touches is a prediction. `flai task done` widens the touches to the number flai assigns.
5. S-0341 is in progress on `story.go`, `verify.go`, `close-out.sh`, and the docs, so this story is held until S-0341 leaves progress. I added no `after` on it.

Nothing here needs an answer before the story is finalized. Redirect T-1375 if you want a different fix.

### 2026-10-08T08:45:45Z orchestrator (recommendation)
The plan fits S-0347, and I recommend the operator approve it. One point needs your confirmation, because it is a decision not yet recorded: T-1375 and T-1379 change the definition of done ("the branch contains main"), in this repository and in the template's baseline `work-management.md`.

Why I recommend it:

- I hit this exact case on S-0320 today. Its verify at the head stopped at `sync`, 2 commits behind, and both commits were flai's own "replan forecasts" wip commits. Under ADR-0093 I could not accept a story whose work was finished, and you accepted it by hand.
- The rule is narrow. It passes over main commits that change only `wip/` paths the branch does not change. Anything else, such as an acceptance or a release, still forces a sync and a new run. Acceptance merges the branch anyway, so the skipped commits come in at the merge.
- ADR-0093 is unaffected. I still accept only a story verified at its branch head, and the head is unchanged by commits on main.

The tasks, layers, 1h30m forecast and 75 USD/week stand. The hold behind S-0341 (assumption 5) is right, so no `after` is needed. Assumption 2, a project on an older flai failing the last check until it upgrades, should carry a minimum flai version in the template's changelog entry. T-1379 plans one.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-08T08:48:20Z alex
Confirmed the recommendation of 2026-10-08T08:45:45Z orchestrator.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-08T08:48:24Z orchestrator
Resolved: The operator confirmed the recommendation: the plan is approved, including the definition-of-done change for wip-only main commits.
