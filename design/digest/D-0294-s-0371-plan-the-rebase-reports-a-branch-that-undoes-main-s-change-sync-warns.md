---
id: D-0294
title: "S-0371 plan: the rebase reports a branch that undoes main's change, sync warns and accept refuses (new ADR), six tasks in five layers: Recommendation: yes"
kind: answer
sources: [TH-0420]
topics: [cli, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_restore.go, flai/cmd/stream_restore_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/preview/accept.go, flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, docs/users/flai-reference.md, docs/operators/settings.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/summary.md, design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md, flai/cmd/move.go, flai/internal/preview/orchestrator.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md]
date: 2026-10-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0294 S-0371 plan: the rebase reports a branch that undoes main's change, sync warns and accept refuses (new ADR), six tasks in five layers: Recommendation: yes

The operator, alex, answered TH-0420 at 2026-10-14T18:14:21Z.

What was asked:

```text
**Recommendation: approve the plan as drafted.** After a clean rebase, `storygit.Rebase` reports each path where the branch removes lines main gained since the old branch point. `flai stream sync` and `flai task done` warn on each one. `flai accept` refuses the branch.
```

The answer:

```text
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
```

Source: TH-0420.
