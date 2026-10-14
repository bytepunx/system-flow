---
id: TH-0420
title: "S-0371 plan: the rebase reports a branch that undoes main's change, sync warns and accept refuses (new ADR), six tasks in five layers"
anchor:
  path: wip/kanban/stories/S-0371-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md
  item: S-0371
status: resolved
participants: [planner-S-0371, orchestrator, alex]
created: 2026-10-12T17:17:16Z
updated: 2026-10-14T22:18:09Z
---

# TH-0420 S-0371 plan: the rebase reports a branch that undoes main's change, sync warns and accept refuses (new ADR), six tasks in five layers

On wip/kanban/stories/S-0371-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md.

## Entries

### 2026-10-12T17:17:16Z planner-S-0371
**Recommendation: approve the plan as drafted.** After a clean rebase, `storygit.Rebase` reports each path where the branch removes lines main gained since the old branch point. `flai stream sync` and `flai task done` warn on each one. `flai accept` refuses the branch.

## Remedy

- **Cause (I-0124):** git drops the commit that became empty, then replays the revert cleanly. No conflict shows. Before the rebase the branch's net change is empty, so the check runs after the rebase, against the old merge base read before it.
- **Check:** for each path that main changed since the old base and that the rebased branch still changes, it compares main's added lines with the lines the branch removes. It reports `undoes_main` (path, line count, old base). This covers paths inside the claim, which the outside-touches note skips.
- **Sync and task done:** they warn and name the restore, `git checkout <main> -- <path>` plus a commit. The sync does not fail, as ADR-0046's checks do not.
- **Accept:** it refuses and merges nothing, as it does for conflict markers, because it runs the rebase with no checks.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1476 ADR, refining ADR-0069 and ADR-0098 | `design/adrs/README.md` |
| 2 | T-1477 the check in `storygit` and the I-0124 reproduction | `undone.go`, `undone_test.go`, `sync.go`, `sync_test.go` |
| 3 | T-1478 sync and task done warning | `cmd/stream_sync.go`, `stream_sync_test.go`, `stream.go` |
| 3 | T-1479 accept refusal | `cmd/branch.go`, `accept.go`, `accept_undone_test.go` |
| 4 | T-1480 design and user docs, reference regenerated | `flai-cli.md`, `workflow.md`, `docs/users/flai.md`, `flai-reference.md` |
| 5 | T-1481 close I-0124 | the issue file, `summary.md` |

T-1478 and T-1479 run together. All touches are files.

## Figures

- Forecast 45m, delivery 2026-10-13T02:21Z. flai gave 26m. I raised it because the work is a new line-level git check, two callers that behave differently, and git-fixture tests in two packages.
- Cost of delay: 25 USD a week, from flai's 10m input, unadjusted.

## Assumptions

- **A warning, not a stop, at sync.** A story may mean to remove main's lines. Accept's refusal is the hard gate. If you want the sync to stop as well, T-1478 changes only its exit code.
- **The dry run cannot name the path.** It does not rebase. The refusal comes at the real acceptance, as the second conflict-marker check does.
- **No convention change.** The warning and the refusal say what to do.
- **Held behind S-0367.** S-0367 is in progress and claims `flai/cmd` and `sync.go`, so this story is held until S-0367 is accepted.

I propose no split, merge, or drop.

### 2026-10-12T17:17:40Z orchestrator (recommendation)
**Recommendation:** approve the plan with one addition. `flai accept` refuses a branch that undoes main's change unless it is given an explicit override naming the paths, say `--keep-removal <path>`. The ADR records the override, and the story's evidence names why the removal is meant.

## Why this needs you

- **It adds a refusal to acceptance.** The plan makes `flai accept` refuse, merging nothing, every branch whose rebase removes lines main gained since its old base. ADR-0093 sets the conditions under which I accept, so a new one is a decision for you.
- **Without an override, intended removals block.** Your plan itself says "A story may mean to remove main's lines." An example: a story that rewrites a doc section another story just added. As drafted, neither you nor I could accept it short of restoring the lines and editing them again in a later commit.
- **An override keeps I-0124's protection.** The silent case, a revert replayed over main's identical change, is still caught. Only a removal someone names on the command line passes.

## The rest fits

- The cause analysis and the check after the rebase against the old merge base match I-0124.
- Warning at sync and refusing at accept is the right split.
- Six tasks in five layers, all file touches.
- The forecast raises flai's 26m to 45m, with its reason stated. The value is 25 USD a week from flai's 10m input.
- It is held behind S-0367, as you say.

I leave S-0371 a draft until you answer. If you take the addition, T-1476 and T-1479 carry it.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-13T17:42:13Z alex
could we instead add a mechanism for resolving this through the dashboard? pushing work down into low level git commands or hand-edited files is an anti-pattern. Give the operator tools to resolve the situation and make suggestions for what that might look like for this scenario.

### 2026-10-13T17:42:38Z orchestrator (recommendation)
**Recommendation:** yes. Make "a branch undoes main's change" a blocker of the acceptance preview that you resolve on the dashboard, path by path, with flai doing the git work. Here is what that could look like.

## On the review page

- **Each path is a card in the preview's blockers.** It shows the lines main gained since the branch's old base that the branch now removes, side by side with the branch's version. ADR-0093 § Decision already makes each failed condition a preview blocker, so this joins them.
- **Three actions per path:**
  1. **Restore main's lines.** flai commits main's version of that path onto `story/S-nnnn` (a new `flai stream restore S-nnnn --path <p>`, behind a host write). The card clears, and the preview runs again.
  2. **Keep the removal.** You give one line saying why. flai records it in the acceptance evidence, and accept passes that path (`flai accept --keep-removal <p> --reason …`). This is the override I recommended, now as a button.
  3. **Send back.** This opens a thread on the story naming the paths and moves it back to in-progress, so its agent fixes it.
- **At sync, a badge, not a stop.** The board card and the story page show "undoes main on N paths", linking to the same view, so you see it before review.
- **The orchestrator stays out of it.** It never chooses Keep: an undone path stays a blocker until you resolve it, so it waits for you.

## How it fits the board

- **S-0371** keeps the check, the accept refusal, and the two flai commands. Its tasks swap the `git checkout` advice for `flai stream restore` and `--keep-removal`, exposed as host writes, and its ADR records both.
- **A new story** builds the review-page cards and the badge, after S-0371. It can reuse what S-0378 builds now: refusal findings as data through `accept.run` (T-1516), and findings listed on the review page with Accept enabled again (T-1520). So it waits for S-0378 too.

If you confirm, edit S-0371 or say so here, and its planner will re-plan it to this shape and draft the dashboard story. S-0371 stays a draft meanwhile.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-14T18:14:21Z alex
Confirmed the recommendation of 2026-10-13T17:42:38Z orchestrator.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-14T18:14:43Z orchestrator
Noted: you confirmed the dashboard-resolution shape for S-0371 (restore, keep with a reason, or send back, per path, as acceptance-preview blockers).

**What happens next needs you.** S-0371 already has touches, a forecast, a cost of delay, and tasks, so `flai plan --candidates` does not list it, and under `plan_backlog_stories` I may start the planner only for a story it lists. A thread confirmation is not one of the replan triggers in ADR-0084 either.

**My recommendation:** start the planner on S-0371 yourself, with Plan on the dashboard or `flai plan S-0371`. Point it at this thread, and it re-plans S-0371 to this shape and drafts the review-page story after S-0371 and S-0378. Editing S-0371's goal or criteria would also trigger a replan, with the `plan` host action on.

S-0371 stays a draft until then. I settle the replan's thread when it comes.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-14T21:02:55Z planner-S-0371
Replanned S-0371 to the shape you confirmed here: the operator restores the path, keeps the removal with a reason, or sends the story back, with each undone path a blocker in the acceptance preview. The plan is on TH-0448: eight tasks in six layers, and the new tasks T-1574 (`flai stream restore`) and T-1575 (the host writes). The review-page cards and the board badge are drafted as S-0388, after S-0371. TH-0448 also asks whether to add three acceptance criteria to S-0371, which you finalized while I planned.

### 2026-10-14T22:18:09Z orchestrator
Resolved: The operator confirmed the dashboard-resolution shape (2026-10-14T18:14Z); the planner replanned S-0371 to it and drafted S-0388 (finalized). The remaining question, S-0371's three added criteria, is on TH-0448 with my recommendation awaiting the operator.
