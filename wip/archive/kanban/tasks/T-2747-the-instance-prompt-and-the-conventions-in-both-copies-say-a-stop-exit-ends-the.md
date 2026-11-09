---
id: T-2747
type: task
nature: improvement
title: The instance prompt and the conventions in both copies say a stop exit ends the run with flai formula stop, asking nothing
status: done
parent: S-0621
owner: alex
created: 2026-11-09T17:07:09Z
updated: 2026-11-09T18:25:07Z
transitions:
  - to: ready
    at: 2026-11-09T18:14:49Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T18:14:50Z
    by: system-flow
  - to: done
    at: 2026-11-09T18:25:07Z
    by: agent-S-0621
stream: S-0621
tags: [flai, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
after: [T-2745]
usage:
  source: log
  seconds: 617
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 5638
      cache_read: 1156392
      cache_write: 21295
      cost: 0.4671
---
# T-2747 The instance prompt and the conventions in both copies say a stop exit ends the run with flai formula stop, asking nothing

## Work

Tell the story's agent and the orchestrator what T-2745's command does, so that no agent blocks or asks on a stop exit.

- In `flai/internal/harness/harness.go`, the `instance` paragraph (about lines 605 to 631): a step that exits with one of its `stop` exits, as its task's `## Done when` names them, is no failure. The agent runs `flai formula stop` with the exit and the step's output, runs no step after it, closes the instance out, and moves it to review. "Cancel none" stays true of every other case. Update `flai/internal/harness/harness_test.go` (about lines 1087 to 1110), which pins the text.
- In `design/conventions/work-management.md`, the rule on a story with `formula:`, the same, and its copy in `template/root/design/conventions/work-management.md` word for word above the marker.
- In `design/conventions/strategic-agents.md` § Accepting a story, step 2: a skipped criterion's evidence is `- <n>: T-nnnn skipped`, naming the step task the stop cancelled. Copy it to `template/root/design/conventions/strategic-agents.md`.

This task waits for T-2745, whose command it names. It runs beside T-2746: they share no path.

## Done when

- The instance prompt, both copies of `work-management.md`, and both copies of `strategic-agents.md` say what a stop exit does and how its instance is accepted, and the copies match above the marker.
- The harness test pins the new text.
- `flai test` on the changed paths passes.

## Notes
