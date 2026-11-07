---
id: T-2507
type: task
nature: remediation
title: Cap the project switcher's width and let it shrink inside the top bar
status: done
parent: S-0557
owner: alex
created: 2026-11-07T03:45:47Z
updated: 2026-11-07T03:52:29Z
transitions:
  - to: ready
    at: 2026-11-07T03:49:36Z
    by: agent-S-0557
  - to: in-progress
    at: 2026-11-07T03:49:37Z
    by: agent-S-0557
  - to: done
    at: 2026-11-07T03:52:29Z
    by: agent-S-0557
stream: S-0557
tags: [flaiover]
touches: [flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts]
usage:
  source: log
  seconds: 172
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 2397
      cache_read: 1702547
      cache_write: 29248
      cost: 0.6188
---
# T-2507 Cap the project switcher's width and let it shrink inside the top bar

## Work

In `flaiover/src/lib/components/ProjectSwitcher.svelte`, give the select a maximum width that leaves the top bar room at 390px and a wider one from Tailwind's `sm`, and let the wrapping `span` and `label` shrink and wrap (`min-w-0`, `max-w-full`, `flex-wrap`), so the "n not connected: why?" link goes under the select rather than past the window. Keep each project's full label in its option, and put the current project's label in the select's `title` beside what the title says now, so a truncated name can still be read. Waits for no task: it is the fix, and the test in the next layer checks it.

## Done when

- The select can no longer grow wider than its cap, whatever its options say.
- A test in `flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts` shows the current project's full label in the select's `title`.
- `flai test flaiover/src/lib/components/ProjectSwitcher.svelte flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts` passes.

## Notes
