---
id: S-0564
type: story
nature: improvement
title: The import and analysis-cadence formulas are each run once in this repository, each instance worked to review and accepted
status: review
parent: E-0024
owner: alex
created: 2026-11-02T18:08:35Z
updated: 2026-11-07T09:13:45Z
transitions:
  - to: ready
    at: 2026-11-07T08:02:12Z
    by: alex
  - to: in-progress
    at: 2026-11-07T08:02:18Z
    by: agent-S-0564
  - to: review
    at: 2026-11-07T09:13:45Z
    by: system-flow
blocked:
  - from: 2026-11-07T08:23:28Z
    until: 2026-11-07T08:26:54Z
    reason: awaiting the operator on TH-0853
  - from: 2026-11-07T08:32:14Z
    until: 2026-11-07T08:33:39Z
    reason: awaiting the operator on TH-0853
  - from: 2026-11-07T08:37:54Z
    until: 2026-11-07T08:45:33Z
    reason: awaiting the operator on TH-0860
tags: [template]
topics: [orchestration]
touches: [wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md, design/issues/summary.md, design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md, design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md, design/issues/I-0302-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md]
after: [S-0424]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3524
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 5
      test_runs: 1
      hand_edits: 2
      work: 143
  models:
    - model: claude-opus-5-5
      input: 326
      output: 42537
      cache_read: 20321318
      cache_write: 205064
      cost: 6.5568
  strategic:
    - kind: orchestrator
      seconds: 397
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 90
          output: 573
          cache_read: 3520303
          cache_write: 438214
          cost: 0.9734
        - model: claude-sonnet-5-5
          input: 14
          output: 951
          cache_read: 198717
          cache_write: 41489
          cost: 0.1741
cost_of_delay:
  inputs:
    penalty_per_week: 25
    by: orchestrator
    at: 2026-11-07T03:54:40Z
  value: 25
  by: planner-S-0564
  at: 2026-11-07T03:56:02Z
forecast:
  duration: 1h
  delivery: 2026-11-07T14:00:00Z
  basis: "Its own forecast of 1h; 13th in the pull order with an in-progress limit of 5, behind S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0560, S-0561 and S-0562."
  by: flai
  at: 2026-11-07T07:55:54Z
finalized:
  by: alex
  at: 2026-11-07T08:01:35Z
---
# S-0564 The import and analysis-cadence formulas are each run once in this repository, each instance worked to review and accepted

## Goal

Run each baseline formula S-0424 shipped once in this repository, so that the epic gets its live runs before it is accepted. One run imports a scratch repository from the board. The other is one analyzer cadence. Each instance is worked by a story's agent to review and accepted.

## Acceptance criteria

- [x] Each formula is run once in this repository: an import of a scratch repository under a folder named for import, and one analyzer cadence, each instance worked to review and accepted, with its Notes naming what it produced.

## Tasks
- T-2514 Check the import run's preconditions, make a scratch repository under a folder named for import, and import it from the board's Import
- T-2520 Check the cadence run's preconditions and start one instance of the analysis-cadence formula from the main checkout
- T-2523 Follow the import instance in the new project to review and acceptance
- T-2526 Follow the analysis-cadence instance to review and acceptance
- T-2532 Record both runs in S-0564's Notes: each instance's ID and step exits, the scratch repository, the analyzer's report, the issues it names, and each run's cost
- T-2534 File each fault the import and cadence runs found in flai as an issue, and say in this task's Notes when they found none

## Notes

- Moved from S-0424 (criterion 4, T-1787 and T-1788) on TH-0810, which the operator confirmed on 2026-11-02, as D-0461 moved S-0423's live run.
- `after: S-0424`: both runs need what S-0424 adds on the main branch: the two formula files, `flai import --formula`, `flai issue list --report`, and the analyzer role step (ADR-0245). `flai formula run` reads the main checkout. The board's Import runs the installed flai, so the import run also waits until S-0424 is published.
- **The import run (from T-1787).** Make a scratch git repository under a folder the host serves for import (`flai serve import add`), outside this repository. Import it from the board's Import, which runs `flai import --formula`. The instance is a story in the imported project, worked by a story's agent on its branch and worktree, so the host needs the `agent` host action on for the new project. Watch the instance reach review, and have it accepted. Record the instance's ID and each step's exit in this story's Notes. Remove the scratch repository afterwards only with the operator's word.
- **The cadence run (from T-1788).** With the `analyze` host action on, run `flai formula run analysis-cadence --input focus=<focus> --ready --autocommit` in the main checkout, with the focus the operator chooses when finalizing this story, or none for all three. The instance's agent starts the analyzer step with `role_run`, waits for the run's `agent` event on the instance, then runs the listing step. Watch the instance reach review, and have it accepted. Record in this story's Notes the instance's ID, the report it wrote, the issues it names, and the run's cost.

### The runs, 2026-11-07

**The import run: S-0001 in `/home/alex/git/flai-scratch-import`, accepted.**

- **Scratch repository.** `/home/alex/git/flai-scratch-import`, below `/home/alex/git`, a folder `flai serve import list` names. One commit, `aedf008`, with a `README.md` and a `Makefile` whose `test` target passes. It stays until the operator says to remove it.
- **Import.** The operator pressed Import on the board at 08:25:46Z (TH-0853). `flai import --formula` committed the import as `1e7935b` and made the instance S-0001, ready, at `f52c1d3`. The operator turned on `agent` for the project at 08:26:06Z.
- **No agent until started by hand.** S-0001 had entered ready before `agent` was on, so flai serve never started it (I-0300). The operator started it on TH-0860.
- **Step 1, T-0001, `flai check --strict`.** It exited 1 on one warning: `layout.lock: template 1.1.24 is recorded but there is no system-flow.lock.yaml`. The import writes no lock (I-0302). The agent blocked on TH-0001 from 08:42:57Z to 08:48:56Z. The operator allowed `flai upgrade --relock` as a one-off, which relocked 74 files in `f71d76d`. The check then exited 0.
- **Step 2, T-0002, `flai import test`.** It exited 0: `make test` passed.
- **Acceptance.** S-0001 was accepted and archived at `049cf36`. Its branch changed `system-flow.yaml` and `system-flow.lock.yaml`, through the relock, though the formula says the branch changes nothing.
- **Cost.** 0.27 USD and 0.43 USD over S-0001's two agent sessions (30 s and 60 s), from their logs.

**The cadence run: S-0584, accepted.**

- **Instance.** `flai formula run analysis-cadence --ready --autocommit` ran in the main checkout with the focus empty, so all three areas. The operator finalized this story without choosing a focus, and the Notes say none in that case. It made S-0584 (T-2545 analyze, T-2547 issues) in ready at `aa01853e5`.
- **Held by its own starter.** S-0584 entered ready with no touches, and was held (no-touches) by S-0564 itself. `flai touches S-0584 --add` on its own story file cleared it (I-0282).
- **Step 1, T-2545, the analyzer role step.** The analyzer started at 08:04:52Z and ended with its work done at 08:22:44Z. The end came late: the stdio `flai mcp` held the run as a zombie, so S-0584 waited on TH-0859 (I-0299, recorded by S-0584).
  - It wrote `design/analysis/2026-11-07-all.md`, with 18 findings.
  - It filed I-0284 to I-0298, and bumped I-0230, I-0246, I-0257, and I-0268.
- **Step 2, T-2547, `flai issue list --report latest --json`.** It exited 0 in the main checkout. It named 19 open issues: I-0230, I-0246, I-0257, I-0268, and I-0284 to I-0298. S-0584's Notes list each with its title and class.
- **Acceptance.** The orchestrator accepted S-0584 under `accept_reviews`, archived at `07fa45be7`.
- **Cost.** The analyzer: 15.50 USD over 796 s, from its log. S-0584's agent: 5.20 USD in its first session. Its second session's log reports 6.56 USD, which looks cumulative across both.

**Faults the runs found in flai, filed by T-2534:**

| Issue | Class | Fault |
|-------|-------|-------|
| I-0282 | defect | A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress |
| I-0300 | defect | A story already in ready when `agent` is turned on gets no agent until someone presses Start agent |
| I-0302 | defect | `flai import` writes no `system-flow.lock.yaml`, so the import formula's `flai check --strict` step fails on `layout.lock` |
| I-0299 | defect | Recorded by S-0584: a role run that the stdio MCP server starts through `role_run` is never reaped |

### Planning

- **Plan thread.** TH-0845, 2026-11-07. It asks the operator three questions, each with a recommendation. Question 1 is the cadence's focus: none, for all three, is recommended. Question 2 is where the scratch repository goes and who presses the board's Import: the story's agent makes it below a folder already named for import, and the operator presses Import and turns on `agent` for the new project. Question 3 is the cost of delay's basis, which the orchestrator settled with an input. T-2514 and T-2520 check the first two answers before they run anything.
- **`after: S-0424`, kept.** This story reads S-0424's output: the `import` and `analysis-cadence` formula files, `flai import --formula`, `flai import test`, `flai issue list --report`, and the analyzer role step (ADR-0245). S-0424 is accepted and archived, so the wait is met. The installed flai 1.83.0 has `--formula` and `--report`, so the board's Import no longer waits on a publish either.
- **Layers.** Three layers, with the two runs side by side.
  - Layer 1: T-2514 starts the import, T-2520 starts the cadence.
  - Layer 2: T-2523 follows the import instance, after T-2514. T-2526 follows the cadence instance, after T-2520.
  - Layer 3: T-2532 records both runs in this story's Notes, and T-2534 files the runs' faults. Both wait for T-2523 and T-2526.
  - No two tasks of one layer share a path. The first four write only their own Notes.
- **Touches.** Two files, both from the tasks.
  - **Design and goal:** this story's own file, where T-2532 records the runs, as the Notes above ask.
  - **Layout:** `design/issues/summary.md`, which T-2534 changes when it files or bumps an issue. A new issue's file has a name no task can know yet, and `flai task done` widens the touches by it.
  - **Co-change:** `flai touches suggest` from the two formula files found no file that changes with them often enough.
  - **Not added:** the formula files and flai's code, which the runs use and do not change. The analyzer's report under `design/analysis`, which the analyzer writes in the main checkout. The scratch repository, which lies outside this repository.
  - No folder touch.
- **Forecast.** flai's 4m, at size 3 (1 criterion, 2 touches), is 80 s a unit over 13 done improvement stories on this model in the small band. I raised it to 1h because the size counts none of the live work: starting two instances, following each through its steps, review, and acceptance, and recording them. The analyzer on all three focuses runs while this story's agent waits on it, and the last single-focus run took 454 s. S-0562's planner raised one publish run to 45m. The delivery is flai's play-out, 2026-11-07T09:09Z (19th in the pull order), plus the 56m added: 2026-11-07T10:05Z.
- **Cost of delay.** 25.00 USD a week, as `flai cod S-0564` gives from the input `penalty_per_week: 25`. The orchestrator set that input on 2026-11-07, as it did for S-0562 on TH-0839: a proof run that holds the epic's acceptance and no other work. It replaces my first value, 342.86, which I worked out by hand from E-0024's 4h per cycle (TH-0471) before S-0564 had an input.
