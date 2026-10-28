---
id: D-0394
title: "Agents do not ask the operator to place stories in the pull order: the operator orders the board by dragging, and an agent pulls the story the operator names or the first ready one"
kind: decision
sources: [S-0058]
topics: [planning, orchestration, conventions, all]
paths: [wip/kanban/board.md, design/system/workflow.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0394 Agents do not ask the operator to place stories in the pull order: the operator orders the board by dragging, and an agent pulls the story the operator names or the first ready one

S-0058 records the operator's ruling of 2026-09-19: "stop asking them to place stories in the pull order; the board gives them no way to change it (S-0057). The agent pulls what they name or move to ready." S-0057 built dragging on the board, so the operator's placement is the order. A story's agent never reorders the list, and never asks the operator to; the orchestrator orders it only by policy with `order_ready` on, which keeps a placement by hand of the last day (ADR-0088).

Source: S-0058.
