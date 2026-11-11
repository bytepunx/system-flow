---
id: T-3038
type: task
nature: remediation
title: Close I-0268, naming the send-back, the kluster blocker, and the tests that hold them
status: done
parent: S-0694
owner: alex
created: 2026-11-10T10:33:56Z
updated: 2026-11-11T14:31:37Z
transitions:
  - to: ready
    at: 2026-11-11T14:31:24Z
    by: agent-S-0694
  - to: in-progress
    at: 2026-11-11T14:31:24Z
    by: agent-S-0694
  - to: done
    at: 2026-11-11T14:31:37Z
    by: agent-S-0694
stream: S-0694
tags: [issues, kluster]
touches: [design/issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md, design/issues/summary.md]
after: [T-3035]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 40
      cache_read: 393934
      cache_write: 1350
      cost: 0.1463
---
# T-3038 Close I-0268, naming the send-back, the kluster blocker, and the tests that hold them

## Work

Run `flai issue close I-0268 --reason "<what fixed it>"` from the story's worktree. The reason names three things:

- the send-back in `flai/internal/serve/kluster.go` (T-3025);
- the `kluster` blocker and `--without-kluster` in `flai accept`, with the queue's reason (T-3028, T-3029);
- the tests that reproduce the cause.

flai regenerates `design/issues/summary.md`. This task waits for T-3035, so that the issue closes only once the build and its docs are done.

## Done when

- I-0268's status is closed, with a reason naming the fix and its tests.
- `design/issues/summary.md` no longer lists I-0268.
- `flai check --strict` passes for the story.

## Notes
