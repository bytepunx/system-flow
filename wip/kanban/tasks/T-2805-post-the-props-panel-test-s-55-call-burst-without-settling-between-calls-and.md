---
id: T-2805
type: task
nature: remediation
title: Post the props-panel test's 55-call burst without settling between calls, and measure the test before and after
status: done
parent: S-0640
owner: alex
created: 2026-11-09T22:07:53Z
updated: 2026-11-10T05:16:08Z
transitions:
  - to: ready
    at: 2026-11-10T05:14:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:14:02Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:16:08Z
    by: agent-S-0640
stream: S-0640
tags: [flaiover, testing, flaky-test]
touches: [flaiover/src/routes/design/components/components.svelte.test.ts]
usage:
  source: log
  seconds: 126
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 140
      cache_read: 1689379
      cache_write: 4776
      cost: 0.6223
---
# T-2805 Post the props-panel test's 55-call burst without settling between calls, and measure the test before and after

## Work

In `flaiover/src/routes/design/components/components.svelte.test.ts`, the test "lists the calls of the component's function props, and clears them on another pick" calls `harness()` 59 times, and each call runs `settle()`: 20 `setTimeout(0)` turns and a flush. Under host load those turns pass 5,000 ms (I-0319).

1. Time the test as it is: run vitest on this file alone and record the test's duration from the reporter.
2. Add a helper beside `harness()` that dispatches several messages and settles once after the last. Use it for the loop of 55 `onpick` calls, so the burst costs one settle in place of 55. Keep every assertion as it is: 50 entries, the newest `[54]` first, and none after the click on `BUTTON`.
3. Read the rest of the file for any other test that settles once per message in a loop, and treat it the same way.
4. Time the test again the same way.
5. Decide whether a test can reproduce the timeout, such as one that runs the burst with every timer turn slowed. Write it if it fits; otherwise say why in the story's `## Notes`.

It waits for no task: it is the first layer.

## Done when

- The test passes, with every assertion it had.
- The task's `## Notes` give the test's duration before and after, from the same vitest command.
- `flai test` on the file passes.
- The story's `## Notes` say whether a reproduction test was written, and why not if it was not.

## Notes
