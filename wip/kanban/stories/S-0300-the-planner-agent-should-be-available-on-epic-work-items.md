---
id: S-0300
type: story
nature: improvement
title: The planner agent should be available on epic work items
status: ready
owner: alex
created: 2026-10-06T21:42:35Z
updated: 2026-10-06T22:31:55Z
transitions:
  - to: ready
    at: 2026-10-06T21:42:36Z
    by: alex
tags: [dashboard, cli]
topics: [planner, planning]
touches: [flaiover/src, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, ".claude/agents/planner.md", template/root/.claude/agents/planner.md, docs/users/flaiover.md, docs/users/flai.md, docs/operators/index.md, "flaiover/src/routes/items/[id]/item.svelte.test.ts"]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 2921
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 202
          output: 41
          cache_read: 1064712
          cache_write: 66559
          cost: 0.23
        - model: claude-opus-5-5
          input: 340
          output: 52689
          cache_read: 15344140
          cache_write: 275748
          cost: 7.2157
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4h
    by: alex
    at: 2026-10-06T21:42:35Z
  value: 600
  by: planner-S-0300
  at: 2026-10-06T21:47:59Z
forecast:
  duration: 40m
  delivery: 2026-10-06T22:55:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 3, behind S-0229 and S-0299."
  by: flai
  at: 2026-10-06T22:12:16Z
---
# S-0300 The planner agent should be available on epic work items

## Goal

The epic page and the epic context menu should provide interfaces for invoking the planner to write its stories and their tasks.

## Acceptance criteria
- [ ] The epic page has a plan button that invokes the planner for the epic
- [ ] The epic item's context menu has a plan option
- [ ] The planner agent is capable of writing draft stories and draft tasks for the epic

## Tasks
- T-1047 The design and the convention say an epic's planner drafts the tasks of each story it drafts
- T-1048 The planner's epic prompt tells it to enrich each story it drafts and draft that story's tasks
- T-1049 The item page test shows Plan on an open epic and not on a done one
- T-1050 The planner agent definitions say an epic's planner drafts the tasks of the stories it drafts
- T-1051 The user and operator docs say Plan on an epic drafts its stories and their tasks

## Notes

### Planning

Planned by planner-S-0300 on 2026-10-06. The plan is in TH-0201. On 2026-10-06 at 22:31Z alex took its recommendations: tasks need no draft flag, and the declared touch `flai/cmd` is dropped.

**What is already there.**

- **Criterion 1:** S-0208's `PlanAction.svelte` already shows **Plan** on an open epic's page (`flaiover/src/routes/items/[id]/+page.svelte`, where `canEdit` includes epics).
- **Criterion 2:** S-0263's `cardMenu` already offers it on an open epic's card (`flaiover/src/lib/cardmenu.ts`, released in flai 1.31.1).
- Both show only while the `plan` host action is on for the project. It was off when the story was written, which is why alex did not see them.
- **Criterion 3 is the new work:** planning an epic drafts its stories but not their tasks. The epic's planner now drafts each drafted story's tasks in the same run, with no draft flag on tasks.

**Touches, and where each came from.**

- **Declared:**
  - `flaiover/src` is kept as a folder touch. T-1049 names its one file, so the story's claim replaces the folder with that file (ADR-0096).
  - `flai/cmd` was dropped with alex's agreement in TH-0201. No task needs it: `flai plan` is unchanged.
- **From the design and the code layout:**
  - `flai/internal/harness/harness.go` and its test hold `planPrompt`.
  - `design/system/strategic-agents.md` holds the planner's "What it is told" table.
  - The convention `design/conventions/strategic-agents.md` and its template copy.
  - `template/CHANGELOG.md`.
  - `.claude/agents/planner.md` and its template copy.
  - `flaiover/src/routes/items/[id]/item.svelte.test.ts`.
- **From co-change (`flai touches suggest`) and the design together:** `docs/users/flaiover.md`, `docs/users/flai.md`, and `docs/operators/index.md`. These are the docs that describe Plan on an epic.
- **Left out of the co-change list:** the host API, the MCP server, and serve. Nothing changes in how the planner is started.

**Forecast: 40m.** flai forecast gave 24m: 84 s per unit of size over 24 large improvement stories, times size 17 from 3 criteria and 14 touches. I raised it to 40m for three reasons:

- five tasks, two of which keep two copies identical;
- a write under `.claude/` that waits on the operator's permission thread;
- the Go and flaiover test suites at close-out.

Two of the three criteria are already met, but that does not offset the rest.

I first set delivery to 2026-10-07T00:30Z, because S-0300 is held behind S-0229 and S-0299 on `flai/internal/harness/harness.go`. flai replans delivery as work ahead completes (ADR-0084). It set 22:55Z at 22:12Z, which counts the pull order but not that hold.

**Cost of delay: 600 USD a week, as flai cod worked it out.** That is the operator's 4h lost per 168h cycle at 150 USD an hour, one cycle a week. There was no reason to adjust it.

**Task layers.**

- **Layer 1:** T-1047 (design and convention), T-1048 (prompt), and T-1049 (item page test). They share no path.
- **Layer 2:** T-1050 (agent definitions) and T-1051 (docs). Both wait for T-1047 so that their words follow the convention's.
