---
id: T-1501
type: task
nature: remediation
title: Opening a stream over an empty narrative writes it afresh from the template
status: done
parent: S-0375
owner: alex
created: 2026-10-13T16:36:32Z
updated: 2026-10-13T16:58:17Z
transitions:
  - to: ready
    at: 2026-10-13T16:39:13Z
    by: agent-S-0375
  - to: in-progress
    at: 2026-10-13T16:39:13Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T16:58:17Z
    by: agent-S-0375
stream: S-0375
tags: [flai]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go]
usage:
  source: log
  seconds: 1144
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 32
      output: 11153
      cache_read: 1467600
      cache_write: 60454
      cost: 0.8965
---
# T-1501 Opening a stream over an empty narrative writes it afresh from the template

## Work

A zero-byte `wip/agents/S-nnnn.md` holds nothing, yet `OpenStream` (`flai/internal/workitem/narrative.go`) refuses it as existing, and `ReopenStream` fails on it in `ReadNarrative` with `ErrNoFrontMatter`. `storystart.openStream` falls through from one to the other, so `flai story start` and `flai stream open` stop on the empty narrative I-0137 describes.

- In `OpenStream`, treat a narrative file of zero bytes, or of whitespace only, as absent: write it afresh from the template, as for a story with none. Log at info that the empty narrative was replaced.
- Leave a narrative with any content to `ReopenStream` as now: never replace one that holds text, even with broken front matter, since that may be the agent's work.
- Test both: a zero-byte narrative is written afresh by `OpenStream`, and a narrative with content still gives `*StreamExistsError`.

Waits for nothing: no other task touches these files, and the story's start code reaches the change through `OpenStream` unchanged.

## Done when

- `OpenStream` writes a narrative over a zero-byte or blank file, and still refuses one with content.
- `flai test flai/internal/workitem` passes.

## Notes

Drafted by the planner from I-0137's instance and `narrative.go` lines 112-160.
