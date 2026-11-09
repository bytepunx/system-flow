---
id: T-2711
type: task
nature: remediation
title: The merge queue's send-back quotes with mdlint.Fenced and, when its thread is still refused, sends the story back with a thread that leaves the findings out
status: done
parent: S-0614
owner: alex
created: 2026-11-09T09:35:51Z
updated: 2026-11-09T10:01:47Z
transitions:
  - to: ready
    at: 2026-11-09T09:52:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:52:28Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:01:47Z
    by: agent-S-0614
stream: S-0614
tags: [flai, merge-queue]
touches: [flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go]
after: [T-2710]
usage:
  source: log
  seconds: 559
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 13770
      cache_read: 5859434
      cache_write: 108140
      cost: 2.2592
---
# T-2711 The merge queue's send-back quotes with mdlint.Fenced and, when its thread is still refused, sends the story back with a thread that leaves the findings out

## Work

In `flai/internal/mergequeue/bisect.go`, replace the package's `fenced` with `mdlint.Fenced` in `failingText`. Then make `SendBack` no longer wedge the queue on a quote: when `openFailing` is refused for the thread's text, open the thread again with a text that leaves the findings out, naming the tier, the summary on one line, and the log folder (`failed.Log`) where the whole output is, so the story is moved to in-progress, dequeued, and the bisection marked sent back as before. A refusal of that second thread stays an error, as now. Log the first refusal with the batch and the story, once, at warn.

Waits for T-2710, whose `mdlint.Fenced` it quotes with.

## Done when

- `TestTheFailingTiersOutputIsQuotedWithoutTabs` still passes, and a new test runs `SendBack` on a repository with the project's `.markdownlint.yaml` and a batch bisected to a story in review whose findings carry go test's summary line with its tabs, CRLF line ends, and ANSI escapes: the thread is written, the story is in-progress, and the bisection is sent back.
- A test where the full thread text is refused (for instance through a findings text the writer refuses, or a seam on the text) shows the story sent back with the thread that leaves the findings out and names the log folder.
- `flai test flai/internal/mergequeue` passes.

## Notes
