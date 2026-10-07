---
id: T-1217
type: task
nature: feature
title: A rebase that stops on conflicts at flai stream sync or flai accept is logged with its story, paths, and time
status: backlog
parent: S-0337
owner: alex
created: 2026-10-07T20:17:22Z
updated: 2026-10-07T20:17:22Z
transitions: []
stream: S-0337
tags: [flai]
touches: [flai/internal/storygit/conflicts.go, flai/internal/storygit/conflicts_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go]
after: [T-1216]
---
# T-1217 A rebase that stops on conflicts at flai stream sync or flai accept is logged with its story, paths, and time

## Work

Record what `flai stats` cannot derive from files today. It waits for T-1216, whose definitions say what to record.

- A log under `.flai-cache`, as `overlaps.jsonl` is, with one line per rebase that stops on conflicts: the story, where (`sync` or `accept`), the paths, and the time. Generated files settled without an agent (ADR-0098) are left out.
- Also one line per trial-merge conflict that S-0332 opens or adds to a conversation for, so a conflict found early is counted apart from one met late.
- Write it from `flai stream sync` and `flai accept`; a line that cannot be written is logged and the command goes on.

## Done when

- Tests cover a stop at sync, a stop at acceptance, a trial-merge conflict, and a generated file left out.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
