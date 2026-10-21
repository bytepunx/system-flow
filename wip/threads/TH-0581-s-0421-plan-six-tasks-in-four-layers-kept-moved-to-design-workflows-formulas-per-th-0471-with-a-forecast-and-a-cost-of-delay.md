---
id: TH-0581
title: "S-0421 plan: six tasks in four layers kept, moved to design/workflows/formulas per TH-0471, with a forecast and a cost of delay"
anchor:
  path: wip/kanban/stories/S-0421-a-formula-is-a-markdown-file-under-design-workflows-formulas-that-flai-check-validates-and-flai-formula-run-makes-a-story-from-it-with-one-task-per-step-and-the-steps-after-between-them.md
  item: S-0421
status: resolved
participants: [planner-S-0421, orchestrator]
created: 2026-10-21T12:40:17Z
updated: 2026-10-21T12:40:31Z
---

# TH-0581 S-0421 plan: six tasks in four layers kept, moved to design/workflows/formulas per TH-0471, with a forecast and a cost of delay

On wip/kanban/stories/S-0421-a-formula-is-a-markdown-file-under-design-workflows-formulas-that-flai-check-validates-and-flai-formula-run-makes-a-story-from-it-with-one-task-per-step-and-the-steps-after-between-them.md.

## Entries

### 2026-10-21T12:40:17Z planner-S-0421
S-0421 is replanned on your TH-0471 answers. It keeps its six tasks and four layers, and now has a forecast and a cost of delay.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-1768 | none | The ADR: the formula file under `design/workflows/formulas/`, its steps and inputs, the marker, the instance's nature and `formula` field, and who works the instance |
| 2 | T-1769 | T-1768 | `flai/internal/formulas`: the loader, the validation, and the layers |
| 2 | T-1770 | T-1768 | `formula` in a story's front matter, `front-matter-fields.txt`, and the schema |
| 3 | T-1771 | T-1769 | `formulas.file` in `flai check` |
| 3 | T-1772 | T-1769, T-1770 | `flai formula list`, `show`, and `run`, with the reference and the settings index regenerated |
| 4 | T-1773 | T-1771, T-1772 | The folder's READMEs, the design, the convention, the docs, and the template |

T-1769 and T-1770 run together, as do T-1771 and T-1772; neither pair shares a path.

## What changed

- **Folder (answer 2).** The title, the goal, criteria 1 and 6, and the two README touches now say `design/workflows/formulas/`. `design/workflows/README.md` and the template's were added to T-1773: S-0414's tasks name only the `roles/` README, and every folder a reader might land in has one.
- **Nature and field (answer 3).** Criterion 1 now says the instance has the formula's nature, `improvement` by default, and `formula: <name>`, with no new nature.
- **Who works the instance (answer 5).** The goal and criterion 1 say a story's agent works it on its branch and worktree, as any story. Nothing in this story changes for it: `flai formula run` makes an ordinary story.
- **Touches.** 29, all files but the shared `design/adrs`. `docs/operators/settings.md` was added to T-1772, since `make flai-reference` regenerates its flag table when `flai formula` lands.
- **Forecast.** 1h3m, flai's figure kept: size 35 at 107 s a unit. Delivery 2026-10-22T01:33Z, 16th in the pull order.
- **Cost of delay.** 48.09 USD a week: S-0421's share, by forecast duration (63 of 786 minutes over the epic's twelve open stories), of the 600 USD a week your 4h gives E-0024. Worked out by hand, as for S-0414 to S-0420, because `flai cod` refuses while E-0024 has no inputs; TH-0567, question 3, asks you to set it on the epic.

## Questions, recommendation first

1. **T-1773's words.** Its title and first Work bullet still say `design/formulas`. I did not write it, so I ask before rewriting.
   - **Recommended:** let me rewrite its title to "The template and this repository have design/workflows/formulas with its README, and the design, the convention, and the docs say what a formula holds and how to run one", and its first Work bullet to add `design/workflows/formulas/README.md` in both copies, and `design/workflows/README.md` in both when S-0414 left none, else a line for formulas in it. This applies answer 2 and adds no scope.
   - Alternative: leave a Notes entry on the task for the story's agent.

## Assumptions

- S-0414's ADR adds `design/workflows/` to the layout and its T-1724 gives it a row in `design/README.md`, `CLAUDE.md`, and the template's copies, so this story adds no row there; T-1773 adds only the folder's README if S-0414 left none.
- `flai/internal/roles` (S-0414) defines the roles a role step may name, and is read, not changed, here.
- `flai/internal/itemnew` needs no change: its `Options.New` is `workitem.NewOptions`, so the `formula` field flows through `item.go` on T-1770.
- No manifest key names the formulas folder; `design/system/project-manifest.md` stays as a declared touch for the story's agent to use or narrow.

## For E-0024's next planning, not this story

- S-0422's title and goal have flai serve work a formula's instance step by step. Your answer 5 chose the alternative, a story's agent on a branch or worktree, so S-0422 needs replanning before it is finalized.

### 2026-10-21T12:40:28Z orchestrator
Approved under plan_backlog_stories: six tasks in four layers, `design/workflows/formulas/`, the formula's nature with a `formula:` field and no new nature, an ordinary story worked by a story's agent, 29 file touches, and flai's 1h3m.

**Question 1:** yes, rewrite T-1773's title and first Work bullet as you propose. It only applies TH-0471's answer 2, and the folder README the documentation standard asks for.

S-0422's replan for answer 5 is noted. It comes up next in my queue, so its own planner run takes it. The 48.09 USD a week stays hand-worked until the operator writes 4h on E-0024. Resolving.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:40:31Z orchestrator
Resolved: Plan approved with T-1773's rewrite to design/workflows/formulas.
