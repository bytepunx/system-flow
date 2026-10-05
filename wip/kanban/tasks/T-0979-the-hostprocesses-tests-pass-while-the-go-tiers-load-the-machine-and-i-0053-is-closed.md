---
id: T-0979
type: task
nature: remediation
title: The HostProcesses tests pass while the Go tiers load the machine, and I-0053 is closed
status: backlog
parent: S-0264
owner: alex
created: 2026-10-05T05:49:44Z
updated: 2026-10-05T05:49:44Z
transitions: []
stream: S-0264
tags: [flaiover, testing, issues]
touches: [design/issues/I-0053-flaiover-s-hostprocesses-test-that-waits-past-the-old-host-after-an-upgrade-fails-under-machine-load.md, design/issues/summary.md]
after: [T-0977]
---
# T-0979 The HostProcesses tests pass while the Go tiers load the machine, and I-0053 is closed

## Work

Show the cause no longer occurs under the load I-0053's instance names, then close the issue.

- Run the flaiover unit tier (`scripts/flaiover-unit.sh`) while the Go tiers (`scripts/flai-test.sh`) run at the same time, as the close-out runs them, at least twice, and record in the narrative's log that `HostProcesses.svelte.test.ts` passed each time.
- Close the issue in the story's worktree: `flai issue close I-0053 --reason "..."`, the reason saying the reconnect tests asserted after a fixed 40 ms wall-clock wait against a 50 ms give-up and now wait for the outcome with a margin of seconds, with the reproduction test S-0264 added.
- Check the story's acceptance criteria that this verifies.

It waits for T-0977: the run under load and the issue's reason are about the fix it makes.

## Done when

- [ ] The unit tier passed at least twice alongside the Go tiers, logged in the narrative.
- [ ] I-0053 is closed with a reason naming what fixed it, and `design/issues/summary.md` shows it closed.
- [ ] `flai check --strict` is clean for the story.

## Notes

Drafted by the planner for S-0264. Run `flai issue close` from the story's worktree, not the main checkout: issue commands write under the checkout they run in.
