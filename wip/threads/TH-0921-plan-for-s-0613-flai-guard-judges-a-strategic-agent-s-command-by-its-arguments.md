---
id: TH-0921
title: "Plan for S-0613: flai guard judges a strategic agent's command by its arguments, redirections left out"
anchor:
  path: wip/kanban/stories/S-0613-flai-guard-refuses-the-orchestrator-a-read-or-a-permitted-flai-command-when-the.md
  item: S-0613
status: resolved
participants: [planner-S-0613, orchestrator]
created: 2026-11-09T08:14:32Z
updated: 2026-11-09T08:17:27Z
---

# TH-0921 Plan for S-0613: flai guard judges a strategic agent's command by its arguments, redirections left out

On wip/kanban/stories/S-0613-flai-guard-refuses-the-orchestrator-a-read-or-a-permitted-flai-command-when-the.md.

## Entries

### 2026-11-09T08:14:32Z planner-S-0613
S-0613 is planned: three tasks in two layers, touches on eight files, forecast 40m, cost of delay value 150 USD a week.

## Proposed solution

The cause is in `flai/internal/guard/guard.go`. `commands` keeps a redirection such as `2>&1` as one of a command's words. `Guard.refuse` judges those words for the orchestrator and for every role a role file holds (`Guard.asRole`). So `ordersByPolicy` counts `2>&1` as a story and takes `flai order --by throughput --apply 2>&1` for a placement by hand. `subcommands` reads `2>/dev/null` as the subcommand of `flai board --json 2>/dev/null`. `parse` already splits each command's arguments with redirections and their targets left out (S-0368), and the fix judges those.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2706 the orchestrator's and every role file's checks judge a command by its arguments, with the tests that reproduce I-0230's instances | `flai/internal/guard/guard.go`, `roles.go`, `orchestrate_test.go`, `roles_test.go` | none |
| 2 | T-2707 the design (§ Its permissions and the guard) and the user guide (§ Role files) say it | `design/system/strategic-agents.md`, `docs/users/flai.md` | T-2706 |
| 2 | T-2708 close I-0230 with `flai issue close --reason` | I-0230's file, `design/issues/summary.md` | T-2706 |

T-2707 and T-2708 share no path, so they run together.

## Assumptions

- The fix covers every role a role file holds (planner, analyzer, curator, and any a project defines), not only the orchestrator. They share `Guard.refuse`, and the planner would be refused the same way.
- The write check for read-only sub-agents keeps reading redirection targets. A redirection to a file is still a write there.
- `authors.go`, `shared.go`, and `hostconfig.go` keep the words. Their checks read `--by`, `FLAI_AGENT`, and `flai serve`'s subcommands, which a trailing redirection does not change. The story's agent may widen the fix to them if a test shows otherwise.
- No flag or help text changes, so `flai/cmd/guard.go` and `docs/users/flai-reference.md` are not touched.
- The forecast is raised from flai's 18m to 40m. Three comparable guard remediations (S-0368, S-0568, S-0556) took 31 to 38 minutes.
- The cost of delay value is flai's 150 USD a week, from the operator's input of 1h lost per cycle. It is unchanged.
- The story names no `after`.

I propose no split, merge, or drop.

### 2026-11-09T08:16:26Z orchestrator
Approved: the plan fits S-0613.

- **Tasks fit the criteria.** Criterion 1 (the cause no longer occurs, with a reproducing test) is T-2706. Criterion 2 (I-0230 closed with a reason) is T-2708. T-2707 carries the design and user guide, as `documentation.md` asks.
- **The layers follow the dependency.** T-2707 and T-2708 both wait for T-2706, share no path, and run together.
- **The touches are files, not folders.** They cover the code the cause lives in (`commands`, `refuse`, `orchestrate`, `asRole`, `ordersByPolicy`, `subcommands` in `flai/internal/guard/guard.go` and `roles.go`), its tests, the two docs, and the issue files.
- **The figures fit.** The 40m forecast over flai's 18m gives its basis (S-0368, S-0568, S-0556). The value of 150 USD a week is flai's, from the input of 1h lost per cycle already on the story.
- **The assumptions are sound.** The fix covers every role through the shared `Guard.refuse`. The read-only sub-agents' write check keeps reading redirection targets.

I am finalizing S-0613 under `finalize_drafts`, since `flai promote --drafts` lists it as complete.

Source: design/conventions/work-management.md

### 2026-11-09T08:17:27Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, forecast 40m, and value 150 fit S-0613.
