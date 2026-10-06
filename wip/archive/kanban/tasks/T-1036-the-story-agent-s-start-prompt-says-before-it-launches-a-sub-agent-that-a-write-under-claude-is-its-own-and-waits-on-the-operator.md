---
id: T-1036
type: task
nature: remediation
title: The story agent's start prompt says, before it launches a sub-agent, that a write under .claude/ is its own and waits on the operator
status: done
parent: S-0299
owner: alex
created: 2026-10-06T21:08:13Z
updated: 2026-10-06T22:37:56Z
transitions:
  - to: ready
    at: 2026-10-06T22:32:32Z
    by: agent-S-0299
  - to: in-progress
    at: 2026-10-06T22:32:33Z
    by: agent-S-0299
  - to: done
    at: 2026-10-06T22:37:56Z
    by: agent-S-0299
stream: S-0299
tags: [flai, harness]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 323
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3751
      cache_read: 519375
      cache_write: 16866
      cost: 0.2845
---
# T-1036 The story agent's start prompt says, before it launches a sub-agent, that a write under .claude/ is its own and waits on the operator

## Work

I-0093's second half: the warning about `.claude/` writes reached S-0223's agent on a thread (TH-0190) only after it had launched the layer. The start prompt is the one text every story agent reads before it launches anything, so the warning goes there.

- In `delegation()` in `flai/internal/harness/harness.go`, the text every story agent `flai serve` starts gets, say:
  - A write to a file in a `.claude/` folder is never a sub-agent's. `flai guard` refuses it unless the operator has turned on `auto-approve`.
  - Say so in the prompt of every sub-agent whose task changes such a file, and have the sub-agent return the file's whole new content in its final message.
  - Make those writes yourself, once the layer's sub-agents are back, never while a layer runs. `permission_prompt` opens a thread on the story and holds the call until the owner answers.
  - When the owner may be away and nothing else is left, write each whole file into the worktree's ignored `.flai-cache/` folder instead. Open one thread on the story with the exact `cp` commands, and end rather than wait.
- Update the doc comment above `delegation()` to name S-0299 with the others.
- This task waits for no other: it shares no path with the guard task and runs beside it.

## Done when

- A test in `flai/internal/harness/harness_test.go` checks that the prompt for a story agent says that writes under `.claude/` are not a sub-agent's, that the story's agent makes them after the layer, and that the `.flai-cache/` and `cp` route is open when the owner is away.
- The existing prompt tests still pass: `go test ./internal/harness/...` in `flai/`.

## Notes

Drafted by the planner for S-0299. The wording follows the board watch's notes on TH-0194 and TH-0195, which agents already follow by hand.
