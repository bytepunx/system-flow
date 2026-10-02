---
id: T-0712
type: task
nature: improvement
title: The prompt and delegation.md say to name the task's ID in a task sub-agent's description
status: done
parent: S-0230
owner: arobson
created: 2026-10-02T17:15:36Z
updated: 2026-10-02T17:22:41Z
transitions:
  - to: ready
    at: 2026-10-02T17:16:14Z
    by: agent-S-0230
  - to: in-progress
    at: 2026-10-02T17:16:14Z
    by: agent-S-0230
  - to: done
    at: 2026-10-02T17:22:41Z
    by: agent-S-0230
stream: S-0230
tags: []
touches: [flai/internal/harness, design/conventions/delegation.md, template/root/design/conventions/delegation.md]
usage:
  source: log
  seconds: 387
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 76
      output: 24962
      cache_read: 2989907
      cache_write: 97119
      cost: 1.718
---
# T-0712 The prompt and delegation.md say to name the task's ID in a task sub-agent's description

## Work

T-0711 measures a task by the calls of the sub-agents started for it, and knows a sub-agent is a task's when the story's agent's `Agent` call names the task's ID in its description (or, failing that, in its prompt). Tell the story's agent so:

- `harness.Prompt` for `claude-code` (`flai/internal/harness/harness.go`, the paragraph that says to hand each task to a task sub-agent): name the task's ID in the sub-agent's description, so that flai measures the task by its calls. Update the prompt's test if it pins the words.
- `delegation.md` § Tasks worked by sub-agents, in `template/root/design/conventions/` first and then copied above the marker into `design/conventions/` so the two baselines stay identical: the same rule, and that a sub-agent (explorer or verifier) started for one task names that task too, while one started for the story as a whole names no task.

Keep the rest of the paragraph and section as they are: the plan with `after`, a task sub-agent per task, a layer at once only when its tasks are long. Bump `updated:` on both convention files. Do not touch `template/CHANGELOG.md` or `template/template.yaml`: the release writes them.

## Done when

- [x] `harness.Prompt` and both copies of `delegation.md` say to name the task's ID in a task sub-agent's description, and the baselines above the marker are identical
- [x] `go test ./internal/harness/...` passes in `flai/`

## Notes

Waits for nothing: layer 1, beside T-0711, with which it shares no path.
