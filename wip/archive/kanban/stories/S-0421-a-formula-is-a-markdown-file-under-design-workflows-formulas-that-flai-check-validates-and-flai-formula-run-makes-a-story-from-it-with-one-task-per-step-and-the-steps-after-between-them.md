---
id: S-0421
type: story
nature: feature
title: A formula is a markdown file under design/workflows/formulas that flai check validates, and flai formula run makes a story from it with one task per step and the steps' after between them
status: done
parent: E-0024
owner: alex
created: 2026-10-16T05:40:40Z
updated: 2026-11-02T03:02:37Z
transitions:
  - to: ready
    at: 2026-10-28T13:10:51Z
    by: alex
  - to: backlog
    at: 2026-10-28T22:20:42Z
    by: alex
  - to: ready
    at: 2026-11-01T23:43:52Z
    by: alex
  - to: in-progress
    at: 2026-11-02T02:07:42Z
    by: agent-S-0421
  - to: review
    at: 2026-11-02T02:53:39Z
    by: system-flow
  - to: done
    at: 2026-11-02T03:02:37Z
    by: orchestrator
tags: [cli, template]
topics: [cli, orchestration, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/workflows/formulas/README.md, design/workflows/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/system/work-hierarchy.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/root.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/workitem/fields_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/item.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/formulas/README.md, template/root/design/workflows/README.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md, flai/internal/guard/formula_test.go, flai/internal/guard/guard.go, flai/internal/workitem/create.go, CLAUDE.md, design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md, design/issues/summary.md, design/issues/I-0249-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2776
  priced_by: harness
  turns:
    - day: 2026-11-02
      ceremony: 5
      hand_edits: 2
      work: 58
  models:
    - model: claude-opus-5-5
      input: 476
      output: 178214
      cache_read: 26790079
      cache_write: 839113
      cost: 13.9303
  strategic:
    - kind: orchestrator
      seconds: 485
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 143
          output: 2235
          cache_read: 29278286
          cache_write: 214297
          cost: 6.4966
        - model: claude-sonnet-5-5
          input: 12
          output: 72
          cache_read: 172697
          cache_write: 55761
          cost: 0.1592
cost_of_delay:
  value: 48.09
  by: planner-S-0421
  at: 2026-10-21T12:39:49Z
forecast:
  duration: 1h3m
  delivery: 2026-11-02T03:14:00Z
  basis: "Its own forecast of 1h3m; 2nd in the pull order with an in-progress limit of 5, behind S-0420, S-0553 and S-0423."
  by: flai
  at: 2026-11-02T01:33:08Z
finalized:
  by: alex
  at: 2026-10-28T13:10:45Z
---
# S-0421 A formula is a markdown file under design/workflows/formulas that flai check validates, and flai formula run makes a story from it with one task per step and the steps' after between them

## Goal

A repeatable process is a file, not prose. A formula under `design/workflows/formulas/` is a named sequence of steps with `after` between them, each a flai command, a script under `scripts/`, or a role run with a prompt and inputs, and the inputs the formula takes. `flai formula run <name>` makes a story from it, its steps as tasks in the same layers, so that it is planned, worked, measured, and archived as any story is: a story's agent works the instance on its branch and worktree, as TH-0471 settled. `flai check` refuses a formula that could not run.

## Acceptance criteria

- [x] An ADR records the formula file: where it lives (`design/workflows/formulas/<name>.md`, in the `design/workflows/` folder TH-0471 settled and S-0414's ADR adds to the layout), its front matter (name, description, nature, inputs, and steps with `id`, `title`, `run` as a flai command, a script, or a role with its prompt, and `after`), its body as the formula's description, the baseline marker, the instance's nature (the formula's, `improvement` by default, no new nature) and its `formula: <name>` field, and that a story's agent works the instance on its branch as any story, as TH-0471 settled them.
- [x] `flai/internal/formulas` reads a project's formulas, and `flai formula list [--json]` and `flai formula show <name> [--json]` print each with its steps and layers.
- [x] `flai formula run <name> [--input key=value]... [--epic E-nnnn] [--json]` makes a story with the formula's nature and `formula: <name>`, one task per step with `## Work` naming what the step runs and `## Done when` its exit, the steps' `after` as the tasks' `after`, through `itemnew` so that `flai check --strict` passes on it, and prints the story's ID. A missing input, or one the formula does not declare, is refused before anything is written.
- [x] `flai check` reports `formulas.file` for a formula that does not parse or misses a field, a step whose `run` is none of the three kinds, a script not under `scripts/`, a role no file defines, an `after` naming no step, a cycle, and an input a step uses that the formula does not declare.
- [x] `formula` is listed in `front-matter-fields.txt` and in `design/system/work-hierarchy.md`'s schema, and an older flai reads past it.
- [x] `design/system/repository-layout.md`, a new section of `design/system/strategic-agents.md` or a design document of its own, `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md`, the convention `strategic-agents.md`, and `template/CHANGELOG.md` with the template's `design/workflows/formulas/README.md` say what a formula holds and how to run one.

## Tasks
- T-1768 An ADR records the formula file: where it lives, its steps and inputs, its body, its marker, and the instance's nature and formula field
- T-1769 flai/internal/formulas reads and validates a project's formulas into steps, inputs, and layers
- T-1770 A story carries formula in its front matter, listed in front-matter-fields.txt and the work hierarchy's schema, and an older flai reads past it
- T-1771 flai check reports formulas.file on each formula that could not run
- T-1772 flai formula list, show, and run: run makes a story with the formula's nature and one task per step, through itemnew, and prints its ID
- T-1773 The template and this repository have design/workflows/formulas with its README, and the design, the convention, and the docs say what a formula holds and how to run one

## Notes

- Drafted by the planner for E-0024. Follows the epic's order: after the dashboard's role page, though it reads only the role file story's loader for the role steps it validates.
- This story makes the instance; the next one works it.

### Planning

- **Plan thread.** TH-0581, replanned 2026-10-21 on TH-0471's answers as S-0414 and S-0420 were, approved and resolved by the orchestrator.
  - Answer 2 moves the folder from `design/formulas/` to `design/workflows/formulas/`. The title, the goal, and criteria 1 and 6 now say so, and the two README touches moved with it. `design/workflows/` itself is S-0414's ADR's to add to the layout; its `README.md`, in both copies, is on T-1773, which adds the folder's README when S-0414 left none and a line for formulas otherwise.
  - Answer 3 stands: the instance has the formula's nature, `improvement` by default, and `formula: <name>`; no new nature. Criterion 1 now says so.
  - Answer 5 takes the alternative: a story's agent works the instance on its branch and worktree, as any story. The goal and criterion 1 say so. Nothing in this story changes for it, since `flai formula run` makes an ordinary story. S-0422's title and goal, which have flai serve work the instance, no longer fit; its own planner run takes that, as the orchestrator said on TH-0581.
  - Answer 1's input, `time_lost_per_cycle: 4h` on E-0024, is not yet on the epic (TH-0567, question 3).
  - The tasks' words predate the answers. T-1768 and T-1772 read the answers through TH-0471 and need no rewrite. T-1773's title and first Work bullet were rewritten to `design/workflows/formulas/`, as approved on TH-0581. Layers and `after` are unchanged.
- **Touches.** 29, all files but `design/adrs`. The story declared none of its own at first; every touch comes from its six tasks.
  - **Declared, kept:** 24 of the 26 earlier touches.
  - **Declared, moved:** `design/formulas/README.md` and the template's, now under `design/workflows/formulas/` (answer 2), on T-1773.
  - **From layout, added:** `design/workflows/README.md` and `template/root/design/workflows/README.md` on T-1773. Every folder a reader might land in has a README, and S-0414's tasks name only the `roles/` one.
  - **Co-change, added:** `docs/operators/settings.md` (16% of the seed commits) on T-1772. Its flag table is generated by `make flai-reference` and gains `flai formula run`'s flags.
  - **Co-change, not added:** `flaiover-dashboard.md`, `docs/users/flaiover.md`, `workflow.md`, `work-management.md`, `flai/internal/harness`, `flai/internal/mcpserver`, and `flai/internal/hostapi`. No page, workflow rule, agent start, or host method changes here.
  - **Checked, not added:** `flai/internal/itemnew`. Its `Options.New` is `workitem.NewOptions`, so the `formula` field reaches creation through `item.go` (T-1770).
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
- **Forecast.** flai's 1h3m is kept: 107 s a unit over 85 done feature stories on this model in the large band, times size 35 (6 criteria, 29 touches). The three added touches are one-line changes that overstate the size a little; the command task, three subcommands and an instance that must pass `flai check --strict`, is heavier than its five touches say. The two about cancel. The delivery is flai's play-out, 16th in the pull order: 2026-10-22T01:33Z.
- **Cost of delay.** 48.09 USD a week, worked out by hand from TH-0471's answer 1, as the planners of S-0414 to S-0420 did, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0421's share is its 1h3m of the 13h6m forecast for E-0024's twelve open stories: 63 / 786 × 600.
  - Once the input is on E-0024, `flai cod S-0421` gives the same figure while the forecasts stand.

### Accepted by the orchestrator

- Verified: cc6ac24be5eff4ba9a6134bb0d47115665974936
- Batch: 20261102-025823
- Enqueued verified: cc6ac24be5eff4ba9a6134bb0d47115665974936
- At: 2026-11-02T03:02:37Z

Verdict: meets all criteria (verified at commit cc6ac24be5eff4ba9a6134bb0d47115665974936)
- 1: design/adrs/0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md, design/adrs/README.md
- 2: flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go
- 3: flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/root.go, flai/internal/workitem/create.go, flai/internal/guard/guard.go, flai/internal/guard/formula_test.go
- 4: flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go
- 5: flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/item.go, flai/internal/workitem/fields_test.go, design/system/work-hierarchy.md
- 6: design/system/repository-layout.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/root/design/workflows/formulas/README.md, template/root/design/workflows/README.md, design/workflows/formulas/README.md, design/workflows/README.md
