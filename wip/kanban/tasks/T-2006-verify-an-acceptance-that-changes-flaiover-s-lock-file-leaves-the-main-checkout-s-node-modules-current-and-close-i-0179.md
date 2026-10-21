---
id: T-2006
type: task
nature: remediation
title: Verify an acceptance that changes flaiover's lock file leaves the main checkout's node_modules current, and close I-0179
status: done
parent: S-0483
owner: alex
created: 2026-10-21T12:11:39Z
updated: 2026-10-21T13:28:14Z
transitions:
  - to: ready
    at: 2026-10-21T13:28:07Z
    by: agent-S-0483
  - to: in-progress
    at: 2026-10-21T13:28:07Z
    by: agent-S-0483
  - to: done
    at: 2026-10-21T13:28:14Z
    by: agent-S-0483
stream: S-0483
tags: [flai]
touches: [design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md, design/issues/summary.md]
after: [T-2003, T-2004]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1233
      cache_read: 221346
      cache_write: 11157
      cost: 0.1365
---
# T-2006 Verify an acceptance that changes flaiover's lock file leaves the main checkout's node_modules current, and close I-0179

## Work

Check that the cause I-0179 records no longer occurs, then close it, the story's second criterion. It waits for T-2003 and T-2004 because it checks the step built and documented there.

- Run `flai test --all` on the story's paths in the worktree, so the integration tier runs `TestEveryListedFlaioverComponentCompilesInTheHarness` and T-2003's tests together.
- Confirm in `accept_install_test.go` that a lock file change accepted from a story branch runs flaiover's `install` in the main checkout, as I-0179's instance needed.
- Close the issue with `flai issue close I-0179 --reason`, naming the ADR, `projects[].install`, and the acceptance step, and that S-0482 covers the story worktree's side.

## Done when

- The integration tier passes in the worktree.
- `flai issue close I-0179` has run, and `design/issues/summary.md` no longer lists I-0179.
- `flai check --strict` passes.

## Notes
