---
id: T-1077
type: task
nature: improvement
title: flai issue bump, new, and close and flai adr new take --commit
status: in-progress
parent: S-0275
owner: alex
created: 2026-10-06T22:52:34Z
updated: 2026-10-07T08:21:01Z
transitions:
  - to: ready
    at: 2026-10-07T08:21:01Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:21:01Z
    by: agent-S-0275
stream: S-0275
tags: [cli]
touches: [flai/cmd/issue.go, flai/cmd/issue_test.go, flai/cmd/adr.go, flai/cmd/adr_test.go, flai/internal/adr/adr.go]
after: [T-1072, T-1074]
---
# T-1077 flai issue bump, new, and close and flai adr new take --commit

## Work

Criterion 1. It waits for T-1072, the commit-and-widen helper it calls, and for T-1074, the paths the issues package now returns.

- Add `--commit` to `flai issue new`, `bump`, and `close` and to `flai adr new`. With it, run in a story's worktree, the command writes as it does today, commits exactly the files it wrote on the story branch with the story's prefix, and widens the story's touches to them, through T-1072's helper. The story comes from `recordingStory()` (`FLAI_STORY`, `FLAI_AGENT`, or the `story/S-nnnn` branch).
- Refuse `--commit` together with `--autocommit`. Refuse `--commit` outside a story worktree, with a message naming `--autocommit` for the main checkout.
- In text, print the commit and the touches added. With `--json`, add `commit` (hash and paths) and `touches_added` to the existing output.
- `flai adr new` numbers as S-0245 leaves it; this task adds only the single-call shape. Touch `flai/internal/adr/adr.go` only if `Result` needs to expose more for the commit.

## Done when

- Tests in `issue_test.go` and `adr_test.go` run each command with `--commit` in a temporary story worktree. They assert one commit on the story branch holding only the files written, the prefix, the widened touches, and the `--json` fields.
- Tests cover the refusals outside a worktree and with `--autocommit`.
- `scripts/flai-test.sh` passes.

## Notes

Written by the planner.
