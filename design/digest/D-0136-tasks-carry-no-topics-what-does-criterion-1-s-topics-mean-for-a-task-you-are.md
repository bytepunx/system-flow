---
id: D-0136
title: "Tasks carry no topics: what does criterion 1's \"topics\" mean for a task: you are correct, that is an error to have included them in the description (you may change it with an edit)"
kind: answer
sources: [TH-0098]
topics: [planner-agent, orchestration, planning, analysis, flai, cli, go, template, conventions, code, all]
paths: [flai/internal/harness, flai/internal/guard, flai/cmd/guard.go, flai/internal/itemnew, flai/internal/mcpserver/items_write_test.go, flai/internal/serve/plan_test.go, design/conventions, template, design/system/strategic-agents.md, docs/users/flai.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flaiover.md, flai/cmd/plan.go, design/system/work-hierarchy.md, design/system/workflow.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
date: 2026-10-04
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0136 Tasks carry no topics: what does criterion 1's "topics" mean for a task: you are correct, that is an error to have included them in the description (you may change it with an edit)

The operator, alex, answered TH-0098 at 2026-10-04T04:04:34Z.

What was asked:

```text
Criterion 1 asks for tasks with "work, done when, nature, tags, topics, touches, and after". flai keeps topics on stories and epics only: item_new documents topics as story/epic, and item_edit refuses them on a task ("topics belong to stories and epics", flai/internal/itemedit/itemedit.go). A story's topics are what select its prime pack, and a task is worked inside its story's pack, so a task's topics would select nothing.
```

The answer:

```text
you are correct, that is an error to have included them in the description (you may change it with an edit)
```

Source: TH-0098.
