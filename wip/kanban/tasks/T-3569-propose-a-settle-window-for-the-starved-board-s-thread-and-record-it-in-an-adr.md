---
id: T-3569
type: task
nature: improvement
title: Propose a settle window for the starved board's thread and record it in an ADR refining ADR-0285
status: backlog
parent: S-0847
owner: alex
created: 2026-11-17T12:29:29Z
updated: 2026-11-17T12:29:29Z
transitions: []
stream: S-0847
tags: [adr, design]
touches: [design/adrs/README.md, design/system/workflow.md]
---
# T-3569 Propose a settle window for the starved board's thread and record it in an ADR refining ADR-0285

## Work

The story's goal asks for a fix proposed from I-0389's instances before it is built, so this task comes first and waits for no other.

I-0389's instances show the cause. ADR-0285 §5 resolves the thread at the first look that finds the board not starved, and "the next starving opens a new thread". §4 adds an entry at every change of the list's text. A board that flips while drafts are finalized and promoted one at a time got 48 threads in five days. TH-1166 lived 6 s and TH-1135 51 s. TH-1134 gained 12 entries in four minutes and TH-1196 17 in sixteen.

Open a thread on S-0847 with the proposal, the recommended answer first:

- **Open late.** `flai serve` keeps in `.flai-cache/starved.json` when a look first found the board starved. It opens the thread only at a look that finds the board still starved a settle window later. A flip back before then opens nothing and forgets that time.
- **Resolve late.** Once the thread is open, a look that finds the board not starved keeps when that began. The thread is resolved only at a look that finds the board still not starved a settle window later. A starving again within the window keeps the same thread open, with an entry only when the list changed.
- **Coalesce entries.** A changed list is written at most once per settle window since the thread's last entry by `flai serve`, the newest list then.
- **The window.** A fixed 15 minutes in the code, with no manifest key. TH-1196's sixteen minutes of churn would have made one entry.

Name the alternatives and why they are not recommended:

- A manifest key for the window. It adds a setting, its row in `docs/operators/settings.md`, and its check, for a figure no operator has asked to tune.
- Reopening the last resolved thread rather than keeping it open. The thread still flips between open and resolved, and the inbox shows each flip.
- Leaving it to the orchestrator. ADR-0285's alternatives already reject this.

Name as out of scope the 48 `ready-empty` events. They start the orchestrator, not a thread to the operator, and ADR-0279 counts them.

Once the thread is answered, write the decision as an ADR with `flai adr new`. It refines ADR-0285 §4 and §5 and keeps the rest. Give it topics with `flai adr topics`: cli, orchestration, ready, threads, as ADR-0285 has. Accept it with `flai adr accept`. Then rewrite the starved-board paragraph of `design/system/workflow.md` § Cadence to match the ADR, linking it, and bump the file's `updated`. `flai task done` widens this task's touches with the ADR's own file, which cannot be named before `flai adr new` numbers it.

## Done when

- The operator answered the proposal on a thread on S-0847.
- An accepted ADR refining ADR-0285 records the settle window, when the thread opens, when it is resolved, and how often it gains an entry, and `design/adrs/README.md` lists it.
- `design/system/workflow.md` § Cadence describes the starved board's thread as that ADR decides, and links it.

## Notes
