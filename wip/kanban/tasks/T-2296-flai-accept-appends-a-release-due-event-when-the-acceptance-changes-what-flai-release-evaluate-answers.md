---
id: T-2296
type: task
nature: feature
title: flai accept appends a release-due event when the acceptance changes what flai release --evaluate answers
status: done
parent: S-0535
owner: alex
created: 2026-10-28T21:38:22Z
updated: 2026-10-28T23:23:30Z
transitions:
  - to: ready
    at: 2026-10-28T22:45:30Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:45:31Z
    by: system-flow
  - to: done
    at: 2026-10-28T23:23:30Z
    by: agent-S-0535
stream: S-0535
tags: [flai, events, release, accept]
touches: [flai/cmd/accept.go, flai/cmd/accept_release_due_test.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, flai/cmd/accept_test.go, docs/users/flai-reference.md]
after: [T-2292]
usage:
  source: log
  seconds: 2257
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 96
      output: 586
      cache_read: 4919904
      cache_write: 257308
      cost: 1.7636
---
# T-2296 flai accept appends a release-due event when the acceptance changes what flai release --evaluate answers

## Work

`release.EvaluateRepo` (`flai/internal/release/evaluate.go`) answers whether `orchestration.release` is met from what is pending in git. `acceptItem` in `flai/cmd/accept.go` is the one path every acceptance takes: the operator's, the orchestrator's, `flai move S-nnnn done`, and the merge queue's `--batch`. The acceptance commit is what changes the pending batch.

- Give `release.Evaluation` a comparison of the parts that make the answer: whether it is due, and the figures the policy reads. Leave out anything that changes with the clock alone.
- In `acceptItem`, evaluate before the acceptance. Evaluate again after `commitAcceptance`, and after `finishCommit` on the second path. When the two differ, append a `release-due` event through T-2292's log with the story, both answers, and `by` the acceptor.
- An evaluation that fails, such as a project with no release policy, appends nothing and does not fail the acceptance.

It waits for T-2292, whose log it appends to. It shares no path with T-2294 or T-2295, so the three run together.

## Done when

- A test accepts a story that makes the release due, and reads one `release-due` event with both answers.
- A test accepts a story that changes nothing in the answer, and reads none.
- A test shows a failed evaluation leaves the acceptance as it was.
- `flai test` on the changed paths passes.

## Notes

The story's criterion 3 makes `release-due` an event of its own, where ADR-0209 §1 wrote `{event: accepted, when: release-due}`. T-2291's ADR records the choice.
