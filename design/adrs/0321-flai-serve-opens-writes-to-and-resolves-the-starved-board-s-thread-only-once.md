---
id: ADR-0321
title: "flai serve opens, writes to, and resolves the starved board's thread only once the board has held its state for a settle window of 15 minutes"
status: accepted
date: 2026-11-27
supersedes: []
superseded_by: []
refines: [ADR-0285]
---

# ADR-0321 flai serve opens, writes to, and resolves the starved board's thread only once the board has held its state for a settle window of 15 minutes

## Context

[I-0389](../issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md), from the analysis of 2026-11-16, records 48 threads titled "The board is starved: nothing in ready can start and no backlog story can be promoted". flai serve opened them on `wip/kanban/board.md` between 2026-11-12 and 2026-11-16, 28 of them on 2026-11-13. TH-1166 was open for 6 s and TH-1135 for 51 s. TH-1134 gained 12 entries in four minutes, and TH-1196 gained 17 in sixteen. In W46 the operator was sent 244 threads, and these 48 were a fifth of them.

The cause is in [ADR-0285](0285-when-nothing-can-be-made-ready-flai-serve-opens-one-thread-to-the-operator-on.md) §4 and §5. §5 resolves the thread at the first look that finds the board not starved, and says "the next starving opens a new thread". §4 adds an entry at every change of the list's text. While drafts are finalized and promoted one at a time, the board flips between starved and not starved from one look to the next. Each flip opens or resolves a thread, and each change to the list adds an entry.

The operator accepted this proposal on TH-1323.

## Decision

**flai serve opens, writes to, and resolves the starved board's thread only once the board has held its state for a settle window of 15 minutes.** The window is fixed in the code, with no manifest key. This refines ADR-0285 §4 and §5. The rest of ADR-0285 stands: when the board is starved (§1), the anchor, the author, and the title (§2), the list and the reminder-channel line (§3), and who answers the thread (§6).

1. **Open late.** `.flai-cache/starved.json` keeps when a look first found the board starved. flai serve opens the thread only at a look that finds the board still starved at least a window after that time. A look that finds the board not starved before then opens nothing and forgets that time.
2. **Resolve late.** While the thread is open, `.flai-cache/starved.json` keeps when a look first found the board not starved. flai serve resolves the thread only at a look that finds the board still not starved at least a window after that time. Its reason is the one that look gives. A look that finds the board starved again within the window forgets that time and keeps the same thread open.
3. **Coalesce entries.** flai serve writes a changed list on the open thread only when at least a window has passed since its own last entry there. It then writes the newest list. A changed list that comes sooner waits for a later look. A list that changes back to the one last written writes nothing.
4. **Kept across a restart.** The times in 1 to 3 are kept in `.flai-cache/starved.json`, beside the thread and the list last written, so a restart neither opens nor resolves the thread early. When the file cannot be read, flai serve finds the open thread by its anchor, author, and title, as before, and the window starts from that look.
5. **A thread the operator resolved.** ADR-0285's rule stays. When the operator resolves the thread while the board stays starved, flai serve opens a new one only when the list changes. The board has already been starved for the window by then, so the new thread opens at that look.

## Consequences

- A board that flips faster than the window opens no thread. TH-1166, TH-1135, and the threads like them would not have been opened.
- A board that stays starved still gets one thread, at most 15 minutes after it starved. The thread names each decision that would start work again, as ADR-0285 meant.
- A board that comes back to life for less than the window keeps its one thread. The operator sees no flip in the inbox.
- A list that keeps changing gains at most one entry per window. TH-1196's sixteen minutes of churn would have made two entries, not 17.
- The thread is resolved up to 15 minutes after a story can start. flai serve's agent starts do not wait for that: the thread starts nothing and holds nothing.
- The `ready-empty` event is unchanged. It starts the orchestrator, not a thread to the operator, and ADR-0279 counts it.

## Alternatives considered

- **A manifest key for the window.** It adds a setting, a row in `docs/operators/settings.md`, and a check, for a figure no operator has asked to tune. A key can follow if one does.
- **Reopen the last resolved thread instead of keeping it open.** The thread would still flip between open and resolved, and the inbox would show every flip.
- **Leave it to the orchestrator.** ADR-0285's alternatives already reject this. The orchestrator runs only while the `orchestrate` host action is on, and it would spend a model run on what flai can work out itself.
