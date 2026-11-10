---
id: ADR-0277
title: "Repo.Save refuses to write an active work item whose file was archived since it was read, and Archive moves items under the items lock"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0184]
---

# ADR-0277 Repo.Save refuses to write an active work item whose file was archived since it was read, and Archive moves items under the items lock

## Context

`flai check` found `item.duplicate-id` (I-0323): a task was both in `wip/kanban/tasks` and in `wip/archive/kanban/tasks`. `flai serve` measures a task's usage by reading it with `Repo.Get` and saving the copy with `Repo.Save`. Meanwhile `flai accept` archived the task: `Repo.Archive` renamed its file into `wip/archive` without the items lock of [ADR-0184](0184-every-work-item-write-takes-one-items-lock-and-merges-a-copy-read-before.md). The save then found no file at the copy's path. ADR-0184 writes a missing file as before, so the save wrote the whole task back at its old path and made it again.

A missing file is not always a moved one. `flai edit` saves a retitled item under its new name while the old file is still beside it, and removes the old file after. Tests and fixtures place an item in the archive by saving it there. Archiving moves an item one way only, from the active folders to the archive.

## Decision

`Repo.Save` refuses to write an active work item whose file was archived since it was read, and `Repo.Archive` moves items under the items lock.

- A copy read from a file whose path now holds no file is refused when it is not archived, no file of its ID is left in its folder, and a file of its ID is in the archive folder of its type. The error wraps `workitem.ErrItemMoved`, names the item and both folders, and says to read the item again where it now is. Nothing is written and no folder is made.
- Every other missing file is written as ADR-0184 says, so a retitle and a save into the archive work as before.
- `Repo.Archive` takes the items lock around its item moves, so a save waiting on the lock reads after the move. The narratives move after it lets go.

## Consequences

- A write from a copy read before an archive no longer makes the item again in `wip/kanban`. Its caller gets `ErrItemMoved` and decides what to do. `flai serve`'s usage write reads the item again where it now is and saves the measure to the archived copy, so the harness's figure is kept.
- Callers that save a copy read before an archive must handle `ErrItemMoved`; until they do, they report it as an error rather than writing a duplicate.
- An active item whose file was removed outright, not archived, is still written as before. No flai command removes an item file that way, save `flai edit`'s retitle, which leaves the new one.
- `Archive` holds the items lock for the moves of every item in its plan, so other item writes wait for an archive. The lock is not re-entrant, and no caller of `Archive` holds it.

## Alternatives considered

- Re-home the save into `wip/archive` inside `Save`: it writes a change into an archived item that its caller did not mean for the archive, and hides the race from the caller.
- Lock `Archive` only and leave `Save` alone: a copy read before the move still finds no file after it and writes the item again at its old path.
- Refuse any copy whose file is missing: refuses `flai edit`'s retitle and every save that places an item in a new folder, unless each item records the path it was read from.
