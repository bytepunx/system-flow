---
id: T-1997
type: task
nature: remediation
title: Verify a worktree's integration tier runs on its own flaiover install, and close I-0178
status: backlog
parent: S-0482
owner: alex
created: 2026-10-21T12:08:57Z
updated: 2026-10-21T12:10:10Z
transitions: []
stream: S-0482
tags: [issues, tests]
touches: [design/issues/I-0178-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md, design/issues/summary.md]
after: [T-1994, T-1996]
---
# T-1997 Verify a worktree's integration tier runs on its own flaiover install, and close I-0178

## Work

It waits for both other tasks: it checks the fix they make together, and closes the issue only once both are in.

- In the story's worktree, with no `flaiover/node_modules` of its own, run `scripts/integration.sh` and confirm it installs from the worktree's lockfile and that `TestEveryListedFlaioverComponentCompilesInTheHarness` takes the worktree's install, not the main checkout's. Do not change the main checkout's `node_modules` to stage this; the reproducing test of the helper task covers the stale case.
- Close the issue: `flai issue close I-0178 --reason "<what fixed it>"`, naming the script change, the lockfile check in `flaioverModules`, and this story. Run it in the story's worktree, as issue commands write under the checkout they run in.
- Tick the story's criteria you verified with `flai criteria tick`.

## Done when

- The worktree's integration run used its own install and passed.
- I-0178 is closed with its reason, and `design/issues/summary.md` no longer lists it.
- Both acceptance criteria of S-0482 are ticked.

## Notes
