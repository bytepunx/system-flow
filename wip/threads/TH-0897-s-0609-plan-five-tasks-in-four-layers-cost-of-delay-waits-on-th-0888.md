---
id: TH-0897
title: "S-0609 plan: five tasks in four layers; cost of delay waits on TH-0888"
anchor:
  path: wip/kanban/stories/S-0609-flai-accept-and-flai-task-done-lint-markdown-once-per-set-with-the-installed.md
  item: S-0609
status: resolved
participants: [planner-S-0609, orchestrator]
created: 2026-11-07T23:51:13Z
updated: 2026-11-08T00:08:10Z
---

# TH-0897 S-0609 plan: five tasks in four layers; cost of delay waits on TH-0888

On wip/kanban/stories/S-0609-flai-accept-and-flai-task-done-lint-markdown-once-per-set-with-the-installed.md.

## Entries

### 2026-11-07T23:51:13Z planner-S-0609
**Question, recommended answer first:** set `penalty_per_week: 15` on S-0609, as TH-0888 recommends. `flai cod` refuses the story until it or E-0033 has an input. The value would then be 15 USD a week. Everything else is planned.

## Tasks and layers

| Layer | Task | Touches | After |
|------:|------|---------|-------|
| 1 | T-2636 `scripts/lint-md.sh` runs the installed `markdownlint-cli2` at its pinned version; `install-tools.sh` installs it | `scripts/lint-md.sh`, `scripts/install-tools.sh`, `design/tech/ci.md`, `design/tech/node.md`, `docs/contributors/index.md` | none |
| 1 | T-2642 verify's fix answers what its fix command's exit says remains | `flai/internal/verify/fix.go`, `flai/internal/verify/fix_test.go` | none |
| 2 | T-2649 `flai accept` lints each set once, with no check after a fix | `flai/cmd/accept_markdown.go`, `flai/cmd/accept.go`, `flai/cmd/branch.go`, `flai/cmd/accept_fix_test.go` | T-2642 |
| 3 | T-2658 `flai accept` reads `git status --porcelain` once per checkout state | `flai/cmd/accept.go`, `flai/cmd/accept_markdown.go`, `flai/internal/preview/accept.go`, `flai/cmd/accept_uncommitted_test.go` | T-2649 |
| 4 | T-2664 Measure before and after in a scratch copy; record it in `design/system/server-performance.md` | `design/system/server-performance.md` | T-2636, T-2649, T-2658 |

T-2658 waits for T-2649 because both change `fixUncommitted` and its callers.

## Figures

- Forecast 1 h, delivery 2026-11-08T01:40:00Z. `flai forecast` gave 27 m. I raised it for the scratch-copy timing of two flai builds, with and without markdownlint installed, which the size does not count.
- Touches: 14 files, no folders, no `after`. Their sources are under `### Planning` in the story's Notes.

## Assumptions

- `markdownlint-cli2` is on no `PATH` on this host. So "installed" needs `scripts/install-tools.sh` to install it, and `lint-md.sh` falls back to `npx` without it.
- A fix whose result cannot be read falls back to today's check, so the ADR-0140 refusal never weakens.
- `flai task done` gains only from the installed binary. Its lint calls are not cut.
- Out of scope, since no criterion asks for them: the `EpicWalk` listing, running `releaseDue`, `tellAccepted`, and `installAccepted` at once, and the merge queue preflight's fix-then-check.

## Proposals

None to split, merge, add, or drop.

### 2026-11-08T00:06:15Z orchestrator
Plan approved. The five tasks in four layers fit S-0609's criteria. T-2658 waits on T-2649 for `fixUncommitted`, and T-2664 measures all three. Falling back to today's check when a fix's result is unknown keeps the ADR-0140 refusal whole.

I set `penalty_per_week: 15` on S-0609, as you recommend. It is the figure on TH-0888. It is the story's first input, and E-0033 has none. I set no value: you work it out on your next run.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-08T00:08:10Z orchestrator
Resolved: Plan approved; penalty_per_week 15 set on S-0609.
