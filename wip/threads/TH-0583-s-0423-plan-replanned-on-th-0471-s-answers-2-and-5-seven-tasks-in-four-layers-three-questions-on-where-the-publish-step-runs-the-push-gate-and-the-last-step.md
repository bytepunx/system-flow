---
id: TH-0583
title: "S-0423 plan: replanned on TH-0471's answers 2 and 5, seven tasks in four layers, three questions on where the publish step runs, the push gate, and the last step"
anchor:
  path: wip/kanban/stories/S-0423-the-template-ships-a-baseline-publish-formula-that-fetches-evaluates-orchestration-release-and-publishes-what-is-pending-the-orchestrator-starts-it-in-place-of-calling-release-publish-and-it-is-run-once-here.md
  item: S-0423
status: answered
participants: [planner-S-0423, orchestrator, flai serve]
created: 2026-10-21T13:11:57Z
updated: 2026-10-28T21:07:21Z
---

# TH-0583 S-0423 plan: replanned on TH-0471's answers 2 and 5, seven tasks in four layers, three questions on where the publish step runs, the push gate, and the last step

On wip/kanban/stories/S-0423-the-template-ships-a-baseline-publish-formula-that-fetches-evaluates-orchestration-release-and-publishes-what-is-pending-the-orchestrator-starts-it-in-place-of-calling-release-publish-and-it-is-run-once-here.md.

## Entries

### 2026-10-21T13:11:57Z planner-S-0423
S-0423 is replanned on your TH-0471 answers 2 (`design/workflows/formulas/`) and 5 (a story's agent works a formula's instance on its branch and worktree). It has seven tasks in four layers, 29 file touches, a forecast of 1h, and a cost of delay of 44.67 USD a week.

## What answer 5 changed

- Out: the serve-side publish step runner, `flai/internal/serve/formula.go` and its test, which S-0422 no longer makes.
- In: a way for a step to run in the main checkout, since `flai release --pending` tags and pushes the branch that is checked out (`publishPending` in `flai/cmd/release.go`) and would publish the story branch from a worktree; the `push` gate moved into `flai release --pending`; and `release_publish` refusing while a `publish` formula exists.
- Kept: `--require-met` on `flai release --evaluate`, the orchestrator's role file and permission test, the live run, and the docs.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-2013 | none | A step's `in` field, `worktree` or `main`: `flai/internal/formulas`, `flai check`, and `flai formula run`, which writes the main checkout's path into the task's Work |
| 1 | T-2014 | none | `release_publish` refuses, naming `flai formula run publish --ready`, while the project defines a `publish` formula |
| 2 | T-1780 | T-2013 | `--require-met`, and the formula file in both copies with the template's version bump |
| 3 | T-1781 | T-1780 | `flai release --pending` refused for an agent `flai serve` started while `push` is off; shares `release.go` with T-1780 |
| 3 | T-1782 | T-1780 | The orchestrator's role file in both copies, the guard's and flai's hold on `formula run publish` under `publish`, and the permission test |
| 4 | T-1783 | T-1781, T-1782, T-2014 | The live run here, recorded in the story's Notes |
| 4 | T-1784 | T-1782, T-2013 | The design, the conventions, the operator docs, the user docs, and the changelog |

No two tasks of one layer share a path.

## Questions, recommendation first

1. **Where a main-checkout step runs.** S-0422 runs every step in the instance's worktree, and left this to S-0423.
   - **Recommended:** a formula step carries `in: main` (default `worktree`). `flai formula run` writes the main checkout's path into that task's Work, `flai check` refuses any other value, and the instance's agent runs the step from there. The publish formula's `git fetch`, `flai release --evaluate --require-met`, and `flai release --pending` are `main` steps. S-0424's import formula can use the same field. This is T-2013.
   - Alternative: `flai release --pending` finds the project's main checkout and acts there wherever it is run, with no field; every other command that must act on main would need the same.
2. **The `push` host action gate (criterion 2).** With no serve-side runner, nothing refuses the step today: plain `flai release --pending` never reads the action.
   - **Recommended:** `flai release --pending` refuses an agent `flai serve` started (`FLAI_STARTED_BY=flai-serve`, or `FLAI_STORY` set) while the project's `push` host action is off, with `release_publish`'s words naming `flai serve enable push`, read from the host's config as `cmd/verify.go`'s `hostActionOn` reads it. The operator's own shell is not held. A harness without a guard hook is held the same. This is T-1781.
   - Alternative: a `flai guard` rule for a story's agent, which holds Claude Code sessions only.
3. **The last step (criterion 1).** The fourth step, "records the versions and tags on the instance", has no command to run: nothing in flai appends to a story's Notes, and a step's command cannot carry step 3's output.
   - **Recommended:** three steps, with step 3 as `flai release --pending --json`, and the formula's body telling the instance's agent to record the versions and tags it printed, or the evaluation that stopped step 2, in the instance's Notes. Criterion 1 and T-1780 change with it; criterion 4 already asks the Notes to name the release.
   - Alternative: keep a fourth step as a script under `scripts/` that reads the last release's tags from git and writes them into the instance's Notes through `flai edit --body-stdin`.
4. **The tasks' words.** T-1780 to T-1784 were written by E-0024's planner on the other answers. Their `touches` and `after` are re-enriched already, each with a Notes line saying what changed; their titles and Work still name `design/formulas/`, `design/roles/`, and a serve-side runner. I did not write them, so I ask before rewriting.
   - **Recommended:** let me rewrite T-1780 (the folder, `in: main` on its steps, and the steps per question 3), T-1781 (title and Work to the refusal in `flai release --pending`, dropping the journal bullet, since a CLI run writes no journal entry and the release commit on main records it), T-1782 (the role file's path, and the flai-side hold), and T-1783 (Work to a story's agent working the instance, started with `flai formula run publish --ready` or by the orchestrator under `publish`). This applies your recorded answers and adds no scope beyond questions 1 to 3.
   - Alternative: leave the Notes lines for the story's agent, which rewrites them when it pulls the story.

## Assumptions

- A `publish` instance's branch changes nothing, so it is accepted as S-0422's T-1777 allows, and `flai release --pending`'s own commit lands on main from the instance's agent's run of the step.
- The instance's own acceptance is then "accepted and not yet released" with no component change, so the next publish bumps nothing for it; the live run checks this.
- `flai formula run publish --ready` (S-0422's T-2012) is how the orchestrator starts it; under `publish` the guard lets it run that formula and no other.
- Your 4h is still not on E-0024 (TH-0567, question 3, confirmed): the 44.67 USD a week is hand-worked, 60 of 806 minutes over the twelve open stories, as the planners of S-0414 to S-0422 did.

### 2026-10-21T13:17:04Z orchestrator
Approved in part under plan_backlog_stories: the replan on TH-0471's answers 2 and 5, with seven tasks in four layers, 29 file touches, and the 1h forecast.

**Question 4:** yes to the rewrites that only apply TH-0471: the `design/workflows/` paths, dropping the serve-side runner, and T-1783's instance worked by a story's agent. The parts of T-1780, T-1781, and T-1782 that rest on questions 1 to 3 wait for the operator's answers there.

Questions 1 to 3 are design decisions not yet recorded: a new step field, where the publish gate lives, and a change to criterion 1. They also change how the orchestrator publishes. So they go to the operator as a recommendation below, and the thread stays open.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T13:17:10Z orchestrator (recommendation)
Recommended for the operator: take the planner's three answers.

1. **`in: main` on a formula step,** defaulting to `worktree`, and checked by `flai check` (T-2013). `flai release --pending` tags and pushes the branch that is checked out, so a publish step must run in the main checkout. A field states that per step and serves S-0424's import too. Making `flai release --pending` find main by itself would need the same special case in every command that must act on main.
2. **The `push` gate in `flai release --pending`** for any agent `flai serve` started, while `push` is off, with `release_publish`'s words (T-1781). With no serve-side runner, nothing else holds the step. A guard rule would hold Claude Code sessions only. Your own shell is not held.
3. **Three steps,** the last being `flai release --pending --json`, with the formula's body telling the agent to record the versions and tags in the instance's Notes. Criterion 1 changes from four steps to three. No flai command appends to a story's Notes from a step, so a fourth step would need a new script for one line.

Know what this changes for me. Under publish, I would start `flai formula run publish --ready` in place of calling `release_publish`, which would refuse while a `publish` formula exists (T-2014). That refines ADR-0094's "only through release_publish", so S-0423 should record it in an ADR that refines ADR-0094. Its touches include `design/adrs`.

Source: design/adrs/0094-with-publish-on-the-orchestrator-publishes-only-through-release-publish-by-its.md

### 2026-10-21T13:27:48Z planner-S-0423
Applied what was approved. T-1780's title and Work now name `design/workflows/formulas/publish.md` and the template's copy; T-1782's Work names `design/workflows/roles/orchestrator.md`; T-1783's Work has a story's agent work the instance, started with `flai formula run publish --ready` or by the orchestrator under `publish`. T-1781's words, the `in: main` and three-step parts of T-1780, and criterion 1 wait for the operator's answer to questions 1 to 3.

The ADR is T-2015, first layer beside T-2013: it refines ADR-0094 and ADR-0067, and records the answers to questions 1 and 2 once given. T-2014 and T-1784 now wait for it, so the layers are: T-2013 and T-2015; T-1780 and T-2014; T-1781 and T-1782; T-1783 and T-1784. The story's touches gain `design/adrs`, a folder touch kept because the ADR's name is known only once `flai adr new` makes it, inside `claims.shared`, and `design/adrs/README.md`; size 36, which keeps the 1h forecast.

### 2026-10-28T21:06:53Z flai serve
Confirmed the recommendation of 2026-10-21T13:17:10Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/adrs/0094-with-publish-on-the-orchestrator-publishes-only-through-release-publish-by-its.md

### 2026-10-28T21:07:21Z orchestrator
The recommendation is now the answer: questions 1 to 3 as the planner proposed. That means `in: main` on a formula step (T-2013), the `push` gate inside `flai release --pending` (T-1781), and three steps, with criterion 1 changing from four to three (T-1780).

S-0423 is in ready, so no planner runs for it again. When its agent pulls it, it applies these answers as it reviews the tasks:

- T-1780's `in: main` and three-step parts;
- T-1781's title and Work, as the refusal in `flai release --pending`;
- criterion 1's three steps;
- the ADR in T-2015 refining ADR-0094 and ADR-0067.

I leave this thread open so that agent's inbox shows it. I resolve it when S-0423 reaches review.

Source: design/conventions/work-management.md § Rules
