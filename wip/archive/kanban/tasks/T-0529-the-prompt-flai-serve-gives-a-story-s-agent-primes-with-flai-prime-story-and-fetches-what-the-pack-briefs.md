---
id: T-0529
type: task
nature: feature
title: The prompt flai serve gives a story's agent primes with flai prime --story and fetches what the pack briefs
status: done
parent: S-0148
owner: alex
created: 2026-09-29T05:24:52Z
updated: 2026-09-29T05:25:25Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:01Z
    by: agent-S-0148
  - to: in-progress
    at: 2026-09-29T05:25:01Z
    by: agent-S-0148
  - to: done
    at: 2026-09-29T05:25:25Z
    by: agent-S-0148
stream: S-0148
tags: []
touches: [flai/internal/harness]
---
# T-0529 The prompt flai serve gives a story's agent primes with flai prime --story and fetches what the pack briefs

## Work

Change `harness.Prompt` so a story's agent primes with `flai prime --story <id>` (or the MCP tool `prime`), is told the pack is a brief within a budget, and reads a briefed document's body, or the section that bears on the story, with `doc_get` and a heading before relying on it or changing what it describes; `doc_search` finds sections. Pin the text in `harness_test.go`, including that `--cat` is gone from the story prompt.

## Done when

- [x] The prompt names `flai prime --story <id>` and `doc_get` with a heading; the test pins both and fails without them.
- [x] `go test ./internal/harness/...` passes.

## Notes
