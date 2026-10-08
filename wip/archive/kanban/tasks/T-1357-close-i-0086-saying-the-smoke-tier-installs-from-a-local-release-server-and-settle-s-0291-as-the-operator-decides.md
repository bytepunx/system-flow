---
id: T-1357
type: task
nature: remediation
title: Close I-0086 saying the smoke tier installs from a local release server, and settle S-0291 as the operator decides
status: done
parent: S-0340
owner: alex
created: 2026-10-08T08:06:23Z
updated: 2026-10-08T08:53:14Z
transitions:
  - to: ready
    at: 2026-10-08T08:50:12Z
    by: agent-S-0340
  - to: in-progress
    at: 2026-10-08T08:50:13Z
    by: agent-S-0340
  - to: done
    at: 2026-10-08T08:53:14Z
    by: agent-S-0340
stream: S-0340
tags: [cli]
touches: [design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md]
after: [T-1354, T-1352]
usage:
  source: log
  seconds: 181
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 17921
      cache_read: 2616640
      cache_write: 77976
      cost: 1.378
---
# T-1357 Close I-0086 saying the smoke tier installs from a local release server, and settle S-0291 as the operator decides

## Work

Once T-1354's run with GitHub unreachable has passed, close I-0086 with `flai issue close I-0086 --reason`. The reason says the smoke tier installs and self-upgrades from a release built from the tree and served locally, and that the check against GitHub runs in its own workflow, outside every close-out. Run it in the story's worktree: `flai issue` writes under the checkout it runs in.

Then settle S-0291, which S-0340 takes over. Ask the operator on a thread on S-0340 whether S-0291 is cancelled or closed as taken over, unless the plan's thread has settled it. Do what the answer says, and record it in the narrative's `## Decisions`.

Waits for T-1354 and T-1352: the issue closes when the fix is shown, and its reason names the GitHub check.

## Done when

- I-0086 is closed with a reason naming the local release server and the separate GitHub check, and `design/issues/summary.md` no longer lists it.
- S-0291 is cancelled or closed as the operator answered, and the answer is in the narrative.

## Notes

Drafted by the planner from S-0340's criterion 6.
