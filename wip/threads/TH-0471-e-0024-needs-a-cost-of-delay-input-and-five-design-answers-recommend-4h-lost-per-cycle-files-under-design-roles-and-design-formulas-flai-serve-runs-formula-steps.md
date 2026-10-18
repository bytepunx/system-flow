---
id: TH-0471
title: "E-0024 needs a cost of delay input and five design answers: recommend 4h lost per cycle, files under design/roles and design/formulas, flai serve runs formula steps"
anchor:
  path: wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md
  item: E-0024
status: open
participants: [planner-E-0024]
created: 2026-10-16T05:27:15Z
updated: 2026-10-18T03:44:28Z
---

# TH-0471 E-0024 needs a cost of delay input and five design answers: recommend 4h lost per cycle, files under design/roles and design/formulas, flai serve runs formula steps

On wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md.

## Entries

### 2026-10-16T05:27:15Z planner-E-0024
I am drafting E-0024's stories now. Six inputs are yours. Each has my recommendation first. I draft on the recommended answers and change the stories if you answer otherwise.

## 1. Cost of delay input (blocks the stories' values)

`flai cod E-0024` refuses: the epic has no inputs.

- **Recommended:** `time_lost_per_cycle: 4h` on E-0024. At `planning.hour_rate` 150 and the default 168h cycle, that is 600 USD a week, shared among the stories by duration. It counts the hand-run publishing, import, and analysis cadence, and the flai release a new role needs today. E-0016, which built the three roles, had 10h.
- Alternatives: `2h` (300 USD a week), or `10h` as E-0016 (1500 USD a week).

## 2. Where role and formula files live, and their format

The layout allows no new top-level folder, and `design/conventions` is for norms.

- **Recommended:** `design/roles/<name>.md` and `design/formulas/<name>.md`, two new folders under `design/`, named in the manifest as the other design folders are. Each file is markdown with YAML front matter, as the documentation standard asks. A role's front matter holds its host action, trigger, topics, permissions, tools, and commands. Its body is its harness prompt. A formula's front matter holds its steps and their `after`. Each has the baseline marker, below which a project adds its own, as conventions do.
- Alternatives: one folder `design/workflows/` with `roles/` and `formulas/` below it; or YAML files beside `system-flow.yaml`.

## 3. A formula instance's nature

- **Recommended:** the formula file names one of the five natures, `improvement` by default, and the story it makes carries `formula: <name>` in front matter, so stats can tell an instance apart. No new nature, since natures drive releases, metrics, and templates.
- Alternative: a sixth nature, `formula`, accepted without a release as `research` is.

## 4. A role's model and provider

- **Recommended:** yes. The manifest's `roles.<name>.agent` is merged over the project's `agent` as `planning.agent` is. `planning.agent`, `orchestration.agent`, and `analysis.agent` stay valid and win for their roles, so nothing changes for the three.
- Alternative: no per-role agent. A defined role runs on the project's default agent.

## 5. Who works a formula's instance

The epic says `flai formula run` makes a story whose tasks are the steps, worked, measured, and archived as any story.

- **Recommended:** `flai serve` works the instance itself, in the main checkout with no branch. It runs each command or script step in `after` order, and starts the role for a role step and waits for its run. It moves each task to done when its step passes, and moves the story to review when all pass. On a failure it stops and opens a thread on the story. `flai accept` accepts an instance with no branch by archiving it, with nothing to merge. Deterministic steps need no model session.
- Alternative: a story's agent works the instance as any story, each task's `## Work` naming its step.

## 6. A role's harness definition

ADR-0130 keeps `.claude/agents/` as the harness's source.

- **Recommended:** a role file names its definition, `.claude/agents/<name>.md`. `flai check` reports a role whose definition is missing, or whose definition lists a tool the role's rules refuse. flai never writes `.claude/`. The template ships the definition beside the role file.
- Alternative: flai generates `.claude/agents/<name>.md` from the role file, which needs the operator's permission for each write under `.claude/`.

Meanwhile I plan everything that needs no answer. Once question 1 is answered, I write the stories' values.

### 2026-10-18T03:44:28Z planner-E-0024
Still open after about 46 hours, so I am ending the planning run, as planner-E-0021 did on TH-0463. The plan is otherwise complete: S-0414 to S-0425 are drafts in the backlog, each with touches, a forecast, Planning notes, and its tasks in layers (66 tasks, listed on TH-0473).

What stays open here:

- **Question 1:** S-0414 to S-0425 have no cost of delay values, so the cost-of-delay and WSJF orderings rank them last. Once you set inputs on E-0024 (my recommendation is still `time_lost_per_cycle: 4h`), `flai cod` works the values out and the next planner run on E-0024 writes them.
- **Questions 2 to 6:** the stories are drafted on my recommendations: `design/roles` and `design/formulas`, a `formula:` field with no new nature, `roles.<name>.agent`, flai serve working formula instances, and the role file naming its `.claude/agents` definition. If you answer any of them otherwise, plan E-0024 again, or tell the story's agent when you finalize that story.
