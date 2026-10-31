---
id: TH-0774
title: "S-0547 plan: six tasks in three layers; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0547-the-close-out-runs-the-queue-tiers-and-the-tests-of-the-packages-that-import-a.md
  item: S-0547
status: open
participants: [planner-S-0547, orchestrator]
created: 2026-10-31T07:42:05Z
updated: 2026-10-31T07:46:20Z
---

# TH-0774 S-0547 plan: six tasks in three layers; cost of delay inputs needed

On wip/kanban/stories/S-0547-the-close-out-runs-the-queue-tiers-and-the-tests-of-the-packages-that-import-a.md.

## Entries

### 2026-10-31T07:42:05Z planner-S-0547
**Question: cost of delay inputs.** Neither S-0547 nor E-0031 has cost of delay inputs, so `flai cod S-0547` refuses and I wrote no value.

Recommended, on S-0547:

| Input | Value | Why |
|---|---|---|
| `time_lost_per_cycle` | 3h | about 7 send-backs a week (14 in the fortnight), each about 25 minutes: a restarted agent's fix, another close-out, a re-enqueue, and another 10-minute verification |
| `penalty_per_week` | none | |
| `revenue_per_week` | none | |

At `planning.hour_rate` 150 and a weekly cycle that is about 450 USD a week. Set it with `flai edit S-0547 --time-lost-per-cycle 3h`, or give inputs on E-0031 so that every open story of the epic takes a share, or answer with other figures, and I write the value.

**Plan.**

| Layer | Task | Criteria | Waits for |
|---|---|---|---|
| 1 | T-2375 ADR refining ADR-0174 and ADR-0221 | 2 | none |
| 1 | T-2376 `when` key: manifest, `flai check`, settings index | 2 | none |
| 1 | T-2377 go-test on `{importers}`, caller-break test | 1 | none |
| 2 | T-2378 verify and MCP verify run a queue tier whose `when` matches | 2 | T-2376 (`Tier.When`) |
| 3 | T-2379 `when` on integration and smoke, proving run, close I-0221 | 3 | T-2377 (`system-flow.yaml`), T-2378 |
| 3 | T-2380 convention, `flai-cli.md`, manifest design, user docs | 2, 4 | T-2375, T-2378 |

Forecast 4h, raised from `flai forecast`'s 18m; delivery 22:00 UTC, after S-0548 and S-0546 reach review. The story's `### Planning` says why, and where each touch came from.

**Assumptions.**

- Criterion 1's "`{packages}` expands to" means this repository's go-test tier takes `{importers}`, which flai has filled since S-0364 (the vet tier uses it). Changing what `{packages}` means would change every project's tiers and golangci-lint's.
- `when` matches the story's touches (its claim) or the paths its branch changes, so a story whose touches are narrower than its change still runs the tier.
- `when` is a `manifest.tests` error on a tier without `queue: true`.
- A queue tier the close-out ran at the head counts for ADR-0221's hand acceptance, and `flai verify --all` at the same patch reuses it.
- Recommended `when` for integration and smoke: `["flai/**", "template/**"]`. That runs both in every flai or template story's close-out, about twelve minutes until S-0548 shortens them, and covers go-test there. So go-test on `{importers}` matters mostly for `flai test` between tasks. If you want a shorter close-out, name a narrower set.
- I-0203 is closed already (S-0529). T-2379 adds a note naming this story's proving run and leaves it closed.
- `design/conventions/testing.md` does not exist. The tiers are in the project additions of `code-quality.md` and `tooling.md`. I kept the declared touch.

**Overlaps.** S-0548 (in progress) changes the same go-test line in `system-flow.yaml` and `flai/internal/check/tests_test.go`. S-0546 (ready, ahead in the pull order) shares `flai/cmd`, `flai-cli.md`, and `testing.md`. No `after`: neither change is needed on main.

**Defect.** As TH-0771 found, this run's start prompt named a thread "asked" that does not exist. The orchestrator passes the trigger `asked` as a thread name.

### 2026-10-31T07:46:20Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 3h` on S-0547, as the planner recommends:

```bash
flai edit S-0547 --time-lost-per-cycle 3h
```

## The figure

- **The goal is time lost.** About 7 send-backs a week, each about 25 minutes of a restarted agent's fix, another close-out, a re-enqueue, and another verification. That is agent time, not revenue or a penalty.
- **On the story, not E-0031.** The epic's stories each measure a different loss (TH-0772, TH-0773). Inputs on the epic would also stop flai taking inputs on any of its stories.

## Why this needs you

I tried to set it, and flai refused: "S-0547 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0547 to ready at 2026-10-31T07:24:13Z. A confirmation of this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.** T-2377 covers criterion 1, T-2376 and T-2378 criterion 2, T-2379 criterion 3, T-2375 and T-2380 criterion 4. The three layers follow the real dependencies.
- **Criterion 1 through `{importers}`: right.** Changing what `{packages}` means would change every project's tiers. The criterion's outcome, changed packages and their importers tested, is met.
- **`when: ["flai/**", "template/**"]` for integration and smoke: right.** Criterion 3 asks for exactly the flai and template paths.
- **Touches and forecast fit.** The 4h forecast over flai's 18m states its reason. No `after` is needed: the overlap hold on S-0548 and S-0552 already orders it.

The thread stays open for the cost of delay input only.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
