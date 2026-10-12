---
id: T-1217
type: task
nature: feature
title: A rebase that stops on conflicts at flai stream sync or flai accept is logged with its story, paths, and time
status: done
parent: S-0337
owner: alex
created: 2026-10-07T20:17:22Z
updated: 2026-10-12T15:46:25Z
transitions:
  - to: ready
    at: 2026-10-12T15:35:41Z
    by: agent-S-0337
  - to: in-progress
    at: 2026-10-12T15:35:42Z
    by: agent-S-0337
  - to: done
    at: 2026-10-12T15:46:25Z
    by: agent-S-0337
stream: S-0337
tags: [flai]
touches: [flai/internal/storygit/conflicts.go, flai/internal/storygit/conflicts_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go]
usage:
  source: log
  seconds: 635
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 64
      output: 321
      cache_read: 3156037
      cache_write: 119830
      cost: 2.8271
---
# T-1217 A rebase that stops on conflicts at flai stream sync or flai accept is logged with its story, paths, and time

## Work

Record what `flai stats` cannot derive from files today. It waits for nothing: the record shape is settled in the narrative's Decisions, so it runs beside T-1216.

- A log under `.flai-cache`, as `overlaps.jsonl` is, with one line per rebase that stops on conflicts: the story, where (`sync` or `accept`), the paths, and the time. Generated files settled without an agent (ADR-0098) are left out.
- Also one line per trial-merge conflict that S-0332 opens or adds to a conversation for, so a conflict found early is counted apart from one met late.
- Write it from `flai stream sync` and `flai accept`; a line that cannot be written is logged and the command goes on.

## Done when

- Tests cover a stop at sync, a stop at acceptance, a trial-merge conflict, and a generated file left out.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
