---
id: T-1296
type: task
nature: remediation
title: Say in the operator docs and the design how the upgrade waits for the name and falls back, and close I-0094
status: done
parent: S-0316
owner: alex
created: 2026-10-07T23:55:23Z
updated: 2026-10-08T00:11:16Z
transitions:
  - to: ready
    at: 2026-10-08T00:09:37Z
    by: agent-S-0316
  - to: in-progress
    at: 2026-10-08T00:09:38Z
    by: agent-S-0316
  - to: done
    at: 2026-10-08T00:11:16Z
    by: agent-S-0316
stream: S-0316
tags: [docs, dashboard]
touches: [docs/operators/index.md, docs/operators/runbooks/update.md, design/system/flai-cli.md, design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md, design/issues/summary.md]
after: [T-1295]
usage:
  source: log
  seconds: 98
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 19
      output: 5423
      cache_read: 1173681
      cache_write: 52239
      cost: 0.7212
---
# T-1296 Say in the operator docs and the design how the upgrade waits for the name and falls back, and close I-0094

## Work

Write down the behaviour T-1295 delivered, as it was built. It waits for T-1295 because the wait's limit and what the command reports on a failed start are settled there.

- `docs/operators/index.md`: in the paragraph on the upgrade under the dashboard host action, say that after the running container stops the upgrade waits for Docker to release its name, and that when the new container still fails to start the previous image is started again.
- `docs/operators/runbooks/update.md`: in step 2 of updating the dashboard, say what the operator sees in that case.
- `design/system/flai-cli.md`: say the same for builders where it describes how `flai dashboard upgrade` swaps the container.

Then close the issue in the story's worktree: `flai issue close I-0094 --reason` naming the fix, the bounded wait for the name, and the rollback to the previous image. It writes the issue and `design/issues/summary.md`.

## Done when

- The three documents say what the upgrade does when the name lingers and when the new container fails to start, and match the code
- I-0094 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it as open
- The markdown lint passes on the changed files

## Notes
