---
id: T-3411
type: task
nature: remediation
title: The hold thread and flai queue name the installed flai, main's unreleased flai commits, and the release to install, and I-0382 is closed
status: done
parent: S-0797
owner: alex
created: 2026-11-16T01:45:07Z
updated: 2026-11-16T03:54:29Z
transitions:
  - to: ready
    at: 2026-11-16T03:40:35Z
    by: agent-S-0797
  - to: in-progress
    at: 2026-11-16T03:40:35Z
    by: agent-S-0797
  - to: done
    at: 2026-11-16T03:54:29Z
    by: agent-S-0797
stream: S-0797
tags: [flai, merge-queue]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_held_test.go]
after: [T-3406]
usage:
  source: log
  seconds: 834
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 49
      output: 18115
      cache_read: 3027343
      cache_write: 116275
      cost: 1.6845
---
# T-3411 The hold thread and flai queue name the installed flai, main's unreleased flai commits, and the release to install, and I-0382 is closed

## Work

Give a hold with the `flai_behind` cause its own words, as ADR-0309's rule 3 says.

- `holdText` in `flai/internal/serve/mergequeue.go`, around the sentence that says "Fix them in the main checkout": under the cause, say the findings may come from the running flai, giving its version and commit, being older than main's flai sub-project, which has n commits beyond it. Give the remedy: publish with `flai release --pending`, install the release, and restart `flai serve`. The first look after that clears the hold. Keep today's words for a hold with no cause. Give the hold a title of its own under the cause only if `followHold` and `resolveHoldThreads` still find its one thread.
- `printQueue` in `flai/cmd/queue.go`: under a held preflight with the cause, print one line naming the running flai, main's commits beyond it, and the remedy. `--json` carries the cause from the state already.

It waits for T-3406, which records the cause on the hold.

Once this task's tests pass and T-3406 is done, close the issue with `flai issue close I-0382 --reason`, naming ADR-0309 and the preflight's second check with `flai.command`.

## Done when

- A test in `flai/internal/serve/mergequeue_test.go` checks that a hold with the cause opens one thread whose text names the running version, main's commit count, and the remedy, and does not say "Fix them in the main checkout". A hold without the cause keeps today's text.
- A test in `flai/cmd/queue_test.go` checks the cause's line in `flai queue`.
- `flai test flai/internal/serve flai/cmd` passes.
- I-0382 is closed with a reason naming what fixed it.

## Notes
