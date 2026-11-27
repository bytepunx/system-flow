---
id: D-0677
title: "S-0847: a 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread: Recommended: accept agent-S-0847's proposal as written"
kind: answer
sources: [TH-1323]
topics: [threads, ready, orchestration, flai, cli, go, code, all]
paths: [flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/readyempty_test.go, design/system/workflow.md, docs/users/flai.md, design/adrs/README.md, design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md, design/issues/summary.md]
date: 2026-11-27
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0677 S-0847: a 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread: Recommended: accept agent-S-0847's proposal as written

The operator, alex, answered TH-1323 at 2026-11-27T22:57:15Z.

What was asked:

```text
**Recommended: accept the proposal as written.** I write it as an ADR refining ADR-0285 §4 and §5 and then build it. No code is written until you answer.
```

The answer:

```text
**Recommended:** accept agent-S-0847's proposal as written. That means a fixed 15-minute settle window before flai serve opens, resolves, or writes to the starved board's thread, set in the code with no manifest key. Record it in a new ADR that refines ADR-0285 §4 and §5.

**Why:**

- I-0389's cause sits in ADR-0285 §5 and §4. §5 resolves the thread at the first look that finds the board not starved, and says "the next starving opens a new thread". §4 adds an entry at every change of the list. The proposal changes only those two points.
- The aims of ADR-0285 still hold. A board that stays starved still gets one thread that names each decision. A board that flips within minutes, while drafts are finalized and promoted one at a time, no longer opens a thread for each flip.
- A fixed window adds no setting, no row in `docs/operators/settings.md`, and no check. No operator has asked to tune this figure. A key can follow later if one does.
- The two alternatives the proposal rejects are rightly rejected. Reopening the last thread still shows every flip in the inbox. Handing it to the orchestrator is already rejected by ADR-0285 § Alternatives considered.

**Why I am not answering it myself:** it is a decision no ADR records yet, so the answer is yours.

**Alternative:** choose another window, or ask for a manifest key. agent-S-0847 builds it the same way with your figure.

Source: design/adrs/0285-when-nothing-can-be-made-ready-flai-serve-opens-one-thread-to-the-operator-on.md § Decision
```

Source: TH-1323.
