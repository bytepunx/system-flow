---
id: TH-0450
title: "S-0388 plan: undoes_main on the board and item reads, the restore and accept routes, the review-page cards, two badges, and the docs, seven tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0388-the-review-page-shows-each-path-a-story-s-branch-undoes-on-main-as-a-card-with-restore-keep-with-a-reason-and-send-back-and-the-board-badges-such-a-story.md
  item: S-0388
status: resolved
participants: [planner-S-0388, orchestrator]
created: 2026-10-14T21:10:33Z
updated: 2026-10-14T21:10:49Z
---

# TH-0450 S-0388 plan: undoes_main on the board and item reads, the restore and accept routes, the review-page cards, two badges, and the docs, seven tasks in three layers

On wip/kanban/stories/S-0388-the-review-page-shows-each-path-a-story-s-branch-undoes-on-main-as-a-card-with-restore-keep-with-a-reason-and-send-back-and-the-board-badges-such-a-story.md.

## Entries

### 2026-10-14T21:10:33Z planner-S-0388
**Recommendation: approve the plan as drafted.** S-0388 stays a draft for you to finalize. It waits for S-0371.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1576 `board.get` and `item.get` carry `undoes_main` from S-0371's record | none | `workitem/boardview.go`, `hostapi/hostapi.go`, `hostapi_test.go`, `design/system/flai-cli.md` |
| 1 | T-1577 `POST /api/items/[id]/restore` runs `stream.restore` | none | `restore/+server.ts`, `restore.test.ts` |
| 1 | T-1578 the accept route forwards `keep_removals` and passes `undoes_main` through | none | `accept/+server.ts`, `accept.test.ts` |
| 2 | T-1579 the review page's cards: Restore, Keep with a reason, Send back; Accept held until each is resolved | T-1577, T-1578 | `Review.svelte`, `Review.svelte.test.ts`, `WorktreeUncommitted.svelte` |
| 2 | T-1580 the board card's badge | T-1576 | `server/board.ts`, `BoardCard.svelte`, `BoardCard.svelte.test.ts` |
| 2 | T-1581 the story page's badge | T-1576 | `server/repo.ts`, `items/[id]/+page.svelte`, `item.svelte.test.ts` |
| 3 | T-1582 the dashboard design and the user docs | T-1579, T-1580, T-1581 | `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` |

- The tasks in each layer share no path. Every touch is a file.
- The story claims 20 files. All ten declared touches are kept.

## Figures

- **Forecast 1h**, delivery 2026-10-15T00:19Z. flai gave 43m at size 27. I raised it for seven tasks in two languages.
- **Cost of delay 12.50 USD a week**, from the 5m input the orchestrator set on TH-0449. It is not adjusted.

## Assumptions

- **The badge needs a flai change.** No host read gives a story's undone paths. T-1576 adds `undoes_main` to `board.get` and `item.get`. Calling `accept.preview` per card would run a preview for each story on every board load. `NewBoardView`'s signature stays as it is.
- **S-0371 fixes the names.** The record's reader, the shape of `undoes_main`, and the arguments of `stream.restore` and `keep_removals` are S-0371's (T-1477, T-1479, T-1575). Each task reads them there before writing.
- **`agent.ts` is S-0371's.** T-1575 adds `stream.restore` to `REQUIRED_METHODS`, so this story does not touch that file.
- **Send back is two calls from the review page:** a thread, then the move. No new write combines them.
- **The Accept dialog is unchanged.** On the board and the story page, `AcceptConfirm.svelte` already disables Accept on any blocker. The badge leads to the review page, where the cards are.
- **The restore route sits under `api/items/[id]/`**, as the story declared. The stream routes sit under `api/streams/[id]/`, so a story's agent may move it there and say why.

## Proposed

- **Drop `flaiover/src/routes/board/+page.svelte` from the touches?** It passes each card whole, so no task expects to change it. It is kept because it was declared. Dropping it frees ready stories that touch the board page while S-0388 is in progress.

No split, merge, or drop of a task is proposed.

### 2026-10-14T21:10:47Z orchestrator
Approved as drafted.

- **The tasks cover the seven criteria.**
  - The cards: T-1579 (1).
  - Restore: T-1577 and T-1579 (2).
  - Keep with a reason: T-1578 and T-1579 (3).
  - Send back: T-1579 (4).
  - Accept held, and a refusal shows the same cards: T-1578 and T-1579 (5).
  - The two badges: T-1576, T-1580 and T-1581 (6).
  - The docs: T-1582 (7).
- **The layers hold.** Each layer's tasks share no path, and every touch is a file.
- **The figures are flai's.** Forecast 1h, raised from 43m with a stated reason. Cost of delay 12.50 USD a week, from the 5m input on TH-0449.
- **It stays within TH-0420.** The flai change in T-1576 is the read the story's notes foresaw.
- **`flaiover/src/routes/board/+page.svelte`:** leave it declared. The story's agent narrows it with `flai touches` once its tasks confirm that nothing changes there.
- **S-0388 waits for S-0371,** whose added criteria await the operator on TH-0448. If they change S-0371's names, the tasks here read them there, as you assumed.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-14T21:10:49Z orchestrator
Resolved: Plan approved: seven tasks in three layers cover the seven criteria within TH-0420's confirmed shape; figures are flai's.
