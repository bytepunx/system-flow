---
id: T-1777
type: task
nature: feature
title: flai accept accepts an instance whose branch changes nothing, and the orchestrator's acceptance takes each step task's done transition as its criterion's evidence
status: done
parent: S-0422
owner: alex
created: 2026-10-16T05:47:54Z
updated: 2026-11-02T03:19:38Z
transitions:
  - to: ready
    at: 2026-11-02T03:04:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T03:04:29Z
    by: system-flow
  - to: done
    at: 2026-11-02T03:19:38Z
    by: agent-S-0422
stream: S-0422
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
usage:
  source: log
  seconds: 909
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 65
      output: 18736
      cache_read: 4759831
      cache_write: 100253
      cost: 1.9109
---
# T-1777 flai accept accepts an instance whose branch changes nothing, and the orchestrator's acceptance takes each step task's done transition as its criterion's evidence

## Work

- An instance's steps may change no file, so its branch may hold no commit beyond main. `flai accept` merges it as a fast-forward no-op, moves it to done, archives, commits, and follows the release rules of its nature; the preview names no blocker for the empty diff. Test it in a new `cmd/accept_formula_test.go` with a story that has `formula:` and an empty branch, and check with `preview/accept_test.go` that a story without `formula` and an empty branch is accepted the same way, since nothing refused it before either.
- Under `--by orchestrator`, `unevidenced` in `preview/orchestrator.go` refuses an evidence line that names no changed file. For a story with `formula:`, a criterion's evidence is instead the done transition of the task of its step, `- <n>: T-nnnn done`; a criterion whose task is not done is a blocker, named. The four orchestrator blockers hold as for any story.
- The convention `strategic-agents.md`, here and in the template, says in "As the orchestrator" step 2 what the evidence of an instance is.
- Tests: the orchestrator's acceptance of an instance with every task done passes; one with a task not done is refused, naming it; an ordinary story's evidence is checked as before.

First layer. It shares no path with the other tasks of the layer.

## Done when

`go test ./cmd/... ./internal/preview/...` passes, and the existing accept tests are unchanged.

## Notes

Rewritten by the planner on 2026-10-21, on TH-0471's answer 5 and TH-0582's approval. It accepted a story with no branch before; an instance now has one.
