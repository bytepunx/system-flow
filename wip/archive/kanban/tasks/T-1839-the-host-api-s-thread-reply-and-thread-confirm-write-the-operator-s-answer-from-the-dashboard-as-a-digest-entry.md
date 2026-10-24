---
id: T-1839
type: task
nature: feature
title: The host API's thread.reply and thread.confirm write the operator's answer from the dashboard as a digest entry
status: done
parent: S-0433
owner: alex
created: 2026-10-16T06:07:44Z
updated: 2026-10-24T21:05:04Z
transitions:
  - to: ready
    at: 2026-10-24T20:58:19Z
    by: agent-S-0433
  - to: in-progress
    at: 2026-10-24T20:58:19Z
    by: agent-S-0433
  - to: done
    at: 2026-10-24T21:05:04Z
    by: agent-S-0433
stream: S-0433
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/thread.go, flai/cmd/thread_test.go]
after: [T-1840]
usage:
  source: log
  seconds: 405
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 9178
      cache_read: 1418487
      cache_write: 78370
      cost: 0.924
---
# T-1839 The host API's thread.reply and thread.confirm write the operator's answer from the dashboard as a digest entry

## Work

The host API's `thread.reply` and `thread.confirm` in `flai/internal/hostapi/writes.go` run `flai thread reply --by=<owner> --json` and `flai thread confirm --by=<owner> --json` and pass the CLI's JSON through, so the dashboard's answer is written as a digest entry by T-1840's code. Show that path in `writes_test.go`: the dashboard's reply and confirmation run the CLI as the owner, and the CLI's `digest` field reaches the dashboard's response unchanged. Change `writes.go` only if the response does not carry it.

Waits for T-1840, whose JSON it passes through.

## Done when

`writes_test.go` covers a reply and a confirmation from the dashboard whose response carries the CLI's `digest`; `flai test flai/internal/hostapi` passes.

## Notes

The dashboard need not show the entry; the response carries it for a later page.
