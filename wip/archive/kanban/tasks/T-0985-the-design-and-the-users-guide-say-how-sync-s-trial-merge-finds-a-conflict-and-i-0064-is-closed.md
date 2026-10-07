---
id: T-0985
type: task
nature: remediation
title: The design and the users' guide say how sync's trial merge finds a conflict, and I-0064 is closed
status: done
parent: S-0251
owner: alex
created: 2026-10-05T05:50:53Z
updated: 2026-10-07T00:04:24Z
transitions:
  - to: ready
    at: 2026-10-07T00:04:00Z
    by: agent-S-0251
  - to: in-progress
    at: 2026-10-07T00:04:00Z
    by: agent-S-0251
  - to: done
    at: 2026-10-07T00:04:24Z
    by: agent-S-0251
stream: S-0251
tags: [flai, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0064-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md, design/issues/summary.md]
after: [T-0984]
usage:
  source: log
  seconds: 24
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 13
      output: 3796
      cache_read: 481979
      cache_write: 22038
      cost: 0.3237
---
# T-0985 The design and the users' guide say how sync's trial merge finds a conflict, and I-0064 is closed

## Work

Describe the fix T-0984 built in two places, in the words of its `## Decisions` entry:

- the `flai stream sync` row of `design/system/flai-cli.md`, where `checkSync` trial-merges with `git merge-tree --write-tree --name-only`;
- the **Conflicts** bullet of `flai stream sync` in `docs/users/flai.md`.

Each should say that a conflict names only the paths both stories changed since they left main. A path where only main and the other story's stale branch disagree is left to that story's own sync.

Then close the issue in the story's worktree with `flai issue close I-0064 --reason`. The reason names the fix and the tests that reproduce both instances. The command updates `design/issues/summary.md`.

It waits for T-0984 because it describes what that task built, and the issue closes only once that task's reproduction tests pass.

## Done when

- `design/system/flai-cli.md` and `docs/users/flai.md` say how the trial merge decides that a path conflicts, and match the help text of `flai stream sync`.
- I-0064 is closed, its reason names the fix, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` and the markdown lint pass.

## Notes
