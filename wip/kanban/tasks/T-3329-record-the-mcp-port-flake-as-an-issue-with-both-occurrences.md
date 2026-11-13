---
id: T-3329
type: task
nature: remediation
title: Record the MCP port flake as an issue with both occurrences
status: done
parent: S-0779
owner: alex
created: 2026-11-13T20:33:25Z
updated: 2026-11-13T20:40:12Z
transitions:
  - to: ready
    at: 2026-11-13T20:39:52Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T20:39:52Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T20:40:12Z
    by: agent-S-0779
stream: S-0779
tags: [flai]
touches: [design/issues/summary.md]
usage:
  source: log
  seconds: 20
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 33
      cache_read: 311169
      cache_write: 967
      cost: 0.1172
---
# T-3329 Record the MCP port flake as an issue with both occurrences

## Work

Record the flake with `flai issue new`, class `defect`, before the fix lands, so the issue holds what was seen. Give it both occurrences on `main`, the second with `flai issue bump`:

- Run 43495183748, commit `9b63bacf3`, 2026-11-13: `TestTheHostMovesAProjectOffAPortAnotherProcessHolds` at `serve_mcp_test.go:516`, the first start given `127.0.0.1:4243` in place of the remembered port.
- Run 43455794271, flai 1.93.1: `TestMCPListenFailureClosesWhatItOpened` at `mcp_http_test.go:805`, `bind: address already in use` on `127.0.0.2`.

Say the cause as the story's Goal does: the tests count on a port they closed, or never held, staying free. Give the cost per occurrence as the time to read the failing job's log and run it again. Name the issue's ID in the story's Notes, above `### Planning`, so that `issues.no-story` finds the story that remediates it.

Waits for nothing. It touches only `design/issues/`, which no other task of this story changes before the last.

## Done when

- An open issue under `design/issues/` names both runs, their tests, and their errors, with a count of 2.
- `design/issues/summary.md` lists it.
- The story's Notes name the issue by its ID.

## Notes
