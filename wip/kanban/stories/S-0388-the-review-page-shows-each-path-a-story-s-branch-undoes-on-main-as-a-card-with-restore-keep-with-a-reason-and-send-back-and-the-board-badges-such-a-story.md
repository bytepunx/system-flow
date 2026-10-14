---
id: S-0388
type: story
nature: feature
title: The review page shows each path a story's branch undoes on main as a card, with Restore, Keep with a reason, and Send back, and the board badges such a story
status: backlog
owner: alex
created: 2026-10-14T21:02:31Z
updated: 2026-10-14T21:52:28Z
transitions: []
tags: [dashboard]
topics: [dashboard, cli]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/WorktreeUncommitted.svelte, "flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts", "flaiover/src/routes/api/items/[id]/restore/+server.ts", "flaiover/src/routes/api/items/[id]/restore/restore.test.ts", flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/server/board.ts, flaiover/src/lib/server/repo.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flai/internal/workitem/boardview.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flaiover.md]
after: [S-0371]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 460
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 1092
          cache_read: 7875603
          cache_write: 31932
          cost: 1.9824
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-14T21:03:59Z
  value: 12.5
  by: planner-S-0388
  at: 2026-10-14T21:10:13Z
forecast:
  duration: 1h
  delivery: 2026-10-15T01:03:00Z
  basis: "Its own forecast of 1h; 10th in the pull order with an in-progress limit of 5, behind S-0371, S-0374, S-0306, S-0241, S-0381, S-0382, S-0383, S-0384, S-0385, S-0386 and S-0387."
  by: flai
  at: 2026-10-14T21:52:28Z
finalized:
  by: orchestrator
  at: 2026-10-14T21:10:51Z
---
# S-0388 The review page shows each path a story's branch undoes on main as a card, with Restore, Keep with a reason, and Send back, and the board badges such a story

## Goal

The operator resolves a story branch that undoes main's change from the dashboard, path by path, without git commands or hand-edited files, as they chose on TH-0420. S-0371 gives flai the check, the acceptance blocker with its data (`undoes_main`), `flai stream restore` behind the host write `stream.restore`, and `accept.run`'s `keep_removals`. This story builds the dashboard's side on them.

## Acceptance criteria
- [ ] The review page shows each path in the acceptance preview's `undoes_main` as a card: the path, and main's lines that the branch removes beside the branch's version
- [ ] Restore on a card runs `stream.restore` for its path, and the preview runs again, so the card clears
- [ ] Keep on a card asks for a one-line reason, and Accept then passes the kept paths and their reasons as `keep_removals`, so the acceptance records them in the story's Notes
- [ ] Send back on a card opens a thread on the story naming the undone paths and moves the story back to in-progress with that reason
- [ ] Accept stays disabled while a card is neither restored nor kept, and a refused acceptance that carries `undoes_main` as data shows the same cards
- [ ] The board's card and the story page show "undoes main on N paths" for a story whose record lists any, linking to the review page
- [ ] The dashboard design and the user docs describe the cards and the badge

## Tasks
- T-1576 The host reads board.get and item.get carry undoes_main, the paths S-0371's record lists for a story
- T-1577 POST /api/items/[id]/restore runs the host write stream.restore for one path
- T-1578 POST /api/items/[id]/accept forwards keep_removals to accept.run and passes a refusal's undoes_main through
- T-1579 The review page shows each undone path as a card with Restore, Keep with a reason, and Send back, and holds Accept until each is resolved
- T-1580 The board's card shows "undoes main on N paths" for a story whose record lists any, linking to its review page
- T-1581 The story page shows "undoes main on N paths" for a story whose record lists any, linking to its review page
- T-1582 The dashboard design and the user docs describe the undone-path cards on the review page and the badge on the board and the story page

## Notes

Drafted by S-0371's planner, as the operator confirmed on TH-0420 (2026-10-14T18:14:21Z): "A new story builds the review-page cards and the badge, after S-0371." It reuses what S-0378 built: refusal findings as data through `accept.run` (T-1516), and findings listed on the review page with Accept enabled again (T-1520).

- It waits for S-0371, whose data and writes it calls. S-0378 is done.
- The badge needs the record's paths on the board or the story page. If no read gives them, its planner adds one to flai, or reads them from `accept.preview`.
- The orchestrator never chooses Keep: `accept.run` refuses `keep_removals` with `by: orchestrator` (S-0371, T-1575).
- Its touches are predicted from S-0378's files; its planner predicts them again with its tasks.

### Planning

Planned on 2026-10-14. The plan thread is TH-0450; the cost of delay input was asked and set on TH-0449.

**The badge's data.** No read gives a story's undone paths today. `board.get` builds its cards in `workitem.NewBoardView`, and `item.get` answers `ItemWithChildren`. Neither reads `.flai-cache/`. Calling `accept.preview` for every card would run a preview per story on each board load. So T-1576 adds `undoes_main` to both host reads, read from S-0371's record in `hostapi.go`, and leaves `NewBoardView`'s signature and its four callers alone.

**Layers:**

| Layer | Tasks | Why |
|-------|-------|-----|
| 1 | T-1576, T-1577, T-1578 | wait for no task; no path in common |
| 2 | T-1579 after T-1577 and T-1578; T-1580 and T-1581 after T-1576 | the cards call the two endpoints; the badges read T-1576's field |
| 3 | T-1582 after T-1579, T-1580, T-1581 | describes what they built |

**Touches**, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flaiover/src/lib/components/Review.svelte` | declared | the cards, the three actions, Accept's hold (T-1579) |
| `flaiover/src/lib/components/Review.svelte.test.ts` | declared | their tests |
| `flaiover/src/lib/components/WorktreeUncommitted.svelte` | layout | `otherBlockers` must leave out the undone paths' blockers |
| `flaiover/src/routes/api/items/[id]/accept/+server.ts` | declared | forwards `keep_removals` (T-1578) |
| `flaiover/src/routes/api/items/[id]/accept/accept.test.ts` | declared | its tests |
| `flaiover/src/routes/api/items/[id]/restore/+server.ts` | declared | new route behind `stream.restore` (T-1577) |
| `flaiover/src/routes/api/items/[id]/restore/restore.test.ts` | declared | its tests |
| `flaiover/src/routes/board/+page.svelte` | declared | kept as declared; it passes the card whole, so no task expects to change it |
| `flaiover/src/lib/components/BoardCard.svelte` | co-change, layout | the board's badge (T-1580) |
| `flaiover/src/lib/components/BoardCard.svelte.test.ts` | co-change | its test |
| `flaiover/src/lib/server/board.ts` | layout | `FlaiCard`, `Card`, and the mapping copy named fields only |
| `flaiover/src/lib/server/repo.ts` | co-change, layout | `itemById` passes named fields only (T-1581) |
| `flaiover/src/routes/items/[id]/+page.svelte` | declared | the story page's badge |
| `flaiover/src/routes/items/[id]/item.svelte.test.ts` | co-change | its test |
| `flai/internal/workitem/boardview.go` | layout | `BoardCard.UndoesMain` (T-1576) |
| `flai/internal/hostapi/hostapi.go` | layout | `board.get` and `item.get` fill it from the record |
| `flai/internal/hostapi/hostapi_test.go` | layout | `TestBoardGet` and `TestItemsListAndGet` sit there |
| `design/system/flaiover-dashboard.md` | declared, design | § Views and § Writes (T-1582) |
| `design/system/flai-cli.md` | co-change, design | the two host reads' rows (T-1576) |
| `docs/users/flaiover.md` | declared, design | § Reviewing a story, § Board, § The item page (T-1582) |

All ten declared touches are kept. `flai touches suggest S-0388` over them listed paths that are not predicted:

- `docs/users/flai.md` (36%) and `docs/operators/index.md` (28%): no command or setting changes.
- `flaiover/src/lib/server/agent.ts` (7%) and `flai/internal/hostapi/writes.go` (6%): S-0371's T-1575 adds `stream.restore` to both.
- `flaiover/src/lib/components/Threads.svelte` (6%): Send back posts to `/api/threads` directly.

`flaiover/src/lib/components/AcceptConfirm.svelte`, the Accept dialog on the board and the story page, is not predicted either. It disables Accept on any blocker, and the badge sends the operator to the review page, where the cards are.

**Overlaps.** `flaiover/src/lib/server/repo.ts` and `flai/internal/hostapi/hostapi.go` are busy files. While this story is in progress they hold any ready story that touches them. `design/system/flai-cli.md` is in the manifest's shared paths and holds nothing.

**Figures:**

- **Forecast 1h**, delivery 2026-10-15T00:19Z. `flai forecast` gave 43m: the median of 95 s per unit of size over 54 done large-band feature stories on claude-opus-5-5, times size 27 (7 criteria, 20 touches). I raised it because the work is seven tasks in three layers, in two languages: a Go change to two host reads, two endpoints, and component tests on three pages. S-0371's eight tasks were raised to 1h15m. The delivery moves by the 17m added to flai's 00:02Z.
- **Cost of delay 12.50 USD a week**, from `flai cod`: the orchestrator's input of 5m lost per 168h cycle (TH-0449, the planner's recommendation) at 150 USD an hour. It is not adjusted. The story saves the operator's time on each occurrence; the loss of main's change itself is S-0371's to prevent.
