---
id: TH-0895
title: "Plan for S-0601: rules at once, a kept props read, then the measurement"
anchor:
  path: wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md
  item: S-0601
status: resolved
participants: [planner-S-0601, orchestrator]
created: 2026-11-07T23:50:24Z
updated: 2026-11-07T23:53:34Z
---

# TH-0895 Plan for S-0601: rules at once, a kept props read, then the measurement

On wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md.

## Entries

### 2026-11-07T23:50:24Z planner-S-0601
**Question: cost of delay input.** S-0601 has none, and neither has E-0033, so `flai cod S-0601` refuses.

Recommended: `penalty_per_week: 50` on S-0601, as TH-0888 proposes. Every `flai check`, `flai verify`, close-out check step, and `flai task done` pays 1.5 to 6 s of it. With that input the value is 50 USD a week.

Alternative: the input on E-0033 that TH-0888 offers. Answering TH-0888 answers this too.

#### Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2633 `check.Run` runs its independent rules at once and gives the same findings in the same order | `check.go`, `concurrent_test.go` | none |
| 1 | T-2643 Keep the props `node` reads under `.flai-cache` and start `node` only for a component whose key changed | `props.go`, `propscache.go`, `propscache_test.go`, `props.mjs` | none |
| 2 | T-2648 Measure a whole `flai check --strict` before and after, and record it in server-performance.md and component-library.md | `server-performance.md`, `component-library.md` | T-2633, T-2643 |

T-2633 and T-2643 share no file, so they run together. T-2648 measures what both change.

#### Figures

- Touches: 8 files, no folder. The sources are under `### Planning` in the story's Notes.
- Forecast: 40m, raised from `flai forecast`'s 19m. The reason is in the basis and the Notes.
- Cost of delay value: waits on the question above.

#### Assumptions

- `workItems` fills `byID`, which later rules read, so it runs first and alone. `vocab` is computed once before the rules start. Each rule writes its own `Result`, merged in today's rule order, so the stable sort gives the same order.
- The criterion keys the kept read by the component file's stat and the svelte version. T-2643 adds `scriptHash`, and the stamps of the files `props.mjs` read for the component's skeletons. Without the stamps, a change to an imported type would keep a stale skeleton. That is why `props.mjs` is touched. Say if you want the key to be exactly the criterion's.
- The fixture for criterion 1 is built in a temporary folder by the test, not under `testdata/`, so no task needs a folder touch.
- No `after` on S-0598 or S-0600. They also change `check.go` or its callers, but S-0601 reads nothing they write.

S-0601 has no draft flag any more, so it was finalized while I planned. I changed only its touches, its forecast, and the `### Planning` heading of its Notes.

### 2026-11-07T23:52:49Z orchestrator
Plan approved. T-2633 and T-2643 share no file, so they run together, and T-2648 measures both. The wider key is fine: `scriptHash` and the stamps of the files `props.mjs` reads. It starts `node` only when a key changes, as criterion 2 asks, and it keeps a stale skeleton out.

I set `penalty_per_week: 50` on S-0601, as you recommend. It is the figure on TH-0888. It is the story's first input, and E-0033 has none. I set no value: you work it out on your next run.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T23:53:34Z orchestrator
Resolved: Plan approved; penalty_per_week 50 set on S-0601.
