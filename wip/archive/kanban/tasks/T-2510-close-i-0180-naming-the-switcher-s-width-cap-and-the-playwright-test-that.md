---
id: T-2510
type: task
nature: remediation
title: Close I-0180 naming the switcher's width cap and the Playwright test that reproduces it
status: done
parent: S-0557
owner: alex
created: 2026-11-07T03:46:26Z
updated: 2026-11-07T04:01:19Z
transitions:
  - to: ready
    at: 2026-11-07T04:01:05Z
    by: agent-S-0557
  - to: in-progress
    at: 2026-11-07T04:01:08Z
    by: agent-S-0557
  - to: done
    at: 2026-11-07T04:01:19Z
    by: agent-S-0557
stream: S-0557
tags: [issues]
touches: [design/issues/I-0180-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects-connected-the-page-is-wider-than-a-390px-viewport.md, design/issues/summary.md]
after: [T-2509]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 683
      cache_read: 485109
      cache_write: 8334
      cost: 0.1763
---
# T-2510 Close I-0180 naming the switcher's width cap and the Playwright test that reproduces it

## Work

Run `flai issue close I-0180 --reason "<what fixed it>"` in the story's worktree, the reason naming the width cap T-2507 gave the select in `flaiover/src/lib/components/ProjectSwitcher.svelte` and the test T-2509 added in `flaiover/e2e/project-switcher.e2e.ts`, which reproduces the overflow with two projects connected. Waits for T-2509: the reason cites that test, which must pass first.

## Done when

- I-0180's status is `closed` and its reason names the fix and the test.
- `design/issues/summary.md` no longer lists I-0180.
- `flai check --strict` is clean for the story.

## Notes
