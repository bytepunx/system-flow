---
id: T-1071
type: task
nature: improvement
title: flai/internal/verify runs the close-out's steps for a story's worktree and answers one result with each step's state, duration, and findings
status: done
parent: S-0270
owner: alex
created: 2026-10-06T22:52:15Z
updated: 2026-10-07T03:36:07Z
transitions:
  - to: ready
    at: 2026-10-07T03:27:06Z
    by: agent-S-0270
  - to: in-progress
    at: 2026-10-07T03:27:06Z
    by: agent-S-0270
  - to: done
    at: 2026-10-07T03:36:07Z
    by: agent-S-0270
stream: S-0270
tags: [flai, cli]
touches: [flai/internal/verify/verify.go, flai/internal/verify/verify_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/proc.go, flai/internal/verify/proc_unix.go, flai/internal/verify/proc_windows.go, flai/internal/verify/proc_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go]
usage:
  source: log
  seconds: 541
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 25993
      cache_read: 4499646
      cache_write: 130600
      cost: 2.235
---
# T-1071 flai/internal/verify runs the close-out's steps for a story's worktree and answers one result with each step's state, duration, and findings

## Work

Criterion 1's core: one implementation of what `scripts/close-out.sh` checks, which the command, the MCP tool, the host method, and the close-out all call (criterion 3), so that they never disagree.

- Add `Verify` to `flai/internal/verify`, the package S-0273 adds with its tier runner and findings. Given a story and its worktree, it runs these steps in order, cheapest first, and stops at the first that fails:
  1. No rebase is unfinished in the worktree.
  2. The branch contains the main branch.
  3. The narrative's `## Current state` and `## Next steps` are written, in the main checkout's `wip/agents/<story>.md`.
  4. `flai check --strict` scoped to the story through `flai/internal/check`'s scope (ADR-0085). A finding outside the story is a note, not a failure (S-0249).
  5. The test and lint tiers that the branch's diff selects (committed against the main branch, pending, and untracked, as the close-out's `touched` does), each tier a step of its own. Run them through S-0273's runner, which reads them from the project's tier declaration. The selection is `flai test`'s, plus every `all_only` tier whose `paths` the diff selects, plus every `all_only` tier with no `paths`. The tiers run with `CLOSE_OUT_STORY` set to the story in their environment.
- The cheap steps come first, so that an empty narrative or an unsynced branch stops the run before minutes of tiers (narrative, Decisions).
- Answer one result: the story, the commit verified, the outcome, and per step its name, its state (`passed`, `failed`, or `not reached`), its duration, and its findings (path, line, message), with the notes apart. Cap the findings per step as S-0273 caps a tier's.
- Write the last result to `.flai-cache/verify/<story>.json`, so that the host channel and the dashboard can show it without running again.
- Do not commit and do not record issues here. Committing stays with the close-out, and recording the notes stays with the command's `--record-issues`.
- This task waits for nothing in this story. It needs S-0273's package and tier declaration, which the story's `after` holds it for.

## Done when

- [ ] `Verify` runs the steps in order and stops at the first failure, with the later steps `not reached`.
- [ ] Tests cover each step's pass and fail, the tier selection by diff, an outside-the-story finding answered as a note, and the stored last result.
- [ ] `scripts/flai-test.sh` passes for `flai/internal/verify`.

## Notes

Drafted by the planner. Assumes S-0273's tiers carry the paths each applies to, as this repository's close-out selects the flai, template, and flaiover tiers by path.
