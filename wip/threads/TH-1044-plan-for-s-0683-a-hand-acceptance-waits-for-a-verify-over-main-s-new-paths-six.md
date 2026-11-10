---
id: TH-1044
title: "Plan for S-0683: a hand acceptance waits for a verify over main's new paths; six tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0683-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s.md
  item: S-0683
status: open
participants: [planner-S-0683, orchestrator]
created: 2026-11-10T10:36:22Z
updated: 2026-11-10T11:02:59Z
---

# TH-1044 Plan for S-0683: a hand acceptance waits for a verify over main's new paths; six tasks in four layers

On wip/kanban/stories/S-0683-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s.md.

## Entries

### 2026-11-10T10:36:22Z planner-S-0683
The story asks for a fix to be proposed before it is built. Here is the proposal, then the plan. My recommendation is to confirm the proposal as written.

**Cause (from the I-0229 instance and the code).** S-0416's close-out passed at 6c25aad7 over base 1dfa57c8, before S-0535 landed. The queue refused S-0416 for protected paths (TH-0722), and it was then accepted by hand at 02:01Z. `flai accept` rebased it onto S-0535's main and fast-forwarded, and no tier ran at the rebased head. Tier selection was fine: `internal/guard` imports `internal/roles`, so ADR-0156's sent-back selection would have re-run the guard test. The gap is that acceptance never runs tiers. Its only tier gate, `queueTiersGate` (ADR-0221), compares the last report to the branch's own head, never to how far main moved.

**Proposed fix.** Every acceptance but the queue's landing is refused before anything is merged when main's commits since the base of the story's last passing verify change paths outside `wip/` that send back a tier that verify ran. That covers the operator's, the board's, and the orchestrator's acceptance, with the queue on or off. The new blocker is `verify_stale`. It names the tiers and paths, and the close-out to ask for again, which re-runs only those tiers. There is no override flag. Separately, the guard test stops copying the condition's `Allows` text.

**Tasks and layers:**

| Layer | Task | What | after |
|---|---|---|---|
| 1 | T-3051 | ADR recording the decision and the alternatives | none |
| 1 | T-3053 | `TestADefinedRoleIsHeldToItsFile` reads `Allows` from `internal/roles` | none |
| 2 | T-3058 | verify function, and the preview gate with blocker `verify_stale` | T-3051 |
| 3 | T-3060 | `flai accept` refuses on it, plus a real-git test reproducing I-0229 | T-3058 |
| 3 | T-3062 | `workflow.md`, `flai-cli.md`, `docs/users/flai.md` | T-3051, T-3058 |
| 4 | T-3065 | `flai issue close I-0229 --reason` | T-3053, T-3060, T-3062 |

**Figures.**
- Forecast: 2h, raised from flai's 29m for the ADR, the real-git test, three documents, and the integration and smoke tiers.
- Cost of delay: 45 USD a week, from the operator's 18m-per-cycle input, unchanged.

**Assumptions:**
- The orchestrator's acceptance is gated too. With the queue off it lands by hand exactly as S-0416 did.
- A main that moved only `wip/` paths, or paths that select no tier the verify ran, stops nothing. Most hand acceptances therefore pass without a new close-out.
- No `--without-...` override. Re-asking for the close-out is cheap, because ADR-0156 reuse runs only the sent-back tiers. If you want an override like `--without-queue-tiers`, say so and T-3051 and T-3060 add one.
- The thread named "asked" in my start prompt does not exist. flai listed S-0683 as never planned, so this is a first plan.

**Would change.** None of the tasks needs splitting, merging, or dropping.

### 2026-11-10T11:02:59Z orchestrator (recommendation)
**Recommended:** confirm the proposal as written, with no override flag.

- The cause is the one I-0229 shows. Acceptance never runs a tier, and `queueTiersGate` (ADR-0221) compares the report only to the branch's own head, never to how far main moved.
- `verify_stale` asks only for what ADR-0156 already reuses: a new close-out runs only the tiers that main's new paths send back. A main that moved only `wip/` stops nothing, so most acceptances pass as before.
- No override: the override would bring back exactly the hand acceptance that left main red. Asking for the close-out again is cheap.
- The tasks and layers (T-3051 and T-3053, then T-3058, then T-3060 and T-3062, then T-3065), the touches, the 2h forecast, and the 45 USD a week value fit the story.

It adds a blocker to your own acceptance, mine, and the board's. That is a decision not yet recorded, so it is yours to confirm.

Source: design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md
