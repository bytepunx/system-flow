---
id: T-2950
type: task
nature: remediation
title: The host channel's stream.state passes context, decisions, and questions to flai stream state
status: done
parent: S-0668
owner: alex
created: 2026-11-10T08:51:03Z
updated: 2026-11-10T10:21:25Z
transitions:
  - to: ready
    at: 2026-11-10T10:13:05Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:13:06Z
    by: system-flow
  - to: done
    at: 2026-11-10T10:21:25Z
    by: agent-S-0668
stream: S-0668
tags: [flai, hostapi, narrative]
touches: [flai/internal/hostapi/writes_streams.go, flai/internal/hostapi/writes_test.go]
after: [T-2947]
usage:
  source: log
  seconds: 499
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 23
      output: 173
      cache_read: 2991057
      cache_write: 19971
      cost: 1.1083
---
# T-2950 The host channel's stream.state passes context, decisions, and questions to flai stream state

## Work

In `flai/internal/hostapi/writes_streams.go`, the `stream.state` write builds a `flai stream state` command line from `current` and `next`. Make it take `context`, `decisions`, and `questions` too:

- Pass `--context=<text>`.
- Pass one `--decision=<text>` per entry of `decisions`.
- Pass one `--question=<text>` per entry of `questions`.

Trim each text as it does now, refuse a bare `-` in each, and refuse a call with none of the five. Exit 4 stays `Refused`. So the CLI, the MCP tool, and the host channel stay one call, as `design/system/agent-narrative.md` says.

Extend `TestStreamStateWritesTheTextsGiven` in `flai/internal/hostapi/writes_test.go` to cover the three new fields, and keep the write's row in the table `TestEveryWriteIsCovered` reads.

Waits for T-2947, since the command line it builds needs the CLI's new flags.

## Done when

- `flai test flai/internal/hostapi/writes_streams.go flai/internal/hostapi/writes_test.go` passes.
- A `stream.state` with `decisions: ["x", "y"]` appends both bullets to the narrative's `## Decisions`.

## Notes
