---
id: TH-1323
title: "S-0847: a 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread"
anchor:
  path: wip/kanban/stories/S-0847-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md
  item: S-0847
status: open
participants: [agent-S-0847, orchestrator]
created: 2026-11-27T22:35:12Z
updated: 2026-11-27T22:36:06Z
---

# TH-1323 S-0847: a 15-minute settle window before flai serve opens, writes to, or resolves the starved board's thread

On wip/kanban/stories/S-0847-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md.

## Entries

### 2026-11-27T22:35:12Z agent-S-0847
**Recommended: accept the proposal as written.** I write it as an ADR refining ADR-0285 §4 and §5 and then build it. No code is written until you answer.

**Cause (I-0389).** ADR-0285 §5 resolves the thread at the first look that finds the board not starved, and "the next starving opens a new thread". §4 adds an entry at every change of the list's text. A board that flips while drafts are finalized and promoted one at a time got 48 threads in five days. TH-1166 lived 6 s and TH-1135 51 s. TH-1134 gained 12 entries in four minutes, and TH-1196 gained 17 in sixteen.

**Proposal.**

1. **Open late.** `.flai-cache/starved.json` keeps when a look first found the board starved. The thread opens only at a look that finds the board still starved a settle window later. A flip back before then opens nothing and forgets that time.
2. **Resolve late.** While the thread is open, a look that finds the board not starved keeps when that began. The thread is resolved only at a look that finds the board still not starved a window later. A starving again within the window keeps the same thread open, and adds an entry only when the list changed.
3. **Coalesce entries.** flai serve writes a changed list at most once per window since its last entry on the thread, and writes the newest list then.
4. **The window.** A fixed 15 minutes in the code, with no manifest key. TH-1166, TH-1135, and the threads like them would not have opened. TH-1196's sixteen minutes of churn would have made one or two entries, not 17.

The rest of ADR-0285 stays: the anchor, the author, the title, the list's text, the reminder-channel line, and the rule that a thread you resolved is not opened again until the list changes.

**Alternatives, not recommended.**

- **A manifest key for the window.** It adds a setting, a row in `docs/operators/settings.md`, and a check, for a figure no operator has asked to tune.
- **Reopen the last resolved thread instead of keeping it open.** The thread still flips between open and resolved, and the inbox shows every flip.
- **Leave it to the orchestrator.** ADR-0285's alternatives already reject this.

**Out of scope.** The 48 `ready-empty` events. They start the orchestrator, not a thread to you, and ADR-0279 counts them.

**What each answer leads to.**

- **Accept:** I write the ADR, update `design/system/workflow.md` § Cadence, then build it in `flai/internal/serve/starved.go` with a test that reproduces I-0389.
- **Another window, or a manifest key:** I build it the same way with your figure, plus the setting and its row in `docs/operators/settings.md` if you want the key.
- **Any other change:** I rewrite the proposal and ask again.

### 2026-11-27T22:36:06Z orchestrator (recommendation)
**Recommended:** accept agent-S-0847's proposal as written. That means a fixed 15-minute settle window before flai serve opens, resolves, or writes to the starved board's thread, set in the code with no manifest key. Record it in a new ADR that refines ADR-0285 §4 and §5.

**Why:**

- I-0389's cause sits in ADR-0285 §5 and §4. §5 resolves the thread at the first look that finds the board not starved, and says "the next starving opens a new thread". §4 adds an entry at every change of the list. The proposal changes only those two points.
- The aims of ADR-0285 still hold. A board that stays starved still gets one thread that names each decision. A board that flips within minutes, while drafts are finalized and promoted one at a time, no longer opens a thread for each flip.
- A fixed window adds no setting, no row in `docs/operators/settings.md`, and no check. No operator has asked to tune this figure. A key can follow later if one does.
- The two alternatives the proposal rejects are rightly rejected. Reopening the last thread still shows every flip in the inbox. Handing it to the orchestrator is already rejected by ADR-0285 § Alternatives considered.

**Why I am not answering it myself:** it is a decision no ADR records yet, so the answer is yours.

**Alternative:** choose another window, or ask for a manifest key. agent-S-0847 builds it the same way with your figure.

Source: design/adrs/0285-when-nothing-can-be-made-ready-flai-serve-opens-one-thread-to-the-operator-on.md § Decision
