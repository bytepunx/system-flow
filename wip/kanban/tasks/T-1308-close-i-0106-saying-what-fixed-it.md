---
id: T-1308
type: task
nature: remediation
title: Close I-0106 saying what fixed it
status: backlog
parent: S-0320
owner: alex
created: 2026-10-08T00:18:01Z
updated: 2026-10-08T00:18:01Z
transitions: []
stream: S-0320
tags: [flai, serve, tests]
touches: [design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md, design/issues/summary.md]
after: [T-1307]
---
# T-1308 Close I-0106 saying what fixed it

## Work

Run `flai issue close I-0106 --reason "<text>"` in the story's worktree, so the close lands on the story branch. The reason names S-0320 and what fixed it: the test's one-second margin before `orchestrateRetry` against an `Ended` written to the second, widened to twenty seconds by commit 7772bd34 and replaced by T-1307 with a clock taken from the run's recorded end.

It waits for T-1307: the reason states what that task changed, and the issue is not closed before the fix is in.

## Done when

- `design/issues/I-0106-…md` has status `closed` and a closing line naming S-0320 and the fix.
- `design/issues/summary.md` no longer lists I-0106.
- `flai check --strict` is clean for the story.

## Notes
