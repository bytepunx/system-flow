---
id: S-0159
type: story
nature: improvement
title: The dashboard's reads are answered in flai serve's process, not by starting flai
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:29Z
updated: 2026-09-29T07:00:29Z
transitions: []
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/hostapi, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0159 The dashboard's reads are answered in flai serve's process, not by starting flai

## Goal

The dashboard's read methods in the write table answer by starting a flai process that reads the repository again from nothing: `publish.preview` 235 to 276 ms, `stats.get` 282 ms, `stream.diff` 261 ms, `item.show` 167 ms, `push.pending` 157 ms. The board asks for `publish.preview` and `push.pending` at every change and they were asked about every 7.5 s. Cause 4 of `design/system/server-performance.md`. Answer the reads in `flai serve`'s own process from the code the commands print with, as ADR-0029 has the rest of the table do, and keep starting flai for writes only.

## Acceptance criteria
- [ ] `publish.preview`, `push.pending`, `stats.get`, `stream.diff`, `item.show`, `item.move.preview`, and `accept.preview` answer with no `exec.flai` phase in their `request answered` events.
- [ ] Each answers what the command's `--json` prints today: a contract test runs both on a fixture and compares them.
- [ ] `publish.preview` and `push.pending` take under 50 ms on this repository once warm.

## Tasks

## Notes

Measured by S-0152. Every flai process also pays its start, 145 ms on the operator's host (S-0160). Which reads are delegated is the `read(...)` entries and `reads: true` specs in `flai/internal/hostapi/writes.go`.
