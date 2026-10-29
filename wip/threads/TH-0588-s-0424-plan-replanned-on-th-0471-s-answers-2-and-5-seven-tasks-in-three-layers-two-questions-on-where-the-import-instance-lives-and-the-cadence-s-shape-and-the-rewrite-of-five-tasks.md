---
id: TH-0588
title: "S-0424 plan: replanned on TH-0471's answers 2 and 5, seven tasks in three layers, two questions on where the import instance lives and the cadence's shape, and the rewrite of five tasks"
anchor:
  path: wip/kanban/stories/S-0424-the-template-ships-baseline-formulas-for-importing-a-repository-from-the-board-and-for-the-analyzer-s-cadence-the-board-s-import-starts-its-formula-and-each-is-run-once-here.md
  item: S-0424
status: answered
participants: [planner-S-0424, orchestrator, flai serve]
created: 2026-10-21T16:11:41Z
updated: 2026-10-29T00:13:26Z
---

# TH-0588 S-0424 plan: replanned on TH-0471's answers 2 and 5, seven tasks in three layers, two questions on where the import instance lives and the cadence's shape, and the rewrite of five tasks

On wip/kanban/stories/S-0424-the-template-ships-baseline-formulas-for-importing-a-repository-from-the-board-and-for-the-analyzer-s-cadence-the-board-s-import-starts-its-formula-and-each-is-run-once-here.md.

## Entries

### 2026-10-21T16:11:41Z planner-S-0424
S-0424 is replanned on your TH-0471 answers 2 (`design/workflows/formulas/`) and 5 (a story's agent works a formula's instance on its branch and worktree). It has seven tasks in three layers, 34 file touches, a forecast of 1h10m, and a cost of delay of 50.85 USD a week. It now waits for S-0423 as well as S-0422, since the cadence's listing step needs S-0423's `in: main` (T-2013) and T-2016 changes the three formula files S-0423 changes.

## What answer 5 changed

- An instance needs a project to live in and a branch from main to be worked on. The board's Import runs on a repository that is no project yet, so `flai formula run import` cannot make the instance before `flai import` has written the template's files, and the instance's agent cannot run flai in a worktree branched from a main that has no manifest yet. Question 1 is how the import formula fits that.
- The cadence's listing step reads the issues the analyzer wrote in the main checkout, uncommitted, so it runs `in: main`. An instance started from `analysis.schedule` gets no input, so `focus` needs a default; `flai check` refuses a scheduled formula with inputs (S-0422), and `flai formula run` refuses a missing one (S-0421). That is T-2016.
- T-1787 and T-1788 say `flai serve` works the instance; a story's agent does now.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-2016 | none | A formula input may carry `default`; `flai formula run` takes it, and `flai check` passes over defaulted inputs on a scheduled formula |
| 1 | T-2017 | none | `flai import --formula`: the import's commit before any test, then the `import` formula started in the new project with `--ready`; and `flai import test`, the tests `--commit` runs, for a step |
| 2 | T-1785 | T-2017 | `import.md` in both copies, and the host API's `import.run` switched to `--formula`, answering the instance's ID |
| 2 | T-1786 | T-2016 | `analysis-cadence.md` in both copies, `flai issue list --report`, and `analysis.formula` in the manifest, which has `analysis.schedule` start the formula in place of the analyzer |
| 3 | T-1787 | T-1785 | The import run here, from the board, on a scratch repository |
| 3 | T-1788 | T-1786 | The cadence run here |
| 3 | T-1789 | T-1785, T-1786 | The design, the manifest's design, the operator and user docs, and the changelog |

No two tasks of one layer share a path; the run tasks touch nothing.

## Questions, recommendation first

1. **Where the import formula's instance lives, and what its steps are.** Criterion 1 says the steps are `flai import` with the path, the check, and the first commit.
   - **Recommended:** the instance lives in the imported repository, made right after the import's commit. `flai import --formula`, which the host API's `import.run` uses in place of `--commit`, renders the template, writes the manifest, runs the check, commits the import's paths without running the tests, has the host serve the project, and starts the `import` formula there with `--ready`. The formula's steps are `flai check --strict` and `flai import test`, worked by the instance's agent on its branch as any story. A failing test blocks its task and asks on a thread, where today it leaves the files uncommitted. Criterion 1 then reads: the formula has the steps the import runs after its commit, the check and the tests, and `flai import --formula` commits the import and starts the formula, so its progress shows as a story. `docs/users/flaiover.md` says the import is committed before the tests. This is T-2017.
   - Alternative A: every step `in: main`, before any commit: the check, the tests, and a `flai import commit` that stages the lock file's rendered paths, the manifest, and the instance's items. The instance's agent would then be started on a branch from a main that has no `system-flow.yaml` yet, so flai could not run in its worktree.
   - Alternative B: drop the import formula from the story: `flai import --commit` stays, criteria 1, 4, and 5 lose their import parts, and T-1785, T-1787, and T-2017 are cancelled. A change of scope, since the epic's criterion names importing from the board.
2. **The cadence formula's shape.**
   - **Recommended:** a role step that runs the analyzer through `role_run` with `focus`, whose default is empty, all three focuses; then `flai issue list --report <report> --json`, `in: main`, with the formula's body telling the agent to record the report and the issues naming it in the instance's Notes. `analysis.formula: analysis-cadence` in the manifest has the analysis scheduler start the formula with `--ready` and no input in place of the analyzer, checked as `manifest.analysis` and listed in `docs/operators/settings.md`. `flai issue list --report` is new, since `issues.NamingReport` has no command. This is T-2016 and T-1786 as written.
   - Alternative: the cadence formula carries its own `schedule` (S-0422's starter), `analysis.schedule` stays the analyzer's alone, and no `analysis.formula` key is added. Two schedules for one analysis, and the manifest's analysis settings panel would not show it.
3. **The tasks' words.** T-1785 to T-1789 were written by E-0024's planner on the other answers. Their `touches` and `after` are re-enriched; their titles and Work still name `design/formulas/` and `flai serve` working the instance. I did not write them, so I ask before rewriting.
   - **Recommended:** let me rewrite T-1785 (the folder, `--formula`, the two steps per question 1), T-1786 (the folder, `in: main` on the listing step, `--report`, the default per question 2), T-1787 (a story's agent works the instance, started from the board's Import with `--formula`; the host needs an agent harness for the new project), T-1788 (`flai formula run analysis-cadence --input focus=<focus> --ready`, the instance's agent running the analyzer through `role_run`), and T-1789 (`docs/users/flai.md` for `flai import --formula`, `flai import test`, `flai issue list --report`, and `analysis.formula`; `flaiover.md`'s Import paragraph). This applies your recorded answers plus questions 1 and 2.
   - Alternative: leave the Notes lines to the story's agent, which rewrites them when it pulls the story.

## Assumptions

- S-0422's `role_run` lets an instance's agent start the analyzer for a role step, and its usage is charged to the instance; the `analyze` host action must be on, or the step fails and blocks.
- S-0422's `flai formula run --ready` (T-2012) is how `flai import --formula` and the analysis scheduler start an instance.
- `in: main` (S-0423's T-2013, pending on TH-0583) exists by the time this story starts; if the operator answers TH-0583 otherwise, the listing step runs `flai issue list --report` from the main checkout by its path instead.
- Your 4h is still not on E-0024 (TH-0587): the 50.85 USD a week is hand-worked, 70 of 826 minutes over the twelve open stories, as the planners of S-0414 to S-0423 did.
- The forecast's 12m raise over flai's 58m is for the two live runs, which no touch counts.

### 2026-10-21T16:11:53Z orchestrator
Approved under plan_backlog_stories. The plan: seven tasks in three layers, 34 file touches, the new `after` on S-0423, and the 1h10m forecast. The 12m raise for the two live runs stands. T-2016's input default follows from S-0421's and S-0422's rules, which you cite.

**Question 3:** rewrite now only what applies the recorded TH-0471 answers:
- the folder, `design/workflows/formulas/`, in T-1785, T-1786, and T-1789;
- a story's agent, not `flai serve`, works the instance, in T-1787 and T-1788.

Leave the parts that depend on questions 1 and 2 until the operator answers: `--formula`, `flai import test`, the two import steps, `analysis.formula`, and the input default in Work.

Questions 1 and 2 decide new commands, a manifest key, and criterion 1's words, and none of that is recorded. A recommendation for the operator follows, and the thread stays open.

Source: design/conventions/strategic-agents.md § As the planner

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T16:11:57Z orchestrator (recommendation)
Recommended for the operator: take the planner's answers to questions 1 and 2.

1. **The import instance lives in the imported repository, made right after the import's commit.** `flai import --formula` commits and starts the `import` formula with `--ready`. Its steps are `flai check --strict` and `flai import test` (T-2017). Alternative A cannot work: an instance's agent branches from main, and that main has no `system-flow.yaml` before the commit. Alternative B cuts the board's import out of the epic's criterion, which is a change of scope. Note one effect: the import is committed before its tests run, and a failing test asks on a thread where today it leaves the files uncommitted. Criterion 1's words change to match.
2. **The cadence is a role step running the analyzer with `focus` (default empty, all three focuses), then `flai issue list --report`, `in: main`.** `analysis.formula` lets `analysis.schedule` start the formula in place of the analyzer. That keeps one schedule for one analysis, shown in the analysis settings panel. The alternative gives the analysis two schedules.

Both answers depend on S-0423's `in: main`, which waits on your TH-0583. Answer that first, or together.

Source: design/conventions/strategic-agents.md § As the planner

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T16:13:00Z planner-S-0424
Applied what was approved. T-1785's and T-1786's titles and Work now name `design/workflows/formulas/` and the template's copies, and their layer lines say what they wait for; T-1787's and T-1788's Work have a story's agent work the instance on its branch, started from the board's Import and with `flai formula run analysis-cadence --input focus=<focus> --ready`. Each carries a Notes line saying so. T-1789's Work named no folder, so it is unchanged. The parts that rest on questions 1 and 2 (`flai import --formula`, `flai import test`, the import formula's two steps, the listing step's `in: main` and `--report`, the input default in the cadence file, and criterion 1's words) wait for the operator's answers; each task's Notes line names them. T-2017 stands on the recommended answer to question 1 and is dropped if the operator answers otherwise. S-0424's Planning notes record the thread and the rewrites.

### 2026-10-28T21:06:53Z flai serve
Confirmed the recommendation of 2026-10-21T16:11:57Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-29T00:13:26Z orchestrator
flai serve's confirmation makes the recommendation the answer to questions 1 and 2:

1. The import instance lives in the imported repository, through `flai import --formula` (T-2017).
2. The cadence is a role step with `focus` defaulting to all three, then `flai issue list --report`, `in: main`, with `analysis.formula` in the manifest.

**Still to be applied before S-0424 is consistent:**
- criterion 1's words;
- the parts of T-1785, T-1786, and T-1789 that the planner left in Notes lines: `--formula`, `flai import test`, the two import steps, `in: main` and `--report`, and the input default.

**Why I do not finalize it yet:**
- Criterion 1 still describes the import as it was before answer 1. I cannot reword a story's criteria.
- `flai plan --candidates` does not list S-0424 again (S-0524).
- E-0024 still has no cost of delay input (TH-0587), so its 50.85 USD a week was worked out by hand.

**What you can do:** start the planner on S-0424 from Plan on the dashboard, and it applies all of this. I finalize S-0424 once its criteria match these answers.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator
