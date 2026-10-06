---
id: T-1071
type: task
nature: improvement
title: flai/internal/verify runs the close-out's steps for a story's worktree and answers one result with each step's state, duration, and findings
status: backlog
parent: S-0270
owner: alex
created: 2026-10-06T22:52:15Z
updated: 2026-10-06T22:52:15Z
transitions: []
stream: S-0270
tags: [flai, cli]
touches: [flai/internal/verify/verify.go, flai/internal/verify/verify_test.go]
---
# T-1071 flai/internal/verify runs the close-out's steps for a story's worktree and answers one result with each step's state, duration, and findings

## Work

Criterion 1's core: one implementation of what `scripts/close-out.sh` checks, which the command, the MCP tool, the host method, and the close-out all call (criterion 3), so that they never disagree.

- Add `Verify` to `flai/internal/verify`, the package S-0273 adds with its tier runner and findings. Given a story and its worktree, it runs these steps in order, cheapest first, and stops at the first that fails:
  1. No rebase is unfinished in the worktree.
  2. The test and lint tiers that the branch's diff selects (committed against the main branch, pending, and untracked, as the close-out's `touched` does). Run them through S-0273's runner, which reads them from the project's tier declaration.
  3. `flai check --strict` scoped to the story through `flai/internal/check`'s scope (ADR-0085). A finding outside the story is a note, not a failure (S-0249).
  4. The narrative's `## Current state` and `## Next steps` are written, in the main checkout's `wip/agents/<story>.md`.
  5. The branch contains the main branch.
- Answer one result: the story, the commit verified, the outcome, and per step its name, its state (`passed`, `failed`, or `not reached`), its duration, and its findings (path, line, message), with the notes apart. Cap the findings per step as S-0273 caps a tier's.
- Write the last result to `.flai-cache/verify/<story>.json`, so that the host channel and the dashboard can show it without running again.
- Do not commit and do not record issues here. Committing stays with the close-out, and recording the notes stays with the command's `--record-issues`.
- This task waits for nothing in this story. It needs S-0273's package and tier declaration, which the story's `after` holds it for.

## Done when

- [ ] `Verify` runs the five steps in order and stops at the first failure, with the later steps `not reached`.
- [ ] Tests cover each step's pass and fail, the tier selection by diff, an outside-the-story finding answered as a note, and the stored last result.
- [ ] `scripts/flai-test.sh` passes for `flai/internal/verify`.

## Notes

Drafted by the planner. Assumes S-0273's tiers carry the paths each applies to, as this repository's close-out selects the flai, template, and flaiover tiers by path.
