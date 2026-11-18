---
id: T-3479
type: task
nature: feature
title: The workflow package says what a state runs on entry and which transition sends a story back from it
status: done
parent: S-0810
owner: alex
created: 2026-11-16T23:51:01Z
updated: 2026-11-18T09:40:54Z
transitions:
  - to: ready
    at: 2026-11-18T09:21:13Z
    by: agent-S-0810
  - to: in-progress
    at: 2026-11-18T09:21:13Z
    by: agent-S-0810
  - to: done
    at: 2026-11-18T09:40:54Z
    by: agent-S-0810
stream: S-0810
tags: [flai]
touches: [flai/internal/workflow/runs.go, flai/internal/workflow/runs_test.go, flai/internal/workflow/review.go, flai/internal/workflow/review_test.go]
usage:
  source: log
  seconds: 1181
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 6283
      cache_read: 1233332
      cache_write: 42886
      cost: 0.6095
---
# T-3479 The workflow package says what a state runs on entry and which transition sends a story back from it

## Work

Waits for no task. It reads the `runs:` key S-0806 parses onto a state, so it starts once S-0806 is on main, as the story's `after` holds.

1. In a new `flai/internal/workflow/runs.go`, accessors over the parsed `runs:` that `flai serve` reads in place of a state's name or kind: whether a state runs anything on entry; the suites it runs, `all` kept apart from a list; whether it names the `kluster:` steps `deploy` and `checks`, each on its own; and the states that name a `kluster:` step.
2. In `flai/internal/workflow/review.go`, the send-back transition out of a given state: the transition with the `send-back` effect whose `from` is that state, and the state it takes a story to. `SendBack` today finds one for the acceptance's state only; keep it, built on the new one.
3. Tests on the baseline (S-0806 gives `review` `runs: {kluster: [deploy, checks]}`), on a file with an `integration` state between `review` and `done` that runs suites and `deploy` only, and on a state with no `runs:`, which runs nothing.

## Done when

- `flai test` passes on the changed paths.
- On the baseline, `review` is the one state that runs the `kluster:` steps, and its send-back transition takes a story to `in-progress`.

## Notes
