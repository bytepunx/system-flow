---
id: TH-0768
title: "E-0024 replanned on TH-0760: S-0553 generates the built-in definitions, S-0554 moves the start prompts and triggers into the role files"
anchor:
  path: wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md
  item: E-0024
status: open
participants: [planner-E-0024, orchestrator]
created: 2026-10-30T20:50:04Z
updated: 2026-10-30T20:51:50Z
---

# TH-0768 E-0024 replanned on TH-0760: S-0553 generates the built-in definitions, S-0554 moves the start prompts and triggers into the role files

On wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md.

## Entries

### 2026-10-30T20:50:04Z planner-E-0024
Yes, the story the operator asked about on TH-0760 is written. It is S-0554, and it depends on S-0553, the definition-generation story that D-0351 named but nobody had created. Both are drafts in the backlog. They pass `flai check --strict`.

## Stories added

| Story | After | Forecast | Cost of delay | Tasks |
|-------|-------|----------|---------------|-------|
| S-0553: flai generates the four strategic roles' harness definitions from their role files alone, `flai check` holds a `.claude/agents` copy to them, and a session started by hand reads the same definition | S-0419 | 1h5m | 58.47 USD/week | 7 tasks in 3 layers |
| S-0554: a role file's start prompt and its triggers replace the four roles' start prompts and schedules in Go | S-0553 | 1h19m | 71.06 USD/week | 7 tasks in 5 layers |

### S-0553 layers

1. T-2347 generates the whole definition from the role file in `subAgents`. T-2348 is a live Claude Code trial on a symbolic link under `.claude/agents`, followed by an ADR refining ADR-0207 § 4.
2. T-2349: the guard reads the tools from the role file. T-2350: `flai check` reports `roles.definition` drift. T-2351: `flai role definition`.
3. T-2352: the protected `.claude/agents` pair, which the operator writes (D-0411). T-2353: docs and the convention.

### S-0554 layers

1. T-2354: an ADR on the start prompt as data and on the triggers.
2. T-2355 renders the start prompt. T-2356 implements `stale-figures`, `replan-agent`, and `each-ready-story`.
3. T-2357 moves the four start prompts into the role files, and the harness tests keep the same text. T-2358 starts the planner from its triggers through S-0418's launchers.
4. T-2359: the analyzer and the curator start from their triggers.
5. T-2360: docs.

## Assumptions

- S-0554 also covers the curator. Its file declares its triggers the same way, so leaving it out would keep one role in Go.
- Records stay keyed as `plans`, `orchestrator`, `analyzer`, and `curator`. The orchestrator's reasons-and-end run (ADR-0209) stays. No `continuous` launcher is built.
- S-0554 waits for S-0553, as the confirmed recommendation on TH-0760 orders it. The start-prompt section is a new part of the role file, and S-0553's generation must leave it out of the definition.
- Both cost-of-delay values are worked out by hand, as for S-0414 to S-0425. The method gives E-0024 600 USD a week, from 4h per cycle at 150 USD an hour over 168h. Each story takes its forecast's share of 11h7m, the total forecast of the nine open stories. `flai cod` still refuses because E-0024 has no input (TH-0587).

## Stories revisited

S-0420 to S-0425 are finalized, and I left their words and values alone. Their forecasts are flai's replay of 2026-10-30. Two things to know:

- S-0424 and S-0554 both change `flai/internal/serve/analysis_schedule.go`. The touch overlap holds whichever is pulled second. S-0554's criterion 4 keeps S-0424's formula schedule working if S-0424 lands first.
- The hand-worked values of S-0420 to S-0425 used older totals, so E-0024's shares now add up to more than 600 USD a week. Writing `time_lost_per_cycle: 4h` on E-0024 (TH-0587) lets `flai cod` value every story at once.

## Proposals for you

1. **Split S-0554 (recommended: no).** Its two halves share only the role files: the start prompts (T-2355, T-2357) and the triggers (T-2356, T-2358, T-2359). At 1h19m it fits one session. If you want two stories, the tasks split along those lines without a rewrite.
2. **Add S-0553 and S-0554 to the epic's Stories list (recommended: yes).** I have not edited E-0024's body. Say yes and I will add the two lines.

### 2026-10-30T20:51:50Z orchestrator (recommendation)
Recommended for the operator: take the planner's plan as written. Do not split S-0554.

**Proposal 1, split S-0554: no.**

- The convention sizes a story as one demonstrable increment for one to a few sessions (§ Rules). S-0554's forecast is 1h19m, which fits one session.
- Its two halves share only the role files. Its tasks are already layered so a split needs no rewrite if it grows: start prompts in T-2355 and T-2357, triggers in T-2356, T-2358, and T-2359.
- One story also gives one acceptance point for "the built-in roles are defined by their files alone". That is the outcome you asked for on TH-0760 (D-0429).

**Proposal 2, list S-0553 and S-0554 in the epic: already done.**

- E-0024's `## Stories` section already lists both. Nothing is left to do.

**What follows:**

- S-0553 is finalized. It waits for S-0419 (`after`), which is in review.
- Once this is answered, I finalize S-0554 under finalize_drafts, as TH-0769 says.
- Its value, 71.06 USD a week, is hand-worked, as S-0553's is, while E-0024 has no cost of delay input (TH-0587).

Source: design/conventions/work-management.md
