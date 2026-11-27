---
id: D-0677
title: "S-0847: a 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread: Recommended: accept agent-S-0847's proposal as written"
kind: answer
sources: [TH-1323, I-0389, S-0847]
topics: [threads, ready, orchestration, flai, cli, go, code, all]
paths: [flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/readyempty_test.go, design/system/workflow.md, docs/users/flai.md, design/adrs/README.md, design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md, design/issues/summary.md, design/adrs/0321-flai-serve-opens-writes-to-and-resolves-the-starved-board-s-thread-only-once.md]
date: 2026-11-27
status: absorbed
absorbed_by: "design/adrs/0321-flai-serve-opens-writes-to-and-resolves-the-starved-board-s-thread-only-once.md#decision"
asked_again: []
cited: []
---

# D-0677 S-0847: a 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread: Recommended: accept agent-S-0847's proposal as written

The operator, alex, answered TH-1323 at 2026-11-27T22:57:15Z by accepting agent-S-0847's proposal: a fixed 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread, set in the code with no manifest key, recorded as an ADR refining ADR-0285 §4 and §5.

The answer:

```text
**Recommended:** accept agent-S-0847's proposal as written. That means a fixed 15-minute settle window before flai serve opens, resolves, or writes to the starved board's thread, set in the code with no manifest key. Record it in a new ADR that refines ADR-0285 §4 and §5.
```

The reasons given: I-0389's cause sits in ADR-0285 §5, which resolved the thread at the first look that found the board not starved, and §4, which added an entry at every change of the list. A fixed window adds no setting, no row in `docs/operators/settings.md`, and no check. The two rejected alternatives were a manifest key and reopening the last resolved thread.

S-0847 built it as ADR-0321 decides, which closed I-0389. From I-0389's fix:

```text
S-0847 fixed it as ADR-0321 decides. flai serve now opens the starved board's thread only once the board has stayed starved for 15 minutes, and resolves it only once the board has stayed not starved for 15 minutes. A starving within that window keeps the same thread. A changed list is written at most once every 15 minutes, and the times are kept in .flai-cache/starved.json across a restart (flai/internal/serve/starved.go, e7e2bb2).
```

The 48 `ready-empty` events stayed out of scope: they start the orchestrator, not a thread to the operator.

Sources: TH-1323; I-0389; S-0847 (`wip/archive/agents/S-0847.md`, Decisions); ADR-0321.
