---
id: T-1093
type: task
nature: improvement
title: The host channel files, bumps, and closes an issue against a story
status: done
parent: S-0275
owner: alex
created: 2026-10-06T22:53:00Z
updated: 2026-10-07T08:43:25Z
transitions:
  - to: ready
    at: 2026-10-07T08:34:27Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:34:28Z
    by: agent-S-0275
  - to: done
    at: 2026-10-07T08:43:25Z
    by: agent-S-0275
stream: S-0275
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/issue.go, flai/cmd/issue_test.go, flaiover/src/lib/server/agent.ts]
after: [T-1077]
usage:
  source: log
  seconds: 537
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 92
      output: 37916
      cache_read: 5843715
      cache_write: 159583
      cost: 2.8525
---
# T-1093 The host channel files, bumps, and closes an issue against a story

## Work

Criterion 2, on the host channel, so the dashboard can file an issue against a story. It waits for T-1077: the host channel runs the CLI, and this task changes `flai/cmd/issue.go` after T-1077 does.

- Add `issue.new`, `issue.bump`, and `issue.close` to `flai/internal/hostapi/writes.go`, beside `issue.story` and `adr.new`. Each validates its input as those do (`issueID`, a story ID when one is given) and builds the CLI arguments with `--story` and the dashboard's trailer.
- The dashboard has no worktree, so these run in the main checkout and commit there, as `adr.new` does with `--autocommit --trailer`. Give `flai issue new`, `bump`, and `close` the `--autocommit` and `--trailer` flags `flai adr new` has, if T-1077 did not.
- `adr.new` already exists on the host channel: leave its behaviour as it is, and check that it still builds after T-1077.

## Done when

- Tests in `writes_test.go` cover each new operation's arguments, a refusal on a bad issue or story ID, and the exit code mapped to `Refused`.
- Tests in `issue_test.go` cover `--autocommit` with `--trailer` on new, bump, and close in the main checkout.
- `scripts/flai-test.sh` passes.

## Notes

Written by the planner. Assumption, named on the plan thread: the host channel commits on the main branch, and does not widen any story's touches.
