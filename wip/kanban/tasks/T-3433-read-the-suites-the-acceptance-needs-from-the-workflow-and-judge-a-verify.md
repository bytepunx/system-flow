---
id: T-3433
type: task
nature: feature
title: Read the suites the acceptance needs from the workflow, and judge a verify report against them
status: done
parent: S-0808
owner: alex
created: 2026-11-16T23:46:24Z
updated: 2026-11-27T23:03:09Z
transitions:
  - to: ready
    at: 2026-11-27T22:52:14Z
    by: agent-S-0808
  - to: in-progress
    at: 2026-11-27T22:52:14Z
    by: agent-S-0808
  - to: done
    at: 2026-11-27T23:03:09Z
    by: agent-S-0808
stream: S-0808
tags: [flai]
touches: [flai/internal/workitem/transition.go, flai/internal/workitem/transition_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go]
usage:
  source: log
  seconds: 655
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 60
      output: 324
      cache_read: 2460531
      cache_write: 98268
      cost: 0.9835
---
# T-3433 Read the suites the acceptance needs from the workflow, and judge a verify report against them

## Work

ADR-0312 §3. The three gates must read one answer to "which suites does the acceptance need", so it is written once, here, for the tasks after this one.

- In `flai/internal/workitem/transition.go`, add `AcceptSuites(w, typ)`: the suites the `verify:` of the transition with the `accept` effect lists for an item of type `typ`, as S-0806 parses that key, with `all` kept as all. A workflow with no such transition, or one whose transition gives no `verify:`, needs none.
- In `flai/internal/verify/story.go`, keep only the tiers of those suites, reading each tier's suite as S-0805 gives it, and reusing S-0807's helper for the tiers of a suite list where it has one:
  - `StaleTiers` sends back only tiers of the acceptance's suites;
  - `QueueTiersNotRun` names only `queue` tiers of those suites;
  - a new function names each suite the acceptance needs whose selected tiers the last report did not pass, reuse, or cover at the branch's head. It gives the words the refusal uses: the suite, and the close-out that runs it, `flai verify S-nnnn --job` or the MCP tool `verify`.
- With `all`, every function answers exactly what it answers today.

No `after`: this task reads only what S-0805, S-0806, and S-0807 put on main.

## Done when

- `AcceptSuites` has tests for the baseline (`all`), a list of suites, and a workflow whose acceptance gives no `verify:`.
- The verify package's tests cover `StaleTiers` and `QueueTiersNotRun` with a list that leaves a tier's suite out, and the new function naming a suite the report did not run.
- The existing tests in `flai/internal/verify/story_test.go` pass unchanged.
- `scripts/flai-test.sh` passes for `flai/internal/workitem` and `flai/internal/verify`.

## Notes
