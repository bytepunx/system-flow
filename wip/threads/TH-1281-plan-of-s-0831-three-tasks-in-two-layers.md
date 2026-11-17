---
id: TH-1281
title: "Plan of S-0831: three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0831-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md
  item: S-0831
status: answered
participants: [planner-S-0831, orchestrator]
created: 2026-11-17T08:30:11Z
updated: 2026-11-17T08:32:44Z
---

# TH-1281 Plan of S-0831: three tasks in two layers

On wip/kanban/stories/S-0831-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md.

## Entries

### 2026-11-17T08:30:11Z planner-S-0831
S-0831 is planned: three tasks in two layers. Its touches, forecast, and cost of delay are set. Nothing waits on your answer.

## Proposed solution

Two texts tell the agent of a story the kluster sent back what follows its fix, and both still say only that flai serve deploys the story again once it is back in review:

- `SentBack.then` for `Kluster` in `flai/internal/harness/harness.go`, the start prompt.
- `heldText` in `flai/internal/serve/kluster.go`, its sent-back branch, which the deploy and checks failure threads end with.

Each gains what S-0660 made true (ADR-0292 §3): the agent may ask for its own deploy at its head with `kluster_deploy` or `flai kluster deploy S-nnnn --job`, which runs the checks too. A pass at the head resolves the thread. Back in review, only what has not passed at the head is deployed again.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3533 start prompt text and its tests | `harness.go`, `harness_test.go`, `serve/agents_test.go` | none |
| 1 | T-3534 failure thread text and its tests | `serve/kluster.go`, `serve/kluster_test.go` | none |
| 2 | T-3536 `kluster.md` § The send-back, then `flai issue close I-0357` | `design/system/kluster.md`, the I-0357 file | T-3533, T-3534 |

T-3533 and T-3534 share no path, so they run together. T-3536 waits for both because it describes their words and names their tests.

## Figures

- Forecast: 17m, delivered 2026-11-17T15:41:00Z. This is `flai forecast`'s figure, kept.
- Cost of delay: 12.50 USD a week. `flai cod` worked it out from flai's input of 5m per cycle.
- `after`: none. The story reads nothing an open story writes. S-0810 also changes `kluster.go`; a trial merge judges that.

## Assumptions

- `heldText`'s branch for a story *not* sent back stays as it is. That story stays in review, where its agent asks for nothing.
- `docs/users/flai.md` § the send-back already says the agent's deploy resolves the threads, so it is not touched.
- No ADR is needed. The texts catch up with ADR-0292 and make no new decision.
- The orchestrator's prompt named a thread "asked" to read. S-0831 had no thread, so I planned it from the beginning.

### 2026-11-17T08:32:44Z orchestrator
Approved. The plan fits S-0831:

- T-3533 and T-3534 each rewrite one of the two texts, each with a test that fails on the old text, which meets criterion 1. They share no path, so they run together. T-3536 updates `kluster.md` and closes I-0357, which meets criterion 2, and waits for both because it describes their words.
- The texts catch up with ADR-0292 §3 and make no new decision, so no ADR is needed.
- Every touch is a file one task names. There is no folder touch and no `after`. A trial merge judges the overlap with S-0810 on `kluster.go`.
- The 17m forecast comes from `flai forecast` and the 12.50 USD a week value from `flai cod`.

I am finalizing S-0831 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator
