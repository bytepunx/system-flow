---
id: TH-0235
title: "S-0307 plan: retry a held index lock, resume a failed acceptance commit, close I-0100"
anchor:
  path: wip/kanban/stories/S-0307-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md
  item: S-0307
status: resolved
participants: [planner-S-0307, alex]
created: 2026-10-07T01:13:50Z
updated: 2026-10-07T01:19:31Z
---

# TH-0235 S-0307 plan: retry a held index lock, resume a failed acceptance commit, close I-0100

On wip/kanban/stories/S-0307-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md.

## Entries

### 2026-10-07T01:13:50Z planner-S-0307
The plan for S-0307 has four tasks in three layers. It takes all three of I-0100's directions.

Layers:

1. T-1137: a helper in `flai/internal/storygit/indexlock.go`. It runs a git command again while git says another process holds `.git/index.lock`, for about ten seconds with a growing wait. It never removes the lock.
2. T-1138: in `flai/cmd/accept.go` and `flai/internal/preview/accept.go`.
   - The acceptance's `git add -A` and `git commit` go through T-1137's helper.
   - A story that is done and archived, whose acceptance commit is missing, resumes at the commit when `flai accept` runs again.
   - A commit that still fails exits non-zero with git's whole error and says `flai accept <id>` finishes it.
   - `flai/cmd/accept_lock_test.go` reproduces I-0100 with a real lock.
3. Run together, no shared path:
   - T-1139: `design/system/flai-cli.md`, `docs/users/flai.md`, and the regenerated `docs/users/flai-reference.md`.
   - T-1140: `flai issue close I-0100`.

Figures:

- Forecast 40m: flai's 22m, raised for the real-git lock test and the resumed commit's notices.
- Cost of delay 5 USD a week, as `flai cod` computes it from the 2m input flai took from I-0100.

Assumptions:

- The host's journal already records a failed `flai accept` as failed, because `hostapi` takes the entry's outcome from the run's error. T-1138 confirms this by reading the code. It changes `flai/internal/hostapi/writes.go` only if the journal does not, which would overlap S-0269 to S-0275.
- A resumed acceptance is recognised by git: the item is done and archived, and the main checkout has uncommitted changes under wip. A done, archived item with nothing to commit is still refused as already done. The orchestrator still refuses a resumed item.
- A resumed commit cannot know what the merge changed. It takes the paths for the overlap notices from the story's commits (`storygit.Committed`).
- The retry covers the acceptance only. The other commits flai makes in the main checkout could fail the same way: `docedit`, `release`, `import`, `agent`, and the autocommits. They could share the helper in a later story. I propose none now, with one occurrence on record.

S-0307 is held until S-0286 moves to review. Both change `flai/cmd/accept.go` and `flai/internal/preview/accept.go`, which S-0286's T-1120 needs, and no other placement avoids it.

### 2026-10-07T01:19:31Z alex
Resolved.
