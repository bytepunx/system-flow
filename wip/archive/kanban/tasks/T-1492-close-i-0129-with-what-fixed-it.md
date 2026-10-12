---
id: T-1492
type: task
nature: improvement
title: Close I-0129 with what fixed it
status: done
parent: S-0373
owner: alex
created: 2026-10-12T17:29:07Z
updated: 2026-10-12T18:17:41Z
transitions:
  - to: ready
    at: 2026-10-12T18:17:35Z
    by: agent-S-0373
  - to: in-progress
    at: 2026-10-12T18:17:36Z
    by: agent-S-0373
  - to: done
    at: 2026-10-12T18:17:41Z
    by: agent-S-0373
stream: S-0373
tags: [issues]
touches: [design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md, design/issues/summary.md]
after: [T-1491]
usage:
  source: log
  seconds: 5
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 480
      cache_read: 185855
      cache_write: 5392
      cost: 0.0899
---
# T-1492 Close I-0129 with what fixed it

## Work

Close I-0129 with `flai issue close I-0129 --reason`, run in the story's worktree. The reason names S-0373 and says that `WriteIndex` lists a narrative that does not parse as unreadable instead of failing (T-1489), so no story's write depends on another's narrative, and that `flai check` names the file with `narrative.front-matter` alone (T-1490). It names the test that reproduces the instance.

It waits for T-1491, the last of the work, so that the reason describes what was built and documented.

## Done when

- I-0129's status is closed, and its Remediation section carries the reason.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes

flai writes both files. Neither is edited by hand.
