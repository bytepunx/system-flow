---
id: ADR-0238
title: "A formula is a markdown file under design/workflows/formulas whose front matter names its nature, inputs, and steps, each running a flai command, a script, or a role, and flai formula run makes from it an ordinary story with one task per step"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0207]
topics: [cli, orchestration, template, conventions]
---

# ADR-0238 A formula is a markdown file under design/workflows/formulas whose front matter names its nature, inputs, and steps, each running a flai command, a script, or a role, and flai formula run makes from it an ordinary story with one task per step

## Context

A repeatable process is written today as prose and shell scripts. Publishing a release, importing a repository, and the analyzer's cadence are examples. An agent reads the prose and follows it. Nothing checks that the steps could run, and nothing tracks each step as work. E-0024 makes such a process data that flai runs: a formula.

[ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) added `design/workflows/` to the design folder. It holds what flai runs as data: `roles/` now, and `formulas/` when E-0024's formula stories land. This decision fills in `formulas/`.

The operator's answers on TH-0471 ([D-0342](../digest/D-0342-e-0024-needs-a-cost-of-delay-input-and-five-design-answers-recommend-4h-lost.md)) settle three things this decision rests on:

- Formula files live in `design/workflows/formulas/`, beside `roles/` (answer 2).
- A formula's instance is a story of an existing nature. There is no new nature (answer 3).
- A story's agent works the instance on its branch and worktree, as it works any story, as the project is configured (answer 5). `flai serve` does not work it.

## Decision

**A formula is one markdown file, `design/workflows/formulas/<name>.md`. Its front matter names the formula, its nature, its inputs, and its steps. Each step runs a flai command, a script under `scripts/`, or a role with a prompt, after the steps it names. `flai formula run` makes from it one ordinary story, with one task per step and the steps' `after` as the tasks' `after`.** This refines [ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) §1, which named `formulas/` without saying what it holds; the rest of it stands.

### 1. Where it lives

- A formula is `design/workflows/formulas/<name>.md`, in the design folder the manifest names (`layout.design`, `design/` by default).
- Its `name` equals the file's stem. It is lowercase kebab-case.
- `README.md` in the folder is its index, not a formula. flai reads every other `*.md` in it.

### 2. The front matter

Every field is plain YAML. A value that starts with `{`, such as a title that starts with an input, is double-quoted.

| Field | Type | Meaning | Allowed values |
|-------|------|---------|----------------|
| `title` | string | The documentation standard's title. It is also the instance story's title, with inputs filled (§4). | Any. Required. |
| `updated` | date | The day of the last meaningful edit. | `YYYY-MM-DD`. Required. |
| `status` | string | The documentation standard's status. | `draft`, `active`, `deprecated`. Required. |
| `name` | string | The formula's name, as `flai formula` takes it and the instance's `formula` field records it. | Lowercase kebab-case; equal to the file stem. Required. |
| `description` | string | One paragraph: the instance story's goal. | Not empty. Required. |
| `nature` | string | The instance story's nature. | `feature`, `improvement`, `remediation`, `research`, or `experiment`. Default `improvement`. |
| `inputs` | list of mappings | The values a run takes. | Each has `name`, `description`, and optionally `default`. Default `[]`. |
| `steps` | list of mappings | What the instance does, one task per step. | Not empty. Required. |

An input:

- `name` matches `[a-z][a-z0-9_-]*` and is unique in the formula.
- `description` says what the value is. It is required.
- `default` is the value a run takes when it is not given. An input without a default is required.

A step:

- `id` is kebab-case and unique in the formula. Required.
- `title` is the task's title. Required.
- `run` is what the step runs (§3). Required.
- `after` is a list of step ids the step waits for. It may be empty or left out.

flai reads past a key it does not know. A later story, S-0422, can then add a schedule without breaking an older flai.

### 3. What a step runs

`run` is a mapping with exactly one of three kinds:

| Kind | Form | Example |
|------|------|---------|
| A flai command | `flai: <the command line after flai>` | `flai: release --evaluate` |
| A script | `script: <command line>`. Its first word is a path under `scripts/`, and the file exists. | `script: scripts/test.sh` |
| A role | `role: <name>` with `prompt: <text>`. A file under `design/workflows/roles/` defines the role. | `role: analyzer`, `prompt: Report the risk in {{area}}.` |

### 4. Inputs

- An input is used as `{{name}}`.
- It may be used in a step's title, in its run (the flai command line, the script's command line, or the role's prompt), and in the formula's title.
- A use of an input the formula does not declare is a finding (§6).

### 5. The body

- The body describes the formula for a reader. flai does not run it.
- It carries the baseline marker `<!-- system-flow:end-of-baseline -->`. Above it is the template's baseline. Below it, under `## Project additions`, is the project's own.
- For a formula the template ships, `flai upgrade` replaces the baseline: the front matter and the body above the marker come from the template, and the project's text from the marker on stays, as ADR-0207 §8 has it for a role file.

### 6. What `flai check` reports

`flai check` reports `formulas.file` on:

- a formula that does not parse, or misses a required field;
- a step whose `run` is none of the three kinds, or more than one;
- a script not under `scripts/`, or missing;
- a role no file under `design/workflows/roles/` defines;
- an `after` naming no step;
- a cycle among the steps' `after`;
- a duplicate step id;
- an input used but not declared.

It does not check that a flai command exists. The check package cannot import the command tree, which imports it, as [D-0408](../digest/D-0408-flai-check-does-not-yet-report-a-role-file-that-names-an-unknown-flai-mcp-tool.md) records for role files.

### 7. The commands

- `flai formula list` lists the formulas.
- `flai formula show <name>` prints one formula: its inputs, its steps, and the steps' layers.
- `flai formula run <name> [--input key=value]... [--epic E-nnnn]` makes the instance.

`run` first checks the inputs. It refuses a required input that is missing, and an input the formula does not declare, before it writes anything.

It then makes the instance through the same creation path as `flai story new` and `flai task new` (`itemnew`):

- **One story**, in backlog, as `flai story new` makes it. It has:
  - the formula's nature;
  - `formula: <name>` in its front matter;
  - the formula's title, inputs filled, as its title;
  - the description as its goal;
  - one acceptance criterion per step;
  - Notes naming the formula file and the inputs the run took.
- **One task per step.** Each has:
  - the step's title, inputs filled, as its title;
  - a `## Work` that names what the step runs, inputs filled;
  - a `## Done when` that is the step's exit: the command or script exits 0, or the role's run ends with its work done;
  - the step's `after` as its `after`, so the tasks' layers are the steps' layers.

It prints the story's ID. `--epic` is optional, since a story need not belong to an epic.

### 8. The instance is an ordinary story

The instance is ordered, pulled, worked, measured, accepted, and archived as any story. A story's agent works it on its branch and worktree (TH-0471, answer 5).

This decision makes the instance only. Working a step, a role step's run, schedules, and the thread a failed step opens are S-0422's.

### 9. The `formula` field

- `formula` is a story front-matter field. Only `flai formula run` writes it.
- It is listed in `flai/internal/workitem/front-matter-fields.txt`. Publishing this flai therefore raises `flai.minimum`.
- An older flai reads past it.

### Example: `release-check.md`

```markdown
---
title: "Release check for {{version}}"
updated: 2026-11-01
status: active
name: release-check
description: "Check that the work on main is ready to publish as {{version}}: the tests pass, flai's release evaluation agrees, and the analyzer reports the risk."
nature: improvement
inputs:
  - name: version
    description: The version to check, such as 1.4.0.
  - name: area
    description: What the analyzer looks at.
    default: the changes since the last release
steps:
  - id: test
    title: The tests pass
    run: {script: scripts/test.sh}
  - id: evaluate
    title: The release evaluation agrees on {{version}}
    run: {flai: release --evaluate}
    after: [test]
  - id: risk
    title: The analyzer reports the risk of {{version}}
    run:
      role: analyzer
      prompt: Report the risk in {{area}} before {{version}} is published.
    after: [test]
---

# Release check

Run before a release is published. The analyzer's step runs beside the evaluation, after the tests.

<!-- system-flow:end-of-baseline -->

## Project additions
```

`flai formula run release-check --input version=1.4.0` makes one story titled "Release check for 1.4.0" with three tasks. `test` is the first layer; `evaluate` and `risk` are the second.

## Consequences

- A repeatable process is a file that `flai check` validates. A step that could not run is a finding before anyone runs it.
- Each run is a story. It shows on the board, is worked by a story's agent, and is measured and archived as any story is. No new nature, state, or metric is needed.
- A role step names a role that the role files define, so a formula and a role share one definition of a role.
- `flai check` does not catch a flai command that does not exist. A misspelled command fails when its task is worked.
- An input's value is filled in as text. A value with shell characters reaches a script's or a flai command line as written. The formula's author quotes it where it matters, and the story's agent sees the filled command in the task's `## Work` before running it.
- Publishing this flai raises `flai.minimum` for the `formula` field.
- `repository-layout.md`, the strategic-agents design and convention, and the user and operator docs say what a formula holds and how to run one, in the template and here.

## Alternatives considered

- **Steps as plain prose in the body.** Easy to write, but nothing could check that a step could run, or make a task from it.
- **A `formula` nature for the instance.** The operator chose an existing nature (TH-0471, answer 3). The formula picks one, `improvement` by default, and `formula: <name>` marks the instance.
- **`flai serve` working the instance itself.** The operator chose a story's agent on the story's branch and worktree, as any story (TH-0471, answer 5).
- **`${name}` placeholders.** They collide with shell text in scripts and in prompts that quote shell. `{{name}}` does not.
