---
id: ADR-0239
title: "A formula's instance is worked by its story's agent step by step, a role step starts through role_run as the one exception for an agent flai serve started, a failed step blocks its task and asks on a thread, the orchestrator takes each step task's done transition as evidence, and a formula's schedule makes its instance ready while the formula host action is on"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0238]
topics: [cli, orchestration, template, conventions]
---

# ADR-0239 A formula's instance is worked by its story's agent step by step, a role step starts through role_run as the one exception for an agent flai serve started, a failed step blocks its task and asks on a thread, the orchestrator takes each step task's done transition as evidence, and a formula's schedule makes its instance ready while the formula host action is on

## Context

[ADR-0238](0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) makes a formula a file and `flai formula run` make from it an ordinary story, its instance, with one task per step. It left to S-0422 how a step is worked, how a role step runs, what a failed step does, and what starts a formula.

The operator's answer 5 on TH-0471 settles who works an instance: a story's agent, on the story's branch and worktree, as it works any story ([D-0342](../digest/D-0342-e-0024-needs-a-cost-of-delay-input-and-five-design-answers-recommend-4h-lost.md)). `flai serve` does not run the steps. Four things follow from that and need a rule:

- The story's agent must know the story is an instance, and how to work it. Its tasks are written already, and it runs each step itself.
- A role step starts a role. `role_run` and `flai role run` refuse every agent `flai serve` started, and a story's agent is one ([ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md), S-0417).
- A step's work may change no file, so the orchestrator's evidence, a changed file per criterion ([ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)), cannot vouch for an instance.
- A formula that runs on a cadence, such as a publish, needs something that makes its instance without the operator.

## Decision

**A formula's instance is worked by its story's agent, step by step, on its branch. A role step starts through `role_run`, the one exception to the refusal of an agent `flai serve` started. A failed step blocks its task and asks on a thread. The orchestrator takes each step task's done transition as its criterion's evidence. A formula's `schedule` makes its instance ready while the `formula` host action is on.** This refines [ADR-0238](0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) §8, which left these to S-0422; the rest of it stands.

### 1. The instance's agent works the steps

- `harness.Prompt` tells the agent of a story with `formula:` that the story is that formula's instance. `work-management.md` says the same.
- The tasks are the steps. The agent writes, reviews, adds, rewrites, and cancels none of them, and changes no step's words.
- It works them in their layers, a layer only once every task of the one before is done. It runs each step itself, never through a task sub-agent.
- A flai command or a script under `scripts/` runs from the root of the story's worktree. The agent reads its output.
- When a step exits 0, or its role run ends `worked`, the agent writes the exit in the task's Notes and closes the task with `task_done`.
- When every task is done, it closes the story out and moves it to review, as any story.

### 2. A role step

- The agent starts the role the step names with the MCP tool `role_run`, or `flai role run`, with the step's prompt and inputs and the instance as the run's item. It then holds `wait_for_events`.
- `role_run` and `flai role run` still refuse every agent `flai serve` started, save one: the agent of a formula's instance, on that instance as item, for a role a role step of its formula names. The agent's story is `FLAI_STORY`, else the story in its name `agent-S-nnnn`. Any other role, item, or agent is refused, naming why.
- The run is a defined role's run, with that role's checks: its host action on, its item kind, one run per item. Its trigger is `formula <name> <S-nnnn>`.
- Its end is a change of kind `agent` on the instance, by the role's agent: `failed` for a failed run, `ended` for one that asked or worked. That change is the step's exit, and the agent does not let it pass.
- The run's usage is not yet charged to the instance. It stays in the role's project strategic total, as every defined role's run does, until S-0425 charges a defined role's usage to its item (TH-0800).

### 3. A failed step

- A step fails when its command or script exits other than 0, or its role run ends `failed` or `asked`.
- The agent runs nothing after it. It blocks the task with `flai block`, naming the step and its exit, writes the exit and the last 40 lines of the step's output in the task's Notes, and opens a thread on the story naming the step and its exit.
- It writes `## Current state` and `## Next steps` and ends the session, as for any question.
- `flai serve` starts it again when the thread is answered. Its prompt then tells it to unblock the task, run that step again as the answer says, and go on with the layers.

### 4. An instance's acceptance

- `flai accept` accepts an instance whose branch changes nothing, as it accepts any such story.
- The orchestrator's evidence for an instance is one item per criterion, `- <n>: T-nnnn done`, naming the task of the criterion's step, in place of changed files.
- A criterion with no such item, or one naming a task that is not done or not the instance's, is the blocker `criterion_unevidenced`, naming it.
- `strategic-agents.md` says so.

### 5. What starts a formula

| Starter | What it does |
|---------|--------------|
| `flai formula run <name>` | Makes the instance in backlog. With `--ready` it moves it to ready as `flai move` does, and with `--autocommit` commits the move with the items. |
| The formula's `schedule` | A front-matter field: a five-field cron expression in UTC, or `daily`. While the `formula` host action is on for the project, each time it comes round `flai serve` runs `flai formula run <name> --ready --autocommit --json` in the main checkout, with no inputs. |
| A role | A role file's command entry `flai formula run <name>` lets the role run that formula alone, with no flag but `--ready`, `--input`, `--autocommit`, and `--json`. |

- The `formula` host action is off by default. It is flai's own: a role file may not name it.
- A scheduled instance is journalled under `formula`, method `serve.formula`, started by `schedule <spec>`. A run `flai formula run` refuses is logged, and tried again at the schedule's next time. Times missed between two looks come round once. While the action is off the schedule keeps up and makes nothing.
- `flai check` reports `formulas.file` on a `schedule` that does not parse, and on a schedule on a formula that declares inputs, since nobody is there to give them.

## Consequences

- An instance is worked, measured, and accepted as any story, by the same agent and the same close-out. No step runner in `flai serve` is needed.
- A step costs what the story's agent costs to run it, plus a role run for a role step. A step that fails waits for the operator on a thread, as any question does.
- The exception for a role step is narrow: one agent, one item, the roles its formula names. The role's host action must still be on.
- Until S-0425, a role step's usage shows in the role's project total, not on the instance.
- A formula with a schedule runs on the host, as the operator, without being asked. Whoever can change a formula file on main decides what runs and when, so a formula file changes only through a story the operator accepts.
- `docs/users/flai.md`, `docs/operators/index.md` and `settings.md`, the strategic-agents design and convention, `work-management.md`, and the formulas README say how an instance is worked, how a step fails, and how a formula starts, in the template and here.

## Alternatives considered

- **`flai serve` running the steps itself.** TH-0471's recommended answer 5: a runner in `flai serve` that runs each step, starts a role step's role, opens the thread on a failure, and resumes with `flai formula resume`. The operator chose a story's agent on the story's branch and worktree, as any story, so the runner, its failure path, and `flai formula resume` were left out.
- **The agent working a role step's prompt itself.** No role run, so no exception. But the role's file, its guard, and its agent would not apply, and the step would not be the role's work.
- **Lifting the refusal for every story's agent.** Simpler, but any story's agent could then start any role on the host. The exception is held to an instance and the roles its formula names.
- **A changed file as an instance's evidence.** A step such as a publish or a check may change no file, so the evidence would have nothing to name. A step task's done transition is what the step's passing leaves.
- **A schedule in `system-flow.yaml`, such as `formulas.<name>.schedule`.** A formula is one file, and ADR-0238 §2 left room for a key it reads past. The schedule sits in the file, and changes with it, through a story.
- **A scheduled instance left in backlog.** Nobody would move it, so it would not run. It is made ready, and its story's agent starts as any ready story's does.
