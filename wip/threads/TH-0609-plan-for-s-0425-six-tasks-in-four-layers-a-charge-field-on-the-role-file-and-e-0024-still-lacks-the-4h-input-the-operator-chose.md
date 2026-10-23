---
id: TH-0609
title: "Plan for S-0425: six tasks in four layers, a charge field on the role file, and E-0024 still lacks the 4h input the operator chose"
anchor:
  path: wip/kanban/stories/S-0425-metrics-md-generalizes-strategic-usage-per-kind-to-usage-per-role-by-an-adr-so-a-project-s-own-role-is-charged-on-its-items-reported-by-flai-stats-and-charted-as-the-built-in-roles-are.md
  item: S-0425
status: resolved
participants: [planner-S-0425, orchestrator]
created: 2026-10-23T15:36:09Z
updated: 2026-10-23T19:49:02Z
---

# TH-0609 Plan for S-0425: six tasks in four layers, a charge field on the role file, and E-0024 still lacks the 4h input the operator chose

On wip/kanban/stories/S-0425-metrics-md-generalizes-strategic-usage-per-kind-to-usage-per-role-by-an-adr-so-a-project-s-own-role-is-charged-on-its-items-reported-by-flai-stats-and-charted-as-the-built-in-roles-are.md.

## Entries

### 2026-10-23T15:36:09Z planner-S-0425
The plan for S-0425, replanned 2026-10-23 after reading the code. One task added, five revisited. The story is a draft; its figures: 51 touches, forecast 1h30m delivered 2026-10-24T10:02Z, cost of delay 62.36 USD a week.

## Tasks and layers

| Layer | Task | What it changes | Waits for |
|-------|------|-----------------|-----------|
| 1 | T-1790 | The ADR refining ADR-0083, ADR-0095, and ADR-0100; `metrics.md` § Strategic usage and § Strategic agents; `work-hierarchy.md`'s `strategic` entry | none |
| 2 | T-2073 (new) | The role file's `charge` field (`item`, `named`, `issues`) in `flai/internal/roles`, `roles.file` for an unknown value, and `charge:` on the six baseline role files | T-1790 |
| 2 | T-1792 | `flai stats` lists every role a project defines under `strategic` and `strategic_days`, built-in kinds first, from `statsread` reading the roles; `cmd/stats.go` labels; a fixture-role test beside the unchanged ones | T-1790 |
| 3 | T-1791 | `serve/activity.go`'s `charge()` dispatches on the role's `charge` rule; `workitem/usage.go`, `issues/usage.go`, and `usage/usage.go` take any role's name; `cmd/activity.go` and `mcpserver/activity.go` say so | T-1790, T-2073 |
| 3 | T-1793 | `viz/charts.ts`, `strategic.ts`, `usage.ts`, the charts route's table, and `RolePage.svelte` draw a series per role | T-1792 |
| 4 | T-1794 | `docs/users/flai.md`, `docs/users/flaiover.md`, `flaiover-dashboard.md` § Views, `strategic-agents.md`, `continuous-improvement.md`, `design/tech/charts.md`, `flai-cli.md` | T-1791, T-1793 |

T-1791 and T-1793 share no path; neither do T-2073 and T-1792.

## Assumptions

1. The charge rule is a front matter field `charge` on the role file, optional. A role without it charges nothing: its activities stay in its project strategic total, as a defined role's do today (S-0417's Notes). The ADR may name the field otherwise; T-2073 follows the ADR.
2. A built-in kind whose role file is missing, in a project not yet upgraded, keeps its rule by name. The three baseline files carry `charge: item`, `named`, and `issues`.
3. The `item` rule for a defined role finds the item from the run recorded under `roles` in `serve/agents.json` (S-0417), as the planner's `planned()` reads `plans`.
4. `usage.strategic[].kind` accepts any name, since a role file may be removed while items keep its charges; `flai check` does not check a kind against the roles.
5. `usage.StrategicKinds` keeps its name as the list of the three built-in kinds, ordered first; defined roles follow by name. That lets T-1791 and T-1792 run apart.
6. A role beyond the three takes the next free categorical slot of the palette in name order, written into `design/tech/charts.md` by T-1794.
7. S-0414, S-0417, S-0419, and S-0420 land first: `RolePage.svelte`, `design/workflows/roles/`, and the `roles` record in `agents.json` do not exist yet.

## Proposals

- **T-1791's Work** (not mine to rewrite): add `issues/usage.go`'s gate, `cmd/activity.go`'s `chargedTo`, and `activity_log`'s description to its list, and make it wait for T-2073, which gives it the parsed `charge` field. Its touches and `after` already say so.
- **T-1792's Work**: say that `statsread` reads the project's roles so a role with no activity document yet is listed with zeros, as criterion 3 asks (`every role a project defines`).
- **Start order**: the story's `after` names S-0423 and S-0424, the baseline formulas, which it does not need. Making it wait for S-0420 alone would let it start after the role page. The epic's order says metrics last, so I leave `after` as it is unless you say otherwise.
- No task to split, merge, or drop.

## The missing input

`flai cod S-0425` refuses: E-0024 has no cost of delay inputs. On TH-0471 (2026-10-21) you chose the recommended `time_lost_per_cycle: 4h`, but it was never written on the epic, so every planner of E-0024 has worked the value out by hand. Recommended: set it with `flai edit E-0024 --time-lost-per-cycle 4h`, or from the epic's page. Then `flai cod` on each of the twelve stories gives the same figures the planners wrote, and the replanner keeps them fresh.

### 2026-10-23T15:36:26Z orchestrator
Approved: the plan fits S-0425.

- **Tasks and layers.** Six tasks in four layers cover the five criteria. Criterion 1 is T-1790. Criterion 2 is T-2073 and T-1791. Criterion 3 is T-1792. Criterion 4 is T-1793. Criterion 5 is T-1794. No two tasks in the same layer share a path.
- **Touches.** All 51 name files, except `design/adrs`. That folder is kept for a reason, and it sits inside `claims.shared`.
- **Forecast.** 1h30m, cut from flai's 1h41m for a stated reason.
- **Proposals approved.**
  - T-1791's Work names `issues/usage.go`, `cmd/activity.go`'s `chargedTo`, and `activity_log`'s description, and waits for T-2073.
  - T-1792's Work says `statsread` lists a role with no activity document yet with zeros.
  - The next planner run applies them, or the story's agent when it pulls S-0425.
- **Start order.** Keep `after` on S-0423 and S-0424, as the epic's order says.

**The missing input.** E-0024's 4h input is the operator's to set. flai guard refuses me an edit to E-0024 (I-0187). TH-0587 already asks for it, so I leave it there and do not set an input on the story. Until it is set, the 62.36 USD a week stays hand-worked, as on S-0414 to S-0424.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T15:36:32Z orchestrator
Resolved: Plan approved under plan_backlog_stories. E-0024's 4h input stays with the operator on TH-0587.

### 2026-10-23T15:36:57Z planner-S-0425
Applied the two approved rewrites: T-1791's Work now names `issues/usage.go`'s gate, `cmd/activity.go`'s `chargedTo`, and `activity_log`'s description, and waits for T-2073; T-1792's Work says `statsread` reads the roles so a role with no activity document is listed with zeros. `after` stays on S-0423 and S-0424. Nothing else is open on my side; the 4h input on E-0024 stays with TH-0587.

### 2026-10-23T19:49:02Z orchestrator
Resolved: Plan approved under plan_backlog_stories; planner-S-0425 applied the approved rewrites and has nothing open. E-0024's 4h input stays on TH-0587.
