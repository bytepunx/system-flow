---
id: T-2896
type: task
nature: improvement
title: A start request marked probe gets a prompt that reads and reports and works no story
status: done
parent: S-0682
owner: alex
created: 2026-11-10T08:45:58Z
updated: 2026-11-10T08:58:52Z
transitions:
  - to: ready
    at: 2026-11-10T08:51:04Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T08:51:05Z
    by: agent-S-0682
  - to: done
    at: 2026-11-10T08:58:52Z
    by: agent-S-0682
stream: S-0682
tags: [harness, containers]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 467
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 40
      output: 292
      cache_read: 1913664
      cache_write: 62075
      cost: 0.7271
---
# T-2896 A start request marked probe gets a prompt that reads and reports and works no story

## Work

Waits for no task: it is the first layer, and pilot and the script build on the field it adds.

- Give `harness.Request` in `flai/internal/harness/harness.go` a field `Probe bool` with the JSON tag `probe,omitempty`, documented as a hand-run trial of an agent image (I-0224): set only in a start request `scripts/agent-container.sh --probe` writes, never by `flai serve`.
- In `Prompt`, answer a request with `Probe` set with a probe prompt before any other case. The prompt names the agent (`r.Name`) and says it is a trial of the image, not the story's agent. It tells the agent to do no work on any item: to call only the flai MCP tools that read (`inbox`, `board`, `item_get`, `doc_get`, `doc_search`), to run `flai --version`, `git status`, and `git log -1` in its working directory, to report one line per check with what it saw, and to end. It names the tools that write (`story_start`, `stream_state`, `task_done`, `item_new`, `item_edit`, `item_move`, `thread_open`, `thread_reply`, `verify`) as never to be called.
- Add tests in `flai/internal/harness/harness_test.go`: the probe prompt names its agent and each read check, tells the agent never to call each writing tool, and holds none of the story prompt's instructions (no `story_start S-nnnn`, no `stream_state` to write). A request with `Probe` set and a `Story` still gets the probe prompt, not the story's.

## Done when

- `Prompt` returns the probe prompt for a request with `Probe` set, whatever its `Story`.
- The new tests pass with `flai test flai/internal/harness/harness.go flai/internal/harness/harness_test.go`.

## Notes
