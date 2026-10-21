---
id: ADR-0184
title: "Every work item write takes one items lock and merges a copy read before another writer's save with that save, by front-matter key and body section, refusing a part both changed"
status: accepted
date: 2026-10-21
supersedes: []
superseded_by: []
---

# ADR-0184 Every work item write takes one items lock and merges a copy read before another writer's save with that save, by front-matter key and body section, refusing a part both changed

## Context

flai writes a work item by reading it, changing its copy, and saving the copy whole with `Repo.Save`. `flai serve`, the MCP server, the dashboard, and a shell write items from different processes, and before S-0479 nothing ordered their writes. A writer that holds its copy for long saves back what another writer changed meanwhile. `itemedit.Apply`, behind the replanner's forecast writes, reads the item and runs a whole-repository `flai check` before it saves. On 2026-10-21 a forecast write undid `flai story start`'s move of S-0463 to in-progress (I-0175), and a usage write undid S-0456's (I-0176). Each story then read as ready while its agent worked it. `Repo.Save` has 48 callers.

## Decision

`Repo.Save` writes every work item under one items lock and merges a copy read before another writer's save with that save, part by part, refusing a part both changed.

- An item read from its file keeps the form it was read in, its `Marshal`. A copy carries that form, and a successful save makes what it wrote the copy's form.
- Save takes one exclusive operating-system lock for all item writes, `.flai-cache/items.lock` (`filelock`, waiting at most 30 s), and reads the file again under it.
- If the file still holds what the copy was read as, the copy is written. If another writer saved the item since, the copy is merged with that save part by part. The parts are each top-level front-matter key, the body above its first level-two section, and each level-two section.
  - A part only the copy changed takes the copy's value.
  - A part only the other writer changed keeps the file's value.
  - `updated` takes the later stamp.
  - A part both changed differently refuses the save with a `ConflictError` naming the item and the parts, and nothing is written.
- The merged item is written and becomes the caller's copy.
- A missing file, a file that does not parse, or an item made in memory is written as before.
- `Repo.Restore` takes one write back out of an item file under the same lock, with the same merge. `itemedit.Apply` uses it when its check refuses an edit, so the refusal does not undo another writer's save.

## Consequences

- No writer loses another's change to a different part of an item, and no caller changed. This covers a move against a forecast or usage write, a child line under `## Tasks` against a note under `## Notes`, and a move against a block.
- Two writers that change the same part concurrently no longer both appear to succeed: the second gets a `ConflictError` and repeats its change from a fresh read. Two `flai story start` calls on one story are such a pair.
- One lock for every item serializes item writes across the project. Each holds it for one read, merge, and write of one file.
- A write that spans several files, such as a move with its epic and the board, is atomic per file, not as a group.
- The lock file lives in `.flai-cache`, which a project must ignore, as `layout.gitignore` already requires.

## Alternatives considered

- A lock and a fresh read at every writer, the planner's plan on TH-0577. It needs each of 48 call sites rewritten to hold the lock across its read, change, and save. It must not hold the lock across `check.Run`, and it leaves any writer missed as a silent loss.
- One lock per item. It has no ordering problem for a single write, but it leaves a lock file per item in `.flai-cache` and gains nothing while writes are this short.
- Refusing any save of a stale copy outright. That is loud, but it fails the many harmless overlaps, such as a usage write during a move, and every caller would need a retry.
- The body as one part. As one part, 21 existing tests conflicted on a child line and a note written to different sections.
