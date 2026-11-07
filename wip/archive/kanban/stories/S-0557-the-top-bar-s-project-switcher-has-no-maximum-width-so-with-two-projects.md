---
id: S-0557
type: story
nature: remediation
title: The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport
status: done
owner: alex
created: 2026-11-02T00:04:25Z
updated: 2026-11-07T04:16:39Z
transitions:
  - to: ready
    at: 2026-11-07T03:48:43Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T03:48:50Z
    by: agent-S-0557
  - to: review
    at: 2026-11-07T04:14:37Z
    by: agent-S-0557
  - to: done
    at: 2026-11-07T04:16:39Z
    by: orchestrator
tags: [flaiover]
touches: [flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts, flaiover/e2e/project-switcher.e2e.ts, design/issues/I-0180-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects-connected-the-page-is-wider-than-a-390px-viewport.md, design/issues/summary.md, design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1554
  priced_by: harness
  turns:
    - day: 2026-11-07
      hand_edits: 2
      work: 58
  models:
    - model: claude-opus-5-5
      input: 142
      output: 20729
      cache_read: 14722353
      cache_write: 252919
      cost: 5.351
  strategic:
    - kind: orchestrator
      seconds: 61
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 76
          cache_read: 627983
          cache_write: 52355
          cost: 0.1651
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-11-02T00:04:25Z
  value: 15
  by: planner-S-0557
  at: 2026-11-07T03:45:19Z
forecast:
  duration: 40m
  delivery: 2026-11-07T09:24:00Z
  basis: flai forecast's 14m (131 s per unit over 14 done small remediation stories, size 6) raised to 40m for the flaiover build and Playwright run the reproducing test needs, which none of those stories ran
  by: planner-S-0557
  at: 2026-11-07T03:45:19Z
finalized:
  by: orchestrator
  at: 2026-11-07T03:48:06Z
---
# S-0557 The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport

## Goal

This story remediates [I-0180](../../../design/issues/I-0180-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects-connected-the-page-is-wider-than-a-390px-viewport.md), "The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0180 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0180 is closed with `flai issue close I-0180 --reason` saying what fixed it

## Tasks
- T-2507 Cap the project switcher's width and let it shrink inside the top bar
- T-2509 Reproduce I-0180 in Playwright: two long-named projects at 390px leave the page no wider than the window
- T-2510 Close I-0180 naming the switcher's width cap and the Playwright test that reproduces it

## Notes

Cost of delay inputs set by flai from I-0180. time_lost_per_cycle 6m: 10m per occurrence × 1 occurrence ÷ 1.7 cycles of 168h (first reported 2026-10-21T09:26:50Z, 11.6 days before this story).

### Planning

Proposed solution, from I-0180's one instance: the select in `flaiover/src/lib/components/ProjectSwitcher.svelte` takes the width of its longest option, and neither it nor its wrapping `span` and `label` may shrink inside the top bar's wrapping flex row. Cap the select's width (a Tailwind max width that leaves room at 390px, wider from `sm`), let the wrapper shrink (`min-w-0`, `max-w-full`) and wrap the "n not connected: why?" link under it, and keep each project's full label in its option and in the select's `title`, so nothing the operator reads is lost. The reproducing test answers `/api/projects` in the browser with two projects with long names, as `flaiover/e2e/board-layout.e2e.ts` answers `/api/board`, since a fake flai would make a project known to the server for every test file after it.

Touches, each a file:

| Touch | From | Why |
|-------|------|-----|
| `flaiover/src/lib/components/ProjectSwitcher.svelte` | design (I-0180's instance) | the select that grows to its longest option |
| `flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts` | layout | the component's unit tests, where the select's `title` carrying the current label is checked |
| `flaiover/e2e/project-switcher.e2e.ts` | layout | a new Playwright file beside `layout.e2e.ts`, which renders with no project and so cannot intercept `/api/projects` for its own pages |
| `design/issues/I-0180-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects-connected-the-page-is-wider-than-a-390px-viewport.md` | design (criterion 2) | `flai issue close` writes it |
| `design/issues/summary.md` | co-change | `flai issue close` regenerates it |

`flai touches suggest` named nothing else that bears on the story: its co-change list is other issues and documents that change with `design/issues/summary.md`. No `docs/users/flaiover.md` change: the switcher does what it did, only narrower. No folder touch and no `after`: the story reads nothing another open story writes.

Forecast: `flai forecast` gave 14m (131 s per unit of size over 14 done small remediation stories, size 6). Raised to 40m, delivery moved by the same 26m, because the reproducing test needs `scripts/flaiover-e2e.sh`, which builds flaiover and runs Playwright's Chromium, and none of those stories ran it.

Cost of delay: `flai cod` gave 15.00 USD a week from the operator's input, 6m lost per 168h cycle at 150 USD an hour. It stands: no thread or issue instance gives a reason to change it.

### Accepted by the orchestrator

- Verified: af568297e84208f304052f212a53eaa4bc5fe935
- At: 2026-11-07T04:16:39Z

Verdict: pass
- 1: flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/e2e/project-switcher.e2e.ts, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts
- 2: design/issues/I-0180-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects-connected-the-page-is-wider-than-a-390px-viewport.md, design/issues/summary.md
