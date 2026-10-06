---
id: S-0300
type: story
nature: improvement
title: The planner agent should be available on epic work items
status: in-progress
owner: alex
created: 2026-10-06T21:42:35Z
updated: 2026-10-06T23:29:10Z
transitions:
  - to: ready
    at: 2026-10-06T21:42:36Z
    by: alex
  - to: in-progress
    at: 2026-10-06T23:09:18Z
    by: agent-S-0300
tags: [dashboard, cli]
topics: [planner, planning]
touches: [flaiover/src, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml, ".claude/agents/planner.md", template/root/.claude/agents/planner.md, docs/users/flaiover.md, docs/users/flai.md, docs/operators/index.md, "flaiover/src/routes/items/[id]/item.svelte.test.ts", flai/cmd/plan.go, flai/internal/mcpserver/plan.go, docs/users/flai-reference.md, design/system/conventions.md, design/issues/I-0097-a-story-s-agent-that-runs-flai-stream-open-is-left-with-the-story-in-ready-so-permission-prompt-refuses-its-claude-write-until-it-moves-the-story-itself.md, design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1541
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 164
      output: 757
      cache_read: 5817588
      cache_write: 322608
      cost: 2.7054
  strategic:
    - kind: planner
      seconds: 3187
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 202
          output: 41
          cache_read: 1064712
          cache_write: 66559
          cost: 0.23
        - model: claude-opus-5-5
          input: 477
          output: 75075
          cache_read: 21368323
          cache_write: 452306
          cost: 11.0931
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4h
    by: alex
    at: 2026-10-06T21:42:35Z
  value: 600
  by: planner-S-0300
  at: 2026-10-06T21:47:59Z
forecast:
  duration: 45m
  delivery: 2026-10-06T23:50:00Z
  basis: "flai's 29m (83 s per unit of size over 25 large improvement stories, size 21: 3 criteria, 18 touches) raised to 45m for six tasks, two copies kept identical, a .claude/ write waiting on the operator, a regenerated reference, and both test suites; the story agent started at 23:03Z."
  by: planner-S-0300
  at: 2026-10-06T23:10:45Z
---
# S-0300 The planner agent should be available on epic work items

## Goal

The epic page and the epic context menu should provide interfaces for invoking the planner to write its stories and their tasks.

## Acceptance criteria
- [x] The epic page has a plan button that invokes the planner for the epic
- [x] The epic item's context menu has a plan option
- [x] The planner agent is capable of writing draft stories and draft tasks for the epic

## Tasks
- T-1047 The design and the convention say an epic's planner drafts the tasks of each story it drafts
- T-1048 The planner's epic prompt tells it to enrich each story it drafts and draft that story's tasks
- T-1049 The item page test shows Plan on an open epic and not on a done one
- T-1050 The planner agent definitions say an epic's planner drafts the tasks of the stories it drafts
- T-1051 The user and operator docs say Plan on an epic drafts its stories and their tasks
- T-1123 flai plan's help, the MCP plan tool's description, and the conventions design say an epic's planner drafts its stories' tasks

## Notes

### Planning

Planned by planner-S-0300 on 2026-10-06. The plan is in TH-0201. On 2026-10-06 at 22:31Z alex took its recommendations: tasks need no draft flag, and the declared touch `flai/cmd` is dropped.

Revisited by planner-S-0300 at 23:07Z, when the story's agent had finished T-1047 to T-1049. The revisit added T-1123 and four touches, and summarises them in TH-0223. Checked again at 23:12Z: `grep` over `design/system`, `docs`, and the planner's code paths found no other text saying an epic's planner drafts only stories, so no further task is needed.

**What is already there.**

- **Criterion 1:** S-0208's `PlanAction.svelte` already shows **Plan** on an open epic's page (`flaiover/src/routes/items/[id]/+page.svelte`, where `canEdit` includes epics).
- **Criterion 2:** S-0263's `cardMenu` already offers it on an open epic's card (`flaiover/src/lib/cardmenu.ts`, released in flai 1.31.1).
- Both show only while the `plan` host action is on for the project. It was off when the story was written, which is why alex did not see them.
- **Criterion 3 is the new work:** planning an epic drafts its stories but not their tasks. The epic's planner now drafts each drafted story's tasks in the same run, with no draft flag on tasks.

**Touches, and where each came from.**

- **Declared:**
  - `flaiover/src` is kept as a folder touch, because it was declared. T-1049 names its one file, so the story's claim replaces the folder with that file (ADR-0096).
  - `flai/cmd` was dropped with alex's agreement in TH-0201.
  - `template/template.yaml` was added by the story's agent with T-1047, for the template version bump.
- **From the design and the code layout:**
  - `flai/internal/harness/harness.go` and its test hold `planPrompt`.
  - `design/system/strategic-agents.md` holds the planner's "What it is told" table.
  - The convention `design/conventions/strategic-agents.md` and its template copy.
  - `template/CHANGELOG.md`.
  - `.claude/agents/planner.md` and its template copy.
  - `flaiover/src/routes/items/[id]/item.svelte.test.ts`.
- **From co-change (`flai touches suggest`) and the design together:** `docs/users/flaiover.md`, `docs/users/flai.md`, and `docs/operators/index.md`. These are the docs that describe Plan on an epic.
- **Added on the revisit, from the layout (T-1123):**
  - `flai/cmd/plan.go`: `flai plan`'s help says that an epic's planner drafts only the stories.
  - `flai/internal/mcpserver/plan.go`: the MCP tool `plan`'s description says the same.
  - `docs/users/flai-reference.md`: regenerated from that help. It is also in the co-change list.
  - `design/system/conventions.md`: the `plan` row of the roles table says the same.
  - The earlier plan said `flai/cmd` was unchanged. That held for how the planner starts, but not for the help text describing what it writes.
- **Overlap:** S-0261, in progress, claims the `flai/internal/mcpserver` folder and `docs/users/flai-reference.md`. T-1123 changes one constant in `plan.go` and regenerates the reference, so a conflict would be small and settled by regenerating.
- **Left out of the co-change list:** the host API, serve, the rest of the MCP server, `design/system/flai-cli.md`, and `design/system/flaiover-dashboard.md`. Nothing changes in how the planner is started, and their text on the command says only that it starts the planner for an epic or a story.

**Forecast: 45m, delivery 2026-10-06T23:50Z.**

- flai forecast now gives 29m: 83 s per unit of size over 25 large improvement stories, times size 21 (3 criteria, 18 touches). It gave 24m at size 17 before the revisit added four touches.
- The first plan raised flai's figure to 40m for three reasons:
  - two of the five tasks keep two copies identical;
  - a write under `.claude/` waits on the operator's permission thread;
  - the Go and flaiover test suites run at close-out.
- The revisit adds 5m for T-1123: three one-sentence edits and a regenerated reference.
- The story's agent started at 23:03Z. Layer 1 was done by 23:07Z, and T-1051 by 23:12Z. So 23:50Z holds unless the `.claude/` permission thread waits long.

**Cost of delay: 600 USD a week, as flai cod works it out.** That is the operator's 4h lost per 168h cycle at 150 USD an hour, one cycle a week. There is no reason to adjust it.

**Task layers.**

- **Layer 1, done:** T-1047 (design and convention), T-1048 (prompt), and T-1049 (item page test).
- **Layer 2:** T-1050 (agent definitions), T-1051 (docs, done), and T-1123 (help, MCP description, conventions design, and reference). Each waits for T-1047 so that its words follow the convention's. They share no path.
