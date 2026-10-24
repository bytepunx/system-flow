---
id: T-1825
type: task
nature: feature
title: "harness.Prompt gives the curator its prompt and the claude-code adapter starts it as the project's curator definition with FLAI_ROLE=curate"
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:05:44Z
updated: 2026-10-24T00:13:50Z
transitions:
  - to: ready
    at: 2026-10-24T00:01:35Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-24T00:01:36Z
    by: agent-S-0431
  - to: done
    at: 2026-10-24T00:13:50Z
    by: agent-S-0431
stream: S-0431
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
after: [T-1824]
usage:
  source: log
  seconds: 734
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 8324
      cache_read: 1631702
      cache_write: 53794
      cost: 0.8172
---
# T-1825 harness.Prompt gives the curator its prompt and the claude-code adapter starts it as the project's curator definition with FLAI_ROLE=curate

## Work

- `flai/internal/harness/harness.go`: `curatePrompt`, dispatched for role `curate`: keep the digest and do nothing else; prime with `--role curate`; call `inbox`; read `flai digest harvest --since <its last activity entry>`, or everything on a first run; for each candidate, write an entry with `flai digest new` when no entry holds the fact, or merge it into the one that does with `flai digest edit`; record in an entry's `asked_again` a thread that asked what the entry already answered; never write a secret's value, only where the operator said it lives; edit no file; ask the operator with `thread_open` only when an entry would need their judgement; end with a one-line summary naming the entries written and merged. Validate the request: no item, no story, no focus.
- `flai/internal/harness/adapters.go`: `claudeCodeStrategic` gains the curator (`--agent curator`, the session named `<key> curate`), and the environment carries `FLAI_ROLE=curate`.

Waits for T-1824, which makes the role known. Runs beside the guard task.

## Done when

`harness_test.go` covers the prompt's wording, a refused request with an item, and the adapter's arguments and environment; `flai test flai/internal/harness` passes.

## Notes

S-0432 adds absorption and compaction to the prompt.
