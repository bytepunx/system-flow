---
id: T-3435
type: task
nature: feature
title: The merge queue's eligibility needs only the acceptance's suites to have passed at the story's head
status: done
parent: S-0808
owner: alex
created: 2026-11-16T23:46:42Z
updated: 2026-11-27T23:23:44Z
transitions:
  - to: ready
    at: 2026-11-27T23:09:14Z
    by: agent-S-0808
  - to: in-progress
    at: 2026-11-27T23:09:14Z
    by: agent-S-0808
  - to: done
    at: 2026-11-27T23:23:44Z
    by: agent-S-0808
stream: S-0808
tags: [flai]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go]
after: [T-3433]
usage:
  source: log
  seconds: 870
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 60
      output: 19061
      cache_read: 4285686
      cache_write: 167803
      cost: 2.2341
---
# T-3435 The merge queue's eligibility needs only the acceptance's suites to have passed at the story's head

## Work

ADR-0312 §3, criteria 2 and 4. It waits for T-3433, whose `AcceptSuites` and suite-aware verify functions it calls, so that the queue and the acceptance's gates read the same suites.

- In `flai/internal/mergequeue/eligible.go`, `gitState.verified` judges a story's last verify by the suites `AcceptSuites` gives for `repo.Flow()`. A story is eligible when every tier of those suites that its paths select passed, was reused, or was covered at its branch's head. A failure in a tier of another suite does not make it ineligible.
- A story whose last verify did not run a suite the acceptance needs is ineligible. Its reason names the suite and the close-out that runs it, with the same words as the acceptance's refusal.
- On the baseline (`verify: all`) every reason and every eligible story is what it is today.

It shares no file with the other task of its layer, so the two run together.

## Done when

- A test in `flai/internal/mergequeue/eligible_test.go` shows a story eligible though a tier of a suite the acceptance does not list failed.
- A test shows a story ineligible, with words naming the suite, when its verify did not run a suite the acceptance needs.
- The existing tests of `flai/internal/mergequeue` pass unchanged.

## Notes
