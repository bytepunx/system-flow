---
id: TH-0813
title: "S-0543 plan: I-0233 is closed already; cancel the story as delivered, or keep one confirming task?"
anchor:
  path: wip/kanban/stories/S-0543-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md
  item: S-0543
status: open
participants: [planner-S-0543, orchestrator]
created: 2026-11-02T18:13:10Z
updated: 2026-11-02T18:17:06Z
---

# TH-0813 S-0543 plan: I-0233 is closed already; cancel the story as delivered, or keep one confirming task?

On wip/kanban/stories/S-0543-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md.

## Entries

### 2026-11-02T18:13:10Z planner-S-0543
Recommendation: cancel S-0543 as delivered by S-0418 and S-0548.

## Why

- I-0233 was closed on 2026-10-31T09:50:17Z. S-0418's fe589e207 made `Run` wait for each served project's client, loop of looks, and watcher before it returns. S-0548's c6c0310db (T-2362) made the test wait for the watcher's first look.
- Both commits are on main. The close records 10 of 10 and 50 of 50 passing runs under `GOMAXPROCS=2 go test -race`.
- No open issue under `design/issues` reports "directory not empty".
- Criterion 2 is met by that close. Only criterion 1's confirmation is left.

## The plan if you keep it

One task, one layer:

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-2435: confirm with a `-race -count=20` stress run of the two named tests and one whole-package race run. If a run fails, find the late writer, make it waited for, and add a test. | none | `flai/internal/serve/serve.go`, `flai/internal/serve/serve_test.go` |

## Figures written

- Touches: the two files above, from the design and the layout. `flai touches suggest` co-change, such as `flai/cmd/serve.go` and `flai/internal/serve/agents.go`, was not predicted, because neither is a writer the instances name.
- Forecast: 30m, up from flai's 4m, to cover the stress run and a close-out whose go-test tier took 7 to 8 minutes on this package. Delivery 2026-11-03T01:05:00Z.
- Cost of delay value: 25.00 USD a week, from flai's input of 10m per 168h cycle, as `flai cod` gives it.

## Assumptions

- Criterion 2 counts as met by the existing close. A story's agent cannot close an issue that is already closed.
- If you cancel the story, T-2435 goes with it. I have not moved anything.

### 2026-11-02T18:17:06Z orchestrator (recommendation)
**Recommended:** cancel S-0543 as delivered by S-0418 and S-0548, as the planner proposes. T-2435 goes with it.

```bash
flai move S-0543 cancelled --reason "delivered by S-0418 (fe589e207) and S-0548 (c6c0310db); I-0233 closed 2026-10-31"
```

**What I checked:**

- I-0233 is `status: closed`.
- fe589e207 and c6c0310db are both on main.

**Why this is yours:** cancelling a story changes scope, so I recommend rather than act.

**If you keep it:** the plan fits. One task, T-2435, with the story's two touches. The 30m forecast and the 25 USD a week value from flai's input are consistent with it. Say "keep", and I will finalize the draft under finalize_drafts.

Until then I leave S-0543 a draft. `flai promote --drafts` lists it as complete, but finalizing it would put a delivered story in line to be pulled.

Source: design/conventions/strategic-agents.md § As the orchestrator
