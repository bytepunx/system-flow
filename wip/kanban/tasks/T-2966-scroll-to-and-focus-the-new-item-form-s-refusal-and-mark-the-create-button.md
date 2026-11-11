---
id: T-2966
type: task
nature: remediation
title: Scroll to and focus the new item form's refusal, and mark the Create button refused
status: done
parent: S-0700
owner: alex
created: 2026-11-10T10:28:13Z
updated: 2026-11-11T08:23:16Z
transitions:
  - to: ready
    at: 2026-11-11T08:20:45Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T08:20:45Z
    by: system-flow
  - to: done
    at: 2026-11-11T08:23:16Z
    by: agent-S-0700
stream: S-0700
tags: [dashboard]
touches: [flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts]
usage:
  source: log
  seconds: 151
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 4348
      cache_read: 2258477
      cache_write: 60037
      cost: 1.019
---
# T-2966 Scroll to and focus the new item form's refusal, and mark the Create button refused

## Work

In `flaiover/src/lib/components/NewItemForm.svelte`, after any refusal sets `error` (the 422 with findings, a refused move, the agent config error, a thrown fetch), wait for the DOM with `tick()`, then call `scrollIntoView?.({ block: 'center' })` on the `refusal` alert and focus it; give the alert `tabindex="-1"` so it can take focus. Beside the Create button, while `error` is set, say `Refused: see above` and give the button `aria-describedby` pointing at the alert. Clear the mark when the next submit starts. Keep the text in the form, as now.

In `flaiover/src/lib/components/NewItemForm.svelte.test.ts`, add a test that reproduces I-0263: submit the form from the title field with a backend that answers 422 with an MD034 finding, and assert that the alert's `scrollIntoView` was called, that it holds focus, and that the button says it was refused. Waits for nothing: it is the first layer.

## Done when

- A refusal scrolls the alert into view, focuses it, and marks the Create button, and the next submit clears the mark.
- The new test fails without the change and passes with it.
- `flai test` on the two files passes.

## Notes
