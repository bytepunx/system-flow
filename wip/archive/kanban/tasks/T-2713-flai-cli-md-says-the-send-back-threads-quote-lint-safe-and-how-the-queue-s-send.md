---
id: T-2713
type: task
nature: remediation
title: flai-cli.md says the send-back threads quote lint-safe and how the queue's send-back goes on when its thread is refused, and I-0242 is closed
status: done
parent: S-0614
owner: alex
created: 2026-11-09T09:36:04Z
updated: 2026-11-09T10:08:22Z
transitions:
  - to: ready
    at: 2026-11-09T10:07:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T10:07:36Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:08:22Z
    by: agent-S-0614
stream: S-0614
tags: [flai, docs, merge-queue]
touches: [design/system/flai-cli.md, design/issues/I-0242-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md, design/issues/summary.md]
after: [T-2711, T-2712]
usage:
  source: log
  seconds: 46
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 3431
      cache_read: 1460145
      cache_write: 26948
      cost: 0.563
---
# T-2713 flai-cli.md says the send-back threads quote lint-safe and how the queue's send-back goes on when its thread is refused, and I-0242 is closed

## Work

In `design/system/flai-cli.md`, under § The merge queue pass and § Reading the pull request back, say that a send-back thread quotes a tier's findings or a check's output through `mdlint.Fenced` (line ends, tabs, control characters, trailing spaces, and blank runs made lint-safe), and that the queue's send-back, when its thread is still refused, opens it with the findings left out and the log folder named, so the story is sent back and the batch goes on (I-0242). Then close the issue: `flai issue close I-0242 --reason "<what fixed it, naming S-0614, mdlint.Fenced, and the fallback thread>"`, which updates `design/issues/summary.md`.

Waits for T-2711 and T-2712: it describes what they built, and the issue closes only once both are done.

## Done when

- `design/system/flai-cli.md` says both, with its `updated` date bumped.
- I-0242 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
