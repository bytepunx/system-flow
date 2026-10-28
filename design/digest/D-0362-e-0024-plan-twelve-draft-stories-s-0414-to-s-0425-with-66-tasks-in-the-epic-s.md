---
id: D-0362
title: "E-0024 plan: twelve draft stories, S-0414 to S-0425, with 66 tasks, in the epic's order; values wait on TH-0471: Recommended: approve the plan, with your TH-0471 answers applied story by story…"
kind: answer
sources: [TH-0473]
topics: [cli, orchestration, planning, template, conventions, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/serve, flai/internal/guard, flai/internal/context, flai/internal/hostapi, flai/cmd, flaiover/src, design/system/strategic-agents.md, design/system/metrics.md, design/conventions/strategic-agents.md, ".claude/agents", docs/operators/index.md, docs/users/flai.md, template]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: [S-0414]
---

# D-0362 E-0024 plan: twelve draft stories, S-0414 to S-0425, with 66 tasks, in the epic's order; values wait on TH-0471: Recommended: approve the plan, with your TH-0471 answers applied story by story…

The operator, alex, answered TH-0473 at 2026-10-21T12:58:20Z.

What was asked:

```text
I drafted twelve stories for E-0024, all in the backlog as drafts. Each has its tasks, file-level touches, and a forecast. The cost of delay values wait on TH-0471, question 1. Every story carries its own docs and template changes, so there is no docs story. Each story's `### Planning` notes say where every touch came from and why each figure stands.
```

The answer:

```text
Recommended: approve the plan, with your TH-0471 answers applied story by story through the story planner. Take the planner's addition about the curator.

- Your TH-0471 answers change the drafts:
  - `design/workflows/roles` and `design/workflows/formulas`, in place of `design/roles` and `design/formulas`;
  - a formula instance gets a branch or worktree, as a story does (question 5);
  - a generated role file, with the Claude agent file symlinked where Claude allows it (question 6).
- I started the planner on S-0414 at 11:03Z under plan_backlog_stories, and S-0415 to S-0425 follow one at a time. Each run brings its story's touches, words, and tasks in line with those answers.
- Each story's value needs the 4h on E-0024 that you chose on TH-0471, written on the epic: `flai edit E-0024` with its cost of delay flags.
- The curator: TH-0476 recommends building it now on the analyzer machinery. So if E-0025 lands first, add one task to S-0414 and one to S-0419 to move it onto the role path, as this thread proposes.

Source: design/conventions/strategic-agents.md § As the planner
```

Source: TH-0473.
