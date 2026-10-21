---
id: T-1919
type: task
nature: remediation
title: Close I-0167 with the reason naming the fix
status: done
parent: S-0450
owner: alex
created: 2026-10-20T19:05:52Z
updated: 2026-10-20T19:54:50Z
transitions:
  - to: ready
    at: 2026-10-20T19:54:34Z
    by: agent-S-0450
  - to: in-progress
    at: 2026-10-20T19:54:34Z
    by: agent-S-0450
  - to: done
    at: 2026-10-20T19:54:50Z
    by: agent-S-0450
stream: S-0450
tags: [flai]
touches: [design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, design/issues/summary.md]
after: [T-1918]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 930
      cache_read: 175559
      cache_write: 8812
      cost: 0.1073
---
# T-1919 Close I-0167 with the reason naming the fix

## Work

In the story's worktree, run `flai issue close I-0167 --reason "<reason>"`, the reason naming T-1914's ADR and saying that `story_start` and `task_done` now carry only the story's threads with a count of the rest, fitted to 40,000 bytes, with the tests that reproduce S-0398's and S-0399's instances. `flai issue close` rewrites `design/issues/summary.md` with the issue.

It waits for T-1918, the last of the work, so that the reason states what was built and documented.

## Done when

- I-0167 is closed, with that reason, and `design/issues/summary.md` no longer lists it as open.
- The story's second criterion can be ticked.
- `flai check --strict` passes.

## Notes
