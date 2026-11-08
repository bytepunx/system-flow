---
id: TH-0905
title: "S-0609 replan: cost of delay 15 USD a week set; tasks unchanged"
anchor:
  path: wip/kanban/stories/S-0609-flai-accept-and-flai-task-done-lint-markdown-once-per-set-with-the-installed.md
  item: S-0609
status: resolved
participants: [planner-S-0609, orchestrator]
created: 2026-11-08T00:11:14Z
updated: 2026-11-08T00:23:27Z
---

# TH-0905 S-0609 replan: cost of delay 15 USD a week set; tasks unchanged

On wip/kanban/stories/S-0609-flai-accept-and-flai-task-done-lint-markdown-once-per-set-with-the-installed.md.

## Entries

### 2026-11-08T00:11:14Z planner-S-0609
The plan approved on TH-0897 stands. This run set the cost of delay value and moved the delivery. Nothing awaits an answer.

## Changes

- Cost of delay value: 15 USD a week, from the input `penalty_per_week: 15` set on TH-0897. `flai cod` gives 15. Kept unadjusted.
- Forecast: 1 h, unchanged. Delivery moved from 01:40Z to 2026-11-08T03:29Z. `flai forecast` now puts the story 9th in the pull order, at 02:56Z for 27 m, plus the same 33 m raise as before.
- `### Planning` in the Notes now records both figures.

## Tasks and layers, revisited, unchanged

| Layer | Task | After |
|------:|------|-------|
| 1 | T-2636 `scripts/lint-md.sh` runs the installed `markdownlint-cli2`; `install-tools.sh` installs it | none |
| 1 | T-2642 verify's fix answers what its fix command's exit says remains | none |
| 2 | T-2649 `flai accept` lints each set once, with no check after a fix | T-2642 |
| 3 | T-2658 `flai accept` reads `git status --porcelain` once per checkout state | T-2649 |
| 4 | T-2664 Measure before and after in a scratch copy; record it in `design/system/server-performance.md` | T-2636, T-2649, T-2658 |

Every task's touches name files. The 14 story touches are unchanged, with no folders and no `after`.

## Assumptions

These are the same as on TH-0897:

- `lint-md.sh` falls back to `npx` when `markdownlint-cli2` is not installed.
- When a fix's result is unknown, the code falls back to today's check, so the ADR-0140 refusal stays whole.
- `flai task done` gains only from the installed binary.

## Proposals

None to split, merge, add, or drop.

### 2026-11-08T00:22:21Z orchestrator
Plan approved. It is unchanged from TH-0897. The five tasks in four layers fit S-0609's three criteria, and the value of 15 is what `flai cod` gives from the input I set from TH-0888. I finalized S-0609: `flai promote --drafts` lists it as complete, and its criteria, touches, forecast, and value agree.

One point for the story's agent, from TH-0903 proposal 4. ADR-0140 § Decision says `flai accept` runs each fixed tier's `command` once the fix has run. Criterion 1 drops that check when the fix's exit already says what remains. Under `decisions.md`, that change to a tool-enforced rule needs an ADR refining ADR-0140, written in T-2642 or T-2649.

Source: design/adrs/0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md § Decision

Source: design/adrs/0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md § Decision

### 2026-11-08T00:23:27Z orchestrator
Resolved: Plan approved, unchanged from TH-0897. Value 15 from flai cod. S-0609 finalized. The ADR refining ADR-0140 is noted for the story's agent.
