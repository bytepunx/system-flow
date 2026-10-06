---
id: T-0907
type: task
nature: feature
title: flai accept and hostapi accept.run take the orchestrator's identity under accept_reviews, refuse on any blocker, and record the evidence
status: done
parent: S-0221
owner: alex
created: 2026-10-05T04:47:55Z
updated: 2026-10-06T11:34:17Z
transitions:
  - to: ready
    at: 2026-10-06T11:24:27Z
    by: agent-S-0221
  - to: in-progress
    at: 2026-10-06T11:24:28Z
    by: agent-S-0221
  - to: done
    at: 2026-10-06T11:34:17Z
    by: agent-S-0221
stream: S-0221
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/hostapi, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-0900, T-0903]
usage:
  source: log
  seconds: 589
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 124
      output: 707
      cache_read: 7055921
      cache_write: 152316
      cost: 3.1011
---
# T-0907 flai accept and hostapi accept.run take the orchestrator's identity under accept_reviews, refuse on any blocker, and record the evidence

## Work

In `flai/cmd/accept.go`, an acceptance with `--by orchestrator`:

- is refused unless the manifest's `orchestration.permissions.accept_reviews` is on. The refusal names the permission. The check sits in flai as well as the guard, so a shell outside the guard is held to it too;
- takes the verified commit and the evidence as the ADR from T-0898 names them, such as `--verified <commit>` and `--evidence <file>`;
- runs `preview.Accept` with T-0900's orchestrator option before anything is merged, and refuses on any blocker, naming each one;
- on success, records the done transition `by: orchestrator`. It writes the evidence under an `### Accepted by the orchestrator` heading in the story's Notes, in the acceptance commit, and returns it in `--json`.

`flai move <S-nnnn> done --by orchestrator` runs the same flow, since it shares `acceptOptions`.

In `flai/internal/hostapi/writes.go`, `accept.run` takes an optional `by`. Given `orchestrator`, it is refused unless `accept_reviews` is on. Given nothing, it stays the operator's `owner(p)`.

Write `flai/cmd/accept_orchestrator_test.go` with:

- An acceptance: a story in review with every criterion ticked, the diff within its touches, no open thread, and the head commit verified. It is accepted `by: orchestrator`, and its Notes carry the evidence.
- A refusal on an open thread: the same story with a thread open is refused, the blocker names the thread, and nothing is merged.
- The permission off: refused, naming `accept_reviews`, and nothing is merged.
- The `--dry-run` blockers from T-0900: an unticked criterion, a file outside the touches, and a verified commit that is not the head.
- `accept.run` with `by: orchestrator` and the permission off is refused.

This task waits for T-0900, whose preview option it calls, and for T-0903, whose guard rule it must agree with.

## Done when

- `flai accept --by orchestrator` accepts only with the permission on and no blocker, and records the transition and the evidence.
- `accept.run` honours `by: orchestrator` only with the permission on.
- `accept_orchestrator_test.go` covers the cases above, and it and `scripts/flai-test.sh` pass.

## Notes
