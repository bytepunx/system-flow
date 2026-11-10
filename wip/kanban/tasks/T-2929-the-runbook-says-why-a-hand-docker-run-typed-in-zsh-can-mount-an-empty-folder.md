---
id: T-2929
type: task
nature: improvement
title: The runbook says why a hand docker run typed in zsh can mount an empty folder, and I-0276 is closed
status: done
parent: S-0667
owner: alex
created: 2026-11-10T08:49:05Z
updated: 2026-11-10T11:03:38Z
transitions:
  - to: ready
    at: 2026-11-10T11:03:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T11:03:21Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:03:38Z
    by: agent-S-0667
stream: S-0667
tags: [docker, agent-container]
touches: [docs/operators/runbooks/agent-container.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md, design/issues/summary.md]
after: [T-3078]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1442
      cache_read: 668011
      cache_write: 11426
      cost: 0.2539
---
# T-2929 The runbook says why a hand docker run typed in zsh can mount an empty folder, and I-0276 is closed

## Work

Waits for T-3078: the runbook points to the rule T-3078 writes, and I-0276 is closed on both.

In `docs/operators/runbooks/agent-container.md`, say why a hand `docker run` typed in zsh can mount the main checkout as an empty folder, as T-2926 found it in S-0667's Notes under `### Cause`: `$ROOT:r` in `-v "$ROOT:$ROOT:ro"` is zsh's modifier. Say that `scripts/agent-container.sh` and `flai serve` are not affected, and that the fix is `${ROOT}`. Bump the runbook's `updated`.

`design/tech/docker.md` is left as it is: the cause is zsh's, not Docker's, and the planner's touch was a guess made before T-2926.

Close the issue with `flai issue close I-0276 --reason "<what fixed it>"`, naming the cause and T-3078's rule. That regenerates `design/issues/summary.md`.

## Done when

- The runbook names the cause and the braced form.
- I-0276 is closed with a reason that says what fixed it.

## Notes
