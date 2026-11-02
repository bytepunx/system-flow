---
id: T-1781
type: task
nature: improvement
title: flai release --pending refuses an agent flai serve started while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message
status: done
parent: S-0423
owner: alex
created: 2026-10-16T05:48:33Z
updated: 2026-11-02T04:45:35Z
transitions:
  - to: ready
    at: 2026-11-02T04:27:37Z
    by: agent-S-0423
  - to: in-progress
    at: 2026-11-02T04:27:38Z
    by: agent-S-0423
  - to: done
    at: 2026-11-02T04:45:35Z
    by: agent-S-0423
stream: S-0423
tags: [cli]
touches: [flai/cmd/release.go, flai/cmd/release_push_action_test.go, docs/users/flai-reference.md]
after: [T-1780]
usage:
  source: log
  seconds: 1077
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 9936
      cache_read: 2608968
      cache_write: 58360
      cost: 1.0792
---
# T-1781 flai release --pending refuses an agent flai serve started while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message

## Work

In `flai release --pending` (`publishPending` in `flai/cmd/release.go`), before anything is applied or tagged:

- Refuse an agent `flai serve` started (`FLAI_STARTED_BY=flai-serve`, or `FLAI_STORY` set) while the project's `push` host action is off, read from the host's config as `cmd/verify.go`'s `hostActionOn` reads it. Use `release_publish`'s words, naming `flai serve enable push`, with a non-zero exit, so that the formula step fails down its failure path. The operator's own shell is not held, and `--dry-run` is not held.
- The `conflict:` refusal for a moved remote stays as it is, exit 3: the instance's agent writes it on the task and in a thread as any failed step's output (S-0422). Check that its message reaches stderr unchanged.

Waits for T-1780, since both change `flai/cmd/release.go`.

## Done when

`flai test` on `flai/cmd` passes, with `release_push_action_test.go` covering: refused for a served agent with `push` off, allowed with it on, allowed for the operator's shell, and `--dry-run` allowed.

## Notes

Replanned 2026-10-21 by planner-S-0423 on TH-0471's answer 5: there is no serve-side step runner. Rewritten 2026-11-02 by agent-S-0423 to TH-0583's answer 2: the gate is in `flai release --pending`. No journal entry: a CLI run writes none, and the release commit on main records it.
