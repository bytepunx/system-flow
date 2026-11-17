---
id: ADR-0306
title: "Every write of board.md takes one board lock, reads the board again inside it, and normalises the pull order against the items before it writes"
status: accepted
date: 2026-11-17
supersedes: []
superseded_by: []
---

# ADR-0306 Every write of board.md takes one board lock, reads the board again inside it, and normalises the pull order against the items before it writes

## Context

`wip/kanban/board.md` holds the WIP limits, the pull order, and who placed each story in it. Four writers change it: a move of a story (`TransitionAll` in `flai/internal/workitem/transition.go`, behind `flai move`, `flai story start`, and the MCP tools), `flai order` (a placement, and `--by <policy> --apply`), `flai accept`, and `flai board limit`. Each one reads the board with `LoadBoard`, changes its copy, and writes the copy whole with `Board.Save`. `Board.Save` takes no lock. A work item's `Repo.Save` takes the items lock and merges a stale copy (ADR-0184), but the board has nothing like it.

I-0375, ``flai check finds `board.order` outside the story at close-out``, has two instances, from the close-outs of S-0738 and S-0746 on 2026-11-15. Both found "order lists S-0738 which is in-progress". The cause is in the history of `board.md` and the archived items. S-0738, S-0746, and S-0739 moved to in-progress in the same second, 2026-11-14T23:49:15Z, when `flai serve` started their agents after S-0737's acceptance. Each start's `TransitionAll` read the board and removed its own story from `order`. Then it saved its items and wrote the board back. One start wrote back a board it had read before another start removed S-0738. So S-0738 stayed in the order while it was in progress, until S-0746's acceptance wrote the board again at 01:56Z.

`design/system/workflow.md` § The pull order says that on every write the list is normalised and names that are no longer ready or backlog stories are dropped. Only `flai order`'s placement does that today, through `Board.Place`. A move, an acceptance, and a limit write keep whatever stale name an earlier write left.

## Decision

Every write of `board.md` is made under one board lock: the board is read again inside the lock, the change is made to that copy, the pull order is normalised against the items read inside the lock, and the board is written before the lock is let go.

1. **One board lock.** It is an operating system lock on `.flai-cache/board.lock`, taken through `filelock` as the items lock is, with the same wait. It is apart from the items lock (ADR-0184) and the checked-write lock (ADR-0272).
2. **The lock order.** The board lock is taken before the items lock and is never taken while the items lock is held. A writer that holds the board lock may save items, and each such save takes and lets go of the items lock inside it. No holder of the items lock waits on the board lock: `Repo.Save`, `SaveAll`, and `Restore` hold the items lock for one read, merge, and write and call nothing back. No checked edit or checked creation writes the board, so the board lock is never taken under the checked-write lock either.
3. **Read again inside the lock.** A writer reads the board after it takes the lock, never before, and makes its change to that copy. A move of a story or an epic holds the board lock from its reading of the items and the board until it has saved the items and the board. So two moves started together write the board one after the other, and each sees the other's change. `flai order`, `flai accept`, and `flai board limit` make their change in the same way. `flai order --by <policy> --apply` works its order out from the board and the items it reads inside the lock.
4. **Normalised on every write.** Before it writes, every write drops from `order` each name that is not an active story in a state the workflow orders, ready or backlog on the baseline. It drops the name's `placed` entry with it, and keeps the other names in their order. A name whose item file does not parse is kept, since its state cannot be told. So a stale name left by any earlier write, or by a hand edit, goes at the next write.
5. **`board.order` stays recorded.** A close-out keeps recording `board.order` outside the story in an issue (ADR-0085). Unlike `board.wip-limit` (ADR-0133), it names a defect in how the board was written, not a state only the operator's acceptances clear. An instance found after this decision points at a writer the lock misses.

## Consequences

- Two story starts in the same second no longer lose one another's removal from the pull order, and I-0375's cause no longer occurs. A test in `flai/internal/workitem` moves several stories in the order from goroutines at once and finds none of them left in it.
- A stale name from before this decision goes at the next write of the board, whoever makes it.
- Moves of stories and epics are serialised across the project for the time a move takes: one read of the items, the saves, and one write of the board. A task's move does not take the board lock.
- `Board.Save` stays, as the write a holder of the board lock makes. Its doc comment says so, and nothing outside the locked update calls it.
- The lock file lives in `.flai-cache`, which a project must ignore, as `layout.gitignore` already requires.

## Alternatives considered

- **Leave `board.order` out of a check scoped to a story, as ADR-0133 does with `board.wip-limit`.** The close-out would go quiet, but the lost update would stay. The finding names a real defect that a story branch cannot cause.
- **Take the items lock for the board too.** A move holds it across its item saves, and each save takes it again. That would need a reentrant lock or every save rewritten, and `filelock` is not reentrant.
- **Merge a stale board copy with the file, as ADR-0184 merges an item.** Two removals of different stories from one list both change the key `order`, so a merge by key refuses one of them. A merge of the list itself needs rules for every kind of change to `order`. A lock and a fresh read need neither.
- **Normalise only in `flai order`, as today.** A move or an acceptance could still write back a stale name, and nothing would drop it until the next placement.
