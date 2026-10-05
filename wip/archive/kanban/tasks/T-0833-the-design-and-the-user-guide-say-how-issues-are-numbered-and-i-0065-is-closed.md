---
id: T-0833
type: task
nature: remediation
title: The design and the user guide say how issues are numbered, and I-0065 is closed
status: done
parent: S-0252
owner: alex
created: 2026-10-04T23:27:43Z
updated: 2026-10-05T00:04:38Z
transitions:
  - to: ready
    at: 2026-10-05T00:04:12Z
    by: agent-S-0252
  - to: in-progress
    at: 2026-10-05T00:04:13Z
    by: agent-S-0252
  - to: done
    at: 2026-10-05T00:04:38Z
    by: agent-S-0252
stream: S-0252
tags: [docs, issues]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/issues/I-0065-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md, design/issues/summary.md]
after: [T-0832]
usage:
  source: log
  seconds: 25
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 4127
      cache_read: 516419
      cache_write: 23567
      cost: 0.3412
---
# T-0833 The design and the user guide say how issues are numbered, and I-0065 is closed

## Work

Describe the new numbering where `flai issue new` is described:

- `design/system/flai-cli.md`, in the `flai issue` row: the number is one past the highest on main, in any story worktree, and on any story branch, and it comes from the storygit helper.
- `design/system/continuous-improvement.md`, beside the line saying that issues live under the checkout's `design/` and are committed on the story's branch.
- `docs/users/flai.md`, in its issues section.

Then close I-0065 with `flai issue close I-0065 --reason`, naming S-0252 and what fixed it, in the story's worktree, so that `summary.md` is regenerated on the branch.

It waits for T-0832: the documents describe what that task built, and the issue closes only once its reproduction test passes.

## Done when

- The three documents say how the next issue number is found, and the markdown lint passes on them.
- I-0065 is closed with a reason naming S-0252, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes
