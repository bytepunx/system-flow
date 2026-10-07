---
id: T-1112
type: task
nature: improvement
title: "The design and the user guide describe the loop: end on a question, wait_for_events' end answer, the restart on the answer, and the empty-wake count"
status: done
parent: S-0272
owner: alex
created: 2026-10-06T22:53:37Z
updated: 2026-10-07T01:01:32Z
transitions:
  - to: ready
    at: 2026-10-07T00:59:07Z
    by: agent-S-0272
  - to: in-progress
    at: 2026-10-07T00:59:08Z
    by: agent-S-0272
  - to: done
    at: 2026-10-07T01:01:32Z
    by: agent-S-0272
stream: S-0272
tags: [docs]
touches: [design/system/workflow.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1087, T-1090, T-1094, T-1101, T-1104]
usage:
  source: log
  seconds: 144
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 55
      output: 19370
      cache_read: 2391241
      cache_write: 91775
      cost: 1.3952
---
# T-1112 The design and the user guide describe the loop: end on a question, wait_for_events' end answer, the restart on the answer, and the empty-wake count

## Work

Criterion 3. There is no `design/system/flai-serve.md`. Serve's loop is described in these places:

- `design/system/workflow.md`, under **An answered agent**
- `design/system/flai-cli.md`, in the `flai mcp` row and the `flai serve` rows

Describe the loop as built:

- **`workflow.md`, An answered agent:** a story's agent with an open question writes its state and ends, and `wait_for_events` answers `end: true` when the story has an open question and no task in progress. flai serve starts it again on the answer, with the answer in its first `inbox`.
- **`agent-narrative.md`:** the `wait_for_events` lines (about 121 and 122) say when it ends at once rather than holding.
- **`flai-cli.md`:** the `wait_for_events` part of the `flai mcp` row, and the empty-wake count in `flai stats`.
- **`docs/users/flai.md`:** the `wait_for_events` row of the MCP tools table (about line 851), and the paragraph on an agent that ended asking (about line 624).
- **`docs/users/flai-reference.md`:** `flai stats`' empty-wake count, text and `--json`.

It waits for every other task, because it describes what they built: T-1087's answer, T-1090's prompt, T-1094's rules, T-1101's count, and T-1104's tested restart.

## Done when

- Each place above describes the loop and the count as built, with nothing left that tells a story's agent to hold for an answer.
- The markdown lint and `flai check --strict` are clean.

## Notes

Drafted by the planner. `metrics.md` and the ADR are T-1097's.
