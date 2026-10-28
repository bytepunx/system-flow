---
id: TH-0701
title: "Plan for S-0529: refuse a hand acceptance while the queue can batch until the queue tiers ran at the story's head"
anchor:
  path: wip/kanban/stories/S-0529-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md
  item: S-0529
status: open
participants: [planner-S-0529, orchestrator]
created: 2026-10-28T18:18:12Z
updated: 2026-10-28T18:18:24Z
---

# TH-0701 Plan for S-0529: refuse a hand acceptance while the queue can batch until the queue tiers ran at the story's head

On wip/kanban/stories/S-0529-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md.

## Entries

### 2026-10-28T18:18:12Z planner-S-0529
**The fix I propose.** The issue leaves the fix to the story. I recommend this one:

- While the merge queue can batch (ADR-0174), `flai accept` by hand is refused for a story whose last `flai verify` at its branch's head did not run the `queue` tiers its paths select. That covers `flai accept`, `flai move <S> done`, the board, and the dashboard.
- The refusal names the tiers and `flai verify <S-nnnn> --all`.
- `--dry-run` and the dashboard's preview list it as a blocker.
- One CLI flag lands the story anyway, and the story's Notes record which tiers did not run.
- `flai verify --all` is made to reuse the tiers a close-out at the same patch passed, so it runs only `integration` and `smoke`.
- A new ADR refines ADR-0174, whose consequences say a hand acceptance runs no tier.

Please confirm this, or pick an alternative. The tasks are planned for it.

| Alternative | Why I did not pick it |
|---|---|
| `flai accept` runs the queue tiers itself before it merges | Takes about 12 minutes here. That is past the dashboard's 10-minute timeout on `accept.run` and `item.move`, so the host API and flaiover would need changes too. |
| An acceptance by hand enqueues the story, and the queue lands it as the operator | Changes what "done" from the board means, and a red batch would undo the operator's acceptance. |

**Tasks and layers.**

| Layer | Task | Touches | After |
|---|---|---|---|
| 1 | T-2261: the ADR refining ADR-0174 | `design/adrs/README.md` and the new ADR | none |
| 1 | T-2262: `verify --all` reuses the passed tiers; a helper names the queue tiers a report did not run | `flai/internal/verify/story.go`, `story_test.go` | none |
| 2 | T-2263: the accept blocker, the refusal, and the override flag, with a test that reproduces I-0203; the reference is regenerated | `flai/cmd/accept.go`, `verify.go`, `accept_queue_tiers_test.go` (new), `move_refused_test.go`, `flai/internal/preview/accept.go`, `accept_test.go`, `docs/users/flai-reference.md` | T-2262 |
| 3 | T-2264: the docs, and `flai issue close I-0203` | `docs/users/flai.md`, `docs/operators/index.md`, `design/system/workflow.md`, `design/system/flai-cli.md`, the I-0203 file, `design/issues/summary.md` | T-2261, T-2263 |

**Figures.**

- Forecast: 1h30m. flai gave 32m; I raised it for the two Go tasks with real-git tests.
- Cost of delay: 37.50 USD a week, from `flai cod` on the 15m input. It stands as computed.

**Assumptions.**

- The override is a CLI flag only. The host API and flaiover do not change: the refusal reaches the board and the dashboard through `acceptItem`.
- With the queue off, and for the orchestrator's acceptance and the queue's landing, nothing changes.
- This instance cannot happen again in the same way. `TestReferenceIsCurrent` now runs in the `go-test` tier, since S-0490 removed `-short`. The class remains for any check only `integration` or `smoke` runs.

**A story to add, apart from this one.** A batch that fails is bisected, and a batch of one blames its story with no step (`RunBisect`, `flai/internal/mergequeue/bisect.go:105-109`). So when main itself is red, every later batch sends an innocent story back. I propose a separate story: before the queue sends a story back, it runs the failing tier on the batch's base alone. When main fails too, the queue holds and opens one thread naming main. I have not created it.

### 2026-10-28T18:18:24Z orchestrator (recommendation)
**Recommendation:** confirm the proposed fix. While the queue can batch, a hand acceptance is refused until the `queue` tiers have run at the story's head. The refusal names `flai verify <S-nnnn> --all`, and one CLI flag lands the story anyway, recorded in its Notes. Leave the bisect story for now, unless you want it drafted.

## Why this needs you

- **It reverses a recorded consequence.** ADR-0174 § Consequences says a hand acceptance is unchanged and runs no test tier. The new ADR refines it.
- **It limits your own acceptance.** The refusal covers the board, the dashboard, `flai move done`, and `flai accept`. Only you can choose to bind yourself there.

## Why I recommend this fix

- It keeps acceptance fast and the dashboard's 10-minute timeout untouched. Running the tiers inside `flai accept` would take about 12 minutes here.
- It keeps "done from the board" meaning what it means today. Enqueueing a hand acceptance would let a red batch undo it.
- The override flag keeps you able to land a story when you judge the tiers irrelevant.

## The plan: it fits

- **Tasks:** four tasks in three layers. T-2261 and T-2262 share no path. T-2263 waits for the reuse helper, and T-2264 waits for the code and the ADR.
- **Touches:** file by file, except `design/adrs`. That folder is kept for the ADR, whose name is not known yet, and it lies inside `claims.shared`.
- **Figures:** 1h30m against flai's 32m, with the reason stated. 37.50 USD a week from flai's 15m input.

## The proposed bisect story

It is a change of scope apart from I-0203. A batch of one blames its story while main itself is red. Say so here if you want the planner to draft it.

The thread stays open for your answer.

Source: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md § Consequences
