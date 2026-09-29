---
id: S-0159
type: story
nature: improvement
title: The dashboard's reads are answered in flai serve's process, not by starting flai
status: done
parent: E-0012
owner: alex
created: 2026-09-29T07:00:29Z
updated: 2026-09-29T20:12:52Z
transitions:
  - to: ready
    at: 2026-09-29T19:18:32Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:55:26Z
    by: agent-S-0159
  - to: review
    at: 2026-09-29T20:12:08Z
    by: agent-S-0159
  - to: done
    at: 2026-09-29T20:12:52Z
    by: alex
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/hostapi, flai/cmd, flai/internal/metrics, flai/internal/preview, flai/internal/storygit, design/system/server-performance.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1026
  models:
    - model: claude-opus-5-5
      input: 204
      output: 91348
      cache_read: 19809639
      cache_write: 289321
      cost: 8.1043
---
# S-0159 The dashboard's reads are answered in flai serve's process, not by starting flai

## Goal

The dashboard's read methods in the write table answer by starting a flai process that reads the repository again from nothing: `publish.preview` 235 to 276 ms, `stats.get` 282 ms, `stream.diff` 261 ms, `item.show` 167 ms, `push.pending` 157 ms. The board asks for `publish.preview` and `push.pending` at every change and they were asked about every 7.5 s. Cause 4 of `design/system/server-performance.md`. Answer the reads in `flai serve`'s own process from the code the commands print with, as ADR-0029 has the rest of the table do, and keep starting flai for writes only.

## Acceptance criteria
- [x] `publish.preview`, `push.pending`, `stats.get`, `stream.diff`, `item.show`, `item.move.preview`, and `accept.preview` answer with no `exec.flai` phase in their `request answered` events.
- [x] Each answers what the command's `--json` prints today: a contract test runs both on a fixture and compares them.
- [x] `publish.preview` and `push.pending` take under 50 ms on this repository once warm.

## Tasks
- T-0573 The seven reads' work is done by internal packages the commands print from
- T-0574 flai serve answers the seven reads in its own process, held to the commands by a contract test
- T-0575 Measure the reads and record cause 4 in the design

## Notes

Measured by S-0152. Every flai process also pays its start, 145 ms on the operator's host (S-0160). Which reads are delegated is the `read(...)` entries and `reads: true` specs in `flai/internal/hostapi/writes.go`.

Done on 2026-09-29: the seven are in `flai/internal/hostapi/reads.go`; the contract test is `flai/cmd/hostapi_reads_test.go`; warm, `publish.preview` 7.2 to 8.3 ms and `push.pending` 17 to 19 ms on this repository (`design/system/server-performance.md`, cause 4). The other `read(...)` entries (`item.template`, `doc.show`, `adr.template`, `dashboard.*`, `host.*`, `checks.*`) still start flai; the story names only these seven.
